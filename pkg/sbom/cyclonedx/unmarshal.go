package cyclonedx

import (
	"bytes"
	"sort"
	"strconv"
	"strings"

	"github.com/aquasecurity/trivy/pkg/log"

	cdx "github.com/CycloneDX/cyclonedx-go"
	"github.com/samber/lo"
	"golang.org/x/xerrors"

	ftypes "github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/purl"
	"github.com/aquasecurity/trivy/pkg/sbom"
)

type CycloneDX struct {
	*sbom.SBOM

	dependencies map[string][]string
	components   map[string]cdx.Component
}

func (c *CycloneDX) UnmarshalJSON(b []byte) error {
	log.Logger.Debug("Unmarshaling CycloneDX JSON...")
	if c.SBOM == nil {
		c.SBOM = &sbom.SBOM{}
	}
	bom := cdx.NewBOM()
	decoder := cdx.NewBOMDecoder(bytes.NewReader(b), cdx.BOMFileFormatJSON)
	if err := decoder.Decode(bom); err != nil {
		return xerrors.Errorf("CycloneDX decode error: %w", err)
	}

	c.dependencies = dependencyMap(bom.Dependencies)
	c.components = componentMap(bom.Metadata, bom.Components)

	var seen = make(map[string]struct{})
	for bomRef := range c.dependencies {
		component := c.components[bomRef]
		switch component.Type {
		case cdx.ComponentTypeOS: // OS info and OS packages
			c.OS = toOS(component)
			pkgInfo, err := c.parseOSPkgs(component, seen)
			if err != nil {
				return xerrors.Errorf("failed to parse os packages: %w", err)
			}
			c.Packages = append(c.Packages, pkgInfo)
		case cdx.ComponentTypeApplication: // It would be a lock file in a CycloneDX report generated by Trivy
			if lookupProperty(component.Properties, PropertyType) == "" {
				continue
			}
			app, err := c.parseLangPkgs(component, seen)
			if err != nil {
				return xerrors.Errorf("failed to parse language packages: %w", err)
			}
			c.Applications = append(c.Applications, *app)
		case cdx.ComponentTypeLibrary:
			// It is an individual package not associated with any lock files and should be processed later.
			// e.g. .gemspec, .egg and .wheel
			continue
		}
	}

	var libComponents []cdx.Component
	for ref, component := range c.components {
		if _, ok := seen[ref]; ok {
			continue
		}
		if component.Type == cdx.ComponentTypeLibrary {
			libComponents = append(libComponents, component)
		}
	}

	aggregatedApps, err := aggregateLangPkgs(libComponents)
	if err != nil {
		return xerrors.Errorf("failed to aggregate packages: %w", err)
	}
	c.Applications = append(c.Applications, aggregatedApps...)

	sort.Slice(c.Applications, func(i, j int) bool {
		if c.Applications[i].Type != c.Applications[j].Type {
			return c.Applications[i].Type < c.Applications[j].Type
		}
		return c.Applications[i].FilePath < c.Applications[j].FilePath
	})

	var metadata ftypes.Metadata
	if bom.Metadata != nil {
		metadata.Timestamp = bom.Metadata.Timestamp
		if bom.Metadata.Component != nil {
			metadata.Component = toTrivyCdxComponent(lo.FromPtr(bom.Metadata.Component))
		}
	}

	var components []ftypes.Component
	for _, component := range lo.FromPtr(bom.Components) {
		components = append(components, toTrivyCdxComponent(component))
	}

	// Keep the original SBOM
	c.CycloneDX = &ftypes.CycloneDX{
		BOMFormat:    bom.BOMFormat,
		SpecVersion:  bom.SpecVersion,
		SerialNumber: bom.SerialNumber,
		Version:      bom.Version,
		Metadata:     metadata,
		Components:   components,
	}
	return nil
}

func (c *CycloneDX) parseOSPkgs(component cdx.Component, seen map[string]struct{}) (ftypes.PackageInfo, error) {
	components := c.walkDependencies(component.BOMRef)
	pkgs, err := parsePkgs(components, seen)
	if err != nil {
		return ftypes.PackageInfo{}, xerrors.Errorf("failed to parse os package: %w", err)
	}

	return ftypes.PackageInfo{
		Packages: pkgs,
	}, nil
}

func (c *CycloneDX) parseLangPkgs(component cdx.Component, seen map[string]struct{}) (*ftypes.Application, error) {
	components := c.walkDependencies(component.BOMRef)
	components = lo.UniqBy(components, func(c cdx.Component) string {
		return c.BOMRef
	})

	app := toApplication(component)
	pkgs, err := parsePkgs(components, seen)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse language-specific packages: %w", err)
	}
	app.Libraries = pkgs

	return app, nil
}

func parsePkgs(components []cdx.Component, seen map[string]struct{}) ([]ftypes.Package, error) {
	var pkgs []ftypes.Package
	for _, com := range components {
		seen[com.BOMRef] = struct{}{}
		_, pkg, err := toPackage(com)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse language package: %w", err)
		}
		pkgs = append(pkgs, *pkg)
	}
	return pkgs, nil
}

// walkDependencies takes all nested dependencies of the root component.
//
// e.g. Library A, B, C, D and E will be returned as dependencies of Application 1.
// type: Application 1
//   - type: Library A
//     - type: Library B
//   - type: Application 2
//     - type: Library C
//     - type: Application 3
//       - type: Library D
//       - type: Library E
func (c *CycloneDX) walkDependencies(rootRef string) []cdx.Component {
	var components []cdx.Component
	for _, dep := range c.dependencies[rootRef] {
		component, ok := c.components[dep]
		if !ok {
			continue
		}

		// Take only 'Libraries'
		if component.Type == cdx.ComponentTypeLibrary {
			components = append(components, component)
		}

		components = append(components, c.walkDependencies(dep)...)
	}
	return components
}

func componentMap(metadata *cdx.Metadata, components *[]cdx.Component) map[string]cdx.Component {
	cmap := make(map[string]cdx.Component)

	for _, component := range lo.FromPtr(components) {
		cmap[component.BOMRef] = component
	}
	if metadata != nil && metadata.Component != nil {
		cmap[metadata.Component.BOMRef] = *metadata.Component
	}
	return cmap
}

func dependencyMap(deps *[]cdx.Dependency) map[string][]string {
	depMap := make(map[string][]string)

	for _, dep := range lo.FromPtr(deps) {
		if _, ok := depMap[dep.Ref]; ok {
			continue
		}

		var refs []string
		for _, d := range lo.FromPtr(dep.Dependencies) {
			refs = append(refs, d.Ref)
		}

		depMap[dep.Ref] = refs
	}
	return depMap
}

func aggregateLangPkgs(libs []cdx.Component) ([]ftypes.Application, error) {
	pkgMap := map[string][]ftypes.Package{}
	for _, lib := range libs {
		appType, pkg, err := toPackage(lib)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse purl to package: %w", err)
		}

		pkgMap[appType] = append(pkgMap[appType], *pkg)
	}

	var apps []ftypes.Application
	for appType, pkgs := range pkgMap {
		sort.Slice(pkgs, func(i, j int) bool {
			return pkgs[i].Name < pkgs[j].Name
		})
		apps = append(apps, ftypes.Application{
			Type:      appType,
			Libraries: pkgs,
		})
	}
	return apps, nil
}

func toOS(component cdx.Component) *ftypes.OS {
	return &ftypes.OS{
		Family: component.Name,
		Name:   component.Version,
	}
}

func toApplication(component cdx.Component) *ftypes.Application {
	return &ftypes.Application{
		Type:     lookupProperty(component.Properties, PropertyType),
		FilePath: component.Name,
	}
}

func toPackage(component cdx.Component) (string, *ftypes.Package, error) {
	p, err := purl.FromString(component.PackageURL)
	if err != nil {
		return "", nil, xerrors.Errorf("failed to parse purl: %w", err)
	}

	pkg := p.Package()
	pkg.Ref = component.BOMRef

	for _, license := range lo.FromPtr(component.Licenses) {
		pkg.Licenses = append(pkg.Licenses, license.Expression)
	}

	for _, prop := range lo.FromPtr(component.Properties) {
		if strings.HasPrefix(prop.Name, Namespace) {
			switch strings.TrimPrefix(prop.Name, Namespace) {
			case PropertySrcName:
				pkg.SrcName = prop.Value
			case PropertySrcVersion:
				pkg.SrcVersion = prop.Value
			case PropertySrcRelease:
				pkg.SrcRelease = prop.Value
			case PropertySrcEpoch:
				pkg.SrcEpoch, err = strconv.Atoi(prop.Value)
				if err != nil {
					return "", nil, xerrors.Errorf("failed to parse source epoch: %w", err)
				}
			case PropertyModularitylabel:
				pkg.Modularitylabel = prop.Value
			case PropertyLayerDiffID:
				pkg.Layer.DiffID = prop.Value
			}
		}
	}

	return p.AppType(), pkg, nil
}

func toTrivyCdxComponent(component cdx.Component) ftypes.Component {
	return ftypes.Component{
		BOMRef:     component.BOMRef,
		MIMEType:   component.MIMEType,
		Type:       ftypes.ComponentType(component.Type),
		Name:       component.Name,
		Version:    component.Version,
		PackageURL: component.PackageURL,
	}
}

func lookupProperty(properties *[]cdx.Property, key string) string {
	for _, p := range lo.FromPtr(properties) {
		if p.Name == Namespace+key {
			return p.Value
		}
	}
	return ""
}

package spdx

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/hashstructure/v2"
	"github.com/samber/lo"
	"github.com/spdx/tools-golang/spdx"
	"github.com/spdx/tools-golang/spdx/v2/common"
	"golang.org/x/exp/maps"
	"golang.org/x/xerrors"
	"k8s.io/utils/clock"

	"github.com/deepfactor-io/trivy/pkg/digest"
	ftypes "github.com/deepfactor-io/trivy/pkg/fanal/types"
	"github.com/deepfactor-io/trivy/pkg/licensing"
	"github.com/deepfactor-io/trivy/pkg/licensing/expression"
	"github.com/deepfactor-io/trivy/pkg/log"
	"github.com/deepfactor-io/trivy/pkg/purl"
	"github.com/deepfactor-io/trivy/pkg/scanner/utils"
	"github.com/deepfactor-io/trivy/pkg/types"
)

const (
	SPDXVersion         = "SPDX-2.2"
	DataLicense         = "CC0-1.0"
	SPDXIdentifier      = "DOCUMENT"
	DocumentNamespace   = "https://deepfactor.io"
	CreatorOrganization = "Deepfactor"
	CreatorTool         = "dfctl"
)

const (
	DocumentSPDXIdentifier = "DOCUMENT"
	noneField              = "NONE"
)

const (
	CategoryPackageManager = "PACKAGE-MANAGER"
	RefTypePurl            = "purl"

	PropertySchemaVersion = "SchemaVersion"

	NoAssertion = "NOASSERTION"

	// Image properties
	PropertySize       = "Size"
	PropertyImageID    = "ImageID"
	PropertyRepoDigest = "RepoDigest"
	PropertyDiffID     = "DiffID"
	PropertyRepoTag    = "RepoTag"

	// Package properties
	PropertyPkgID       = "PkgID"
	PropertyLayerDiffID = "LayerDiffID"
	PropertyLayerDigest = "LayerDigest"
	// Package Purpose fields
	PackagePurposeOS          = "OPERATING-SYSTEM"
	PackagePurposeContainer   = "CONTAINER"
	PackagePurposeSource      = "SOURCE"
	PackagePurposeApplication = "APPLICATION"
	PackagePurposeLibrary     = "LIBRARY"

	PackageSupplierNoAssertion  = "NOASSERTION"
	PackageSupplierOrganization = "Organization"

	RelationShipContains  = common.TypeRelationshipContains
	RelationShipDescribe  = common.TypeRelationshipDescribe
	RelationShipDependsOn = common.TypeRelationshipDependsOn

	ElementOperatingSystem = "OperatingSystem"
	ElementApplication     = "Application"
	ElementPackage         = "Package"
	ElementFile            = "File"
)

var (
	SourcePackagePrefix = "built package from"
)

type Marshaler struct {
	format     spdx.Document
	clock      clock.Clock
	newUUID    newUUID
	hasher     Hash
	appVersion string // Trivy version. It needed for `creator` field
}

type Hash func(v interface{}, format hashstructure.Format, opts *hashstructure.HashOptions) (uint64, error)

type newUUID func() uuid.UUID

type marshalOption func(*Marshaler)

func WithClock(clock clock.Clock) marshalOption {
	return func(opts *Marshaler) {
		opts.clock = clock
	}
}

func WithNewUUID(newUUID newUUID) marshalOption {
	return func(opts *Marshaler) {
		opts.newUUID = newUUID
	}
}

func WithHasher(hasher Hash) marshalOption {
	return func(opts *Marshaler) {
		opts.hasher = hasher
	}
}

func NewMarshaler(version string, opts ...marshalOption) *Marshaler {
	m := &Marshaler{
		format:     spdx.Document{},
		clock:      clock.RealClock{},
		newUUID:    uuid.New,
		hasher:     hashstructure.Hash,
		appVersion: version,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// The function augmentSpdxData updates each package in packages key,
// ensuring the spdx json is valid as per https://tools.spdx.org/app/validate/
// The following keys are being updated
//  1. licenseConcluded (incorrect delimiter and string value throws error)
//  2. licenseDeclared (incorrect delimiter and string value throws error)
//  3. copyrightText (throws a warning if the value is empty)
//  4. downloadLocation (throws a warning if the value is empty)
func augmentSpdxData(p *spdx.Package) {
	p.PackageLicenseConcluded = NoAssertion
	p.PackageLicenseDeclared = NoAssertion
	p.PackageCopyrightText = NoAssertion
	p.PackageDownloadLocation = NoAssertion
}

func (m *Marshaler) Marshal(r types.Report) (*spdx.Document, error) {
	var relationShips []*spdx.Relationship
	packages := make(map[spdx.ElementID]*spdx.Package)
	pkgDownloadLocation := getPackageDownloadLocation(r.ArtifactType, r.ArtifactName)

	// Root package contains OS, OS packages, language-specific packages and so on.
	rootPkg, err := m.rootPackage(r, pkgDownloadLocation)
	if err != nil {
		return nil, xerrors.Errorf("failed to generate a root package: %w", err)
	}
	packages[rootPkg.PackageSPDXIdentifier] = rootPkg
	relationShips = append(relationShips,
		relationShip(DocumentSPDXIdentifier, rootPkg.PackageSPDXIdentifier, RelationShipDescribe),
	)

	for _, result := range r.Results {
		parentPackage, err := m.resultToSpdxPackage(result, r.Metadata.OS, pkgDownloadLocation)
		if err != nil {
			return nil, xerrors.Errorf("failed to parse result: %w", err)
		}
		packages[parentPackage.PackageSPDXIdentifier] = &parentPackage
		relationShips = append(relationShips,
			relationShip(rootPkg.PackageSPDXIdentifier, parentPackage.PackageSPDXIdentifier, RelationShipContains),
		)

		for _, pkg := range result.Packages {
			spdxPackage, err := m.pkgToSpdxPackage(result.Type, pkgDownloadLocation, result.Class, r.Metadata, pkg, r.ArtifactType)
			if err != nil {
				return nil, xerrors.Errorf("failed to parse package: %w", err)
			}
			packages[spdxPackage.PackageSPDXIdentifier] = &spdxPackage
			relationShips = append(relationShips,
				relationShip(parentPackage.PackageSPDXIdentifier, spdxPackage.PackageSPDXIdentifier, RelationShipContains),
			)
		}
	}

	if len(relationShips) > 1 {
		// for consistent report generation accross UI and CLI
		// sort relationships except for the first item
		// the first relationship will reference the DOCUMENT (root-package/image)
		// eg: Relationship: SPDXRef-DOCUMENT DESCRIBES SPDXRef-ContainerImage-7eef1cebfe0e056b
		// and is expected to be at start
		from := 1
		sort.Slice(relationShips[from:], func(i, j int) bool {
			r1 := relationShips[i+from]
			r2 := relationShips[j+from]

			s1 := string(r1.RefA.ElementRefID) + r1.Relationship + string(r1.RefB.ElementRefID)
			s2 := string(r2.RefA.ElementRefID) + r2.Relationship + string(r2.RefB.ElementRefID)

			return s1 < s2
		})
	}

	// Augment SPDX data
	for _, val := range packages {
		augmentSpdxData(val)
	}

	return &spdx.Document{
		SPDXVersion:       spdx.Version,
		DataLicense:       spdx.DataLicense,
		SPDXIdentifier:    DocumentSPDXIdentifier,
		DocumentName:      r.ArtifactName,
		DocumentNamespace: getDocumentNamespace(r, m),
		CreationInfo: &spdx.CreationInfo{
			Creators: []common.Creator{
				{
					Creator:     CreatorOrganization,
					CreatorType: "Organization",
				},
				{
					Creator:     fmt.Sprintf("%s-%s", CreatorTool, m.appVersion),
					CreatorType: "Tool",
				},
			},
			Created: r.DfScanMeta.Created.UTC().Format(time.RFC3339),
		},
		Packages:      toPackages(packages),
		Relationships: relationShips,
	}, nil
}

func toPackages(packages map[spdx.ElementID]*spdx.Package) []*spdx.Package {
	ret := maps.Values(packages)
	sort.Slice(ret, func(i, j int) bool {
		if ret[i].PackageName != ret[j].PackageName {
			return ret[i].PackageName < ret[j].PackageName
		}
		return ret[i].PackageSPDXIdentifier < ret[j].PackageSPDXIdentifier
	})
	return ret
}

func (m *Marshaler) resultToSpdxPackage(result types.Result, os *ftypes.OS, pkgDownloadLocation string) (spdx.Package, error) {
	switch result.Class {
	case types.ClassOSPkg:
		osPkg, err := m.osPackage(os, pkgDownloadLocation)
		if err != nil {
			return spdx.Package{}, xerrors.Errorf("failed to parse operating system package: %w", err)
		}
		return osPkg, nil
	case types.ClassLangPkg:
		langPkg, err := m.langPackage(result.Target, result.Type, pkgDownloadLocation)
		if err != nil {
			return spdx.Package{}, xerrors.Errorf("failed to parse application package: %w", err)
		}
		return langPkg, nil
	default:
		// unsupported packages
		return spdx.Package{}, nil
	}
}

func (m *Marshaler) parseFile(filePath string, digest digest.Digest) (spdx.File, error) {
	pkgID, err := calcPkgID(m.hasher, filePath)
	if err != nil {
		return spdx.File{}, xerrors.Errorf("failed to get %s package ID: %w", filePath, err)
	}
	file := spdx.File{
		FileSPDXIdentifier: spdx.ElementID(fmt.Sprintf("File-%s", pkgID)),
		FileName:           filePath,
		Checksums:          digestToSpdxFileChecksum(digest),
	}
	return file, nil
}

func (m *Marshaler) rootPackage(r types.Report, pkgDownloadLocation string) (*spdx.Package, error) {
	var externalReferences []*spdx.PackageExternalReference
	attributionTexts := []string{attributionText(PropertySchemaVersion, strconv.Itoa(r.SchemaVersion))}

	// When the target is a container image, add PURL to the external references of the root package.
	if p, err := purl.NewPackageURL(purl.TypeOCI, r.Metadata, ftypes.Package{}); err != nil {
		return nil, xerrors.Errorf("failed to new package url for oci: %w", err)
	} else if p.Type != "" {
		externalReferences = append(externalReferences, purlExternalReference(p.ToString()))
	}

	if r.Metadata.ImageID != "" {
		attributionTexts = appendAttributionText(attributionTexts, PropertyImageID, r.Metadata.ImageID)
	}
	if r.Metadata.Size != 0 {
		attributionTexts = appendAttributionText(attributionTexts, PropertySize, strconv.FormatInt(r.Metadata.Size, 10))
	}

	for _, d := range r.Metadata.RepoDigests {
		attributionTexts = appendAttributionText(attributionTexts, PropertyRepoDigest, d)
	}

	// sort diffIDs for consistency
	sort.Slice(r.Metadata.DiffIDs, func(i, j int) bool {
		return r.Metadata.DiffIDs[i] < r.Metadata.DiffIDs[j]
	})

	for _, d := range r.Metadata.DiffIDs {
		attributionTexts = appendAttributionText(attributionTexts, PropertyDiffID, d)
	}
	for _, t := range r.Metadata.RepoTags {
		attributionTexts = appendAttributionText(attributionTexts, PropertyRepoTag, t)
	}

	pkgID, err := calcPkgID(m.hasher, fmt.Sprintf("%s-%s", r.ArtifactName, r.ArtifactType))
	if err != nil {
		return nil, xerrors.Errorf("failed to get %s package ID: %w", err)
	}

	pkgPurpose := PackagePurposeSource
	if r.ArtifactType == ftypes.ArtifactContainerImage {
		pkgPurpose = PackagePurposeContainer
	}

	return &spdx.Package{
		PackageName:               r.ArtifactName,
		PackageSPDXIdentifier:     elementID(camelCase(string(r.ArtifactType)), pkgID),
		PackageDownloadLocation:   pkgDownloadLocation,
		PackageAttributionTexts:   attributionTexts,
		PackageExternalReferences: externalReferences,
		PrimaryPackagePurpose:     pkgPurpose,
	}, nil
}

func (m *Marshaler) osPackage(osFound *ftypes.OS, pkgDownloadLocation string) (spdx.Package, error) {
	if osFound == nil {
		return spdx.Package{}, nil
	}

	pkgID, err := calcPkgID(m.hasher, osFound)
	if err != nil {
		return spdx.Package{}, xerrors.Errorf("failed to get os metadata package ID: %w", err)
	}

	return spdx.Package{
		PackageName:             osFound.Family,
		PackageVersion:          osFound.Name,
		PackageSPDXIdentifier:   elementID(ElementOperatingSystem, pkgID),
		PackageDownloadLocation: pkgDownloadLocation,
		PrimaryPackagePurpose:   PackagePurposeOS,
	}, nil
}

func (m *Marshaler) langPackage(target, appType, pkgDownloadLocation string) (spdx.Package, error) {
	pkgID, err := calcPkgID(m.hasher, fmt.Sprintf("%s-%s", target, appType))
	if err != nil {
		return spdx.Package{}, xerrors.Errorf("failed to get %s package ID: %w", target, err)
	}

	return spdx.Package{
		PackageName:             appType,
		PackageSourceInfo:       target, // TODO: Files seems better
		PackageSPDXIdentifier:   elementID(ElementApplication, pkgID),
		PackageDownloadLocation: pkgDownloadLocation,
		PrimaryPackagePurpose:   PackagePurposeApplication,
	}, nil
}

// Create a pkg object that will be common for cli and deepfactor portal
func createDFPkgObject(pkg ftypes.Package, artifactType ftypes.ArtifactType) ftypes.Package {
	pkgObj := ftypes.Package{
		ID:         pkg.ID,
		Arch:       pkg.Arch,
		Name:       pkg.Name,
		Version:    pkg.Version,
		SrcName:    pkg.SrcName,
		SrcVersion: pkg.SrcVersion,
		SrcRelease: pkg.SrcRelease,
		SrcEpoch:   pkg.SrcEpoch,
		Licenses:   pkg.Licenses,
		FilePath:   pkg.FilePath,
		Release:    pkg.Release,
		Ref:        pkg.Ref,
		Epoch:      pkg.Epoch,
		DependsOn:  pkg.DependsOn,
		Maintainer: pkg.Maintainer,
		// BuildInfo: pkg.BuildInfo,
		Modularitylabel: pkg.Modularitylabel,
		Indirect:        pkg.Indirect,
		// Locations:       pkg.Locations,
	}

	if artifactType == ftypes.ArtifactContainerImage {
		pkgObj.Layer = ftypes.Layer{
			Digest:    pkg.Layer.Digest,
			DiffID:    pkg.Layer.DiffID,
			CreatedBy: pkg.Layer.CreatedBy,
		}
	}

	return pkgObj
}

func (m *Marshaler) pkgToSpdxPackage(t, pkgDownloadLocation string, class types.ResultClass, metadata types.Metadata, pkg ftypes.Package, artifactType ftypes.ArtifactType) (spdx.Package, error) {
	license := GetLicense(pkg)

	// Create a pkg object that will be common for cli and deepfactor portal
	dfPkgObj := createDFPkgObject(pkg, artifactType)

	pkgID, err := calcPkgID(m.hasher, dfPkgObj)
	if err != nil {
		return spdx.Package{}, xerrors.Errorf("failed to get %s package ID: %w", pkg.Name, err)
	}

	var pkgSrcInfo string
	if class == types.ClassOSPkg && pkg.SrcName != "" {
		pkgSrcInfo = fmt.Sprintf("%s: %s %s", SourcePackagePrefix, pkg.SrcName, utils.FormatSrcVersion(pkg))
	}

	packageURL, err := purl.NewPackageURL(t, metadata, pkg)
	if err != nil {
		return spdx.Package{}, xerrors.Errorf("failed to parse purl (%s): %w", pkg.Name, err)
	}
	pkgExtRefs := []*spdx.PackageExternalReference{purlExternalReference(packageURL.String())}

	var attrTexts []string
	attrTexts = appendAttributionText(attrTexts, PropertyPkgID, pkg.ID)
	attrTexts = appendAttributionText(attrTexts, PropertyLayerDigest, pkg.Layer.Digest)
	attrTexts = appendAttributionText(attrTexts, PropertyLayerDiffID, pkg.Layer.DiffID)

	files, err := m.pkgFiles(pkg)
	if err != nil {
		return spdx.Package{}, xerrors.Errorf("package file error: %w", err)
	}

	supplier := &spdx.Supplier{Supplier: PackageSupplierNoAssertion}
	if pkg.Maintainer != "" {
		supplier = &spdx.Supplier{
			SupplierType: PackageSupplierOrganization, // Always use "Organization" at the moment as it is difficult to distinguish between "Person" or "Organization".
			Supplier:     pkg.Maintainer,
		}
	}
	return spdx.Package{
		PackageName:             pkg.Name,
		PackageVersion:          utils.FormatVersion(pkg),
		PackageSPDXIdentifier:   elementID(ElementPackage, pkgID),
		PackageDownloadLocation: pkgDownloadLocation,
		PackageSourceInfo:       pkgSrcInfo,

		// The Declared License is what the authors of a project believe govern the package
		PackageLicenseConcluded: license,

		// The Concluded License field is the license the SPDX file creator believes governs the package
		PackageLicenseDeclared: license,

		PackageExternalReferences: pkgExtRefs,
		PackageAttributionTexts:   attrTexts,
		PrimaryPackagePurpose:     PackagePurposeLibrary,
		PackageSupplier:           supplier,
		Files:                     files,
	}, nil
}

func (m *Marshaler) pkgFiles(pkg ftypes.Package) ([]*spdx.File, error) {
	if pkg.FilePath == "" {
		return nil, nil
	}

	file, err := m.parseFile(pkg.FilePath, pkg.Digest)
	if err != nil {
		return nil, xerrors.Errorf("failed to parse file: %w")
	}
	return []*spdx.File{
		&file,
	}, nil
}

func elementID(elementType, pkgID string) spdx.ElementID {
	return spdx.ElementID(fmt.Sprintf("%s-%s", elementType, pkgID))
}

func relationShip(refA, refB spdx.ElementID, operator string) *spdx.Relationship {
	ref := spdx.Relationship{
		RefA:         common.MakeDocElementID("", string(refA)),
		RefB:         common.MakeDocElementID("", string(refB)),
		Relationship: operator,
	}
	return &ref
}

func appendAttributionText(attributionTexts []string, key, value string) []string {
	if value == "" {
		return attributionTexts
	}
	return append(attributionTexts, attributionText(key, value))
}

func attributionText(key, value string) string {
	return fmt.Sprintf("%s: %s", key, value)
}

func purlExternalReference(packageURL string) *spdx.PackageExternalReference {
	return &spdx.PackageExternalReference{
		Category: CategoryPackageManager,
		RefType:  RefTypePurl,
		Locator:  packageURL,
	}
}

func GetLicense(p ftypes.Package) string {
	if len(p.Licenses) == 0 {
		return noneField
	}

	license := strings.Join(lo.Map(p.Licenses, func(license string, index int) string {
		// e.g. GPL-3.0-with-autoconf-exception
		license = strings.ReplaceAll(license, "-with-", " WITH ")
		license = strings.ReplaceAll(license, "-WITH-", " WITH ")

		return fmt.Sprintf("(%s)", license)
	}), " AND ")
	s, err := expression.Normalize(license, licensing.Normalize, expression.NormalizeForSPDX)
	if err != nil {
		// Not fail on the invalid license
		log.Logger.Warnf("Unable to marshal SPDX licenses %q", license)
		return ""
	}
	return s
}

func getDocumentNamespace(r types.Report, m *Marshaler) string {
	return fmt.Sprintf("%s/%s/%s-%s",
		DocumentNamespace,
		string(r.ArtifactType),
		strings.ReplaceAll(strings.ReplaceAll(r.ArtifactName, "https://", ""), "http://", ""), // remove http(s):// prefix when scanning repos
		r.DfScanMeta.ScanID, // overriden for consistency
	)
}

func calcPkgID(h Hash, v interface{}) (string, error) {
	f, err := h(v, hashstructure.FormatV2, &hashstructure.HashOptions{
		ZeroNil:      true,
		SlicesAsSets: true,
	})
	if err != nil {
		return "", xerrors.Errorf("could not build package ID for %+v: %w", v, err)
	}

	return fmt.Sprintf("%x", f), nil
}

func camelCase(inputUnderScoreStr string) (camelCase string) {
	isToUpper := false
	for k, v := range inputUnderScoreStr {
		if k == 0 {
			camelCase = strings.ToUpper(string(inputUnderScoreStr[0]))
		} else {
			if isToUpper {
				camelCase += strings.ToUpper(string(v))
				isToUpper = false
			} else {
				if v == '_' {
					isToUpper = true
				} else {
					camelCase += string(v)
				}
			}
		}
	}
	return
}

func getPackageDownloadLocation(t ftypes.ArtifactType, artifactName string) string {
	location := noneField
	// this field is used for git/mercurial/subversion/bazaar:
	// https://spdx.github.io/spdx-spec/v2.2.2/package-information/#77-package-download-location-field
	if t == ftypes.ArtifactRemoteRepository {
		// Trivy currently only supports git repositories. Format examples:
		// git+https://git.myproject.org/MyProject.git
		// git+http://git.myproject.org/MyProject
		location = fmt.Sprintf("git+%s", artifactName)
	}
	return location
}

func digestToSpdxFileChecksum(d digest.Digest) []common.Checksum {
	if d == "" {
		return nil
	}

	var alg spdx.ChecksumAlgorithm
	switch d.Algorithm() {
	case digest.SHA1:
		alg = spdx.SHA1
	case digest.SHA256:
		alg = spdx.SHA256
	default:
		return nil
	}

	return []spdx.Checksum{
		{
			Algorithm: alg,
			Value:     d.Encoded(),
		},
	}
}

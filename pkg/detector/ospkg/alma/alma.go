package alma

import (
	"strings"
	"time"

	version "github.com/knqyf263/go-rpm-version"
	"golang.org/x/xerrors"
	"k8s.io/utils/clock"

	ftypes "github.com/aquasecurity/fanal/types"
	dbTypes "github.com/aquasecurity/trivy-db/pkg/types"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/alma"
	"github.com/aquasecurity/trivy-db/pkg/vulnsrc/epel"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/scanner/utils"
	"github.com/aquasecurity/trivy/pkg/types"
)

var (
	eolDates = map[string]time.Time{
		// Source:
		// https://wiki.almalinux.org/FAQ.html#how-long-will-cloudlinux-support-almalinux
		"8": time.Date(2029, 12, 31, 23, 59, 59, 0, time.UTC),
	}
)

type options struct {
	clock clock.Clock
}

type option func(*options)

func WithClock(clock clock.Clock) option {
	return func(opts *options) {
		opts.clock = clock
	}
}

// Scanner implements the AlmaLinux scanner
type Scanner struct {
	osVS   alma.VulnSrc
	epelVS epel.VulnSrc
	*options
}

// NewScanner is the factory method for Scanner
func NewScanner(opts ...option) *Scanner {
	o := &options{
		clock: clock.RealClock{},
	}

	for _, opt := range opts {
		opt(o)
	}
	return &Scanner{
		osVS:    alma.NewVulnSrc(),
		epelVS:  epel.NewVulnSrc(),
		options: o,
	}
}

// Detect vulnerabilities in package using AlmaLinux scanner
func (s *Scanner) Detect(osVer string, pkgs []ftypes.Package) ([]types.DetectedVulnerability, error) {
	log.Logger.Info("Detecting AlmaLinux vulnerabilities...")
	if strings.Count(osVer, ".") > 0 {
		osVer = osVer[:strings.Index(osVer, ".")]
	}
	log.Logger.Debugf("AlmaLinux: os version: %s", osVer)
	log.Logger.Debugf("AlmaLinux: the number of packages: %d", len(pkgs))

	var vulns []types.DetectedVulnerability
	var skipPkgs []string
	for _, pkg := range pkgs {
		pkgName := addModularNamespace(pkg.Name, pkg.Modularitylabel)

		var advisories []dbTypes.Advisory
		var err error
		// https://docs.fedoraproject.org/en-US/epel/epel-faq/#how_can_i_find_out_if_a_package_is_from_epel
		if pkg.Vendor == "Fedora Project" {
			advisories, err = s.epelVS.Get(osVer, pkgName)
			if err != nil {
				return nil, xerrors.Errorf("failed to get EPEL advisories: %w", err)
			}
		} else {
			if strings.Contains(pkg.Release, ".module_el") {
				skipPkgs = append(skipPkgs, pkg.Name)
				continue
			}
			advisories, err = s.osVS.Get(osVer, pkgName)
			if err != nil {
				return nil, xerrors.Errorf("failed to get AlmaLinux advisories: %w", err)
			}
		}

		installed := utils.FormatVersion(pkg)
		installedVersion := version.NewVersion(installed)

		for _, adv := range advisories {
			fixedVersion := version.NewVersion(adv.FixedVersion)
			if installedVersion.LessThan(fixedVersion) {
				vuln := types.DetectedVulnerability{
					VulnerabilityID:  adv.VulnerabilityID,
					PkgName:          pkg.Name,
					InstalledVersion: installed,
					FixedVersion:     fixedVersion.String(),
					Layer:            pkg.Layer,
					DataSource:       adv.DataSource,
				}
				vulns = append(vulns, vuln)
			}
		}
	}
	if len(skipPkgs) > 0 {
		log.Logger.Infof("Skipped detection of these packages: %q because modular packages cannot be detected correctly due to a bug in AlmaLinux. See also: https://bugs.almalinux.org/view.php?id=173", skipPkgs)
	}

	return vulns, nil
}

// IsSupportedVersion checks the OSFamily can be scanned using AlmaLinux scanner
func (s *Scanner) IsSupportedVersion(osFamily, osVer string) bool {
	if strings.Count(osVer, ".") > 0 {
		osVer = osVer[:strings.Index(osVer, ".")]
	}

	eol, ok := eolDates[osVer]
	if !ok {
		log.Logger.Warnf("This OS version is not on the EOL list: %s %s", osFamily, osVer)
		return false
	}

	return s.clock.Now().Before(eol)
}

func addModularNamespace(name, label string) string {
	// e.g. npm, nodejs:12:8030020201124152102:229f0a1c => nodejs:12::npm
	var count int
	for i, r := range label {
		if r == ':' {
			count++
		}
		if count == 2 {
			return label[:i] + "::" + name
		}
	}
	return name
}

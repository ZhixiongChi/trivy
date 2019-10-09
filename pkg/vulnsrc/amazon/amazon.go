package amazon

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"go.uber.org/zap"

	"github.com/aquasecurity/trivy/pkg/db"
	"github.com/aquasecurity/trivy/pkg/utils"
	"github.com/aquasecurity/trivy/pkg/vulnsrc/vulnerability"
	"github.com/aquasecurity/vuln-list-update/amazon"
	bolt "github.com/etcd-io/bbolt"
	"golang.org/x/xerrors"
)

const (
	amazonDir      = "amazon"
	platformFormat = "amazon linux %s"
)

var (
	targetVersions = []string{"1", "2"}
	fileWalker     = utils.FileWalk // TODO: Remove once utils.go exposes an interface
)

type Operations interface {
	Update(string, map[string]struct{}) error
	Get(string, string) ([]vulnerability.Advisory, error)
}

type Config struct {
	lg  *zap.SugaredLogger
	dbc db.Operations
}

type alas struct {
	Version string
	amazon.ALAS
}

func (ac Config) Update(dir string, updatedFiles map[string]struct{}) error {
	rootDir := filepath.Join(dir, amazonDir)
	targets, err := utils.FilterTargets(amazonDir, updatedFiles)
	if err != nil {
		return xerrors.Errorf("failed to filter target files: %w", err)
	} else if len(targets) == 0 {
		ac.lg.Debug("amazon: no updated file")
		return nil
	}
	ac.lg.Debugf("Amazon Linux AMI Security Advisory updated files: %d", len(targets))

	bar := utils.PbStartNew(len(targets))
	defer bar.Finish()

	var alasList []alas
	err = fileWalker(rootDir, targets, func(r io.Reader, path string) error { // TODO: untested
		paths := strings.Split(path, string(filepath.Separator))
		if len(paths) < 2 {
			return nil
		}
		version := paths[len(paths)-2]
		if !utils.StringInSlice(version, targetVersions) {
			ac.lg.Debugf("unsupported amazon version: %s", version)
			return nil
		}

		var vuln amazon.ALAS
		if err = json.NewDecoder(r).Decode(&vuln); err != nil {
			return xerrors.Errorf("failed to decode amazon JSON: %w", err)
		}

		alasList = append(alasList, alas{
			Version: version,
			ALAS:    vuln,
		})
		bar.Increment()
		return nil
	})
	if err != nil {
		return xerrors.Errorf("error in amazon walk: %w", err)
	}

	if err = ac.save(alasList); err != nil {
		return xerrors.Errorf("error in amazon save: %w", err)
	}

	return nil
}

func (ac Config) save(alasList []alas) error {
	ac.lg.Debug("Saving amazon DB")
	err := ac.dbc.BatchUpdate(func(tx *bolt.Tx) error { // TODO: Untested
		for _, alas := range alasList {
			for _, cveID := range alas.CveIDs {
				for _, pkg := range alas.Packages {
					platformName := fmt.Sprintf(platformFormat, alas.Version)
					advisory := vulnerability.Advisory{
						VulnerabilityID: cveID,
						FixedVersion:    constructVersion(pkg.Epoch, pkg.Version, pkg.Release),
					}
					if err := ac.dbc.PutNestedBucket(tx, platformName, pkg.Name, cveID, advisory); err != nil {
						return xerrors.Errorf("failed to save amazon advisory: %w", err)
					}

					var references []string
					for _, ref := range alas.References {
						references = append(references, ref.Href)
					}

					vuln := vulnerability.Vulnerability{
						Severity:    severityFromPriority(alas.Severity),
						References:  references,
						Description: alas.Description,
						// TODO
						Title: "",
					}
					if err := vulnerability.Put(tx, cveID, vulnerability.Amazon, vuln); err != nil {
						return xerrors.Errorf("failed to save amazon vulnerability: %w", err)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return xerrors.Errorf("error in batch update: %w", err)
	}
	return nil
}

// Get returns a security advisory
func (ac Config) Get(version string, pkgName string) ([]vulnerability.Advisory, error) {
	bucket := fmt.Sprintf(platformFormat, version)
	advisories, err := db.ForEach(bucket, pkgName)
	if err != nil {
		return nil, xerrors.Errorf("error in amazon foreach: %w", err)
	}
	if len(advisories) == 0 {
		return nil, nil
	}

	var results []vulnerability.Advisory
	for _, v := range advisories {
		var advisory vulnerability.Advisory
		if err = json.Unmarshal(v, &advisory); err != nil {
			return nil, xerrors.Errorf("failed to unmarshal amazon JSON: %w", err)
		}
		results = append(results, advisory)
	}
	return results, nil
}

func severityFromPriority(priority string) vulnerability.Severity {
	switch priority {
	case "low":
		return vulnerability.SeverityLow
	case "medium":
		return vulnerability.SeverityMedium
	case "important":
		return vulnerability.SeverityHigh
	case "critical":
		return vulnerability.SeverityCritical
	default:
		return vulnerability.SeverityUnknown
	}
}

func constructVersion(epoch, version, release string) string {
	verStr := ""
	if epoch != "0" && epoch != "" {
		verStr += fmt.Sprintf("%s:", epoch)
	}
	verStr += version

	if release != "" {
		version += fmt.Sprintf("-%s", release)

	}
	return version
}

package report

import (
	"fmt"
	"io"

	"github.com/owenrumney/go-sarif/sarif"
)

// SarifWriter implements result Writer
type SarifWriter struct {
	Output  io.Writer
	Version string
	run     *sarif.Run
}

type sarifData struct {
	vulnerabilityId string
	title           string
	description     string
	severity        string
	pkgName         string
	fixedVersion    string
	url             string
	resourceType    string
	target          string
}

func (sw *SarifWriter) addSarifRule(data *sarifData) {
	description := data.description
	if description == "" {
		description = data.title
	}

	helpText := fmt.Sprintf("Vulnerability %v\\n%v\\nSeverity: %v\\nPackage: %v\\nFixed Version: %v\\nLink: [%v](%v)",
		data.vulnerabilityId, description, data.severity, data.pkgName, data.fixedVersion, data.vulnerabilityId, data.url)
	helpMarkdown := fmt.Sprintf("**Vulnerability %v**\n%v\n| Severity | Package | Fixed Version | Link |\n| --- | --- | --- | --- |\n|%v|%v|%v|[%v](%v)|\n",
		data.vulnerabilityId, description, data.severity, data.pkgName, data.fixedVersion, data.vulnerabilityId, data.url)

	sw.run.AddRule(data.vulnerabilityId).
		WithName(toSarifRuleName(data.resourceType)).
		WithDescription(data.vulnerabilityId).
		WithFullDescription(&sarif.MultiformatMessageString{Text: &description}).
		WithHelp(helpText).
		WithMarkdownHelp(helpMarkdown).
		WithProperties(sarif.Properties{
			"tags": []string{
				"vulnerability",
				data.severity,
			},
			"precision": "very-high",
		})
}

func (sw *SarifWriter) addSarifResult(data *sarifData) {
	sw.addSarifRule(data)

	message := sarif.NewTextMessage(data.description)
	region := sarif.NewSimpleRegion(1, 1)

	location := sarif.NewPhysicalLocation().
		WithArtifactLocation(sarif.NewSimpleArtifactLocation(data.target).WithUriBaseId("ROOTPATH")).
		WithRegion(region)

	sw.run.AddResult(data.vulnerabilityId).
		WithMessage(message).
		WithLevel(toSarifErrorLevel(data.severity)).
		WithLocation(sarif.NewLocation().WithPhysicalLocation(location))
}

func (sw SarifWriter) Write(report Report) error {
	sarifReport, err := sarif.New(sarif.Version210)
	if err != nil {
		return err
	}
	sw.run = sarif.NewRun("Trivy", "https://github.com/aquasecurity/trivy")
	sw.run.Tool.Driver.WithVersion(sw.Version)

	for _, res := range report.Results {
		for _, vuln := range res.Vulnerabilities {
			sw.addSarifResult(&sarifData{
				vulnerabilityId: vuln.VulnerabilityID,
				title:           vuln.Title,
				description:     vuln.Description,
				severity:        vuln.Severity,
				pkgName:         vuln.PkgName,
				fixedVersion:    vuln.FixedVersion,
				url:             vuln.PrimaryURL,
				resourceType:    res.Type,
				target:          res.Target,
			})
		}
		for _, misconf := range res.Misconfigurations {
			sw.addSarifResult(&sarifData{
				vulnerabilityId: misconf.ID,
				title:           misconf.Title,
				description:     misconf.Description,
				severity:        misconf.Severity,
				pkgName:         res.Type,
				fixedVersion:    "",
				url:             misconf.PrimaryURL,
				resourceType:    misconf.Type,
				target:          res.Target,
			})
		}
	}
	sarifReport.AddRun(sw.run)
	return sarifReport.PrettyWrite(sw.Output)
}

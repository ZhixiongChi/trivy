package k8s

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/aquasecurity/table"

	"github.com/liamg/tml"
)

type SummaryWriter struct {
	Output           io.Writer
	Severities       []string
	SeverityHeadings []string
}

func NewSummaryWriter(output io.Writer, severity string) SummaryWriter {
	var severities []string
	var severityHeadings []string
	severities, severityHeadings = getRequiredSeverities(severity)
	return SummaryWriter{
		Output:           output,
		Severities:       severities,
		SeverityHeadings: severityHeadings,
	}
}

// Write writes the results in a summarized table format
func (s SummaryWriter) Write(report Report) error {
	consolidated := report.consolidate()
	_, _ = fmt.Fprintln(s.Output)
	_, _ = fmt.Fprintf(s.Output, "Summary Report for %s\n", consolidated.ClusterName)

	t := table.New(s.Output)
	t.SetRowLines(false)
	configureHeader(s, t)

	sort.Slice(consolidated.Findings, func(i, j int) bool {
		return consolidated.Findings[i].Namespace > consolidated.Findings[j].Namespace
	})

	for _, finding := range consolidated.Findings {
		if !finding.Results.Failed() {
			continue
		}
		vCount, mCount, sCount := accumulateSeverityCounts(finding)
		name := fmt.Sprintf("%s/%s", finding.Kind, finding.Name)
		rowParts := []string{finding.Namespace, name}
		rowParts = append(rowParts, s.generateSummary(vCount)...)
		rowParts = append(rowParts, s.generateSummary(mCount)...)
		rowParts = append(rowParts, s.generateSummary(sCount)...)

		t.AddRow(rowParts...)
	}

	t.Render()

	keyParts := []string{"Severities:"}
	for _, s := range s.Severities {
		keyParts = append(keyParts, fmt.Sprintf("%s=%s", s[:1], colourSeverityValue(s, s)))
	}

	_, _ = fmt.Fprintln(s.Output, strings.Join(keyParts, " "))
	_, _ = fmt.Fprintln(s.Output)
	return nil
}

func (s SummaryWriter) generateSummary(sevCount map[string]int) []string {
	var parts []string

	for _, sev := range s.Severities {
		if count, ok := sevCount[sev]; ok {
			parts = append(parts, colourSeverityValue(strconv.Itoa(count), sev))
		} else {
			parts = append(parts, " ")
		}
	}

	return parts
}

func getRequiredSeverities(severity string) ([]string, []string) {
	requiredSevOrder := []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "UNKNOWN"}
	var severities []string
	var severityHeadings []string
	sevSplit := strings.Split(severity, ",")
	for _, sev := range requiredSevOrder {
		for _, p := range sevSplit {
			if p == sev {
				severities = append(severities, sev)
				severityHeadings = append(severityHeadings, strings.ToUpper(sev[:1]))
				continue
			}
		}
	}
	return severities, severityHeadings
}

func accumulateSeverityCounts(finding Resource) (map[string]int, map[string]int, map[string]int) {
	vCount := make(map[string]int)
	mCount := make(map[string]int)
	sCount := make(map[string]int)
	for _, r := range finding.Results {
		for _, rv := range r.Vulnerabilities {
			vCount[rv.Severity] = vCount[rv.Severity] + 1
		}
		for _, rv := range r.Misconfigurations {
			mCount[rv.Severity] = mCount[rv.Severity] + 1
		}
		for _, rv := range r.Secrets {
			sCount[rv.Severity] = sCount[rv.Severity] + 1
		}
	}
	return vCount, mCount, sCount
}

func configureHeader(s SummaryWriter, t *table.Table) {
	sevCount := len(s.Severities)

	headerRow := []string{"Namespace", "Resource"}
	//  vulnerabilities headings
	headerRow = append(headerRow, s.SeverityHeadings...)
	//  misconfig headings
	headerRow = append(headerRow, s.SeverityHeadings...)
	//  secrets headings
	headerRow = append(headerRow, s.SeverityHeadings...)
	headerAlignment := []table.Alignment{table.AlignLeft, table.AlignLeft}

	for i := 0; i < len(headerRow)-2; i++ {
		headerAlignment = append(headerAlignment, table.AlignCenter)
	}

	t.SetHeaders("Namespace", "Resource", "Vulnerabilities", "Misconfigurations", "Secrets")
	t.AddHeaders(headerRow...)
	t.SetAlignment(headerAlignment...)
	t.SetAutoMergeHeaders(true)
	t.SetHeaderColSpans(0, 1, 1, sevCount, sevCount, sevCount)
}

func colourSeverityValue(value string, severity string) string {
	switch severity {
	case "CRITICAL":
		return tml.Sprintf("<bold><red>%s</red></bold>", value)
	case "HIGH":
		return tml.Sprintf("<red>%s</red>", value)
	case "MEDIUM":
		return tml.Sprintf("<yellow>%s</yellow>", value)
	case "UNKNOWN":
		return tml.Sprintf("<blue>%s</blue>", value)
	default:
		return tml.Sprintf("%s", value)
	}
}

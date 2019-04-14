package nvd

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type Severity int

const (
	SeverityCritical Severity = iota
	SeverityHigh
	SeverityMedium
	SeverityLow
	SeverityUnknown
)

var (
	SeverityNames = []string{
		"CRITICAL",
		"HIGH",
		"MEDIUM",
		"LOW",
		"UNKNOWN",
	}
	SeverityColor = []func(a ...interface{}) string{
		color.New(color.FgRed).SprintFunc(),
		color.New(color.FgHiRed).SprintFunc(),
		color.New(color.FgYellow).SprintFunc(),
		color.New(color.FgBlue).SprintFunc(),
		color.New(color.FgCyan).SprintFunc(),
	}
)

func NewSeverity(severity string) (Severity, error) {
	for i, name := range SeverityNames {
		if severity == name {
			return Severity(i), nil
		}
	}
	return SeverityUnknown, fmt.Errorf("unknown severity: %s", severity)
}

func ColorizeSeverity(severity string) string {
	for i, name := range SeverityNames {
		if severity == name {
			return SeverityColor[i](severity)
		}
	}
	return color.New(color.FgBlue).SprintFunc()(severity)
}

func (s Severity) String() string {
	return SeverityNames[s]
}

type LastUpdated struct {
	Date time.Time
}

type NVD struct {
	CVEItems []Item `json:"CVE_Items"`
}

type Item struct {
	Cve    Cve
	Impact Impact
}

type Cve struct {
	Meta Meta `json:"CVE_data_meta"`
}

type Meta struct {
	ID string
}

type Impact struct {
	BaseMetricV2 BaseMetricV2
	BaseMetricV3 BaseMetricV3
}

type BaseMetricV2 struct {
	Severity string
}

type BaseMetricV3 struct {
	CvssV3 CvssV3
}

type CvssV3 struct {
	BaseSeverity string
}

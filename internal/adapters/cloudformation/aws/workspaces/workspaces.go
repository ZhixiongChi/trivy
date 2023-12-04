package workspaces

import (
	"github.com/aquasecurity/defsec/pkg/providers/aws/workspaces"
	"github.com/aquasecurity/trivy/pkg/iac/scanners/cloudformation/parser"
)

// Adapt adapts workspaces resources
func Adapt(cfFile parser.FileContext) workspaces.WorkSpaces {
	return workspaces.WorkSpaces{
		WorkSpaces: getWorkSpaces(cfFile),
	}
}
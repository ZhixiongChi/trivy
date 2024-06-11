package sbt

import (
	"context"
	"golang.org/x/xerrors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aquasecurity/trivy/pkg/dependency/parser/sbt/lockfile"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/analyzer/language"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
	"github.com/aquasecurity/trivy/pkg/utils/fsutils"
)

func init() {
	analyzer.RegisterPostAnalyzer(analyzer.TypeSbtLock, newSbtDependencyLockAnalyzer)
}

const (
	version      = 1
	requiredFile = "build.sbt.lock"
)

// sbtDependencyLockAnalyzer analyzes '*.sbt.lock'
type sbtDependencyLockAnalyzer struct {
	parser language.Parser
}

func newSbtDependencyLockAnalyzer(_ analyzer.AnalyzerOptions) (analyzer.PostAnalyzer, error) {
	return &sbtDependencyLockAnalyzer{
		parser: lockfile.NewParser(),
	}, nil
}

func (a sbtDependencyLockAnalyzer) PostAnalyze(_ context.Context, input analyzer.PostAnalysisInput) (*analyzer.AnalysisResult, error) {
	required := func(path string, d fs.DirEntry) bool {
		return a.Required(path, nil)
	}

	var apps []types.Application
	var err error
	err = fsutils.WalkDir(input.FS, ".", required, func(filePath string, _ fs.DirEntry, r io.Reader) error {
		var app *types.Application
		app, err = language.Parse(types.Sbt, filePath, r, a.parser)
		if err != nil {
			return xerrors.Errorf("%s parse error: %w", filePath, err)
		}

		if app != nil {

			apps = append(apps, *app)
		}

		return nil
	})

	if err != nil {
		return nil, xerrors.Errorf("sbt walk error: %w", err)
	}

	return &analyzer.AnalysisResult{
		Applications: apps,
	}, nil
}

func (a sbtDependencyLockAnalyzer) Required(filePath string, _ os.FileInfo) bool {
	return requiredFile == filepath.Base(filePath)
}

func (a sbtDependencyLockAnalyzer) Type() analyzer.Type {
	return analyzer.TypeSbtLock
}

func (a sbtDependencyLockAnalyzer) Version() int {
	return version
}

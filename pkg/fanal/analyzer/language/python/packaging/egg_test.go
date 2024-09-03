package packaging

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/aquasecurity/trivy/pkg/fanal/analyzer"
	"github.com/aquasecurity/trivy/pkg/fanal/types"
)

func Test_eggAnalyzer_Analyze(t *testing.T) {
	tests := []struct {
		name            string
		dir             string
		includeChecksum bool
		want            *analyzer.AnalysisResult
		wantErr         string
	}{
		{
			name: "egg zip",
			dir:  "testdata/egg-zip",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "kitchen-1.2.6-py2.7.egg",
						Packages: types.Packages{
							{
								Name:    "kitchen",
								Version: "1.2.6",
								Licenses: []string{
									"GNU Library or Lesser General Public License (LGPL)",
								},
								FilePath: "kitchen-1.2.6-py2.7.egg",
							},
						},
					},
				},
			},
		},
		{
			name:            "egg zip with checksum",
			dir:             "testdata/egg-zip",
			includeChecksum: true,
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "kitchen-1.2.6-py2.7.egg",
						Packages: types.Packages{
							{
								Name:    "kitchen",
								Version: "1.2.6",
								Licenses: []string{
									"GNU Library or Lesser General Public License (LGPL)",
								},
								FilePath: "kitchen-1.2.6-py2.7.egg",
								Digest:   "sha1:4e13b6e379966771e896ee43cf8e240bf6083dca",
							},
						},
					},
				},
			},
		},
		{
			name: "egg zip with license file",
			dir:  "testdata/egg-zip-with-license-file",
			want: &analyzer.AnalysisResult{
				Applications: []types.Application{
					{
						Type:     types.PythonPkg,
						FilePath: "sample_package.egg",
						Packages: types.Packages{
							{
								Name:    "sample_package",
								Version: "0.1",
								Licenses: []string{
									"MIT",
								},
								FilePath: "sample_package.egg",
							},
						},
					},
				},
			},
		},
		{
			name: "egg zip doesn't contain required files",
			dir:  "testdata/no-req-files",
			want: &analyzer.AnalysisResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a, err := newEggAnalyzer(analyzer.AnalyzerOptions{})
			require.NoError(t, err)
			got, err := a.PostAnalyze(context.Background(), analyzer.PostAnalysisInput{
				FS: os.DirFS(tt.dir),
				Options: analyzer.AnalysisOptions{
					FileChecksum: tt.includeChecksum,
				},
			})

			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}

}

func Test_eggAnalyzer_Required(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		want     bool
	}{
		{
			name:     "egg zip",
			filePath: "python2.7/site-packages/cssutils-1.0-py2.7.egg",
			want:     true,
		},
		{
			name:     "egg-info PKG-INFO",
			filePath: "python3.8/site-packages/wrapt-1.12.1.egg-info/PKG-INFO",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := eggAnalyzer{}
			got := a.Required(tt.filePath, nil)
			assert.Equal(t, tt.want, got)
		})
	}
}

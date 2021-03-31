// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package artifact

import (
	"context"
	"github.com/aquasecurity/fanal/analyzer"
	"github.com/aquasecurity/fanal/applier"
	image2 "github.com/aquasecurity/fanal/artifact/image"
	local2 "github.com/aquasecurity/fanal/artifact/local"
	"github.com/aquasecurity/fanal/artifact/remote"
	"github.com/aquasecurity/fanal/cache"
	"github.com/aquasecurity/fanal/image"
	"github.com/aquasecurity/trivy-db/pkg/db"
	"github.com/aquasecurity/trivy/pkg/detector/ospkg"
	"github.com/aquasecurity/trivy/pkg/scanner"
	"github.com/aquasecurity/trivy/pkg/scanner/local"
	"github.com/aquasecurity/trivy/pkg/types"
	"github.com/aquasecurity/trivy/pkg/vulnerability"
	"time"
)

// Injectors from inject.go:

func initializeDockerScanner(ctx context.Context, imageName string, artifactCache cache.ArtifactCache, localArtifactCache cache.LocalArtifactCache, timeout time.Duration, disableAnalyzers []analyzer.Type) (scanner.Scanner, func(), error) {
	applierApplier := applier.NewApplier(localArtifactCache)
	detector := ospkg.Detector{}
	localScanner := local.NewScanner(applierApplier, detector)
	dockerOption, err := types.GetDockerOption(timeout)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	imageImage, cleanup, err := image.NewDockerImage(ctx, imageName, dockerOption)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	artifact := image2.NewArtifact(imageImage, artifactCache, disableAnalyzers)
	scannerScanner := scanner.NewScanner(localScanner, artifact)
	return scannerScanner, func() {
		cleanup()
	}, nil
}

func initializeArchiveScanner(ctx context.Context, filePath string, artifactCache cache.ArtifactCache, localArtifactCache cache.LocalArtifactCache, timeout time.Duration, disableAnalyzers []analyzer.Type) (scanner.Scanner, error) {
	applierApplier := applier.NewApplier(localArtifactCache)
	detector := ospkg.Detector{}
	localScanner := local.NewScanner(applierApplier, detector)
	imageImage, err := image.NewArchiveImage(filePath)
	if err != nil {
		return scanner.Scanner{}, err
	}
	artifact := image2.NewArtifact(imageImage, artifactCache, disableAnalyzers)
	scannerScanner := scanner.NewScanner(localScanner, artifact)
	return scannerScanner, nil
}

func initializeFilesystemScanner(ctx context.Context, dir string, artifactCache cache.ArtifactCache, localArtifactCache cache.LocalArtifactCache, disableAnalyzers []analyzer.Type) (scanner.Scanner, func(), error) {
	applierApplier := applier.NewApplier(localArtifactCache)
	detector := ospkg.Detector{}
	localScanner := local.NewScanner(applierApplier, detector)
	artifact := local2.NewArtifact(dir, artifactCache, disableAnalyzers)
	scannerScanner := scanner.NewScanner(localScanner, artifact)
	return scannerScanner, func() {
	}, nil
}

func initializeRepositoryScanner(ctx context.Context, url string, artifactCache cache.ArtifactCache, localArtifactCache cache.LocalArtifactCache, disableAnalyzers []analyzer.Type) (scanner.Scanner, func(), error) {
	applierApplier := applier.NewApplier(localArtifactCache)
	detector := ospkg.Detector{}
	localScanner := local.NewScanner(applierApplier, detector)
	remoteRemote, err := types.GetGitOption(url)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	artifact, cleanup, err := remote.NewArtifact(remoteRemote, artifactCache, disableAnalyzers)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	scannerScanner := scanner.NewScanner(localScanner, artifact)
	return scannerScanner, func() {
		cleanup()
	}, nil
}

func initializeVulnerabilityClient() vulnerability.Client {
	config := db.Config{}
	client := vulnerability.NewClient(config)
	return client
}

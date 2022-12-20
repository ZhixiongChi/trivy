package walker

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

const (
	opq string = ".wh..wh..opq"
	wh  string = ".wh."
)

type LayerTar struct {
	walker
	threshold int64
}

func NewLayerTar(skipFiles, skipDirs []string, slow bool) LayerTar {
	threshold := defaultSizeThreshold
	if slow {
		threshold = slowSizeThreshold
	}

	return LayerTar{
		walker:    newWalker(skipFiles, skipDirs, slow),
		threshold: threshold,
	}
}

func (w LayerTar) Walk(layer io.Reader, analyzeFn WalkFunc) ([]string, []string, error) {

	var opqDirs, whFiles, skipDirs []string
	tr := tar.NewReader(layer)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, nil, xerrors.Errorf("failed to extract the archive: %w", err)
		}

		filePath := filepath.ToSlash(hdr.Name)
		filePath = strings.TrimLeft(filepath.Clean(filePath), "/")
		fileDir, fileName := filepath.Split(filePath)

		fileDir = filepath.ToSlash(fileDir)
		filePath = filepath.ToSlash(filePath)

		// e.g. etc/.wh..wh..opq
		if opq == fileName {
			opqDirs = append(opqDirs, fileDir)
			continue
		}
		// etc/.wh.hostname
		if strings.HasPrefix(fileName, wh) {
			name := strings.TrimPrefix(fileName, wh)
			fpath := filepath.Join(fileDir, name)
			whFiles = append(whFiles, fpath)
			continue
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if w.shouldSkipDir(filePath) {
				skipDirs = append(skipDirs, filePath)
				continue
			}
		case tar.TypeReg:
			if w.shouldSkipFile(filePath) {
				continue
			}
		// symlinks and hardlinks have no content in reader, skip them
		default:
			continue
		}

		if underSkippedDir(filePath, skipDirs) {
			continue
		}

		// A symbolic/hard link or regular file will reach here.
		if err = w.processFile(filePath, tr, hdr.FileInfo(), analyzeFn); err != nil {
			return nil, nil, xerrors.Errorf("failed to process the file: %w", err)
		}
	}
	return opqDirs, whFiles, nil
}

func (w LayerTar) processFile(filePath string, tr *tar.Reader, fi fs.FileInfo, analyzeFn WalkFunc) error {
	cf := newCachedFile(fi.Size(), tr, w.threshold)
	defer func() {
		// nolint
		_ = cf.Clean()
	}()

	if err := analyzeFn(filePath, fi, cf.Open); err != nil {
		return xerrors.Errorf("failed to analyze file: %w", err)
	}

	return nil
}

func underSkippedDir(filePath string, skipDirs []string) bool {
	for _, skipDir := range skipDirs {
		rel, err := filepath.Rel(skipDir, filePath)
		if err != nil {
			return false
		}
		if !strings.HasPrefix(rel, fmt.Sprintf("..%c", os.PathSeparator)) {
			return true
		}
	}
	return false
}

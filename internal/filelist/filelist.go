package filelist

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type FileMetadata struct {
	ModifiedTime int64  // Unix timestamp in nanoseconds
	Size         int64  // Size of the file in bytes
	Hash         string // SHA-512 hash of the file content
}

type FileList map[string]*FileMetadata

func LoadFileList(filelistPath string) (FileList, error) {
	bt, err := os.ReadFile(filepath.Clean(filelistPath))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // File does not exist, return empty FileList
		}
		return nil, err
	}

	fileList := make(FileList)
	if err := yaml.Unmarshal(bt, &fileList); err != nil {
		return nil, err
	}

	return fileList, nil
}

func (fl FileList) Save(filelistPath string) error {
	bt, err := yaml.Marshal(fl)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Clean(filelistPath), bt, 0644); err != nil {
		return err
	}

	return nil
}

func (fl FileList) Diff(base FileList) (updated, removed FileList) {
	updated = make(FileList)
	removed = make(FileList)

	for relPath, meta := range fl {
		if otherMeta, exists := base[relPath]; !exists || *meta != *otherMeta {
			updated[relPath] = meta
		}
	}

	for relPath, meta := range base {
		if _, exists := fl[relPath]; !exists {
			removed[relPath] = meta
		}
	}

	return updated, removed
}

func (fl FileList) Equal(other FileList) bool {
	if len(fl) != len(other) {
		return false
	}

	for path, meta := range fl {
		if otherMeta, exists := other[path]; !exists || *meta != *otherMeta {
			return false
		}
	}

	return true
}

func MakeFileList(localDir string, filePatterns []*regexp.Regexp) (FileList, error) {
	localDir, err := filepath.Abs(localDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for local directory '%s': %w", localDir, err)
	}

	fileList := make(FileList)

	err = filepath.Walk(localDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		for _, pattern := range filePatterns {
			relpath, err := filepath.Rel(localDir, p)
			if err != nil {
				return fmt.Errorf("failed to get relative path for '%s': %w", p, err)
			}

			if pattern.MatchString(strings.ToLower(relpath)) {
				h, err := calculateFileHash(p)
				if err != nil {
					return fmt.Errorf("failed to calculate hash for file '%s': %w", p, err)
				}

				fileList[relpath] = &FileMetadata{
					ModifiedTime: info.ModTime().UnixNano(),
					Size:         info.Size(),
					Hash:         h,
				}

				return nil // Stop checking other patterns for this file
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk through directory '%s': %w", localDir, err)
	}

	return fileList, nil
}

func calculateFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file '%s': %w", filePath, err)
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("failed to calculate hash for file '%s': %w", filePath, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

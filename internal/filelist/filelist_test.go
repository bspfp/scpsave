package filelist_test

import (
	"fmt"
	"path/filepath"
	"regexp"
	"scpsave/internal/filelist"
	"scpsave/internal/scp"
	"testing"
)

const (
	testfilelistpath = `./testfilelist.yaml`
)

func TestFileList(t *testing.T) {
	localdir, err := filepath.Abs(`..`)
	if err != nil {
		t.Fatalf("Failed to get absolute path for local directory: %v", err)
	}
	fl, err := filelist.MakeFileList(localdir, []*regexp.Regexp{regexp.MustCompile(`.+\.go`)})
	if err != nil {
		t.Fatalf("MakeFileList failed: %v", err)
	}
	for path, meta := range fl {
		fmt.Printf("%s %d %d %s\n", path, meta.ModifiedTime, meta.Size, meta.Hash)
	}
	if err := fl.Save(testfilelistpath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loadedFl, err := filelist.LoadFileList(testfilelistpath)
	if err != nil {
		t.Fatalf("LoadFileList failed: %v", err)
	}
	for path, meta := range loadedFl {
		fmt.Printf("%s %d %d %s\n", path, meta.ModifiedTime, meta.Size, meta.Hash)
	}

	u, r := loadedFl.Diff(fl)
	if len(u) > 0 || len(r) > 0 {
		t.Fatalf("Diff failed: added=%d, removed=%d", len(u), len(r))
	}
	t.Log("Diff passed: no changes detected between original and loaded file lists")

	_ = scp.DeleteLocalFile(testfilelistpath)
}

package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Copier struct {
	source string
	target string
}

func NewCopier(source, target string) *Copier {
	// TODO: Check whether source and target exist

	return &Copier{
		source: filepath.Clean(source),
		target: filepath.Clean(target),
	}
}

// Copy copies a file if it is more recent than the version in the target directory. The files path
// is computed relatively from the source dir and then applied to the target dir in order to mirror the directory
// hierarchy.
// Return values are whether the file has is already present in a (more) recent version, the amount of bytes copied and
// whether an error occurred.
func (c *Copier) Copy(path string, info os.FileInfo) (bool, int, error) {
	// TODO: check that source file is more recent than destination file
	// TODO: copy file over

	return true, 0, nil
}

func (c *Copier) Move(oldPath, newPath string, info os.FileInfo) (bool, int, error) {

}

func (c *Copier) Rename(oldPath, newPath string) error {
	oldPath, err := c.resolveToSourcePath(oldPath)

	newPath, err := c.resolveToSourcePath(newPath)
	if err != nil {
		return err
	}


}

func (c *Copier) resolveToSourcePath(p string) (string, error) {
	p = filepath.Clean(p)

	// Is p a file located somewhere in source?
	if !strings.HasPrefix(p, c.source) {
		return "", fmt.Errorf("file not below source directory")
	}

	relPathToSource, err := filepath.Rel(c.source, oath)
	if err != nil {
		return "", err
	}

	return path.Join(c.target, relPathToSource), nil
}

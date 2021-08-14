package internal

import "os"

type Copier struct {
	source string
	target string
}

func NewCopier(source, target string) *Copier {
	return &Copier{
		source: source,
		target: target,
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

}

package internal

import "os"

type Copier struct {

}

func (c *Copier) Copy(path string, info os.FileInfo) error {
	// TODO: check that source file is more recent than destination file
	// TODO: copy file over

	return nil
}

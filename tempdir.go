package pathlib

import (
	"os"
)

// NewTempDir creates a unique temporary directory in the
// users home directory (~/tmp). The directory is created,
// but the user is responsible for clean up. `namePattern`
// is a string that effects the nameing of the temp directory.
// An empty pattern string will produce a random name.
// A string such as "temp" with create a name like "tempZd93lsk".
// The inclusion of an astrix liek "*temp" will produce
// a name like "83dwilsCtemp".
func NewTempDir(namePattern string) (Path, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return Path{}, err
	}

	fileInfo, err := os.Stat(dir)
	if err != nil {
		return Path{}, err
	}

	homeDir := NewPathWithMode(dir, fileInfo.Mode())
	tmpDir := homeDir.Join("tmp")
	err = tmpDir.MkDirs()
	if err != nil {
		return Path{}, err
	}

	newTmpDir, err := os.MkdirTemp(tmpDir.String(), namePattern)
	if err != nil {
		return Path{}, err
	}

	tmpdir := NewPathWithMode(newTmpDir, fileInfo.Mode())

	err = tmpdir.MkDirs()
	if err != nil {
		return Path{}, err
	}

	return tmpdir, nil
}

type TempDirCleanupFunc func()

func emptyCleanupFunc() {
}

// NewTempDirWithCleanup returns a new temporary directory along with
// a cleanup function. Example:
//
//	pth, clean, err := NewTempDirWithCleanup("")
//	if err != nil {
//	    return err
//	}
//	defer clean()
func NewTempDirWithCleanup(namePattern string) (Path, TempDirCleanupFunc, error) {
	tmpDir, err := NewTempDir(namePattern)
	if err != nil {
		return tmpDir, emptyCleanupFunc, err
	}
	cleanup := func() {
		tmpDir.Remove()
	}
	return tmpDir, cleanup, err
}

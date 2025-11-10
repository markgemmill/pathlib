// pathlib.go contain Path methods that are purely path struct
// operations: manipulating the path, fileMode or err attributes.
// All operation return a new Path object.
package pathlib

import (
	"os"
	"path/filepath"
	"strings"
)

const version = "0.2.0-dev.0"

type jsonPathInfo struct {
	Path     string      `json:"path"`
	FileMode os.FileMode `json:"mode"`
	Err      error       `json:"err"`
}

type Path struct {
	path     string
	fileMode os.FileMode
	err      error
}

// Error returns the Path's error value or nil.
func (pth Path) Error() error {
	return pth.err
}

// HasError returns true if there is an error on the Path object.
func (pth Path) HasError() bool {
	return pth.err != nil
}

// FileMode returns the current Path's FileMode
func (pth Path) FileMode() os.FileMode {
	return pth.fileMode
}

// SetFileMode creates a copy of the current Path object with the given
// os.FileMode. This operation does NOT change the mode of the file on
// the operating system if the file exists. To effect a change in the
// use either Path.Touch() or Path.Chmod() methods.
func (pth Path) SetFileMode(mode os.FileMode) Path {
	p := NewPathWithMode(pth.path, mode)
	p.err = pth.err
	return p
}

// SetPath creates a copy of the current Path object with the given path string.
func (pth Path) SetPath(path string) Path {
	p := pth.Copy()
	p.path = path
	return p
}

// SetErr creates a copy of the current Path object with the given error.
func (pth Path) SetErr(err error) Path {
	p := pth.Copy()
	p.err = err
	return p
}

// Copy creates a copy of the current Path object.
func (pth Path) Copy() Path {
	p := NewPathWithMode(pth.path, pth.fileMode)
	p.err = pth.err
	return p
}

// Parent returns the parent of the current path object.
func (pth Path) Parent() Path {
	return pth.SetPath(filepath.Dir(pth.path))
}

// Resolve expands relative paths to their absolute path string.
func (pth Path) Resolve() Path {
	if pth.HasError() {
		return pth
	}
	// resolve ~
	if strings.HasPrefix(pth.path, "~") {
		home := Home()
		newPath := strings.TrimLeft(pth.path, "~")
		newPath = strings.TrimLeft(newPath, string(filepath.Separator))

		fullPath := home.Join(newPath)
		fullPath.SetFileMode(pth.fileMode)

		return fullPath
	}

	absPth, err := filepath.Abs(pth.path)
	if err != nil {
		pth.err = err
		return pth
	}
	return pth.SetErr(err).SetPath(absPth)
}

// Suffix returns the extension of the final path segment. Examples:
//
//	/some/path/filename.txt   -> .txt
//	/some/path/zipfile.tar.gz -> .gz
//	/some/path/dirname        ->
func (pth Path) Suffix() string {
	return filepath.Ext(pth.path)
}

// Name returns the final path segment. Examples:
//
//	/some/path/filename.txt   -> filename.txt
//	/some/path/zipfile.tar.gz -> zipfile.tar.gz
//	/some/path/dirname        -> dirname
func (pth Path) Name() string {
	return filepath.Base(pth.path)
}

// Stem returns the file segment less any extension. Examples:
//
//	/some/path/filename.txt   -> filename
//	/some/path/zipfile.tar.gz -> zipfile.tar
//	/some/path/dirname        -> dirname
func (pth Path) Stem() string {
	segments := strings.Split(filepath.Base(pth.path), ".")
	t := len(segments)
	switch t {
	case 1:
		return segments[0]
	case 2:
		return segments[0]
	default:
		return strings.Join(segments[0:t-2], ".")
	}
}

// String returns the string representation of the Path object.
func (pth Path) String() string {
	return pth.path
}

// Equal compares the path strings of two Path object.
func (pth Path) Equal(other Path) bool {
	return pth == other
}

// Join returns a copy of current path with the additional path segments appended to it.
func (pth Path) Join(paths ...string) Path {
	return pth.SetPath(filepath.Join(pth.path, filepath.Join(paths...)))
}

func (pth Path) Split() []string {
	return strings.Split(pth.String(), string(filepath.Separator))
}

// Check takes one or more Path objects and returns
// the error of the first path that contains one.
func Check(paths ...Path) error {
	for _, pth := range paths {
		if pth.err != nil {
			return pth.err
		}
	}
	return nil
}

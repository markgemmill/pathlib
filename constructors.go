package pathlib

import (
	"os"
)

// NewPathWithMode creates a new Path object from the given
// path string and file mode.
func NewPathWithMode(path string, fileMode os.FileMode) Path {
	// TODO: for relative paths, should we be stripping any ./ references?
	//    for example filepath.Join("./here", "there") returns
	//    "here/there", stripping off the "./". Maybe we should do the same?
	return Path{
		path:     path,
		fileMode: fileMode,
	}
}

// NewFilePath creates a new Path object with default permission
// mode of 0644. There is no check that the path is actually
// a file instead of a directory. The assumption is the file
// may or may not already exist.
func NewFilePath(path string) Path {
	return NewPathWithMode(path, READONLY_FILE)
}

// NewDirPath creates a new Path object with default permission
// mode of 0755. There is no check that the path is actually
// a directory instead of a file. The assumption is the directory
// may or may not already exist.
func NewDirPath(path string) Path {
	return NewPathWithMode(path, READONLY_DIR)
}

// Home returns the current user's home directory as a Path.
// If home errors for any reason, the Path is returned with
// a default file mode of 0700 and the error.  Otherwise
// the Path is returned with whatever the current users
// home directory permissions are.
func Home() Path {
	pth, err := os.UserHomeDir()
	if err != nil {
		home := NewPathWithMode(pth, PRIVATE_DIR)
		home.err = err
		return home
	}

	stat, err := os.Stat(pth)
	if err != nil {
		home := NewPathWithMode(pth, PRIVATE_DIR)
		home.err = err
		return home
	}

	return NewPathWithMode(pth, stat.Mode())
}

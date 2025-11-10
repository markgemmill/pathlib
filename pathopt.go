package pathlib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func (pth Path) Chown(usr *user.User) error {
	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		return err
	}

	// change ownership of html dir to caddy
	err = os.Chown(pth.String(), uid, gid)
	if err != nil {
		return err
	}

	return nil
}

func (pth Path) ChownTree(usr *user.User) error {
	var err error
	if pth.IsFile() {
		err = pth.Chown(usr)
		if err != nil {
			return err
		}
		return nil
	}

	err = pth.Chown(usr)
	if err != nil {
		return err
	}

	content, err := pth.ReadDir()
	if err != nil {
		return err
	}

	for _, path := range content {
		err = path.ChownTree(usr)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pth Path) CopyTo(dstPath Path) (Path, error) {
	src, err := pth.Open()
	if err != nil {
		return dstPath, err
	}
	defer src.Close()

	dst, err := dstPath.Open()
	if err != nil {
		return dstPath, err
	}
	defer dst.Close()

	// do the copy
	_, err = io.Copy(dst, src)
	if err != nil {
		return dstPath, err
	}

	err = dst.Sync()
	if err != nil {
		return dstPath, err
	}

	return dstPath, nil
}

// MoveTo is the same as Rename, except the provided Path must be a directory, and
// we first attempt to create the directories.
func (pth Path) MoveTo(directory Path) (Path, error) {
	err := directory.MkDirs()
	if err != nil {
		return pth, err
	}

	newPath := directory.Join(pth.Name())
	return pth.Rename(newPath)
}

// Rename moves the current path to the given path and returns the new Path object.
// On error, the returned path is the current Path object.
func (pth Path) Rename(newPth Path) (Path, error) {
	err := os.Rename(pth.String(), newPth.String())
	if err != nil {
		return pth, err
	}
	return newPth, nil
}

func (pth *Path) UnmarshalJSON(b []byte) error {
	pthInfo := new(jsonPathInfo)
	err := json.Unmarshal(b, &pthInfo)
	if err != nil {
		return err
	}
	pth.path = pthInfo.Path
	pth.fileMode = pthInfo.FileMode
	pth.err = pthInfo.Err
	return nil
}

func (pth Path) MarshalJSON() ([]byte, error) {
	pthInfo := jsonPathInfo{
		Path:     pth.path,
		FileMode: pth.fileMode,
		Err:      pth.err,
	}
	return json.Marshal(pthInfo)
}

func (pth Path) MkDir() error {
	return os.Mkdir(pth.path, pth.fileMode)
}

func (pth Path) MkDirs() error {
	return os.MkdirAll(pth.path, pth.fileMode)
}

func (pth Path) Read() ([]byte, error) {
	return os.ReadFile(pth.path)
}

func (pth Path) Write(data []byte) error {
	return os.WriteFile(pth.path, data, pth.fileMode)
}

func (pth Path) Touch() error {
	return os.WriteFile(pth.path, []byte(""), pth.fileMode)
}

func (pth Path) Open() (*os.File, error) {
	if pth.Exists() {
		return os.Open(pth.String())
	}
	return os.Create(pth.String())
}

func (pth Path) Remove() error {
	return os.RemoveAll(pth.path)
}

func (pth Path) RelativeTo(oth Path) (Path, error) {
	relativePth, found := strings.CutPrefix(pth.String(), oth.String())
	if !found {
		return pth, fmt.Errorf("%s is not relative to %s", pth.String(), oth.String())
	}
	relativePth, _ = strings.CutPrefix(relativePth, string(filepath.Separator))

	return NewPathWithMode(relativePth, pth.fileMode), nil
}

func (pth Path) ReadDir(filters ...PathFilterFunc) ([]Path, error) {
	dirToRead := pth.path
	if pth.IsFile() {
		dirToRead = pth.Parent().path
	}
	content, err := os.ReadDir(dirToRead)
	if err != nil {
		return nil, err
	}
	paths := []Path{}
	directory := pth.SetPath(dirToRead)
	for _, info := range content {
		p := directory.Join(info.Name())
		if ApplyPathFilters(p, filters...) {
			paths = append(paths, p)
		}
	}
	return paths, nil
}

type PathFilterFunc func(Path) bool

func ApplyPathFilters(pth Path, filters ...PathFilterFunc) bool {
	for _, filter := range filters {
		if !filter(pth) {
			return false
		}
	}
	return true
}

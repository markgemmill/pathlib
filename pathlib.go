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
	"time"
)

const version = "0.1.1"

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

func NewPath(path string, fileMode os.FileMode) Path {
	return Path{
		path:     path,
		fileMode: fileMode,
	}
}

func Home() Path {
	pth, err := os.UserHomeDir()
	home := NewPath(pth, PRIVATE_DIR)
	if err != nil {
		home.err = err
	}
	return home
}

func Check(paths ...Path) error {
	for _, pth := range paths {
		if pth.err != nil {
			return pth.err
		}
	}
	return nil
}

func NewTempDir(namePattern string) (Path, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return Path{}, err
	}

	fileInfo, err := os.Stat(dir)
	if err != nil {
		return Path{}, err
	}

	homeDir := NewPath(dir, fileInfo.Mode())
	tmpDir := homeDir.Join("tmp")
	err = tmpDir.MkDirs()
	if err != nil {
		return Path{}, err
	}

	newTmpDir, err := os.MkdirTemp(tmpDir.String(), namePattern)
	if err != nil {
		return Path{}, err
	}

	tmpdir := NewPath(newTmpDir, fileInfo.Mode())

	err = tmpdir.MkDirs()
	if err != nil {
		return Path{}, err
	}

	return tmpdir, nil
}

func (pth Path) SetFileMode(mode os.FileMode) Path {
	p := NewPath(pth.path, mode)
	p.err = pth.err
	return p
}

func (pth Path) Owner() (*user.User, error) {
	return GetOwner(pth.String())
}

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

func (pth Path) Copy() Path {
	p := NewPath(pth.path, pth.fileMode)
	p.err = pth.err
	return p
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

func (pth Path) SetPath(path string) Path {
	p := pth.Copy()
	p.path = path
	return p
}

func (pth Path) SetErr(err error) Path {
	p := pth.Copy()
	p.err = err
	return p
}

func (pth Path) HasError() bool {
	return pth.err != nil
}

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

func (pth Path) Suffix() string {
	return filepath.Ext(pth.path)
}

func (pth Path) Name() string {
	return filepath.Base(pth.path)
}

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

func (pth Path) String() string {
	return pth.path
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

func (pth Path) Exists() bool {
	_, err := pth.Stat()
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (pth Path) Equal(other Path) bool {
	return pth.String() == other.String()
}

func (pth Path) IsDir() bool {
	stats, err := os.Stat(pth.path)
	if err != nil {
		return false
	}
	return stats.IsDir()
}

func (pth Path) IsFile() bool {
	return !pth.IsDir()
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

func (pth Path) Join(paths ...string) Path {
	return pth.SetPath(filepath.Join(pth.path, filepath.Join(paths...)))
}

func (pth Path) Parent() Path {
	return pth.SetPath(filepath.Dir(pth.path))
}

func (pth Path) Remove() error {
	return os.RemoveAll(pth.path)
}

func (pth Path) Stat() (os.FileInfo, error) {
	return os.Stat(pth.path)
}

func (pth Path) ModTime() (time.Time, error) {
	stat, err := pth.Stat()
	if err != nil {
		return time.Now(), err
	}
	return stat.ModTime(), nil
}

func (pth Path) RelativeTo(oth Path) (Path, error) {
	relativePth, found := strings.CutPrefix(pth.String(), oth.String())
	if !found {
		return pth, fmt.Errorf("%s is not relative to %s", pth.String(), oth.String())
	}
	relativePth, _ = strings.CutPrefix(relativePth, string(filepath.Separator))

	return NewPath(relativePth, pth.fileMode), nil
}

func (pth Path) Split() []string {
	return strings.Split(pth.String(), string(filepath.Separator))
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

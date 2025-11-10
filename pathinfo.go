package pathlib

import (
	"os"
	"os/user"
	"time"
)

func (pth Path) Owner() (*user.User, error) {
	return GetOwner(pth.String())
}

func (pth Path) Stat() (os.FileInfo, error) {
	return os.Stat(pth.path)
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

func (pth Path) ModTime() (time.Time, error) {
	stat, err := pth.Stat()
	if err != nil {
		return time.Now(), err
	}
	return stat.ModTime(), nil
}

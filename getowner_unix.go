//go:build darwin || linux

package pathlib

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// GetOwner returns the owner of the specified file or directory
func GetOwner(filepath string) (*user.User, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", filepath, err)
	}
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to get syscall.Stat_t for %s", filepath)
	}
	uid := strconv.Itoa(int(stat.Uid))
	return user.LookupId(uid)
}

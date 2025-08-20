package pathlib

import (
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"syscall"
)

var (
	PRIVATE_FILE = os.FileMode(0600)
	PRIVATE_EXE  = os.FileMode(0700)
	PRIVATE_DIR  = os.FileMode(0700)

	READONLY_FILE = os.FileMode(0644)
	READONLY_DIR  = os.FileMode(0755)

	PUBLIC_FILE = os.FileMode(0666)
	PUBLIC_DIR  = os.FileMode(0777)
)

// GetOwner returns the owner of the specified file or directory
func GetOwner(filepath string) (*user.User, error) {
	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info for %s: %w", filepath, err)
	}

	// Get the underlying syscall.Stat_t
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("failed to get syscall.Stat_t for %s", filepath)
	}

	switch runtime.GOOS {
	case "windows":
		return getWindowsFileOwner(filepath)
	case "linux", "darwin", "freebsd", "openbsd", "netbsd":
		return getUnixFileOwner(stat)
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// getUnixFileOwner handles Unix-like systems (Linux, macOS, BSD variants)
func getUnixFileOwner(stat *syscall.Stat_t) (*user.User, error) {
	uid := strconv.Itoa(int(stat.Uid))
	return user.LookupId(uid)
}

// getWindowsFileOwner handles Windows systems
func getWindowsFileOwner(filepath string) (*user.User, error) {
	// TODO: implement this using Win32 API calls using ffi

	return nil, fmt.Errorf("windows get owner api is not implemented")
}

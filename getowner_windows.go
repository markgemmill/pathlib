//go:build windows

package pathlib

import (
	"fmt"
	"os/user"
)

// GetOwner returns the owner of the specified file or directory
func GetOwner(filepath string) (*user.User, error) {
	return nil, fmt.Errorf("not implemented on Windows yet")
}

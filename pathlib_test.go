package pathlib

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	cwd  = ""
	home = ""
)

func init() {
	home, _ = os.UserHomeDir()
	cwd, _ = os.Getwd()
}

func TestPathResolveCwd(t *testing.T) {
	pth := NewPath("./here", PUBLIC_DIR)

	assert.Equal(t, "./here", pth.path)

	pth = pth.Resolve()

	assert.Equal(t, cwd+"/here", pth.path)
}

func TestPathResolveHome(t *testing.T) {
	pth := NewPath("~/here", PUBLIC_DIR)

	assert.Equal(t, "~/here", pth.path)

	pth = pth.Resolve()

	assert.Equal(t, home+"/here", pth.path)
}

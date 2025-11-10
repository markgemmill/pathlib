//go:build linux || darwin
// +build linux darwin

package pathlib

import (
	"fmt"
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

func TestNewPath(t *testing.T) {
	dir := NewDirPath("./here")
	file := NewDirPath("./here/file.txt")

	t.Run("directory path methods and attributes", func(t *testing.T) {
		assert.Nil(t, dir.Error())
		assert.False(t, dir.HasError())
		assert.Equal(t, dir.FileMode(), READONLY_DIR)
		// TODO: for relative paths, should we be stripping any ./ references?
		//    for example filepath.Join("./here", "there") returns
		//    "here/there", stripping off the "./". Maybe we should do the same?
		assert.Equal(t, "./here", dir.String())
		assert.Equal(t, "here", dir.Name())
		assert.Equal(t, "here", dir.Stem())
		assert.Equal(t, "", dir.Suffix())
		assert.Equal(t, []string{".", "here"}, dir.Split())

		parent := dir.Parent()
		assert.Equal(t, ".", parent.String())

		extended := dir.Join("and", "there")
		assert.Equal(t, "here/and/there", extended.String())
	})

	t.Run("file path methods and attributes", func(t *testing.T) {
		assert.Nil(t, file.Error())
		assert.False(t, file.HasError())
		assert.Equal(t, file.FileMode(), READONLY_DIR)
		assert.Equal(t, "./here/file.txt", file.String())
		assert.Equal(t, "file.txt", file.Name())
		assert.Equal(t, "file", file.Stem())
		assert.Equal(t, ".txt", file.Suffix())
		assert.Equal(t, []string{".", "here", "file.txt"}, file.Split())

		parent := file.Parent()
		assert.Equal(t, "here", parent.String())

		extended := file.Join("and", "there")
		// TODO: this should be an error condition
		assert.Equal(t, "here/file.txt/and/there", extended.String())
	})

	t.Run("path manipulation and equality", func(t *testing.T) {
		copy := file.Copy()
		assert.Equal(t, file.String(), copy.String())
		assert.True(t, file.Equal(copy))
		assert.True(t, file == copy)

		newMode := file.SetFileMode(PUBLIC_FILE)
		assert.Equal(t, file.String(), newMode.String())
		assert.False(t, file.Equal(newMode))

		newPath := file.SetPath("./there/file.txt")
		assert.NotEqual(t, file.String(), newPath.String())
		assert.False(t, file.Equal(newPath))

		newErr := file.SetErr(fmt.Errorf("file error"))
		assert.Equal(t, file.String(), newErr.String())
		assert.False(t, file.Equal(newErr))
	})
}

func TestPathManipulationAndEquality(t *testing.T) {
	file := NewDirPath("./here/file.txt")
	copy := file.Copy()
	assert.Equal(t, file.String(), copy.String())
	assert.True(t, file.Equal(copy))
	assert.True(t, file == copy)

	newMode := file.SetFileMode(PUBLIC_FILE)
	assert.Equal(t, file.String(), newMode.String())
	assert.False(t, file.Equal(newMode))

	newPath := file.SetPath("./there/file.txt")
	assert.NotEqual(t, file.String(), newPath.String())
	assert.False(t, file.Equal(newPath))

	newErr := file.SetErr(fmt.Errorf("file error"))
	assert.Equal(t, file.String(), newErr.String())
	assert.False(t, file.Equal(newErr))
}

func TestPathResolveCwd(t *testing.T) {
	pth := NewPathWithMode("./here", PUBLIC_DIR)

	assert.Equal(t, "./here", pth.path)

	pth = pth.Resolve()

	assert.Equal(t, cwd+"/here", pth.path)
}

func TestPathResolveHome(t *testing.T) {
	pth := NewPathWithMode("~/here", PUBLIC_DIR)

	assert.Equal(t, "~/here", pth.path)

	pth = pth.Resolve()

	assert.Equal(t, home+"/here", pth.path)
}

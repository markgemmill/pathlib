package pathlib

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTempDir(t *testing.T) {
	tmpdir, err := NewTempDir("mytemp")
	defer tmpdir.Remove()
	assert.Nil(t, err)
	assert.True(t, tmpdir.Exists())
	assert.True(t, strings.HasPrefix(tmpdir.Name(), "mytemp"))
}

func TestNewTempDirWithSuffixedName(t *testing.T) {
	tmpdir, err := NewTempDir("*mytemp")
	defer tmpdir.Remove()
	assert.Nil(t, err)
	assert.True(t, tmpdir.Exists())
	assert.True(t, strings.HasSuffix(tmpdir.Name(), "mytemp"))
}

func TestNewTempDirWithCleanup(t *testing.T) {
	tmpdir, clean, err := NewTempDirWithCleanup("mytemp")
	assert.Nil(t, err)
	assert.True(t, tmpdir.Exists())
	clean()
	assert.False(t, tmpdir.Exists())
}

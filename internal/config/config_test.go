package config

import (
	"io/fs"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetConfig(t *testing.T) {
	want := Config{
		SkipPaths: []string{".ccls", "node-modules"},
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "test")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	err = ioutil.WriteFile(
		tmpFile.Name(),
		[]byte(`skip_paths = [ ".ccls", "node-modules" ]`),
		fs.ModePerm,
	)
	require.NoError(t, err)

	os.Setenv(envConfigPath, tmpFile.Name())
	assert.EqualValues(t, want, GetConfig())
}

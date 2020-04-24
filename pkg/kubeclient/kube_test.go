package kubeclient

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fileExists(t *testing.T) {
	file := path.Join(".", "test_data", "test.txt")

	exist, err := fileExists(file)
	assert.NoError(t, err, "should not error when looking for the file")
	assert.True(t, exist, "should have found the file")
}

func Test_fileNoExists(t *testing.T) {
	file := path.Join(".", "test_data", "dummy")

	exist, err := fileExists(file)
	assert.NoError(t, err, "should not error when looking for the file")
	assert.False(t, exist, "should not have found the file")
}

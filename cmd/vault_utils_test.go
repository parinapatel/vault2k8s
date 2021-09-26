package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestconvertPath  calls convertPath with a path, checking
// for a valid return value.
func TestConvertPath(t *testing.T) {
	vaultPath := "shared/myTeam/secret/logging/certs"
	want := "shared/data/myTeam/secret/logging/certs"
	msg, err := convertPath(vaultPath)
	assert.NoError(t, err, "vaultPath should not result in errors.")
	assert.Equal(t, want, msg, "The two strings should be the same.")
	vaultPath = "myTeam/secret/logging/certs"
	msg, err = convertPath(vaultPath)
	assert.Error(t, err, "vault path should begin with shared")
}

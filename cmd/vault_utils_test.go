package cmd

import (
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

// TestconvertPath  calls convertPath with a path, checking
// for a valid return value.
func TestConvertPath(t *testing.T) {
	vaultPath := "shared/myTeam/secret/logging/certs"
	want := "shared/data/myTeam/secret/logging/certs"
	msg, err := convertPath(vaultPath)
	assert.NoError(t, err, "vaultPath should not result in errors.")
	assert.Equal(t, want, msg, "The two strings should be the same.")
}

/* HOWTO figure out later
func createTestVault(t *testing.T) (net.Listener, *api.Client) {
	t.Helper()

	// Create an in-memory, unsealed core (the "backend", if you will).
	core, keyShares, rootToken := vault.TestCoreUnsealed(t)
	_ = keyShares

	// Start an HTTP server for the core.
	ln, addr := http.TestServer(t, core)

	// Create a client that talks to the server, initially authenticating with
	// the root token.
	conf := api.DefaultConfig()
	conf.Address = addr

	client, err := api.NewClient(conf)
	if err != nil {
		t.Fatal(err)
	}
	client.SetToken(rootToken)

	// Setup required secrets, policies, etc.
	_, err = client.Logical().Write("shared/foo", map[string]interface{}{
		"secret": "bar",
	})
	if err != nil {
		t.Fatal(err)
	}

	return ln, client
}


func TestVaultStuff(t *testing.T) {
	ln, client := createTestVault(t)
	defer ln.Close()

	// Pass the client to the code under test.
	data  :=getVaultData("shared/foo",client)
	print(data)
}
*/

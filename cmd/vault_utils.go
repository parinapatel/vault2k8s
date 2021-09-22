package cmd

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/sirupsen/logrus"
	"strings"
)

func convertPath(vaultPath string) (string, error) {
	fixPath := ""
	if strings.HasPrefix(vaultPath, "shared") {
		c := strings.Split(vaultPath, "/")
		c = append(c[:2], c[1:]...)
		c[1] = "data"
		fixPath = strings.Join(c, "/")
	} else {
		return fixPath, errors.New("vault path should begin with shared")
	}
	return fixPath, nil
}

func createVaultClient() (*api.Client, error) {
	client, err := api.NewClient(nil)
	// Get the token if it came in from the environment
	vaultToken := client.Token()

	// If we don't have a token, check the token helper
	if vaultToken == "" {
		helper, err := token.NewInternalTokenHelper()
		if err != nil {
			logrus.Errorf("failed to get token helper, error : %s", err)
			return nil, err
		}
		vaultToken, err = helper.Get()
		if err != nil {
			logrus.Errorf("failed to get token from token helper : %s", err)
			return nil, err
		}
	}
	client.SetToken(vaultToken)
	if err != nil {
		logrus.Errorf("failed to set token: %s", err)
		return nil, err
	}
	return client, nil
}

func getVaultData(vaultPath string, client *api.Client) map[string]interface{} {

	vaultPath, err := convertPath(vaultPath)
	if err != nil {
		logrus.Errorf("failed to get vaultPath: %s", err)
		return nil
	}
	secret, err := client.Logical().Read(vaultPath)
	if err != nil {
		logrus.Errorf("failed to get vaultPath: %s", err)
		return nil
	}

	//logrus.Info(secret.Data["data"].(map[string]interface{}))
	return secret.Data["data"].(map[string]interface{})
}

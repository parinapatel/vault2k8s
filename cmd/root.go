package cmd

import (
	"errors"
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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

func getVaultDate(vaultPath string) map[string]interface{} {
	client, err := api.NewClient(nil)
	// Get the token if it came in from the environment
	vaultToken := client.Token()

	// If we don't have a token, check the token helper
	if vaultToken == "" {
		helper, err := token.NewInternalTokenHelper()
		if err != nil {
			logrus.Errorf("failed to get token helper, error : %s", err)
			return nil
		}
		vaultToken, err = helper.Get()
		if err != nil {
			logrus.Errorf("failed to get token from token helper : %s", err)
			return nil
		}
	}
	client.SetToken(vaultToken)
	if err != nil {
		logrus.Errorf("failed to parse verbosity: %s", err)
	}
	vaultPath, err = convertPath(vaultPath)
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

func NewRootCmd() *cobra.Command {
	result := &cobra.Command{
		Use:   "vault2k8s",
		Short: "Read Vault Path and convert to k8s secrets",

		PreRunE: func(_ *cobra.Command, _ []string) error {
			lvl, err := logrus.ParseLevel(viper.GetString("verbosity"))
			if err != nil {
				return fmt.Errorf("failed to parse verbosity: %w", err)
			}

			logrus.SetLevel(lvl)
			return nil
		},
		Run: func(C *cobra.Command, _ []string) {
			logrus.Debugf(viper.GetString("secretName"))
			logrus.Debugf(viper.GetString("vaultPath"))
			if viper.GetString("vaultPath") == "" {
				logrus.Errorf("Missing flags vaultPath")
				os.Exit(1)
			}
			if viper.GetString("secretName") == "" {
				logrus.Errorf("Missing flags secretName")
				os.Exit(1)

			}
			secretMap := getVaultDate(viper.GetString("vaultPath"))
			if secretMap == nil {
				logrus.Errorf("No data recieved from vault")
			}
			output, err := convertK8sSecret(secretMap, viper.GetString("secretName"), viper.GetString("namespace"), viper.GetBool("forceDecode"))
			if err != nil {
				panic(err)
			}
			fmt.Printf("%s", output)
		},
	}

	flags := result.PersistentFlags()

	flags.String("verbosity", "info", "The level to log at [trace, debug, info, warning, error, fatal, panic]")
	flags.StringP("vaultPath", "p", "", "Vault Path. ( Required ) ")
	flags.Bool("forceDecode", false, "Force to decode String b64, By default it will try to encode string to base64")

	flags.StringP("secretName", "s", "", "K8s Secret Name ( Required )")
	flags.StringP("namespace", "n", "istio-system", "K8s Namespace secret is deployed in")
	result.MarkFlagRequired("vaultPath")
	result.MarkFlagRequired("secretName")

	_ = viper.BindPFlags(flags)
	return result
}

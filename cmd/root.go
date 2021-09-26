package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
			client, err := createVaultClient()
			if err != nil {
				panic(err)
			}
			secretMap := getVaultData(viper.GetString("vaultPath"), client)
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

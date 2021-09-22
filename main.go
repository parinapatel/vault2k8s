package main

import (
	"os"
	"vault2k8s/cmd"

	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {
	logrus.SetFormatter(&prefixed.TextFormatter{
		ForceFormatting: true,
		FullTimestamp:   true,
	})

	if err := cmd.NewRootCmd().Execute(); err != nil {
		logrus.Errorf("error : %s", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
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
		fmt.Errorf("error : %w", err)
		os.Exit(1)
	}
}

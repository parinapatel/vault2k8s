package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/creasty/defaults"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type secret struct {
	ApiGroup string `yaml:"apiVersion" default:"v1"`
	Kind     string `default:"Secret" yaml:"kind"`
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace" default:"istio-system"`
	}
	T    string            `yaml:"type" default:"Opaque"`
	Data map[string]string `yaml:"data"`
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func checkCert(data map[string]string) bool {
	keys := make([]string, len(data))

	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	requiredKeys := []string{"tls.key", "tls.crt"}

	if len(keys) < len(requiredKeys) {
		return false
	}
	for _, e := range keys {
		if !contains(requiredKeys, e) {
			return false
		}
	}
	return true
}

func convertK8sSecret(data map[string]interface{}, secretName string, namespace string, forceDecode bool) ([]byte, error) {
	d := secret{}
	_ = defaults.Set(&d)
	d.Metadata.Name = secretName
	d.Metadata.Namespace = namespace
	d.Data = make(map[string]string)
	for k, v := range data {
		_, err := base64.StdEncoding.DecodeString(v.(string))
		if err != nil {
			if forceDecode {
				logrus.Errorf("Unable to decode secret %s with key %s , are they base64 encoded ?", secretName, k)
				return nil, fmt.Errorf("unable to decode secret %s with key %s , are they base64 encoded", secretName, k)
			}
			v = base64.StdEncoding.EncodeToString([]byte(v.(string)))
		}
		// Do some common key updates
		if k == "cert" {
			k = "tls.crt"
		}
		if k == "key" {
			k = "tls.key"
		}
		if k == "csr" {
			continue
		}
		d.Data[k] = v.(string)
	}
	if checkCert(d.Data) {
		d.T = "kubernetes.io/tls"
	}
	return yaml.Marshal(&d)
}

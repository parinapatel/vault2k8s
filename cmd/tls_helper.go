package cmd

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

type certs struct {
	cert   []byte
	key    []byte
	caCert []byte
}

// CertValidForDays : Currently Verify that cert is valid for atleast a week :P ,Configure it later.
const CertValidForDays int = 7

func checkValidity(cert *x509.Certificate) {
	currentTime := time.Now()
	futureTime := currentTime.Add(time.Duration(CertValidForDays*24) * time.Hour)

	if cert.NotBefore.Sub(currentTime).Seconds() > 0 || cert.NotAfter.Sub(futureTime) < 0 {
		logrus.Warnf("Cert will be invalid at Not Before : %s,  Not After : %s", cert.NotBefore.String(), cert.NotAfter.String())
	}
}

func verifyKey(keyPem []byte, publicKey interface{}) (bool, error) {
	k, _ := pem.Decode(keyPem)
	if k == nil || k.Type != "PRIVATE KEY" {
		logrus.Fatal("failed to decode PEM block containing private key")
		return false, errors.New("failed to decode PEM block containing private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(k.Bytes)
	switch priv := key.(type) {
	case *rsa.PrivateKey:
		logrus.Info("pub is of type RSA.")
		if !priv.PublicKey.Equal(publicKey) {
			return false, errors.New("private key does not matches cert")
		}
	case *ecdsa.PrivateKey:
		logrus.Info("pub is of type ECDSA.")
		if !priv.PublicKey.Equal(publicKey) {
			return false, errors.New("private key does not matches cert")
		}
	default:
		logrus.Fatal("unknown type of private key")
		return false, err
	}
	return true, nil
}

func verifyCerts(c certs) ([]string, error) {

	// TODO: Verify Full chain of certs
	pub, _ := pem.Decode(c.cert)
	if pub == nil || pub.Type != "CERTIFICATE" {
		logrus.Fatal("failed to decode PEM block containing public key")
		return nil, errors.New("failed to decode PEM block containing public key")
	}
	cert, err := x509.ParseCertificate(pub.Bytes)
	if err != nil {
		logrus.Errorf("Failed to load Public Key. Error : %s", err)
		return nil, err
	}
	checkValidity(cert)
	KeyMatchesCerts, err := verifyKey(c.key, cert.PublicKey)
	if err != nil || !KeyMatchesCerts {
		logrus.Error(err)
		return nil, err
	}
	if len(c.caCert) > 0 {
		certPool := x509.CertPool{}
		additional_certs := c.caCert
		for len(additional_certs) > 0 {
			pub, additional_certs = pem.Decode(c.caCert)
			if (pub == nil || pub.Type != "CERTIFICATE") && len(additional_certs) > 0 {
				logrus.Warnf("unable to load cert %s", err)
				return nil, errors.New("failed to decode PEM block containing CA CERT")
			}
			caCert, err := x509.ParseCertificate(pub.Bytes)
			if err != nil {
				logrus.Errorf("Failed to load Public Key. Error : %s", err)
				return nil, err
			}
			if !caCert.IsCA {
				logrus.Errorf("Not A CA certificate.")
			}
			checkValidity(caCert)
			certPool.AddCert(caCert)
		}
		//TODO : Verify if cert is exists.
	}
	return cert.DNSNames, nil
}

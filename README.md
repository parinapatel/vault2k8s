# vault2k8s

This is small project that given vault path converts to k8s secrets, It also converts keys like `certs` , `keys` to required 'tls.key' & `tls.crt` format.

Note : It expected that you have `VAULT_ADDR` & vault login done.

# Usage 

```bash
ppatel4@ppatel4-mbp:[Wed Sep 22 00:53:37]:~/Documents/repos/vault2k8s:[master !+]$ bin/OSX/vault2k8s -h
Read Vault Path and convert to k8s secrets

Usage:
  vault2k8s [flags]

Flags:
      --forceDecode         Force to decode String b64, By default it will try to encode string to base64
  -h, --help                help for vault2k8s
  -n, --namespace string    K8s Namespace secret is deployed in (default "istio-system")
  -s, --secretName string   K8s Secret Name ( Required )
  -p, --vaultPath string    Vault Path. ( Required ) 
      --verbosity string    The level to log at [trace, debug, info, warning, error, fatal, panic] (default "info")

```
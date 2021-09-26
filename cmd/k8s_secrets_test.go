package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestContains(t *testing.T) {
	verify := []string{"tls.crt", "tls.key"}
	exist := "tls.key"
	doesNotExist := "ca.crt"
	stringExists := contains(verify, exist)
	stringNotExists := contains(verify, doesNotExist)

	assert.True(t, stringExists)
	assert.False(t, stringNotExists)
}

func TestCheckCert(t *testing.T) {
	verify := map[string]string{"tls.crt": "datat", "tls.key": "SECRETKEY"}
	notValid := map[string]string{"d": "datat", "tls.key": "SECRETKEY"}
	missingKey := map[string]string{"tls.key": "SECRETKEY"}

	stringExists := checkCert(verify)
	stringNotExists := checkCert(notValid)
	missingKeyOutput := checkCert(missingKey)

	assert.True(t, stringExists)
	assert.False(t, stringNotExists)
	assert.False(t, missingKeyOutput)
}

func convertMap2Interface(m map[string]string) map[string]interface{} {
	data := make(map[string]interface{}, len(m))
	for k, v := range m {
		data[k] = v
	}
	return data
}

func TestConvertK8sSecret(t *testing.T) {
	name := "testSecret"
	namespace := "testNamespace"
	temp := map[string]string{
		"tls.crt": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUIrRENDQVoyZ0F3SUJBZ0lVSHlSYy95cXRWMjhVcEYyREw3ZlJQWnJhSGRrd0NnWUlLb1pJemowRUF3SXcKTFRFUk1BOEdBMVVFQXd3SWRHVnpkQzVqYjIweEN6QUpCZ05WQkFZVEFsVlRNUXN3Q1FZRFZRUUlEQUpEUVRBZwpGdzB5TVRBNU1qWXdPRE0zTXpkYUdBOHlNVEl4TURrd01qQTRNemN6TjFvd0xURVJNQThHQTFVRUF3d0lkR1Z6CmRDNWpiMjB4Q3pBSkJnTlZCQVlUQWxWVE1Rc3dDUVlEVlFRSURBSkRRVEJaTUJNR0J5cUdTTTQ5QWdFR0NDcUcKU000OUF3RUhBMElBQkJlMWIvNFhuZS9LTGxOOE03NlJrTUoyZTZrdGRUWStjQ2M5M3M1bjBlK0cwSllMd3FPTQpqY0l2NjJlN2dXekhhVW5QNUNLL09WRDNFQ0YxaEhIcFluMmpnWmd3Z1pVd0hRWURWUjBPQkJZRUZCN2xQd2oyCndsOGE4ZHBqeHhrbjZTM1U2NkpiTUI4R0ExVWRJd1FZTUJhQUZCN2xQd2oyd2w4YThkcGp4eGtuNlMzVTY2SmIKTUE0R0ExVWREd0VCL3dRRUF3SUZvREFnQmdOVkhTVUJBZjhFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSApBd0l3REFZRFZSMFRBUUgvQkFJd0FEQVRCZ05WSFJFRUREQUtnZ2gwWlhOMExtTnZiVEFLQmdncWhrak9QUVFECkFnTkpBREJHQWlFQWwrdmwwMkNZVmFqWGZ0OTM4M1dVb3dEcExxU2ZSbXlpZjkwL0V3M3NHZ2NDSVFESXR5bTIKVXRXdjVvSm1vODJqcmlHVGljcW1yRnIwNGFxb2pVK1pxQkxUSEE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t",
		"tls.key": "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ2NHK053Q0d1WHNYQ0psMUsKbVZwSFViY0lxelJqbk1yOXhnN3RPV3lGYkdLaFJBTkNBQVFYdFcvK0Y1M3Z5aTVUZkRPK2taRENkbnVwTFhVMgpQbkFuUGQ3T1o5SHZodENXQzhLampJM0NMK3RudTRGc3gybEp6K1FpdnpsUTl4QWhkWVJ4NldKOQotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0t",
	}
	data := convertMap2Interface(temp)
	s, err := convertK8sSecret(data, name, namespace, true)
	assert.NoError(t, err)
	d := secret{}
	err = yaml.Unmarshal(s, &d)
	assert.NoError(t, err)
	assert.Equal(t, d.T, "kubernetes.io/tls")
	assert.Equal(t, d.Metadata.Name, name)
	assert.Equal(t, d.Metadata.Namespace, namespace)
	assert.Contains(t, d.Data, "tls.crt")
	assert.Contains(t, d.Data, "tls.key")

	temp = map[string]string{
		"cert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUIrRENDQVoyZ0F3SUJBZ0lVSHlSYy95cXRWMjhVcEYyREw3ZlJQWnJhSGRrd0NnWUlLb1pJemowRUF3SXcKTFRFUk1BOEdBMVVFQXd3SWRHVnpkQzVqYjIweEN6QUpCZ05WQkFZVEFsVlRNUXN3Q1FZRFZRUUlEQUpEUVRBZwpGdzB5TVRBNU1qWXdPRE0zTXpkYUdBOHlNVEl4TURrd01qQTRNemN6TjFvd0xURVJNQThHQTFVRUF3d0lkR1Z6CmRDNWpiMjB4Q3pBSkJnTlZCQVlUQWxWVE1Rc3dDUVlEVlFRSURBSkRRVEJaTUJNR0J5cUdTTTQ5QWdFR0NDcUcKU000OUF3RUhBMElBQkJlMWIvNFhuZS9LTGxOOE03NlJrTUoyZTZrdGRUWStjQ2M5M3M1bjBlK0cwSllMd3FPTQpqY0l2NjJlN2dXekhhVW5QNUNLL09WRDNFQ0YxaEhIcFluMmpnWmd3Z1pVd0hRWURWUjBPQkJZRUZCN2xQd2oyCndsOGE4ZHBqeHhrbjZTM1U2NkpiTUI4R0ExVWRJd1FZTUJhQUZCN2xQd2oyd2w4YThkcGp4eGtuNlMzVTY2SmIKTUE0R0ExVWREd0VCL3dRRUF3SUZvREFnQmdOVkhTVUJBZjhFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSApBd0l3REFZRFZSMFRBUUgvQkFJd0FEQVRCZ05WSFJFRUREQUtnZ2gwWlhOMExtTnZiVEFLQmdncWhrak9QUVFECkFnTkpBREJHQWlFQWwrdmwwMkNZVmFqWGZ0OTM4M1dVb3dEcExxU2ZSbXlpZjkwL0V3M3NHZ2NDSVFESXR5bTIKVXRXdjVvSm1vODJqcmlHVGljcW1yRnIwNGFxb2pVK1pxQkxUSEE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t",
		"key":  "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ2NHK053Q0d1WHNYQ0psMUsKbVZwSFViY0lxelJqbk1yOXhnN3RPV3lGYkdLaFJBTkNBQVFYdFcvK0Y1M3Z5aTVUZkRPK2taRENkbnVwTFhVMgpQbkFuUGQ3T1o5SHZodENXQzhLampJM0NMK3RudTRGc3gybEp6K1FpdnpsUTl4QWhkWVJ4NldKOQotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0t",
	}
	data = convertMap2Interface(temp)
	s, err = convertK8sSecret(data, name, namespace, true)
	assert.NoError(t, err)
	err = yaml.Unmarshal(s, &d)
	assert.NoError(t, err)
	assert.Equal(t, d.T, "kubernetes.io/tls")
	assert.Equal(t, d.Metadata.Name, name)
	assert.Equal(t, d.Metadata.Namespace, namespace)
	assert.Contains(t, d.Data, "tls.crt")
	assert.Contains(t, d.Data, "tls.key")

	// Fail Decode
	temp = map[string]string{
		"tls,crt": "USUZJQ0FURS0tLS0tCk1JSUIrRENDQVoyZ0F3SUJBZ0lVSHlSYy95cXRWMjhVcEYyREw3ZlJQWnJhSGRrd0NnWUlLb1pJemowRUF3SXcKTFRFUk1BOEdBMVVFQXd3SWRHVnpkQzVqYjIweEN6QUpCZ05WQkFZVEFsVlRNUXN3Q1FZRFZRUUlEQUpEUVRBZwpGdzB5TVRBNU1qWXdPRE0zTXpkYUdBOHlNVEl4TURrd01qQTRNemN6TjFvd0xURVJNQThHQTFVRUF3d0lkR1Z6CmRDNWpiMjB4Q3pBSkJnTlZCQVlUQWxWVE1Rc3dDUVlEVlFRSURBSkRRVEJaTUJNR0J5cUdTTTQ5QWdFR0NDcUcKU000OUF3RUhBMElBQkJlMWIvNFhuZS9LTGxOOE03NlJrTUoyZTZrdGRUWStjQ2M5M3M1bjBlK0cwSllMd3FPTQpqY0l2NjJlN2dXekhhVW5QNUNLL09WRDNFQ0YxaEhIcFluMmpnWmd3Z1pVd0hRWURWUjBPQkJZRUZCN2xQd2oyCndsOGE4ZHBqeHhrbjZTM1U2NkpiTUI4R0ExVWRJd1FZTUJhQUZCN2xQd2oyd2w4YThkcGp4eGtuNlMzVTY2SmIKTUE0R0ExVWREd0VCL3dRRUF3SUZvREFnQmdOVkhTVUJBZjhFRmpBVUJnZ3JCZ0VGQlFjREFRWUlLd1lCQlFVSApBd0l3REFZRFZSMFRBUUgvQkFJd0FEQVRCZ05WSFJFRUREQUtnZ2gwWlhOMExtTnZiVEFLQmdncWhrak9QUVFECkFnTkpBREJHQWlFQWwrdmwwMkNZVmFqWGZ0OTM4M1dVb3dEcExxU2ZSbXlpZjkwL0V3M3NHZ2NDSVFESXR5bTIKVXRXdjVvSm1vODJqcmlHVGljcW1yRnIwNGFxb2pVK1pxQkxUSEE9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t",
		"tls,key": "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JR0hBZ0VBTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEJHMHdhd0lCQVFRZ2NHK053Q0d1WHNYQ0psMUsKbVZwSFViY0lxelJqbk1yOXhnN3RPV3lGYkdLaFJBTkNBQVFYdFcvK0Y1M3Z5aTVUZkRPK2taRENkbnVwTFhVMgpQbkFuUGQ3T1o5SHZodENXQzhLampJM0NMK3RudTRGc3gybEp6K1FpdnpsUTl4QWhkWVJ4NldKOQotLS0tLUVORCBQUklWQVRFIEtFWS0tLS0t",
	}
	data = convertMap2Interface(temp)
	s, err = convertK8sSecret(data, name, namespace, true)
	assert.Error(t, err)

}

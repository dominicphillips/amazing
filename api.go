package amazing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var serviceDomains = map[string][]string{
	"CA": []string{"ecs.amazonaws.ca", "xml-ca.amznxslt.com"},
	"CN": []string{"webservices.amazon.cn", "xml-cn.amznxslt.com"},
	"DE": []string{"ecs.amazonaws.de", "xml-de.amznxslt.com"},
	"ES": []string{"webservices.amazon.es", "xml-es.amznxslt.com"},
	"FR": []string{"ecs.amazonaws.fr", "xml-fr.amznxslt.com"},
	"IT": []string{"webservices.amazon.it", "xml-it.amznxslt.com"},
	"JP": []string{"ecs.amazonaws.jp", "xml-jp.amznxslt.com"},
	"UK": []string{"ecs.amazonaws.co.uk", "xml-uk.amznxslt.com"},
	"US": []string{"ecs.amazonaws.com", "xml-us.amznxslt.com"},
}

const (
	resourcePath    = "/onca/xml"
	resourceService = "AWSECommerceService"
	resourceVersion = "2011-08-01"
	envTag          = "AMAZING_ASSOCIATE_TAG"
	envAccess       = "AMAZING_ACCESS_KEY"
	envSecret       = "AMAZING_SECRET_KEY"
)

type Amazing struct {
	Config *AmazingClientConfig
}

type AmazingClientConfig struct {
	ServiceDomain  []string
	AssociateTag   string
	AWSAccessKeyId string
	AWSSecretKey   string
}

func NewAmazing(domain, tag, access, secret string) (*Amazing, error) {
	return newAmazing(domain, tag, access, secret)
}

func NewAmazingFromEnv(domain string) (*Amazing, error) {
	tag := os.Getenv(envTag)
	access := os.Getenv(envAccess)
	secret := os.Getenv(envSecret)

	if tag == "" || access == "" || secret == "" {
		return nil, fmt.Errorf("Can't read configuration from environment variables, are they set? "+
			"%s: %s\n"+
			"%s: %s\n"+
			"%s: %s\n", envTag, tag, envAccess, access, envSecret, secret)
	}

	return newAmazing(domain, tag, access, secret)

}

func newAmazing(domain, tag, access, secret string) (*Amazing, error) {
	if d, ok := serviceDomains[domain]; ok {
		config := &AmazingClientConfig{
			ServiceDomain:  d,
			AssociateTag:   tag,
			AWSAccessKeyId: access,
			AWSSecretKey:   secret,
		}
		return &Amazing{config}, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Service domain does not exist %v", serviceDomains))
	}
}

func (a *Amazing) MergeParamsWithDefaults(extra url.Values) url.Values {
	params := url.Values{
		"AWSAccessKeyId": []string{a.Config.AWSAccessKeyId},
		"AssociateTag":   []string{a.Config.AssociateTag},
		"Service":        []string{resourceService},
		"Timestamp":      []string{time.Now().Format(time.RFC3339Nano)},
		"Version":        []string{resourceVersion},
	}
	for k, v := range extra {
		params[k] = v
	}

	return params

}

func (a *Amazing) ItemLookup(params url.Values) (*AmazonLookupResult, error) {

	var err error
	var result AmazonLookupResult

	httpClient := NewTimeoutClient(time.Duration(3*time.Second), time.Duration(3*time.Second))

	merged := a.MergeParamsWithDefaults(params)
	signThis := fmt.Sprintf("GET\n%s\n%s\n%s", a.Config.ServiceDomain[0], resourcePath, merged.Encode())
	h := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(a.Config.AWSSecretKey))
	h.Write([]byte(signThis))
	signed := base64.StdEncoding.EncodeToString(h.Sum(nil))
	merged.Set("Signature", signed)

	u := url.URL{
		Scheme:   "http",
		Host:     a.Config.ServiceDomain[0],
		Path:     resourcePath,
		RawQuery: merged.Encode(),
	}

	r, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		var errorResponse AmazonItemLookupErrorResponse
		err = xml.Unmarshal(b, &errorResponse)
		if err != nil {
			return nil, err
		}
		return nil, &errorResponse
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(b, &result)

	return &result, err
}

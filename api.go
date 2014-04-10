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
	"time"
)

var SERVICE_DOMAINS = map[string][]string{
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

type Amazing struct {
	Config *AmazingClientConfig
}

type AmazingClientConfig struct {
	ServiceDomain  []string
	AssociateTag   string
	AWSAccessKeyId string
	AWSSecretKey   string
}

type AmazonLookupResult struct {
	XMLName xml.Name `xml:"ItemLookupResponse"`
	//	Item    Product  `xml:"Items>Item"`
	//	OperationRequest
}

func NewAmazing(domain, tag, access, secret string) (*Amazing, error) {

	if d, ok := SERVICE_DOMAINS[domain]; ok {
		config := &AmazingClientConfig{
			ServiceDomain:  d,
			AssociateTag:   tag,
			AWSAccessKeyId: access,
			AWSSecretKey:   secret,
		}
		return &Amazing{config}, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Service domain does not exist %v", SERVICE_DOMAINS))
	}
}

func (a *Amazing) MergeParamsWithDefaults(extra url.Values) url.Values {
	params := url.Values{
		"AWSAccessKeyId": []string{a.Config.AWSAccessKeyId},
		"AssociateTag":   []string{a.Config.AssociateTag},
		"Service":        []string{"AWSECommerceService"},
		"Timestamp":      []string{time.Now().Format(time.RFC3339Nano)},
		"Version":        []string{"2011-08-01"},
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
	signThis := fmt.Sprintf("GET\n%s\n/onca/xml\n%s", a.Config.ServiceDomain[0], merged.Encode())
	h := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(a.Config.AWSSecretKey))

	h.Write([]byte(signThis))
	signed := base64.StdEncoding.EncodeToString(h.Sum(nil))
	merged.Set("Signature", signed)

	u, err := url.ParseRequestURI(a.Config.ServiceDomain[0])
	if err != nil {
		return nil, err
	}
	u.Path = "/onca/xml"
	u.RawQuery = merged.Encode()
	urlStr := fmt.Sprintf("%v", u)

	r, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(r)
	if err != nil || res.StatusCode != http.StatusOK {
		return nil, err
	}

	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(b, &result)

	return &result, err
}

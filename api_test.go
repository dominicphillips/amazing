package amazing

import (
	"fmt"
	"net/url"
	"os"
	_ "reflect"
	"testing"
)

func CheckEnv(t *testing.T) {
	// skip this if env variables are not set
	tag := os.Getenv(envTag)
	access := os.Getenv(envAccess)
	secret := os.Getenv(envSecret)
	if os.Getenv(envTag) == "" || os.Getenv(envAccess) == "" || os.Getenv(envSecret) == "" {
		t.Errorf("Can't read configuration from environment variables, are they set? "+
			"%s: %s\n"+
			"%s: %s\n"+
			"%s: %s\n", envTag, tag, envAccess, access, envSecret, secret)
		t.Skip()
	}
}

func TestNew(t *testing.T) {
	CheckEnv(t)
	_, err := NewAmazing("test", "tag", "access", "secret")
	if err == nil {
		t.Error("Should not be able to create client with wrong domain")
	}

	client, err := NewAmazing("DE", "tag", "access", "secret")
	if err != nil || client == nil {
		t.Errorf("Client is nil or error %s", err)
	}

	client, err = NewAmazingFromEnv("DE")
	if err != nil || client == nil {
		t.Errorf("Client is nil or error %s", err)
	}

}

func TestItemLookup(t *testing.T) {
	CheckEnv(t)
	client, _ := NewAmazingFromEnv("DE")
	params := url.Values{
		"Condition":     []string{"New"},
		"ResponseGroup": []string{"Images,Medium,Offers"},
		"IdType":        []string{"ASIN"},
		"ItemId":        []string{"1234"},
		"Operation":     []string{"ItemLookup"},
	}

	result, err := client.ItemLookup(params)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result.XMLName)

}

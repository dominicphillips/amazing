package amazing

import (
	"fmt"
	"net/url"
	"os"
	_ "reflect"
	"testing"
)

func CheckEnv(t *testing.T) {

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
	// check error handling
	params := url.Values{
		"Condition":     []string{"New"},
		"ResponseGroup": []string{"Images,Medium,Offers"},
		"IdType":        []string{"ASIN"},
		"ItemId":        []string{"1234"},
		"Operation":     []string{"ItemLookup"},
	}

	result, err := client.ItemLookup(params)

	if result == nil || err != nil {
		t.Errorf("Result is nil or error", err)
		t.Skip()
	}

	// verify there is an Error in the Result
	if len(result.AmazonItems.Request.Errors) == 0 {
		t.Errorf("Error list is empty, should be 1")
	}

	fmt.Println(result.XMLName)

	// get PlayStation 4
	params.Set("ItemId", "B00BIYAO3K")
	result, err = client.ItemLookup(params)

	if result == nil || err != nil {
		t.Errorf("Result is nil or error", err)
		t.Skip()
	}

}

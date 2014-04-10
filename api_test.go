package amazing

import (
	_ "reflect"
	"testing"
)

func TestNew(t *testing.T) {

	_, err := NewAmazing("blub", "tag", "access", "secret")
	if err == nil {
		t.Error("Should not be able to create client with wrong domain")
	}
	client, _ := NewAmazing("DE", "tag", "access", "secret")
	if client == nil {
		t.Errorf("Client is nil")
	}

}

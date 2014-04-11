package amazing

import (
	"encoding/xml"
	"fmt"
)

const (
	errorFormat = "ErrorCode: %s\nMessage: %s\nRequest:%s"
)

type AmazonLookupResult struct {
	XMLName xml.Name `xml:"ItemLookupResponse"`
}

type AmazonItemLookupErrorResponse struct {
	XMLName xml.Name `xml:"ItemLookupErrorResponse"`
	AmazonError
}

type AmazonError struct {
	Code      string `xml:"Error>Code"`
	Message   string `xml:"Error>Message"`
	RequestId string
}

type AmazonImage struct {
	XMLName xml.Name `xml:"MediumImage"`
	URL     string
	Height  uint16
	Width   uint16
}

func (a *AmazonError) Error() string {
	return fmt.Sprintf(errorFormat, a.Code, a.Message, a.RequestId)
}

package amazing

import (
	"encoding/xml"
	"fmt"
)

const (
	errorFormat = "ErrorCode: %s\nMessage: %s\nRequest:%s"
)

type AmazonItemLookupResponse struct {
	XMLName          xml.Name `xml:"ItemLookupResponse"`
	OperationRequest AmazonOperationRequest
	AmazonItems      AmazonItems `xml:"Items"`
}

type AmazonItemSearchResponse struct {
	XMLName          xml.Name `xml:"ItemSearchResponse"`
	OperationRequest AmazonOperationRequest
	AmazonItems      AmazonItems `xml:"Items"`
}

type AmazonSimilarityLookupResponse struct {
	XMLName          xml.Name `xml:"SimilarityLookupResponse"`
	OperationRequest AmazonOperationRequest
	AmazonItems      AmazonItems `xml:"Items"`
}

type AmazonOperationRequest struct {
	HTTPHeaders           []AmazonOperationRequestHeader   `xml:"HTTPHeaders>Header"`
	Arguments             []AmazonOperationRequestArgument `xml:"Arguments>Argument"`
	RequestId             string
	RequestProcessingTime float64
}

type AmazonOperationRequestHeader struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

type AmazonOperationRequestArgument struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

type AmazonItems struct {
	Request AmazonRequest
	Items   []AmazonItem `xml:"Item"`
}

type AmazonItem struct {
	ASIN             string
	ParentASIN       string
	DetailPageURL    string
	SalesRank        string
	ItemLinks        []AmazonItemLink `xml:"ItemLinks>ItemLink"`
	SmallImage       AmazonImage
	MediumImage      AmazonImage
	LargeImage       AmazonImage
	ImageSets        []AmazonImageSet `xml:"ImageSets>ImageSet"`
	ItemAttributes   AmazonItemAttributes
	OfferSummary     AmazonItemOfferSummary
	EditorialReviews []AmazonEditorialReview `xml:"EditorialReviews>EditorialReview"`
	BrowseNodes      []AmazonBrowseNode      `xml:"BrowseNodes>BrowseNode"`
}
type AmazonBrowseNode struct {
	BrowseNodeId string            `xml:"BrowseNodeId"`
	Name         string            `xml:"Name"`
	Ancestors    *AmazonBrowseNode `xml:"Ancestors>BrowseNode"`
}

type AmazonEditorialReview struct {
	Source  string
	Content string
}

type AmazonItemAtributes AmazonItemAttributes

type AmazonItemAttributes struct {
	Title     string
	Brand     string
	ListPrice AmazonItemPrice
}

type AmazonItemOfferSummary struct {
	LowestUsedPrice        AmazonItemPrice
	LowestNewPrice         AmazonItemPrice
	LowestCollectiblePrice AmazonItemPrice
}

type AmazonItemPrice struct {
	Amount         int64
	CurrencyCode   string
	FormattedPrice string
}

type AmazonItemLink struct {
	Description string
	URL         string
}

type AmazonImageSet struct {
	Category       string `xml:"Category,attr"`
	SwatchImage    AmazonImage
	SmallImage     AmazonImage
	ThumbnailImage AmazonImage
	TinyImage      AmazonImage
	MediumImage    AmazonImage
	LargeImage     AmazonImage
}

type AmazonRequest struct {
	IsValid           bool
	ItemLookupRequest AmazonItemLookupRequest
	Errors            []AmazonError
}

type AmazonItemLookupRequest struct {
	Condition      string
	IdType         string
	ItemId         string
	ResponseGroups []string `xml:"ResponseGroup"`
	VariationPage  string
}

type AmazonItemLookupErrorResponse struct {
	XMLName xml.Name
	AmazonError
}

type AmazonError struct {
	Code      string `xml:"Error>Code"`
	Message   string `xml:"Error>Message"`
	RequestId string
}

type AmazonImage struct {
	URL    string
	Height uint16
	Width  uint16
}

func (a *AmazonError) Error() string {
	return fmt.Sprintf(errorFormat, a.Code, a.Message, a.RequestId)
}

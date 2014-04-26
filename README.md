## Go Client for the Amazon Product API

A golang client for the amazon product api. This is very much WIP but it's still quite usable for simple tasks such as lookup and search. Handles annoying things such as signature signing and xml -> struct mapping.

MIT license.

### Install

    go get github.com/dominicphillips/amazing

### Configuration

Initialize a client with your ServiceDomain, AssociateTag, AWSAccessKeyId & AWSSecretKey. Service Domain may be one of the following:

    CA
    CN
    DE
    ES
    FR
    IT
    JP
    UK
    US

    client, err := amazing.NewAmazing("DE", "tag", "access", "secret")

Or from environment variables (recommended)

    export AMAZING_ASSOCIATE_TAG=
    export AMAZING_ACCESS_KEY=
    export AMAZING_SECRET_KEY=

    ----------

    client, err := amazing.NewAmazingFromEnv("DE")

Currently these operations are supported:

    ItemLookup
    result, err := client.ItemLookup(params)

    ItemSearch
    result, err := client.ItemSearch(params)

    SimilarityLookup
    result, err := client.SimilarityLookup(params)


Params are of type url.Values, for ItemLookup you would pass them like this:

    params := url.Values{
      "IdType":        []string{"ASIN"},
      "ItemId":        []string{"B00BIYAO3K"},
      "Operation":     []string{"ItemLookup"},
    }

    result, err := client.ItemLookup(params)

Results are defined in response.go, you may also pass in your own structs to the Request() function directly:

    client, _ := amazing.NewAmazingFromEnv("DE")

    params := url.Values{
        "SearchIndex": []string{"All"},
        "Operation":   []string{"ItemSearch"},
        "Keywords":    []string{"Golang"},
    }
    type CustomResult struct {
        XMLName          xml.Name `xml:"ItemSearchResponse"`
        OperationRequest AmazonOperationRequest
        AmazonItems      AmazonItems `xml:"Items"`
    }
    var result CustomResult
    err := client.Request(params, &result)


For a quick reference over the various parameters, response groups etc. check [the quick reference card](http://s3.amazonaws.com/awsdocs/Associates/2011-08-01/prod-adv-api-qrc-2011-08-01.pdf).




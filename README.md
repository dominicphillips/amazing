## Go Client for the Amazon Product API


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

    client, err := NewAmazing("DE", "tag", "access", "secret")

Or from environment variables AMAZING_ASSOCIATE_TAG, AMAZING_ACCESS_KEY, AMAZING_SECRET_KEY:

    client, err := NewAmazingFromEnv("DE")

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

Results are defined in response.go, you may also pass in your own structs to the Do(params url.Values, result interface{)) function directly:

    client, _ := NewAmazingFromEnv("DE")
    params := url.Values{
      "SearchIndex":   []string{"All"},
      "Operation":     []string{"ItemSearch"},
      "Keywords":      []string{"Golang"},
    }
    var result AmazonItemSearchResponse
    err := client.Request(params, &result)


For a quick reference over the various parameters, response groups etc. check [the quick reference card](http://s3.amazonaws.com/awsdocs/Associates/2011-08-01/prod-adv-api-qrc-2011-08-01.pdf).




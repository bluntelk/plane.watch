# \DefaultApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AutocompleteGet**](DefaultApi.md#AutocompleteGet) | **Get** /autocomplete | Autocomplete
[**AutosuggestGet**](DefaultApi.md#AutosuggestGet) | **Get** /autosuggest | Autosuggest
[**BrowseGet**](DefaultApi.md#BrowseGet) | **Get** /browse | Browse
[**DiscoverGet**](DefaultApi.md#DiscoverGet) | **Get** /discover | Discover
[**GeocodeGet**](DefaultApi.md#GeocodeGet) | **Get** /geocode | Geocode
[**LookupGet**](DefaultApi.md#LookupGet) | **Get** /lookup | Lookup By ID
[**RevgeocodeGet**](DefaultApi.md#RevgeocodeGet) | **Get** /revgeocode | Reverse Geocode



## AutocompleteGet

> OpenSearchAutocompleteResponse AutocompleteGet(ctx).Q(q).At(at).In(in).Limit(limit).Types(types).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Autocomplete



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    q := "Berlin Pariser 20" // string | Enter a free-text query  Examples:  * `ber`, `berl`, `berli`, ...  * `berlin+p`, `berlin+paris`, `berlin+parise`, ...  * `berlin+pariser+20`   _Note: Whitespace, urls, email addresses, or other out-of-scope queries will yield no results._ 
    at := "at_example" // string | Specify the center of the search context expressed as coordinates.  Format: `{latitude},{longitude}`  Type: `{decimal},{decimal}`  Example: `-13.163068,-72.545128` (Machu Picchu Mountain, Peru)  (optional)
    in := "in_example" // string | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: `countryCode:{countryCode}[,{countryCode}]*`     Examples:     * `countryCode:USA`     * `countryCode:CAN,MEX,USA`   (optional)
    limit := int32(56) // int32 | Maximum number of results to be returned. (optional) (default to 5)
    types := []string{"Types_example"} // []string | BETA: Limit the result items to the specified types. Currently supported values of the type filter for Autocomplete:  * `city` - restricting results to result type `locality` and locality type `city`  * `postalCode` - restricting results to result type `locality` and locality type `postalCode`,  * `area` - restricting results to result types: `locality` or `administrativeArea` including all the sub-types  Provide one of the supported values or a comma separated list. (optional)
    lang := []string{"Inner_example"} // []string | Select the preferred response language for result rendering from a list of BCP47 compliant Language Codes. The autocomplete endpoint tries to detect the query language based on matching name variants and then chooses the same language for the response.  Therefore the end-user can see and recognize all the entered terms in the same language as in the query. The specified preferred language is used only for not matched address tokens and for matched address tokens in case of ambiguity  (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.AutocompleteGet(context.Background()).Q(q).At(at).In(in).Limit(limit).Types(types).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.AutocompleteGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AutocompleteGet`: OpenSearchAutocompleteResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.AutocompleteGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAutocompleteGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **q** | **string** | Enter a free-text query  Examples:  * &#x60;ber&#x60;, &#x60;berl&#x60;, &#x60;berli&#x60;, ...  * &#x60;berlin+p&#x60;, &#x60;berlin+paris&#x60;, &#x60;berlin+parise&#x60;, ...  * &#x60;berlin+pariser+20&#x60;   _Note: Whitespace, urls, email addresses, or other out-of-scope queries will yield no results._  | 
 **at** | **string** | Specify the center of the search context expressed as coordinates.  Format: &#x60;{latitude},{longitude}&#x60;  Type: &#x60;{decimal},{decimal}&#x60;  Example: &#x60;-13.163068,-72.545128&#x60; (Machu Picchu Mountain, Peru)  | 
 **in** | **string** | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: &#x60;countryCode:{countryCode}[,{countryCode}]*&#x60;     Examples:     * &#x60;countryCode:USA&#x60;     * &#x60;countryCode:CAN,MEX,USA&#x60;   | 
 **limit** | **int32** | Maximum number of results to be returned. | [default to 5]
 **types** | **[]string** | BETA: Limit the result items to the specified types. Currently supported values of the type filter for Autocomplete:  * &#x60;city&#x60; - restricting results to result type &#x60;locality&#x60; and locality type &#x60;city&#x60;  * &#x60;postalCode&#x60; - restricting results to result type &#x60;locality&#x60; and locality type &#x60;postalCode&#x60;,  * &#x60;area&#x60; - restricting results to result types: &#x60;locality&#x60; or &#x60;administrativeArea&#x60; including all the sub-types  Provide one of the supported values or a comma separated list. | 
 **lang** | **[]string** | Select the preferred response language for result rendering from a list of BCP47 compliant Language Codes. The autocomplete endpoint tries to detect the query language based on matching name variants and then chooses the same language for the response.  Therefore the end-user can see and recognize all the entered terms in the same language as in the query. The specified preferred language is used only for not matched address tokens and for matched address tokens in case of ambiguity  | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**OpenSearchAutocompleteResponse**](OpenSearchAutocompleteResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## AutosuggestGet

> OpenSearchAutosuggestResponse AutosuggestGet(ctx).Q(q).At(at).In(in).Limit(limit).Route(route).TermsLimit(termsLimit).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Autosuggest



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    q := "Eismieze Berlin" // string | Enter a free-text query  Examples:  * `res`, `rest`, `resta`, `restau`, ...  * `berlin+bran`, `berlin+brand`, `berlin+branden`, ...  * `New+Yok+Giants`   _Note: Whitespace, urls, email addresses, or other out-of-scope queries will yield no results. 
    at := "52.5308,13.3856" // string | Specify the center of the search context expressed as coordinates  Format: `{latitude},{longitude}`  Type: `{decimal},{decimal}`  Example: `-13.163068,-72.545128` (Machu Picchu Mountain, Peru)  The following constraints apply:   * One of \"at\", \"in=circle\" or \"in=bbox\" is required.   * Parameters \"at\", \"in=circle\" and \"in=bbox\" are mutually exclusive. Only one of them is allowed.  (optional)
    in := "in_example" // string | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: `countryCode:{countryCode}[,{countryCode}]*`     Examples:     * `countryCode:USA`     * `countryCode:CAN,MEX,USA`    * a circular area, provided as latitude, longitude, and radius (in meters)     Format: `circle:{latitude},{longitude};r={radius}`     Type: `circle:{decimal},{decimal};r={integer}`     Example: `circle:52.53,13.38;r=10000`    * a bounding box, provided as _west longitude_, _south latitude_, _east longitude_, _north latitude_     Format: `bbox:{west longitude},{south latitude},{east longitude},{north latitude}`     Example: `bbox:13.08836,52.33812,13.761,52.6755`   The following constraints apply:   * Parameters \"at\", \"in=circle\" and \"in=bbox\" are mutually exclusive. Only one of them is allowed.    * One of \"at\", \"in=circle\" or \"in=bbox\" is required.   * The \"in=countryCode\" parameter must be accompanied by exactly one of \"at\", \"in=circle\" or \"in=bbox\".  (optional)
    limit := int32(56) // int32 | Maximum number of results to be returned. (optional) (default to 20)
    route := "route_example" // string | BETA: Select within a geographic corridor. This is a hard filter. Results will be returned if they are located within the specified area.  A `route` is defined by a [Flexible Polyline Encoding](https://github.com/heremaps/flexible-polyline),  followed by an optional width, represented by a sub-parameter \"w\".  Format: `{route};w={width}`  In regular expression syntax, the values of `route` look like:  `[a-zA-Z0-9_-]+(;w=\\d+)?`  \"[a-zA-Z0-9._-]+\" is the encoded flexible polyline.  \"w=\\d+\" is the optional width. The width is specified in meters from the center of the path. If no width is provided, the default is 1000 meters.  Type: `{Flexible Polyline Encoding};w={integer}`  The following constraints apply:  * A `route` MUST NOT contain more than 2000 points.  Examples:  * `BFoz5xJ67i1B1B7PzIhaxL7Y`  * `BFoz5xJ67i1B1B7PzIhaxL7Y;w=5000`  * `BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC;w=2000`  Note: The last example above can be decoded (using the Python class [here](https://github.com/heremaps/flexible-polyline/tree/master/python) as follows:  ``` >>> import flexpolyline >>> polyline = 'BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC' >>> flexpolyline.decode(polyline) [(52.51994, 13.38663, 1.0), (52.51009, 13.28169, 2.0), (52.43518, 13.19352, 3.0), (52.41073, 13.19645, 4.0), (52.38871, 13.15578, 5.0), (52.37278, 13.1491, 6.0), (52.37375, 13.11546, 7.0), (52.38752, 13.08722, 8.0), (52.40294, 13.07062, 9.0), (52.41058, 13.07555, 10.0)] ```  (optional)
    termsLimit := int32(56) // int32 | Maximum number of Query Terms Suggestions to be returned. (optional)
    lang := []string{"Inner_example"} // []string | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'phonemes': Renders phonemes for address and place names into the results.  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: 'tz': Renders result items with additional time zone information. Please note that this may impact latency significantly. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.AutosuggestGet(context.Background()).Q(q).At(at).In(in).Limit(limit).Route(route).TermsLimit(termsLimit).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.AutosuggestGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `AutosuggestGet`: OpenSearchAutosuggestResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.AutosuggestGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiAutosuggestGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **q** | **string** | Enter a free-text query  Examples:  * &#x60;res&#x60;, &#x60;rest&#x60;, &#x60;resta&#x60;, &#x60;restau&#x60;, ...  * &#x60;berlin+bran&#x60;, &#x60;berlin+brand&#x60;, &#x60;berlin+branden&#x60;, ...  * &#x60;New+Yok+Giants&#x60;   _Note: Whitespace, urls, email addresses, or other out-of-scope queries will yield no results.  | 
 **at** | **string** | Specify the center of the search context expressed as coordinates  Format: &#x60;{latitude},{longitude}&#x60;  Type: &#x60;{decimal},{decimal}&#x60;  Example: &#x60;-13.163068,-72.545128&#x60; (Machu Picchu Mountain, Peru)  The following constraints apply:   * One of \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; or \&quot;in&#x3D;bbox\&quot; is required.   * Parameters \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; and \&quot;in&#x3D;bbox\&quot; are mutually exclusive. Only one of them is allowed.  | 
 **in** | **string** | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: &#x60;countryCode:{countryCode}[,{countryCode}]*&#x60;     Examples:     * &#x60;countryCode:USA&#x60;     * &#x60;countryCode:CAN,MEX,USA&#x60;    * a circular area, provided as latitude, longitude, and radius (in meters)     Format: &#x60;circle:{latitude},{longitude};r&#x3D;{radius}&#x60;     Type: &#x60;circle:{decimal},{decimal};r&#x3D;{integer}&#x60;     Example: &#x60;circle:52.53,13.38;r&#x3D;10000&#x60;    * a bounding box, provided as _west longitude_, _south latitude_, _east longitude_, _north latitude_     Format: &#x60;bbox:{west longitude},{south latitude},{east longitude},{north latitude}&#x60;     Example: &#x60;bbox:13.08836,52.33812,13.761,52.6755&#x60;   The following constraints apply:   * Parameters \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; and \&quot;in&#x3D;bbox\&quot; are mutually exclusive. Only one of them is allowed.    * One of \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; or \&quot;in&#x3D;bbox\&quot; is required.   * The \&quot;in&#x3D;countryCode\&quot; parameter must be accompanied by exactly one of \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; or \&quot;in&#x3D;bbox\&quot;.  | 
 **limit** | **int32** | Maximum number of results to be returned. | [default to 20]
 **route** | **string** | BETA: Select within a geographic corridor. This is a hard filter. Results will be returned if they are located within the specified area.  A &#x60;route&#x60; is defined by a [Flexible Polyline Encoding](https://github.com/heremaps/flexible-polyline),  followed by an optional width, represented by a sub-parameter \&quot;w\&quot;.  Format: &#x60;{route};w&#x3D;{width}&#x60;  In regular expression syntax, the values of &#x60;route&#x60; look like:  &#x60;[a-zA-Z0-9_-]+(;w&#x3D;\\d+)?&#x60;  \&quot;[a-zA-Z0-9._-]+\&quot; is the encoded flexible polyline.  \&quot;w&#x3D;\\d+\&quot; is the optional width. The width is specified in meters from the center of the path. If no width is provided, the default is 1000 meters.  Type: &#x60;{Flexible Polyline Encoding};w&#x3D;{integer}&#x60;  The following constraints apply:  * A &#x60;route&#x60; MUST NOT contain more than 2000 points.  Examples:  * &#x60;BFoz5xJ67i1B1B7PzIhaxL7Y&#x60;  * &#x60;BFoz5xJ67i1B1B7PzIhaxL7Y;w&#x3D;5000&#x60;  * &#x60;BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC;w&#x3D;2000&#x60;  Note: The last example above can be decoded (using the Python class [here](https://github.com/heremaps/flexible-polyline/tree/master/python) as follows:  &#x60;&#x60;&#x60; &gt;&gt;&gt; import flexpolyline &gt;&gt;&gt; polyline &#x3D; &#39;BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC&#39; &gt;&gt;&gt; flexpolyline.decode(polyline) [(52.51994, 13.38663, 1.0), (52.51009, 13.28169, 2.0), (52.43518, 13.19352, 3.0), (52.41073, 13.19645, 4.0), (52.38871, 13.15578, 5.0), (52.37278, 13.1491, 6.0), (52.37375, 13.11546, 7.0), (52.38752, 13.08722, 8.0), (52.40294, 13.07062, 9.0), (52.41058, 13.07555, 10.0)] &#x60;&#x60;&#x60;  | 
 **termsLimit** | **int32** | Maximum number of Query Terms Suggestions to be returned. | 
 **lang** | **[]string** | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;phonemes&#39;: Renders phonemes for address and place names into the results.  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: &#39;tz&#39;: Renders result items with additional time zone information. Please note that this may impact latency significantly. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**OpenSearchAutosuggestResponse**](OpenSearchAutosuggestResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BrowseGet

> OpenSearchBrowseResponse BrowseGet(ctx).At(at).Categories(categories).Chains(chains).FoodTypes(foodTypes).In(in).Limit(limit).Name(name).Route(route).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Browse



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    at := "52.5308,13.3856" // string | Specify the center of the search context expressed as coordinates  Required parameter for endpoints that are expected to rank results by distance from the explicitly  specified search center  Format: `{latitude},{longitude}`  Type: `{decimal},{decimal}`  Example: `-13.163068,-72.545128` (Machu Picchu Mountain, Peru) 
    categories := []string{"Inner_example"} // []string | Category filter consisting of a comma-separated list of category-IDs for Categories defined in the HERE Places Category System, described in the Appendix to the HERE Search Developer Guide. Places with any assigned categories that match any of the requested categories are included in the response.  An exclamation mark \"`!`\" in front of a category ID causes that category to be excluded from the results. It is possible to mix excluded and included categories in the request - e.g. searching for places that are restaurants but not fast food restaurants. An exclusion will always win over an inclusion.  (optional)
    chains := []string{"Inner_example"} // []string | Chain filter consisting of a comma-separated list of chain-IDs for Chains defined in the HERE Places Chain System, described in the Appendix to the HERE Search Developer Guide. Places with any assigned chains that match any of the requested chains are included in the response.  An exclamation mark \"`!`\" in front of a chain ID causes that chain to be excluded from the results. It is possible to mix excluded and included chains in the request - e.g. searching for places that are amazon but not wholefoods. An exclusion will always win over an inclusion.  (optional)
    foodTypes := []string{"Inner_example"} // []string | FoodType filter consisting of a comma-separated list of cuisine-IDs for FoodTypes defined in the HERE Places Cuisine System, described in the Appendix to the HERE Search Developer Guide. Places with any assigned foodTypes that match any of the requested foodTypes are included in the response.  An exclamation mark \"`!`\" in front of a cuisine ID causes that foodType to be excluded from the results. It is possible to mix excluded and included foodTypes in the request - e.g. searching for places that serve italian but not chinese. An exclusion will always win over an inclusion.  (optional)
    in := "in_example" // string | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: `countryCode:{countryCode}[,{countryCode}]*`     Examples:     * `countryCode:USA`     * `countryCode:CAN,MEX,USA`    * a circular area, provided as latitude, longitude, and radius (in meters)     Format: `circle:{latitude},{longitude};r={radius}`     Type: `circle:{decimal},{decimal};r={integer}`     Example: `circle:52.53,13.38;r=10000`    * a bounding box, provided as _west longitude_, _south latitude_, _east longitude_, _north latitude_     Format: `bbox:{west longitude},{south latitude},{east longitude},{north latitude}`     Example: `bbox:13.08836,52.33812,13.761,52.6755`   The following constraints apply:   * Parameters \"in=circle\" and \"in=bbox\" are mutually exclusive. Only one of them is allowed.  (optional)
    limit := int32(56) // int32 | Maximum number of results to be returned. (optional) (default to 20)
    name := "name_example" // string | Full-text filter on POI names/titles. Results with a partial match on the name parameter are included in the response. (optional)
    route := "route_example" // string | BETA: Select within a geographic corridor. This is a hard filter. Results will be returned if they are located within the specified area.  A `route` is defined by a [Flexible Polyline Encoding](https://github.com/heremaps/flexible-polyline),  followed by an optional width, represented by a sub-parameter \"w\".  Format: `{route};w={width}`  In regular expression syntax, the values of `route` look like:  `[a-zA-Z0-9_-]+(;w=\\d+)?`  \"[a-zA-Z0-9._-]+\" is the encoded flexible polyline.  \"w=\\d+\" is the optional width. The width is specified in meters from the center of the path. If no width is provided, the default is 1000 meters.  Type: `{Flexible Polyline Encoding};w={integer}`  The following constraints apply:  * A `route` MUST NOT contain more than 2000 points.  Examples:  * `BFoz5xJ67i1B1B7PzIhaxL7Y`  * `BFoz5xJ67i1B1B7PzIhaxL7Y;w=5000`  * `BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC;w=2000`  Note: The last example above can be decoded (using the Python class [here](https://github.com/heremaps/flexible-polyline/tree/master/python) as follows:  ``` >>> import flexpolyline >>> polyline = 'BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC' >>> flexpolyline.decode(polyline) [(52.51994, 13.38663, 1.0), (52.51009, 13.28169, 2.0), (52.43518, 13.19352, 3.0), (52.41073, 13.19645, 4.0), (52.38871, 13.15578, 5.0), (52.37278, 13.1491, 6.0), (52.37375, 13.11546, 7.0), (52.38752, 13.08722, 8.0), (52.40294, 13.07062, 9.0), (52.41058, 13.07555, 10.0)] ```  (optional)
    lang := []string{"Inner_example"} // []string | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'phonemes': Renders phonemes for address and place names into the results.  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: 'tz': Renders result items with additional time zone information. Please note that this may impact latency significantly. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.BrowseGet(context.Background()).At(at).Categories(categories).Chains(chains).FoodTypes(foodTypes).In(in).Limit(limit).Name(name).Route(route).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.BrowseGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BrowseGet`: OpenSearchBrowseResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.BrowseGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiBrowseGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **at** | **string** | Specify the center of the search context expressed as coordinates  Required parameter for endpoints that are expected to rank results by distance from the explicitly  specified search center  Format: &#x60;{latitude},{longitude}&#x60;  Type: &#x60;{decimal},{decimal}&#x60;  Example: &#x60;-13.163068,-72.545128&#x60; (Machu Picchu Mountain, Peru)  | 
 **categories** | **[]string** | Category filter consisting of a comma-separated list of category-IDs for Categories defined in the HERE Places Category System, described in the Appendix to the HERE Search Developer Guide. Places with any assigned categories that match any of the requested categories are included in the response.  An exclamation mark \&quot;&#x60;!&#x60;\&quot; in front of a category ID causes that category to be excluded from the results. It is possible to mix excluded and included categories in the request - e.g. searching for places that are restaurants but not fast food restaurants. An exclusion will always win over an inclusion.  | 
 **chains** | **[]string** | Chain filter consisting of a comma-separated list of chain-IDs for Chains defined in the HERE Places Chain System, described in the Appendix to the HERE Search Developer Guide. Places with any assigned chains that match any of the requested chains are included in the response.  An exclamation mark \&quot;&#x60;!&#x60;\&quot; in front of a chain ID causes that chain to be excluded from the results. It is possible to mix excluded and included chains in the request - e.g. searching for places that are amazon but not wholefoods. An exclusion will always win over an inclusion.  | 
 **foodTypes** | **[]string** | FoodType filter consisting of a comma-separated list of cuisine-IDs for FoodTypes defined in the HERE Places Cuisine System, described in the Appendix to the HERE Search Developer Guide. Places with any assigned foodTypes that match any of the requested foodTypes are included in the response.  An exclamation mark \&quot;&#x60;!&#x60;\&quot; in front of a cuisine ID causes that foodType to be excluded from the results. It is possible to mix excluded and included foodTypes in the request - e.g. searching for places that serve italian but not chinese. An exclusion will always win over an inclusion.  | 
 **in** | **string** | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: &#x60;countryCode:{countryCode}[,{countryCode}]*&#x60;     Examples:     * &#x60;countryCode:USA&#x60;     * &#x60;countryCode:CAN,MEX,USA&#x60;    * a circular area, provided as latitude, longitude, and radius (in meters)     Format: &#x60;circle:{latitude},{longitude};r&#x3D;{radius}&#x60;     Type: &#x60;circle:{decimal},{decimal};r&#x3D;{integer}&#x60;     Example: &#x60;circle:52.53,13.38;r&#x3D;10000&#x60;    * a bounding box, provided as _west longitude_, _south latitude_, _east longitude_, _north latitude_     Format: &#x60;bbox:{west longitude},{south latitude},{east longitude},{north latitude}&#x60;     Example: &#x60;bbox:13.08836,52.33812,13.761,52.6755&#x60;   The following constraints apply:   * Parameters \&quot;in&#x3D;circle\&quot; and \&quot;in&#x3D;bbox\&quot; are mutually exclusive. Only one of them is allowed.  | 
 **limit** | **int32** | Maximum number of results to be returned. | [default to 20]
 **name** | **string** | Full-text filter on POI names/titles. Results with a partial match on the name parameter are included in the response. | 
 **route** | **string** | BETA: Select within a geographic corridor. This is a hard filter. Results will be returned if they are located within the specified area.  A &#x60;route&#x60; is defined by a [Flexible Polyline Encoding](https://github.com/heremaps/flexible-polyline),  followed by an optional width, represented by a sub-parameter \&quot;w\&quot;.  Format: &#x60;{route};w&#x3D;{width}&#x60;  In regular expression syntax, the values of &#x60;route&#x60; look like:  &#x60;[a-zA-Z0-9_-]+(;w&#x3D;\\d+)?&#x60;  \&quot;[a-zA-Z0-9._-]+\&quot; is the encoded flexible polyline.  \&quot;w&#x3D;\\d+\&quot; is the optional width. The width is specified in meters from the center of the path. If no width is provided, the default is 1000 meters.  Type: &#x60;{Flexible Polyline Encoding};w&#x3D;{integer}&#x60;  The following constraints apply:  * A &#x60;route&#x60; MUST NOT contain more than 2000 points.  Examples:  * &#x60;BFoz5xJ67i1B1B7PzIhaxL7Y&#x60;  * &#x60;BFoz5xJ67i1B1B7PzIhaxL7Y;w&#x3D;5000&#x60;  * &#x60;BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC;w&#x3D;2000&#x60;  Note: The last example above can be decoded (using the Python class [here](https://github.com/heremaps/flexible-polyline/tree/master/python) as follows:  &#x60;&#x60;&#x60; &gt;&gt;&gt; import flexpolyline &gt;&gt;&gt; polyline &#x3D; &#39;BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC&#39; &gt;&gt;&gt; flexpolyline.decode(polyline) [(52.51994, 13.38663, 1.0), (52.51009, 13.28169, 2.0), (52.43518, 13.19352, 3.0), (52.41073, 13.19645, 4.0), (52.38871, 13.15578, 5.0), (52.37278, 13.1491, 6.0), (52.37375, 13.11546, 7.0), (52.38752, 13.08722, 8.0), (52.40294, 13.07062, 9.0), (52.41058, 13.07555, 10.0)] &#x60;&#x60;&#x60;  | 
 **lang** | **[]string** | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;phonemes&#39;: Renders phonemes for address and place names into the results.  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: &#39;tz&#39;: Renders result items with additional time zone information. Please note that this may impact latency significantly. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**OpenSearchBrowseResponse**](OpenSearchBrowseResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DiscoverGet

> OpenSearchSearchResponse DiscoverGet(ctx).Q(q).At(at).In(in).Limit(limit).Route(route).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Discover



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    q := "Eismieze Berlin" // string | Enter a free-text query  Examples:  * `125, Berliner, berlin`  * `Beacon, Boston, Hospital`  * `Schnurrbart German Pub and Restaurant, Hong Kong`   _Note: Whitespace, urls, email addresses, or other out-of-scope queries will yield no results. 
    at := "52.5308,13.3856" // string | Specify the center of the search context expressed as coordinates  Format: `{latitude},{longitude}`  Type: `{decimal},{decimal}`  Example: `-13.163068,-72.545128` (Machu Picchu Mountain, Peru)  The following constraints apply:   * One of \"at\", \"in=circle\" or \"in=bbox\" is required.   * Parameters \"at\", \"in=circle\" and \"in=bbox\" are mutually exclusive. Only one of them is allowed.  (optional)
    in := "in_example" // string | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: `countryCode:{countryCode}[,{countryCode}]*`     Examples:     * `countryCode:USA`     * `countryCode:CAN,MEX,USA`    * a circular area, provided as latitude, longitude, and radius (in meters)     Format: `circle:{latitude},{longitude};r={radius}`     Type: `circle:{decimal},{decimal};r={integer}`     Example: `circle:52.53,13.38;r=10000`    * a bounding box, provided as _west longitude_, _south latitude_, _east longitude_, _north latitude_     Format: `bbox:{west longitude},{south latitude},{east longitude},{north latitude}`     Example: `bbox:13.08836,52.33812,13.761,52.6755`   The following constraints apply:   * Parameters \"at\", \"in=circle\" and \"in=bbox\" are mutually exclusive. Only one of them is allowed.    * One of \"at\", \"in=circle\" or \"in=bbox\" is required.   * The \"in=countryCode\" parameter must be accompanied by exactly one of \"at\", \"in=circle\" or \"in=bbox\".  (optional)
    limit := int32(56) // int32 | Maximum number of results to be returned. (optional) (default to 20)
    route := "route_example" // string | BETA: Select within a geographic corridor. This is a hard filter. Results will be returned if they are located within the specified area.  A `route` is defined by a [Flexible Polyline Encoding](https://github.com/heremaps/flexible-polyline),  followed by an optional width, represented by a sub-parameter \"w\".  Format: `{route};w={width}`  In regular expression syntax, the values of `route` look like:  `[a-zA-Z0-9_-]+(;w=\\d+)?`  \"[a-zA-Z0-9._-]+\" is the encoded flexible polyline.  \"w=\\d+\" is the optional width. The width is specified in meters from the center of the path. If no width is provided, the default is 1000 meters.  Type: `{Flexible Polyline Encoding};w={integer}`  The following constraints apply:  * A `route` MUST NOT contain more than 2000 points.  Examples:  * `BFoz5xJ67i1B1B7PzIhaxL7Y`  * `BFoz5xJ67i1B1B7PzIhaxL7Y;w=5000`  * `BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC;w=2000`  Note: The last example above can be decoded (using the Python class [here](https://github.com/heremaps/flexible-polyline/tree/master/python) as follows:  ``` >>> import flexpolyline >>> polyline = 'BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC' >>> flexpolyline.decode(polyline) [(52.51994, 13.38663, 1.0), (52.51009, 13.28169, 2.0), (52.43518, 13.19352, 3.0), (52.41073, 13.19645, 4.0), (52.38871, 13.15578, 5.0), (52.37278, 13.1491, 6.0), (52.37375, 13.11546, 7.0), (52.38752, 13.08722, 8.0), (52.40294, 13.07062, 9.0), (52.41058, 13.07555, 10.0)] ```  (optional)
    lang := []string{"Inner_example"} // []string | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'phonemes': Renders phonemes for address and place names into the results.  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: 'tz': Renders result items with additional time zone information. Please note that this may impact latency significantly. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.DiscoverGet(context.Background()).Q(q).At(at).In(in).Limit(limit).Route(route).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.DiscoverGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoverGet`: OpenSearchSearchResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.DiscoverGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **q** | **string** | Enter a free-text query  Examples:  * &#x60;125, Berliner, berlin&#x60;  * &#x60;Beacon, Boston, Hospital&#x60;  * &#x60;Schnurrbart German Pub and Restaurant, Hong Kong&#x60;   _Note: Whitespace, urls, email addresses, or other out-of-scope queries will yield no results.  | 
 **at** | **string** | Specify the center of the search context expressed as coordinates  Format: &#x60;{latitude},{longitude}&#x60;  Type: &#x60;{decimal},{decimal}&#x60;  Example: &#x60;-13.163068,-72.545128&#x60; (Machu Picchu Mountain, Peru)  The following constraints apply:   * One of \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; or \&quot;in&#x3D;bbox\&quot; is required.   * Parameters \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; and \&quot;in&#x3D;bbox\&quot; are mutually exclusive. Only one of them is allowed.  | 
 **in** | **string** | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: &#x60;countryCode:{countryCode}[,{countryCode}]*&#x60;     Examples:     * &#x60;countryCode:USA&#x60;     * &#x60;countryCode:CAN,MEX,USA&#x60;    * a circular area, provided as latitude, longitude, and radius (in meters)     Format: &#x60;circle:{latitude},{longitude};r&#x3D;{radius}&#x60;     Type: &#x60;circle:{decimal},{decimal};r&#x3D;{integer}&#x60;     Example: &#x60;circle:52.53,13.38;r&#x3D;10000&#x60;    * a bounding box, provided as _west longitude_, _south latitude_, _east longitude_, _north latitude_     Format: &#x60;bbox:{west longitude},{south latitude},{east longitude},{north latitude}&#x60;     Example: &#x60;bbox:13.08836,52.33812,13.761,52.6755&#x60;   The following constraints apply:   * Parameters \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; and \&quot;in&#x3D;bbox\&quot; are mutually exclusive. Only one of them is allowed.    * One of \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; or \&quot;in&#x3D;bbox\&quot; is required.   * The \&quot;in&#x3D;countryCode\&quot; parameter must be accompanied by exactly one of \&quot;at\&quot;, \&quot;in&#x3D;circle\&quot; or \&quot;in&#x3D;bbox\&quot;.  | 
 **limit** | **int32** | Maximum number of results to be returned. | [default to 20]
 **route** | **string** | BETA: Select within a geographic corridor. This is a hard filter. Results will be returned if they are located within the specified area.  A &#x60;route&#x60; is defined by a [Flexible Polyline Encoding](https://github.com/heremaps/flexible-polyline),  followed by an optional width, represented by a sub-parameter \&quot;w\&quot;.  Format: &#x60;{route};w&#x3D;{width}&#x60;  In regular expression syntax, the values of &#x60;route&#x60; look like:  &#x60;[a-zA-Z0-9_-]+(;w&#x3D;\\d+)?&#x60;  \&quot;[a-zA-Z0-9._-]+\&quot; is the encoded flexible polyline.  \&quot;w&#x3D;\\d+\&quot; is the optional width. The width is specified in meters from the center of the path. If no width is provided, the default is 1000 meters.  Type: &#x60;{Flexible Polyline Encoding};w&#x3D;{integer}&#x60;  The following constraints apply:  * A &#x60;route&#x60; MUST NOT contain more than 2000 points.  Examples:  * &#x60;BFoz5xJ67i1B1B7PzIhaxL7Y&#x60;  * &#x60;BFoz5xJ67i1B1B7PzIhaxL7Y;w&#x3D;5000&#x60;  * &#x60;BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC;w&#x3D;2000&#x60;  Note: The last example above can be decoded (using the Python class [here](https://github.com/heremaps/flexible-polyline/tree/master/python) as follows:  &#x60;&#x60;&#x60; &gt;&gt;&gt; import flexpolyline &gt;&gt;&gt; polyline &#x3D; &#39;BlD05xgKuy2xCCx9B7vUCl0OhnRC54EqSCzpEl-HCxjD3pBCiGnyGCi2CvwFCsgD3nDC4vB6eC&#39; &gt;&gt;&gt; flexpolyline.decode(polyline) [(52.51994, 13.38663, 1.0), (52.51009, 13.28169, 2.0), (52.43518, 13.19352, 3.0), (52.41073, 13.19645, 4.0), (52.38871, 13.15578, 5.0), (52.37278, 13.1491, 6.0), (52.37375, 13.11546, 7.0), (52.38752, 13.08722, 8.0), (52.40294, 13.07062, 9.0), (52.41058, 13.07555, 10.0)] &#x60;&#x60;&#x60;  | 
 **lang** | **[]string** | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;phonemes&#39;: Renders phonemes for address and place names into the results.  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: &#39;tz&#39;: Renders result items with additional time zone information. Please note that this may impact latency significantly. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**OpenSearchSearchResponse**](OpenSearchSearchResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GeocodeGet

> OpenSearchGeocodeResponse GeocodeGet(ctx).At(at).In(in).Limit(limit).Q(q).Qq(qq).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Geocode



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    at := "at_example" // string | Specify the center of the search context expressed as coordinates.  Format: `{latitude},{longitude}`  Type: `{decimal},{decimal}`  Example: `-13.163068,-72.545128` (Machu Picchu Mountain, Peru)  (optional)
    in := "in_example" // string | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: `countryCode:{countryCode}[,{countryCode}]*`     Examples:     * `countryCode:USA`     * `countryCode:CAN,MEX,USA`   (optional)
    limit := int32(56) // int32 | Maximum number of results to be returned. (optional) (default to 20)
    q := "Invalidenstraße 116 Berlin" // string | Enter a free-text query  Examples:  * `125, Berliner, berlin`  * `Beacon, Boston, Hospital`  * `Schnurrbart German Pub and Restaurant, Hong Kong`  _Note: Either `q` or `qq`-parameter is required on this endpoint. Both parameters can be provided in the same request._  (optional)
    qq := "qq_example" // string | Enter a qualified query. A qualified query is similar to a free-text query, but in a structured manner.  It can take multiple _sub-parameters_, separated by semicolon, allowing to specify different aspects of a query.  Currently supported _sub-parameters_ are `country`, `state`, `county`, `city`, `district`, `street`,  `houseNumber`, and `postalCode`.  Format: `{sub-parameter}={string}[;{sub-parameter}={string}]*`  Examples:  * `city=Berlin;country=Germany;street=Friedrichstr;houseNumber=20`  * `city=Berlin;country=Germany`  * `postalCode=10969`  _Note: Either `q` or `qq`-parameter is required on this endpoint. Both parameters can be provided in the same request._  (optional)
    lang := []string{"Inner_example"} // []string | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'countryInfo': For each result item renders additional block with the country info, such as [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) and [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code.  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: 'parsing'  BETA: 'tz': Renders result items with additional time zone information. Please note that this may impact latency significantly. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.GeocodeGet(context.Background()).At(at).In(in).Limit(limit).Q(q).Qq(qq).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.GeocodeGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `GeocodeGet`: OpenSearchGeocodeResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.GeocodeGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiGeocodeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **at** | **string** | Specify the center of the search context expressed as coordinates.  Format: &#x60;{latitude},{longitude}&#x60;  Type: &#x60;{decimal},{decimal}&#x60;  Example: &#x60;-13.163068,-72.545128&#x60; (Machu Picchu Mountain, Peru)  | 
 **in** | **string** | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a country (or multiple countries), provided as comma-separated [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country codes     The country codes are to be provided in all uppercase.     Format: &#x60;countryCode:{countryCode}[,{countryCode}]*&#x60;     Examples:     * &#x60;countryCode:USA&#x60;     * &#x60;countryCode:CAN,MEX,USA&#x60;   | 
 **limit** | **int32** | Maximum number of results to be returned. | [default to 20]
 **q** | **string** | Enter a free-text query  Examples:  * &#x60;125, Berliner, berlin&#x60;  * &#x60;Beacon, Boston, Hospital&#x60;  * &#x60;Schnurrbart German Pub and Restaurant, Hong Kong&#x60;  _Note: Either &#x60;q&#x60; or &#x60;qq&#x60;-parameter is required on this endpoint. Both parameters can be provided in the same request._  | 
 **qq** | **string** | Enter a qualified query. A qualified query is similar to a free-text query, but in a structured manner.  It can take multiple _sub-parameters_, separated by semicolon, allowing to specify different aspects of a query.  Currently supported _sub-parameters_ are &#x60;country&#x60;, &#x60;state&#x60;, &#x60;county&#x60;, &#x60;city&#x60;, &#x60;district&#x60;, &#x60;street&#x60;,  &#x60;houseNumber&#x60;, and &#x60;postalCode&#x60;.  Format: &#x60;{sub-parameter}&#x3D;{string}[;{sub-parameter}&#x3D;{string}]*&#x60;  Examples:  * &#x60;city&#x3D;Berlin;country&#x3D;Germany;street&#x3D;Friedrichstr;houseNumber&#x3D;20&#x60;  * &#x60;city&#x3D;Berlin;country&#x3D;Germany&#x60;  * &#x60;postalCode&#x3D;10969&#x60;  _Note: Either &#x60;q&#x60; or &#x60;qq&#x60;-parameter is required on this endpoint. Both parameters can be provided in the same request._  | 
 **lang** | **[]string** | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;countryInfo&#39;: For each result item renders additional block with the country info, such as [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) and [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code.  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: &#39;parsing&#39;  BETA: &#39;tz&#39;: Renders result items with additional time zone information. Please note that this may impact latency significantly. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**OpenSearchGeocodeResponse**](OpenSearchGeocodeResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## LookupGet

> LookupResponse LookupGet(ctx).Id(id).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Lookup By ID



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    id := "here:pds:place:276u33db-8097f3194e4b411081b761ea9a366776" // string | Location ID, which is the ID of a result item eg. of a Discover request
    lang := []string{"Inner_example"} // []string | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'countryInfo': For each result item renders additional block with the country info, such as [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) and [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code.  'phonemes': Renders phonemes for address and place names into the results.  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: 'tz': Renders result items with additional time zone information. Please note that this may impact latency significantly. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.LookupGet(context.Background()).Id(id).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.LookupGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `LookupGet`: LookupResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.LookupGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiLookupGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string** | Location ID, which is the ID of a result item eg. of a Discover request | 
 **lang** | **[]string** | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;countryInfo&#39;: For each result item renders additional block with the country info, such as [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) and [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code.  &#39;phonemes&#39;: Renders phonemes for address and place names into the results.  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: &#39;tz&#39;: Renders result items with additional time zone information. Please note that this may impact latency significantly. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**LookupResponse**](LookupResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RevgeocodeGet

> OpenSearchReverseGeocodeResponse RevgeocodeGet(ctx).At(at).In(in).Limit(limit).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()

Reverse Geocode



### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    at := "52.5308,13.3856" // string | Specify the center of the search context expressed as coordinates.  Format: `{latitude},{longitude}`  Type: `{decimal},{decimal}`  Example: `-13.163068,-72.545128` (Machu Picchu Mountain, Peru)  The following constraints apply:   * Either \"at\" or \"in=circle\" is required.   * Parameters \"at\" and \"in=circle\" are mutually exclusive. Only one of them is allowed.  (optional)
    in := "in_example" // string | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a circular area, provided as latitude, longitude, and radius (in meters)     Format: `circle:{latitude},{longitude};r={radius}`     Type: `circle:{decimal},{decimal};r={integer}`     Example: `circle:52.53,13.38;r=10000`   The following constraints apply:   * Either \"at\" or \"in=circle\" is required.  (optional)
    limit := int32(56) // int32 | Maximum number of results to be returned. (optional) (default to 1)
    lang := []string{"Inner_example"} // []string | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. (optional)
    politicalView := "politicalView_example" // string | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \"400\" error code. (optional)
    show := []string{"Show_example"} // []string | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  'countryInfo': For each result item renders additional block with the country info, such as [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) and [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code.  'streetInfo': For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: 'tz': Renders result items with additional time zone information. Please note that this may impact latency significantly. (optional)
    xRequestID := "xRequestID_example" // string | Used to correlate requests with their responses within a customer's application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. (optional)

    configuration := openapiclient.NewConfiguration()
    api_client := openapiclient.NewAPIClient(configuration)
    resp, r, err := api_client.DefaultApi.RevgeocodeGet(context.Background()).At(at).In(in).Limit(limit).Lang(lang).PoliticalView(politicalView).Show(show).XRequestID(xRequestID).Execute()
    if err.Error() != "" {
        fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.RevgeocodeGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RevgeocodeGet`: OpenSearchReverseGeocodeResponse
    fmt.Fprintf(os.Stdout, "Response from `DefaultApi.RevgeocodeGet`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiRevgeocodeGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **at** | **string** | Specify the center of the search context expressed as coordinates.  Format: &#x60;{latitude},{longitude}&#x60;  Type: &#x60;{decimal},{decimal}&#x60;  Example: &#x60;-13.163068,-72.545128&#x60; (Machu Picchu Mountain, Peru)  The following constraints apply:   * Either \&quot;at\&quot; or \&quot;in&#x3D;circle\&quot; is required.   * Parameters \&quot;at\&quot; and \&quot;in&#x3D;circle\&quot; are mutually exclusive. Only one of them is allowed.  | 
 **in** | **string** | Search within a geographic area. This is a hard filter. Results will be returned if they are located within the specified area.  A geographic area can be   * a circular area, provided as latitude, longitude, and radius (in meters)     Format: &#x60;circle:{latitude},{longitude};r&#x3D;{radius}&#x60;     Type: &#x60;circle:{decimal},{decimal};r&#x3D;{integer}&#x60;     Example: &#x60;circle:52.53,13.38;r&#x3D;10000&#x60;   The following constraints apply:   * Either \&quot;at\&quot; or \&quot;in&#x3D;circle\&quot; is required.  | 
 **limit** | **int32** | Maximum number of results to be returned. | [default to 1]
 **lang** | **[]string** | Select the language to be used for result rendering from a list of [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) compliant language codes. | 
 **politicalView** | **string** | Toggle the political view.  This parameter accepts single ISO 3166-1 alpha-3 country code. The country codes are to be provided in all uppercase.  Currently the only supported political views are:  * RUS expressing the Russian view on Crimea  * SRB expressing the Serbian view on Kosovo, Vukovar and Sarengrad Islands  * MAR expressing the Moroccan view on Western Sahara  * SUR Suriname view on Courantyne Headwaters and Lawa Headwaters  * KEN Kenya view on Ilemi Triangle  * TZA Tanzania view on Lake Malawi  * URY Uruguay view on Rincon de Artigas  * EGY Egypt view on Bir Tawil  * SDN Sudan view on Halaib Triangle  * SYR Syria view on Golan Heights  * ARG Argentina view on Southern Patagonian Ice Field and Tierra Del Fuego, including Falkland Islands, South Georgia and South Sandwich Islands  For any valid 3 letter country code, for which GS7 does not have dedicated political view, it falls back to the default view.  For not accepted values of the politicalView parameter the GS7 responds with \&quot;400\&quot; error code. | 
 **show** | **[]string** | Select additional fields to be rendered in the response. Please note that some of the fields involve additional webservice calls and can increase the overall response time.  The value is a comma-separated list of the sections to be enabled. For some sections there is a long and a short ID.  Description of accepted values:  &#39;countryInfo&#39;: For each result item renders additional block with the country info, such as [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) and [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code.  &#39;streetInfo&#39;: For each result item renders additional block with the street name decomposed into its parts like the base name, the street type, etc.  BETA: &#39;tz&#39;: Renders result items with additional time zone information. Please note that this may impact latency significantly. | 
 **xRequestID** | **string** | Used to correlate requests with their responses within a customer&#39;s application, for logging and error reporting.  Format: Free string, but a valid UUIDv4 is recommended. | 

### Return type

[**OpenSearchReverseGeocodeResponse**](OpenSearchReverseGeocodeResponse.md)

### Authorization

[ApiKey](../README.md#ApiKey), [Bearer](../README.md#Bearer)

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


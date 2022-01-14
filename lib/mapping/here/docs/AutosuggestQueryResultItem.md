# AutosuggestQueryResultItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | The localized display name of this result item. | 
**Id** | Pointer to **string** | The unique identifier for the result item. This ID can be used for a Look Up by ID search as well. | [optional] 
**ResultType** | Pointer to **string** | WARNING: The resultType values &#39;intersection&#39; and &#39;postalCodePoint&#39; are in BETA state | [optional] 
**Href** | Pointer to **string** | URL of the follow-up query | [optional] 
**Highlights** | Pointer to [**TitleHighlighting**](TitleHighlighting.md) | Describes how the parts of the response element matched the input query | [optional] 

## Methods

### NewAutosuggestQueryResultItem

`func NewAutosuggestQueryResultItem(title string, ) *AutosuggestQueryResultItem`

NewAutosuggestQueryResultItem instantiates a new AutosuggestQueryResultItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAutosuggestQueryResultItemWithDefaults

`func NewAutosuggestQueryResultItemWithDefaults() *AutosuggestQueryResultItem`

NewAutosuggestQueryResultItemWithDefaults instantiates a new AutosuggestQueryResultItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *AutosuggestQueryResultItem) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *AutosuggestQueryResultItem) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *AutosuggestQueryResultItem) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetId

`func (o *AutosuggestQueryResultItem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AutosuggestQueryResultItem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AutosuggestQueryResultItem) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AutosuggestQueryResultItem) HasId() bool`

HasId returns a boolean if a field has been set.

### GetResultType

`func (o *AutosuggestQueryResultItem) GetResultType() string`

GetResultType returns the ResultType field if non-nil, zero value otherwise.

### GetResultTypeOk

`func (o *AutosuggestQueryResultItem) GetResultTypeOk() (*string, bool)`

GetResultTypeOk returns a tuple with the ResultType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResultType

`func (o *AutosuggestQueryResultItem) SetResultType(v string)`

SetResultType sets ResultType field to given value.

### HasResultType

`func (o *AutosuggestQueryResultItem) HasResultType() bool`

HasResultType returns a boolean if a field has been set.

### GetHref

`func (o *AutosuggestQueryResultItem) GetHref() string`

GetHref returns the Href field if non-nil, zero value otherwise.

### GetHrefOk

`func (o *AutosuggestQueryResultItem) GetHrefOk() (*string, bool)`

GetHrefOk returns a tuple with the Href field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHref

`func (o *AutosuggestQueryResultItem) SetHref(v string)`

SetHref sets Href field to given value.

### HasHref

`func (o *AutosuggestQueryResultItem) HasHref() bool`

HasHref returns a boolean if a field has been set.

### GetHighlights

`func (o *AutosuggestQueryResultItem) GetHighlights() TitleHighlighting`

GetHighlights returns the Highlights field if non-nil, zero value otherwise.

### GetHighlightsOk

`func (o *AutosuggestQueryResultItem) GetHighlightsOk() (*TitleHighlighting, bool)`

GetHighlightsOk returns a tuple with the Highlights field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHighlights

`func (o *AutosuggestQueryResultItem) SetHighlights(v TitleHighlighting)`

SetHighlights sets Highlights field to given value.

### HasHighlights

`func (o *AutosuggestQueryResultItem) HasHighlights() bool`

HasHighlights returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



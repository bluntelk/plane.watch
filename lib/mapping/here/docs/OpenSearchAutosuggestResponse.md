# OpenSearchAutosuggestResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | [**[]OneOfAutosuggestEntityResultItemAutosuggestQueryResultItem**](OneOfAutosuggestEntityResultItemAutosuggestQueryResultItem.md) | The results are presented as a JSON list of candidates in ranked order (most-likely to least-likely) based on the matched location criteria. | 
**QueryTerms** | [**[]QueryTermResultItem**](QueryTermResultItem.md) | Suggestions for refining individual query terms | 

## Methods

### NewOpenSearchAutosuggestResponse

`func NewOpenSearchAutosuggestResponse(items []OneOfAutosuggestEntityResultItemAutosuggestQueryResultItem, queryTerms []QueryTermResultItem, ) *OpenSearchAutosuggestResponse`

NewOpenSearchAutosuggestResponse instantiates a new OpenSearchAutosuggestResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenSearchAutosuggestResponseWithDefaults

`func NewOpenSearchAutosuggestResponseWithDefaults() *OpenSearchAutosuggestResponse`

NewOpenSearchAutosuggestResponseWithDefaults instantiates a new OpenSearchAutosuggestResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *OpenSearchAutosuggestResponse) GetItems() []OneOfAutosuggestEntityResultItemAutosuggestQueryResultItem`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *OpenSearchAutosuggestResponse) GetItemsOk() (*[]OneOfAutosuggestEntityResultItemAutosuggestQueryResultItem, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *OpenSearchAutosuggestResponse) SetItems(v []OneOfAutosuggestEntityResultItemAutosuggestQueryResultItem)`

SetItems sets Items field to given value.


### GetQueryTerms

`func (o *OpenSearchAutosuggestResponse) GetQueryTerms() []QueryTermResultItem`

GetQueryTerms returns the QueryTerms field if non-nil, zero value otherwise.

### GetQueryTermsOk

`func (o *OpenSearchAutosuggestResponse) GetQueryTermsOk() (*[]QueryTermResultItem, bool)`

GetQueryTermsOk returns a tuple with the QueryTerms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryTerms

`func (o *OpenSearchAutosuggestResponse) SetQueryTerms(v []QueryTermResultItem)`

SetQueryTerms sets QueryTerms field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



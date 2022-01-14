# QueryTermResultItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Term** | **string** | The term that will be suggested to the user. | 
**Replaces** | **string** | The sub-string of the original query that is replaced by this Query Term. | 
**Start** | **int32** | The start index in codepoints (inclusive) of the text replaced in the original query. | 
**End** | **int32** | The end index in codepoints (exclusive) of the text replaced in the original query. | 

## Methods

### NewQueryTermResultItem

`func NewQueryTermResultItem(term string, replaces string, start int32, end int32, ) *QueryTermResultItem`

NewQueryTermResultItem instantiates a new QueryTermResultItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewQueryTermResultItemWithDefaults

`func NewQueryTermResultItemWithDefaults() *QueryTermResultItem`

NewQueryTermResultItemWithDefaults instantiates a new QueryTermResultItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTerm

`func (o *QueryTermResultItem) GetTerm() string`

GetTerm returns the Term field if non-nil, zero value otherwise.

### GetTermOk

`func (o *QueryTermResultItem) GetTermOk() (*string, bool)`

GetTermOk returns a tuple with the Term field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerm

`func (o *QueryTermResultItem) SetTerm(v string)`

SetTerm sets Term field to given value.


### GetReplaces

`func (o *QueryTermResultItem) GetReplaces() string`

GetReplaces returns the Replaces field if non-nil, zero value otherwise.

### GetReplacesOk

`func (o *QueryTermResultItem) GetReplacesOk() (*string, bool)`

GetReplacesOk returns a tuple with the Replaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplaces

`func (o *QueryTermResultItem) SetReplaces(v string)`

SetReplaces sets Replaces field to given value.


### GetStart

`func (o *QueryTermResultItem) GetStart() int32`

GetStart returns the Start field if non-nil, zero value otherwise.

### GetStartOk

`func (o *QueryTermResultItem) GetStartOk() (*int32, bool)`

GetStartOk returns a tuple with the Start field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStart

`func (o *QueryTermResultItem) SetStart(v int32)`

SetStart sets Start field to given value.


### GetEnd

`func (o *QueryTermResultItem) GetEnd() int32`

GetEnd returns the End field if non-nil, zero value otherwise.

### GetEndOk

`func (o *QueryTermResultItem) GetEndOk() (*int32, bool)`

GetEndOk returns a tuple with the End field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnd

`func (o *QueryTermResultItem) SetEnd(v int32)`

SetEnd sets End field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



# MatchInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Start** | **int32** | First index of the matched range (0-based indexing, inclusive) | 
**End** | **int32** | One past the last index of the matched range (0-based indexing, exclusive); The difference between end and start gives the length of the term | 
**Value** | **string** | Matched term in the input string | 
**Qq** | Pointer to **string** | The matched qualified query field type. If this is not returned, then matched value refers to the freetext query | [optional] 

## Methods

### NewMatchInfo

`func NewMatchInfo(start int32, end int32, value string, ) *MatchInfo`

NewMatchInfo instantiates a new MatchInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMatchInfoWithDefaults

`func NewMatchInfoWithDefaults() *MatchInfo`

NewMatchInfoWithDefaults instantiates a new MatchInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStart

`func (o *MatchInfo) GetStart() int32`

GetStart returns the Start field if non-nil, zero value otherwise.

### GetStartOk

`func (o *MatchInfo) GetStartOk() (*int32, bool)`

GetStartOk returns a tuple with the Start field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStart

`func (o *MatchInfo) SetStart(v int32)`

SetStart sets Start field to given value.


### GetEnd

`func (o *MatchInfo) GetEnd() int32`

GetEnd returns the End field if non-nil, zero value otherwise.

### GetEndOk

`func (o *MatchInfo) GetEndOk() (*int32, bool)`

GetEndOk returns a tuple with the End field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnd

`func (o *MatchInfo) SetEnd(v int32)`

SetEnd sets End field to given value.


### GetValue

`func (o *MatchInfo) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *MatchInfo) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *MatchInfo) SetValue(v string)`

SetValue sets Value field to given value.


### GetQq

`func (o *MatchInfo) GetQq() string`

GetQq returns the Qq field if non-nil, zero value otherwise.

### GetQqOk

`func (o *MatchInfo) GetQqOk() (*string, bool)`

GetQqOk returns a tuple with the Qq field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQq

`func (o *MatchInfo) SetQq(v string)`

SetQq sets Qq field to given value.

### HasQq

`func (o *MatchInfo) HasQq() bool`

HasQq returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



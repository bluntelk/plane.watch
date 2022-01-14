# TitleHighlighting

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | Pointer to [**[]Range**](Range.md) | Ranges of indexes that matched in the title attribute | [optional] 

## Methods

### NewTitleHighlighting

`func NewTitleHighlighting() *TitleHighlighting`

NewTitleHighlighting instantiates a new TitleHighlighting object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTitleHighlightingWithDefaults

`func NewTitleHighlightingWithDefaults() *TitleHighlighting`

NewTitleHighlightingWithDefaults instantiates a new TitleHighlighting object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *TitleHighlighting) GetTitle() []Range`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *TitleHighlighting) GetTitleOk() (*[]Range, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *TitleHighlighting) SetTitle(v []Range)`

SetTitle sets Title field to given value.

### HasTitle

`func (o *TitleHighlighting) HasTitle() bool`

HasTitle returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



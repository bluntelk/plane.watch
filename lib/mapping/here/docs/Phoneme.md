# Phoneme

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Value** | **string** | The actual phonetic transcription in the NT-SAMPA format. | 
**Language** | Pointer to **string** | The [BCP 47](https://en.wikipedia.org/wiki/IETF_language_tag) language code. | [optional] 
**Preferred** | Pointer to **bool** | Whether or not it is the preferred phoneme. | [optional] 

## Methods

### NewPhoneme

`func NewPhoneme(value string, ) *Phoneme`

NewPhoneme instantiates a new Phoneme object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPhonemeWithDefaults

`func NewPhonemeWithDefaults() *Phoneme`

NewPhonemeWithDefaults instantiates a new Phoneme object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetValue

`func (o *Phoneme) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *Phoneme) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *Phoneme) SetValue(v string)`

SetValue sets Value field to given value.


### GetLanguage

`func (o *Phoneme) GetLanguage() string`

GetLanguage returns the Language field if non-nil, zero value otherwise.

### GetLanguageOk

`func (o *Phoneme) GetLanguageOk() (*string, bool)`

GetLanguageOk returns a tuple with the Language field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguage

`func (o *Phoneme) SetLanguage(v string)`

SetLanguage sets Language field to given value.

### HasLanguage

`func (o *Phoneme) HasLanguage() bool`

HasLanguage returns a boolean if a field has been set.

### GetPreferred

`func (o *Phoneme) GetPreferred() bool`

GetPreferred returns the Preferred field if non-nil, zero value otherwise.

### GetPreferredOk

`func (o *Phoneme) GetPreferredOk() (*bool, bool)`

GetPreferredOk returns a tuple with the Preferred field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPreferred

`func (o *Phoneme) SetPreferred(v bool)`

SetPreferred sets Preferred field to given value.

### HasPreferred

`func (o *Phoneme) HasPreferred() bool`

HasPreferred returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



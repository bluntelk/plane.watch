# CountryInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Alpha2** | Pointer to **string** | [ISO 3166-1 alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) country code | [optional] 
**Alpha3** | Pointer to **string** | [ISO 3166-1 alpha-3](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3) country code | [optional] 

## Methods

### NewCountryInfo

`func NewCountryInfo() *CountryInfo`

NewCountryInfo instantiates a new CountryInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCountryInfoWithDefaults

`func NewCountryInfoWithDefaults() *CountryInfo`

NewCountryInfoWithDefaults instantiates a new CountryInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAlpha2

`func (o *CountryInfo) GetAlpha2() string`

GetAlpha2 returns the Alpha2 field if non-nil, zero value otherwise.

### GetAlpha2Ok

`func (o *CountryInfo) GetAlpha2Ok() (*string, bool)`

GetAlpha2Ok returns a tuple with the Alpha2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlpha2

`func (o *CountryInfo) SetAlpha2(v string)`

SetAlpha2 sets Alpha2 field to given value.

### HasAlpha2

`func (o *CountryInfo) HasAlpha2() bool`

HasAlpha2 returns a boolean if a field has been set.

### GetAlpha3

`func (o *CountryInfo) GetAlpha3() string`

GetAlpha3 returns the Alpha3 field if non-nil, zero value otherwise.

### GetAlpha3Ok

`func (o *CountryInfo) GetAlpha3Ok() (*string, bool)`

GetAlpha3Ok returns a tuple with the Alpha3 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAlpha3

`func (o *CountryInfo) SetAlpha3(v string)`

SetAlpha3 sets Alpha3 field to given value.

### HasAlpha3

`func (o *CountryInfo) HasAlpha3() bool`

HasAlpha3 returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



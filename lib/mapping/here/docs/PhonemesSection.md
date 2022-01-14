# PhonemesSection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PlaceName** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the name of the place. | [optional] 
**CountryName** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the county name. | [optional] 
**State** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the state name. | [optional] 
**County** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the county name. | [optional] 
**City** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the city name. | [optional] 
**District** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the district name. | [optional] 
**Subdistrict** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the subdistrict name. | [optional] 
**Street** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the street name. | [optional] 
**Block** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the block. | [optional] 
**Subblock** | Pointer to [**[]Phoneme**](Phoneme.md) | Phonemes for the sub-block. | [optional] 

## Methods

### NewPhonemesSection

`func NewPhonemesSection() *PhonemesSection`

NewPhonemesSection instantiates a new PhonemesSection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPhonemesSectionWithDefaults

`func NewPhonemesSectionWithDefaults() *PhonemesSection`

NewPhonemesSectionWithDefaults instantiates a new PhonemesSection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPlaceName

`func (o *PhonemesSection) GetPlaceName() []Phoneme`

GetPlaceName returns the PlaceName field if non-nil, zero value otherwise.

### GetPlaceNameOk

`func (o *PhonemesSection) GetPlaceNameOk() (*[]Phoneme, bool)`

GetPlaceNameOk returns a tuple with the PlaceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlaceName

`func (o *PhonemesSection) SetPlaceName(v []Phoneme)`

SetPlaceName sets PlaceName field to given value.

### HasPlaceName

`func (o *PhonemesSection) HasPlaceName() bool`

HasPlaceName returns a boolean if a field has been set.

### GetCountryName

`func (o *PhonemesSection) GetCountryName() []Phoneme`

GetCountryName returns the CountryName field if non-nil, zero value otherwise.

### GetCountryNameOk

`func (o *PhonemesSection) GetCountryNameOk() (*[]Phoneme, bool)`

GetCountryNameOk returns a tuple with the CountryName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryName

`func (o *PhonemesSection) SetCountryName(v []Phoneme)`

SetCountryName sets CountryName field to given value.

### HasCountryName

`func (o *PhonemesSection) HasCountryName() bool`

HasCountryName returns a boolean if a field has been set.

### GetState

`func (o *PhonemesSection) GetState() []Phoneme`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *PhonemesSection) GetStateOk() (*[]Phoneme, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *PhonemesSection) SetState(v []Phoneme)`

SetState sets State field to given value.

### HasState

`func (o *PhonemesSection) HasState() bool`

HasState returns a boolean if a field has been set.

### GetCounty

`func (o *PhonemesSection) GetCounty() []Phoneme`

GetCounty returns the County field if non-nil, zero value otherwise.

### GetCountyOk

`func (o *PhonemesSection) GetCountyOk() (*[]Phoneme, bool)`

GetCountyOk returns a tuple with the County field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCounty

`func (o *PhonemesSection) SetCounty(v []Phoneme)`

SetCounty sets County field to given value.

### HasCounty

`func (o *PhonemesSection) HasCounty() bool`

HasCounty returns a boolean if a field has been set.

### GetCity

`func (o *PhonemesSection) GetCity() []Phoneme`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *PhonemesSection) GetCityOk() (*[]Phoneme, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *PhonemesSection) SetCity(v []Phoneme)`

SetCity sets City field to given value.

### HasCity

`func (o *PhonemesSection) HasCity() bool`

HasCity returns a boolean if a field has been set.

### GetDistrict

`func (o *PhonemesSection) GetDistrict() []Phoneme`

GetDistrict returns the District field if non-nil, zero value otherwise.

### GetDistrictOk

`func (o *PhonemesSection) GetDistrictOk() (*[]Phoneme, bool)`

GetDistrictOk returns a tuple with the District field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistrict

`func (o *PhonemesSection) SetDistrict(v []Phoneme)`

SetDistrict sets District field to given value.

### HasDistrict

`func (o *PhonemesSection) HasDistrict() bool`

HasDistrict returns a boolean if a field has been set.

### GetSubdistrict

`func (o *PhonemesSection) GetSubdistrict() []Phoneme`

GetSubdistrict returns the Subdistrict field if non-nil, zero value otherwise.

### GetSubdistrictOk

`func (o *PhonemesSection) GetSubdistrictOk() (*[]Phoneme, bool)`

GetSubdistrictOk returns a tuple with the Subdistrict field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubdistrict

`func (o *PhonemesSection) SetSubdistrict(v []Phoneme)`

SetSubdistrict sets Subdistrict field to given value.

### HasSubdistrict

`func (o *PhonemesSection) HasSubdistrict() bool`

HasSubdistrict returns a boolean if a field has been set.

### GetStreet

`func (o *PhonemesSection) GetStreet() []Phoneme`

GetStreet returns the Street field if non-nil, zero value otherwise.

### GetStreetOk

`func (o *PhonemesSection) GetStreetOk() (*[]Phoneme, bool)`

GetStreetOk returns a tuple with the Street field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreet

`func (o *PhonemesSection) SetStreet(v []Phoneme)`

SetStreet sets Street field to given value.

### HasStreet

`func (o *PhonemesSection) HasStreet() bool`

HasStreet returns a boolean if a field has been set.

### GetBlock

`func (o *PhonemesSection) GetBlock() []Phoneme`

GetBlock returns the Block field if non-nil, zero value otherwise.

### GetBlockOk

`func (o *PhonemesSection) GetBlockOk() (*[]Phoneme, bool)`

GetBlockOk returns a tuple with the Block field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlock

`func (o *PhonemesSection) SetBlock(v []Phoneme)`

SetBlock sets Block field to given value.

### HasBlock

`func (o *PhonemesSection) HasBlock() bool`

HasBlock returns a boolean if a field has been set.

### GetSubblock

`func (o *PhonemesSection) GetSubblock() []Phoneme`

GetSubblock returns the Subblock field if non-nil, zero value otherwise.

### GetSubblockOk

`func (o *PhonemesSection) GetSubblockOk() (*[]Phoneme, bool)`

GetSubblockOk returns a tuple with the Subblock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubblock

`func (o *PhonemesSection) SetSubblock(v []Phoneme)`

SetSubblock sets Subblock field to given value.

### HasSubblock

`func (o *PhonemesSection) HasSubblock() bool`

HasSubblock returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



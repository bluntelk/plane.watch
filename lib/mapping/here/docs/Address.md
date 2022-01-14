# Address

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Label** | Pointer to **string** | Assembled address value built out of the address components according to the regional postal rules. These are the same rules for all endpoints. It may not include all the input terms. For example: \&quot;Schulstraße 4, 32547 Bad Oeynhausen, Germany\&quot; | [optional] 
**CountryCode** | Pointer to **string** | A three-letter country code. For example: \&quot;DEU\&quot; | [optional] 
**CountryName** | Pointer to **string** | The localised country name. For example: \&quot;Deutschland\&quot; | [optional] 
**StateCode** | Pointer to **string** | A state code or state name abbreviation – country specific. For example, in the United States it is the two letter state abbreviation: \&quot;CA\&quot; for California. | [optional] 
**State** | Pointer to **string** | The state division of a country. For example: \&quot;North Rhine-Westphalia\&quot; | [optional] 
**CountyCode** | Pointer to **string** | A county code or county name abbreviation – country specific. For example, for Italy it is the province abbreviation: \&quot;RM\&quot; for Rome. | [optional] 
**County** | Pointer to **string** | A division of a state; typically, a secondary-level administrative division of a country or equivalent. | [optional] 
**City** | Pointer to **string** | The name of the primary locality of the place. For example: \&quot;Bad Oyenhausen\&quot; | [optional] 
**District** | Pointer to **string** | A division of city; typically an administrative unit within a larger city or a customary name of a city&#39;s neighborhood. For example: \&quot;Bad Oyenhausen\&quot; | [optional] 
**Subdistrict** | Pointer to **string** | A subdivision of a district. For example: \&quot;Minden-Lübbecke\&quot; | [optional] 
**Street** | Pointer to **string** | Name of street. For example: \&quot;Schulstrasse\&quot; | [optional] 
**Block** | Pointer to **string** | Name of block. | [optional] 
**Subblock** | Pointer to **string** | Name of sub-block. | [optional] 
**PostalCode** | Pointer to **string** | An alphanumeric string included in a postal address to facilitate mail sorting, such as post code, postcode, or ZIP code. For example: \&quot;32547\&quot; | [optional] 
**HouseNumber** | Pointer to **string** | House number. For example: \&quot;4\&quot; | [optional] 
**Building** | Pointer to **string** | Name of building. | [optional] 

## Methods

### NewAddress

`func NewAddress() *Address`

NewAddress instantiates a new Address object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAddressWithDefaults

`func NewAddressWithDefaults() *Address`

NewAddressWithDefaults instantiates a new Address object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLabel

`func (o *Address) GetLabel() string`

GetLabel returns the Label field if non-nil, zero value otherwise.

### GetLabelOk

`func (o *Address) GetLabelOk() (*string, bool)`

GetLabelOk returns a tuple with the Label field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLabel

`func (o *Address) SetLabel(v string)`

SetLabel sets Label field to given value.

### HasLabel

`func (o *Address) HasLabel() bool`

HasLabel returns a boolean if a field has been set.

### GetCountryCode

`func (o *Address) GetCountryCode() string`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *Address) GetCountryCodeOk() (*string, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *Address) SetCountryCode(v string)`

SetCountryCode sets CountryCode field to given value.

### HasCountryCode

`func (o *Address) HasCountryCode() bool`

HasCountryCode returns a boolean if a field has been set.

### GetCountryName

`func (o *Address) GetCountryName() string`

GetCountryName returns the CountryName field if non-nil, zero value otherwise.

### GetCountryNameOk

`func (o *Address) GetCountryNameOk() (*string, bool)`

GetCountryNameOk returns a tuple with the CountryName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryName

`func (o *Address) SetCountryName(v string)`

SetCountryName sets CountryName field to given value.

### HasCountryName

`func (o *Address) HasCountryName() bool`

HasCountryName returns a boolean if a field has been set.

### GetStateCode

`func (o *Address) GetStateCode() string`

GetStateCode returns the StateCode field if non-nil, zero value otherwise.

### GetStateCodeOk

`func (o *Address) GetStateCodeOk() (*string, bool)`

GetStateCodeOk returns a tuple with the StateCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateCode

`func (o *Address) SetStateCode(v string)`

SetStateCode sets StateCode field to given value.

### HasStateCode

`func (o *Address) HasStateCode() bool`

HasStateCode returns a boolean if a field has been set.

### GetState

`func (o *Address) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Address) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Address) SetState(v string)`

SetState sets State field to given value.

### HasState

`func (o *Address) HasState() bool`

HasState returns a boolean if a field has been set.

### GetCountyCode

`func (o *Address) GetCountyCode() string`

GetCountyCode returns the CountyCode field if non-nil, zero value otherwise.

### GetCountyCodeOk

`func (o *Address) GetCountyCodeOk() (*string, bool)`

GetCountyCodeOk returns a tuple with the CountyCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountyCode

`func (o *Address) SetCountyCode(v string)`

SetCountyCode sets CountyCode field to given value.

### HasCountyCode

`func (o *Address) HasCountyCode() bool`

HasCountyCode returns a boolean if a field has been set.

### GetCounty

`func (o *Address) GetCounty() string`

GetCounty returns the County field if non-nil, zero value otherwise.

### GetCountyOk

`func (o *Address) GetCountyOk() (*string, bool)`

GetCountyOk returns a tuple with the County field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCounty

`func (o *Address) SetCounty(v string)`

SetCounty sets County field to given value.

### HasCounty

`func (o *Address) HasCounty() bool`

HasCounty returns a boolean if a field has been set.

### GetCity

`func (o *Address) GetCity() string`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *Address) GetCityOk() (*string, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *Address) SetCity(v string)`

SetCity sets City field to given value.

### HasCity

`func (o *Address) HasCity() bool`

HasCity returns a boolean if a field has been set.

### GetDistrict

`func (o *Address) GetDistrict() string`

GetDistrict returns the District field if non-nil, zero value otherwise.

### GetDistrictOk

`func (o *Address) GetDistrictOk() (*string, bool)`

GetDistrictOk returns a tuple with the District field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistrict

`func (o *Address) SetDistrict(v string)`

SetDistrict sets District field to given value.

### HasDistrict

`func (o *Address) HasDistrict() bool`

HasDistrict returns a boolean if a field has been set.

### GetSubdistrict

`func (o *Address) GetSubdistrict() string`

GetSubdistrict returns the Subdistrict field if non-nil, zero value otherwise.

### GetSubdistrictOk

`func (o *Address) GetSubdistrictOk() (*string, bool)`

GetSubdistrictOk returns a tuple with the Subdistrict field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubdistrict

`func (o *Address) SetSubdistrict(v string)`

SetSubdistrict sets Subdistrict field to given value.

### HasSubdistrict

`func (o *Address) HasSubdistrict() bool`

HasSubdistrict returns a boolean if a field has been set.

### GetStreet

`func (o *Address) GetStreet() string`

GetStreet returns the Street field if non-nil, zero value otherwise.

### GetStreetOk

`func (o *Address) GetStreetOk() (*string, bool)`

GetStreetOk returns a tuple with the Street field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreet

`func (o *Address) SetStreet(v string)`

SetStreet sets Street field to given value.

### HasStreet

`func (o *Address) HasStreet() bool`

HasStreet returns a boolean if a field has been set.

### GetBlock

`func (o *Address) GetBlock() string`

GetBlock returns the Block field if non-nil, zero value otherwise.

### GetBlockOk

`func (o *Address) GetBlockOk() (*string, bool)`

GetBlockOk returns a tuple with the Block field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlock

`func (o *Address) SetBlock(v string)`

SetBlock sets Block field to given value.

### HasBlock

`func (o *Address) HasBlock() bool`

HasBlock returns a boolean if a field has been set.

### GetSubblock

`func (o *Address) GetSubblock() string`

GetSubblock returns the Subblock field if non-nil, zero value otherwise.

### GetSubblockOk

`func (o *Address) GetSubblockOk() (*string, bool)`

GetSubblockOk returns a tuple with the Subblock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubblock

`func (o *Address) SetSubblock(v string)`

SetSubblock sets Subblock field to given value.

### HasSubblock

`func (o *Address) HasSubblock() bool`

HasSubblock returns a boolean if a field has been set.

### GetPostalCode

`func (o *Address) GetPostalCode() string`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *Address) GetPostalCodeOk() (*string, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *Address) SetPostalCode(v string)`

SetPostalCode sets PostalCode field to given value.

### HasPostalCode

`func (o *Address) HasPostalCode() bool`

HasPostalCode returns a boolean if a field has been set.

### GetHouseNumber

`func (o *Address) GetHouseNumber() string`

GetHouseNumber returns the HouseNumber field if non-nil, zero value otherwise.

### GetHouseNumberOk

`func (o *Address) GetHouseNumberOk() (*string, bool)`

GetHouseNumberOk returns a tuple with the HouseNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumber

`func (o *Address) SetHouseNumber(v string)`

SetHouseNumber sets HouseNumber field to given value.

### HasHouseNumber

`func (o *Address) HasHouseNumber() bool`

HasHouseNumber returns a boolean if a field has been set.

### GetBuilding

`func (o *Address) GetBuilding() string`

GetBuilding returns the Building field if non-nil, zero value otherwise.

### GetBuildingOk

`func (o *Address) GetBuildingOk() (*string, bool)`

GetBuildingOk returns a tuple with the Building field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuilding

`func (o *Address) SetBuilding(v string)`

SetBuilding sets Building field to given value.

### HasBuilding

`func (o *Address) HasBuilding() bool`

HasBuilding returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



# FieldScore

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Country** | Pointer to **float64** | Indicates how good the result country name or [ISO 3166-1 alpha-3] country code matches to the freeform or qualified input. | [optional] 
**CountryCode** | Pointer to **float64** | Indicates how good the result [ISO 3166-1 alpha-3] country code matches to the freeform or qualified input. | [optional] 
**State** | Pointer to **float64** | Indicates how good the result state name matches to the freeform or qualified input. | [optional] 
**StateCode** | Pointer to **float64** | Indicates how good the result state code matches to the freeform or qualified input. | [optional] 
**County** | Pointer to **float64** | Indicates how good the result county name matches to the freeform or qualified input. | [optional] 
**CountyCode** | Pointer to **float64** | Indicates how good the result county code matches to the freeform or qualified input. | [optional] 
**City** | Pointer to **float64** | Indicates how good the result city name matches to the freeform or qualified input. | [optional] 
**District** | Pointer to **float64** | Indicates how good the result district name matches to the freeform or qualified input. | [optional] 
**Subdistrict** | Pointer to **float64** | Indicates how good the result sub-district name matches to the freeform or qualified input. | [optional] 
**Streets** | Pointer to **[]float64** | Indicates how good the result street names match to the freeform or qualified input. If the input contains multiple street names, the field score is calculated and returned for each of them individually. | [optional] 
**Block** | Pointer to **float64** | Indicates how good the result block name matches to the freeform or qualified input. | [optional] 
**Subblock** | Pointer to **float64** | Indicates how good the result sub-block name matches to the freeform or qualified input. | [optional] 
**HouseNumber** | Pointer to **float64** | Indicates how good the result house number matches to the freeform or qualified input. It may happen, that the house number, which one is looking for, is not yet in the map data. For such cases, the /geocode returns the nearest known house number on the same street. This represents the numeric difference between the requested and the returned house numbers. | [optional] 
**PostalCode** | Pointer to **float64** | Indicates how good the result postal code matches to the freeform or qualified input. | [optional] 
**Building** | Pointer to **float64** | Indicates how good the result building name matches to the freeform or qualified input. | [optional] 
**Unit** | Pointer to **float64** | Indicates how good the result unit (such as a micro point address) matches to the freeform or qualified input. | [optional] 
**PlaceName** | Pointer to **float64** | Indicates how good the result place name matches to the freeform or qualified input. | [optional] 
**OntologyName** | Pointer to **float64** | Indicates how good the result ontology name matches to the freeform or qualified input. | [optional] 

## Methods

### NewFieldScore

`func NewFieldScore() *FieldScore`

NewFieldScore instantiates a new FieldScore object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFieldScoreWithDefaults

`func NewFieldScoreWithDefaults() *FieldScore`

NewFieldScoreWithDefaults instantiates a new FieldScore object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCountry

`func (o *FieldScore) GetCountry() float64`

GetCountry returns the Country field if non-nil, zero value otherwise.

### GetCountryOk

`func (o *FieldScore) GetCountryOk() (*float64, bool)`

GetCountryOk returns a tuple with the Country field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountry

`func (o *FieldScore) SetCountry(v float64)`

SetCountry sets Country field to given value.

### HasCountry

`func (o *FieldScore) HasCountry() bool`

HasCountry returns a boolean if a field has been set.

### GetCountryCode

`func (o *FieldScore) GetCountryCode() float64`

GetCountryCode returns the CountryCode field if non-nil, zero value otherwise.

### GetCountryCodeOk

`func (o *FieldScore) GetCountryCodeOk() (*float64, bool)`

GetCountryCodeOk returns a tuple with the CountryCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryCode

`func (o *FieldScore) SetCountryCode(v float64)`

SetCountryCode sets CountryCode field to given value.

### HasCountryCode

`func (o *FieldScore) HasCountryCode() bool`

HasCountryCode returns a boolean if a field has been set.

### GetState

`func (o *FieldScore) GetState() float64`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *FieldScore) GetStateOk() (*float64, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *FieldScore) SetState(v float64)`

SetState sets State field to given value.

### HasState

`func (o *FieldScore) HasState() bool`

HasState returns a boolean if a field has been set.

### GetStateCode

`func (o *FieldScore) GetStateCode() float64`

GetStateCode returns the StateCode field if non-nil, zero value otherwise.

### GetStateCodeOk

`func (o *FieldScore) GetStateCodeOk() (*float64, bool)`

GetStateCodeOk returns a tuple with the StateCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStateCode

`func (o *FieldScore) SetStateCode(v float64)`

SetStateCode sets StateCode field to given value.

### HasStateCode

`func (o *FieldScore) HasStateCode() bool`

HasStateCode returns a boolean if a field has been set.

### GetCounty

`func (o *FieldScore) GetCounty() float64`

GetCounty returns the County field if non-nil, zero value otherwise.

### GetCountyOk

`func (o *FieldScore) GetCountyOk() (*float64, bool)`

GetCountyOk returns a tuple with the County field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCounty

`func (o *FieldScore) SetCounty(v float64)`

SetCounty sets County field to given value.

### HasCounty

`func (o *FieldScore) HasCounty() bool`

HasCounty returns a boolean if a field has been set.

### GetCountyCode

`func (o *FieldScore) GetCountyCode() float64`

GetCountyCode returns the CountyCode field if non-nil, zero value otherwise.

### GetCountyCodeOk

`func (o *FieldScore) GetCountyCodeOk() (*float64, bool)`

GetCountyCodeOk returns a tuple with the CountyCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountyCode

`func (o *FieldScore) SetCountyCode(v float64)`

SetCountyCode sets CountyCode field to given value.

### HasCountyCode

`func (o *FieldScore) HasCountyCode() bool`

HasCountyCode returns a boolean if a field has been set.

### GetCity

`func (o *FieldScore) GetCity() float64`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *FieldScore) GetCityOk() (*float64, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *FieldScore) SetCity(v float64)`

SetCity sets City field to given value.

### HasCity

`func (o *FieldScore) HasCity() bool`

HasCity returns a boolean if a field has been set.

### GetDistrict

`func (o *FieldScore) GetDistrict() float64`

GetDistrict returns the District field if non-nil, zero value otherwise.

### GetDistrictOk

`func (o *FieldScore) GetDistrictOk() (*float64, bool)`

GetDistrictOk returns a tuple with the District field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistrict

`func (o *FieldScore) SetDistrict(v float64)`

SetDistrict sets District field to given value.

### HasDistrict

`func (o *FieldScore) HasDistrict() bool`

HasDistrict returns a boolean if a field has been set.

### GetSubdistrict

`func (o *FieldScore) GetSubdistrict() float64`

GetSubdistrict returns the Subdistrict field if non-nil, zero value otherwise.

### GetSubdistrictOk

`func (o *FieldScore) GetSubdistrictOk() (*float64, bool)`

GetSubdistrictOk returns a tuple with the Subdistrict field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubdistrict

`func (o *FieldScore) SetSubdistrict(v float64)`

SetSubdistrict sets Subdistrict field to given value.

### HasSubdistrict

`func (o *FieldScore) HasSubdistrict() bool`

HasSubdistrict returns a boolean if a field has been set.

### GetStreets

`func (o *FieldScore) GetStreets() []float64`

GetStreets returns the Streets field if non-nil, zero value otherwise.

### GetStreetsOk

`func (o *FieldScore) GetStreetsOk() (*[]float64, bool)`

GetStreetsOk returns a tuple with the Streets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreets

`func (o *FieldScore) SetStreets(v []float64)`

SetStreets sets Streets field to given value.

### HasStreets

`func (o *FieldScore) HasStreets() bool`

HasStreets returns a boolean if a field has been set.

### GetBlock

`func (o *FieldScore) GetBlock() float64`

GetBlock returns the Block field if non-nil, zero value otherwise.

### GetBlockOk

`func (o *FieldScore) GetBlockOk() (*float64, bool)`

GetBlockOk returns a tuple with the Block field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlock

`func (o *FieldScore) SetBlock(v float64)`

SetBlock sets Block field to given value.

### HasBlock

`func (o *FieldScore) HasBlock() bool`

HasBlock returns a boolean if a field has been set.

### GetSubblock

`func (o *FieldScore) GetSubblock() float64`

GetSubblock returns the Subblock field if non-nil, zero value otherwise.

### GetSubblockOk

`func (o *FieldScore) GetSubblockOk() (*float64, bool)`

GetSubblockOk returns a tuple with the Subblock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubblock

`func (o *FieldScore) SetSubblock(v float64)`

SetSubblock sets Subblock field to given value.

### HasSubblock

`func (o *FieldScore) HasSubblock() bool`

HasSubblock returns a boolean if a field has been set.

### GetHouseNumber

`func (o *FieldScore) GetHouseNumber() float64`

GetHouseNumber returns the HouseNumber field if non-nil, zero value otherwise.

### GetHouseNumberOk

`func (o *FieldScore) GetHouseNumberOk() (*float64, bool)`

GetHouseNumberOk returns a tuple with the HouseNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumber

`func (o *FieldScore) SetHouseNumber(v float64)`

SetHouseNumber sets HouseNumber field to given value.

### HasHouseNumber

`func (o *FieldScore) HasHouseNumber() bool`

HasHouseNumber returns a boolean if a field has been set.

### GetPostalCode

`func (o *FieldScore) GetPostalCode() float64`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *FieldScore) GetPostalCodeOk() (*float64, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *FieldScore) SetPostalCode(v float64)`

SetPostalCode sets PostalCode field to given value.

### HasPostalCode

`func (o *FieldScore) HasPostalCode() bool`

HasPostalCode returns a boolean if a field has been set.

### GetBuilding

`func (o *FieldScore) GetBuilding() float64`

GetBuilding returns the Building field if non-nil, zero value otherwise.

### GetBuildingOk

`func (o *FieldScore) GetBuildingOk() (*float64, bool)`

GetBuildingOk returns a tuple with the Building field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuilding

`func (o *FieldScore) SetBuilding(v float64)`

SetBuilding sets Building field to given value.

### HasBuilding

`func (o *FieldScore) HasBuilding() bool`

HasBuilding returns a boolean if a field has been set.

### GetUnit

`func (o *FieldScore) GetUnit() float64`

GetUnit returns the Unit field if non-nil, zero value otherwise.

### GetUnitOk

`func (o *FieldScore) GetUnitOk() (*float64, bool)`

GetUnitOk returns a tuple with the Unit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnit

`func (o *FieldScore) SetUnit(v float64)`

SetUnit sets Unit field to given value.

### HasUnit

`func (o *FieldScore) HasUnit() bool`

HasUnit returns a boolean if a field has been set.

### GetPlaceName

`func (o *FieldScore) GetPlaceName() float64`

GetPlaceName returns the PlaceName field if non-nil, zero value otherwise.

### GetPlaceNameOk

`func (o *FieldScore) GetPlaceNameOk() (*float64, bool)`

GetPlaceNameOk returns a tuple with the PlaceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlaceName

`func (o *FieldScore) SetPlaceName(v float64)`

SetPlaceName sets PlaceName field to given value.

### HasPlaceName

`func (o *FieldScore) HasPlaceName() bool`

HasPlaceName returns a boolean if a field has been set.

### GetOntologyName

`func (o *FieldScore) GetOntologyName() float64`

GetOntologyName returns the OntologyName field if non-nil, zero value otherwise.

### GetOntologyNameOk

`func (o *FieldScore) GetOntologyNameOk() (*float64, bool)`

GetOntologyNameOk returns a tuple with the OntologyName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOntologyName

`func (o *FieldScore) SetOntologyName(v float64)`

SetOntologyName sets OntologyName field to given value.

### HasOntologyName

`func (o *FieldScore) HasOntologyName() bool`

HasOntologyName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



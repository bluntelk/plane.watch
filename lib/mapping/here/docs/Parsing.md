# Parsing

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PlaceName** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Place name matches | [optional] 
**Country** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Country matches | [optional] 
**State** | Pointer to [**[]MatchInfo**](MatchInfo.md) | State matches | [optional] 
**County** | Pointer to [**[]MatchInfo**](MatchInfo.md) | County matches | [optional] 
**City** | Pointer to [**[]MatchInfo**](MatchInfo.md) | City matches | [optional] 
**District** | Pointer to [**[]MatchInfo**](MatchInfo.md) | District matches | [optional] 
**Subdistrict** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Subdistrict matches | [optional] 
**Street** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Street matches | [optional] 
**Block** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Block matches | [optional] 
**Subblock** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Subblock matches | [optional] 
**HouseNumber** | Pointer to [**[]MatchInfo**](MatchInfo.md) | HouseNumber matches | [optional] 
**PostalCode** | Pointer to [**[]MatchInfo**](MatchInfo.md) | PostalCode matches | [optional] 
**Building** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Building matches | [optional] 
**SecondaryUnits** | Pointer to [**[]MatchInfo**](MatchInfo.md) | secondaryUnits matches | [optional] 
**OntologyName** | Pointer to [**[]MatchInfo**](MatchInfo.md) | Ontology name matches | [optional] 

## Methods

### NewParsing

`func NewParsing() *Parsing`

NewParsing instantiates a new Parsing object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewParsingWithDefaults

`func NewParsingWithDefaults() *Parsing`

NewParsingWithDefaults instantiates a new Parsing object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPlaceName

`func (o *Parsing) GetPlaceName() []MatchInfo`

GetPlaceName returns the PlaceName field if non-nil, zero value otherwise.

### GetPlaceNameOk

`func (o *Parsing) GetPlaceNameOk() (*[]MatchInfo, bool)`

GetPlaceNameOk returns a tuple with the PlaceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlaceName

`func (o *Parsing) SetPlaceName(v []MatchInfo)`

SetPlaceName sets PlaceName field to given value.

### HasPlaceName

`func (o *Parsing) HasPlaceName() bool`

HasPlaceName returns a boolean if a field has been set.

### GetCountry

`func (o *Parsing) GetCountry() []MatchInfo`

GetCountry returns the Country field if non-nil, zero value otherwise.

### GetCountryOk

`func (o *Parsing) GetCountryOk() (*[]MatchInfo, bool)`

GetCountryOk returns a tuple with the Country field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountry

`func (o *Parsing) SetCountry(v []MatchInfo)`

SetCountry sets Country field to given value.

### HasCountry

`func (o *Parsing) HasCountry() bool`

HasCountry returns a boolean if a field has been set.

### GetState

`func (o *Parsing) GetState() []MatchInfo`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *Parsing) GetStateOk() (*[]MatchInfo, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *Parsing) SetState(v []MatchInfo)`

SetState sets State field to given value.

### HasState

`func (o *Parsing) HasState() bool`

HasState returns a boolean if a field has been set.

### GetCounty

`func (o *Parsing) GetCounty() []MatchInfo`

GetCounty returns the County field if non-nil, zero value otherwise.

### GetCountyOk

`func (o *Parsing) GetCountyOk() (*[]MatchInfo, bool)`

GetCountyOk returns a tuple with the County field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCounty

`func (o *Parsing) SetCounty(v []MatchInfo)`

SetCounty sets County field to given value.

### HasCounty

`func (o *Parsing) HasCounty() bool`

HasCounty returns a boolean if a field has been set.

### GetCity

`func (o *Parsing) GetCity() []MatchInfo`

GetCity returns the City field if non-nil, zero value otherwise.

### GetCityOk

`func (o *Parsing) GetCityOk() (*[]MatchInfo, bool)`

GetCityOk returns a tuple with the City field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCity

`func (o *Parsing) SetCity(v []MatchInfo)`

SetCity sets City field to given value.

### HasCity

`func (o *Parsing) HasCity() bool`

HasCity returns a boolean if a field has been set.

### GetDistrict

`func (o *Parsing) GetDistrict() []MatchInfo`

GetDistrict returns the District field if non-nil, zero value otherwise.

### GetDistrictOk

`func (o *Parsing) GetDistrictOk() (*[]MatchInfo, bool)`

GetDistrictOk returns a tuple with the District field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistrict

`func (o *Parsing) SetDistrict(v []MatchInfo)`

SetDistrict sets District field to given value.

### HasDistrict

`func (o *Parsing) HasDistrict() bool`

HasDistrict returns a boolean if a field has been set.

### GetSubdistrict

`func (o *Parsing) GetSubdistrict() []MatchInfo`

GetSubdistrict returns the Subdistrict field if non-nil, zero value otherwise.

### GetSubdistrictOk

`func (o *Parsing) GetSubdistrictOk() (*[]MatchInfo, bool)`

GetSubdistrictOk returns a tuple with the Subdistrict field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubdistrict

`func (o *Parsing) SetSubdistrict(v []MatchInfo)`

SetSubdistrict sets Subdistrict field to given value.

### HasSubdistrict

`func (o *Parsing) HasSubdistrict() bool`

HasSubdistrict returns a boolean if a field has been set.

### GetStreet

`func (o *Parsing) GetStreet() []MatchInfo`

GetStreet returns the Street field if non-nil, zero value otherwise.

### GetStreetOk

`func (o *Parsing) GetStreetOk() (*[]MatchInfo, bool)`

GetStreetOk returns a tuple with the Street field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreet

`func (o *Parsing) SetStreet(v []MatchInfo)`

SetStreet sets Street field to given value.

### HasStreet

`func (o *Parsing) HasStreet() bool`

HasStreet returns a boolean if a field has been set.

### GetBlock

`func (o *Parsing) GetBlock() []MatchInfo`

GetBlock returns the Block field if non-nil, zero value otherwise.

### GetBlockOk

`func (o *Parsing) GetBlockOk() (*[]MatchInfo, bool)`

GetBlockOk returns a tuple with the Block field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlock

`func (o *Parsing) SetBlock(v []MatchInfo)`

SetBlock sets Block field to given value.

### HasBlock

`func (o *Parsing) HasBlock() bool`

HasBlock returns a boolean if a field has been set.

### GetSubblock

`func (o *Parsing) GetSubblock() []MatchInfo`

GetSubblock returns the Subblock field if non-nil, zero value otherwise.

### GetSubblockOk

`func (o *Parsing) GetSubblockOk() (*[]MatchInfo, bool)`

GetSubblockOk returns a tuple with the Subblock field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubblock

`func (o *Parsing) SetSubblock(v []MatchInfo)`

SetSubblock sets Subblock field to given value.

### HasSubblock

`func (o *Parsing) HasSubblock() bool`

HasSubblock returns a boolean if a field has been set.

### GetHouseNumber

`func (o *Parsing) GetHouseNumber() []MatchInfo`

GetHouseNumber returns the HouseNumber field if non-nil, zero value otherwise.

### GetHouseNumberOk

`func (o *Parsing) GetHouseNumberOk() (*[]MatchInfo, bool)`

GetHouseNumberOk returns a tuple with the HouseNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumber

`func (o *Parsing) SetHouseNumber(v []MatchInfo)`

SetHouseNumber sets HouseNumber field to given value.

### HasHouseNumber

`func (o *Parsing) HasHouseNumber() bool`

HasHouseNumber returns a boolean if a field has been set.

### GetPostalCode

`func (o *Parsing) GetPostalCode() []MatchInfo`

GetPostalCode returns the PostalCode field if non-nil, zero value otherwise.

### GetPostalCodeOk

`func (o *Parsing) GetPostalCodeOk() (*[]MatchInfo, bool)`

GetPostalCodeOk returns a tuple with the PostalCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPostalCode

`func (o *Parsing) SetPostalCode(v []MatchInfo)`

SetPostalCode sets PostalCode field to given value.

### HasPostalCode

`func (o *Parsing) HasPostalCode() bool`

HasPostalCode returns a boolean if a field has been set.

### GetBuilding

`func (o *Parsing) GetBuilding() []MatchInfo`

GetBuilding returns the Building field if non-nil, zero value otherwise.

### GetBuildingOk

`func (o *Parsing) GetBuildingOk() (*[]MatchInfo, bool)`

GetBuildingOk returns a tuple with the Building field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuilding

`func (o *Parsing) SetBuilding(v []MatchInfo)`

SetBuilding sets Building field to given value.

### HasBuilding

`func (o *Parsing) HasBuilding() bool`

HasBuilding returns a boolean if a field has been set.

### GetSecondaryUnits

`func (o *Parsing) GetSecondaryUnits() []MatchInfo`

GetSecondaryUnits returns the SecondaryUnits field if non-nil, zero value otherwise.

### GetSecondaryUnitsOk

`func (o *Parsing) GetSecondaryUnitsOk() (*[]MatchInfo, bool)`

GetSecondaryUnitsOk returns a tuple with the SecondaryUnits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecondaryUnits

`func (o *Parsing) SetSecondaryUnits(v []MatchInfo)`

SetSecondaryUnits sets SecondaryUnits field to given value.

### HasSecondaryUnits

`func (o *Parsing) HasSecondaryUnits() bool`

HasSecondaryUnits returns a boolean if a field has been set.

### GetOntologyName

`func (o *Parsing) GetOntologyName() []MatchInfo`

GetOntologyName returns the OntologyName field if non-nil, zero value otherwise.

### GetOntologyNameOk

`func (o *Parsing) GetOntologyNameOk() (*[]MatchInfo, bool)`

GetOntologyNameOk returns a tuple with the OntologyName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOntologyName

`func (o *Parsing) SetOntologyName(v []MatchInfo)`

SetOntologyName sets OntologyName field to given value.

### HasOntologyName

`func (o *Parsing) HasOntologyName() bool`

HasOntologyName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



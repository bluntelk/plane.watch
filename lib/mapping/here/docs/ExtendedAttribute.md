# ExtendedAttribute

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EvStation** | Pointer to [**EvChargingAttributes**](EvChargingAttributes.md) | EV charging pool information | [optional] 

## Methods

### NewExtendedAttribute

`func NewExtendedAttribute() *ExtendedAttribute`

NewExtendedAttribute instantiates a new ExtendedAttribute object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExtendedAttributeWithDefaults

`func NewExtendedAttributeWithDefaults() *ExtendedAttribute`

NewExtendedAttributeWithDefaults instantiates a new ExtendedAttribute object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEvStation

`func (o *ExtendedAttribute) GetEvStation() EvChargingAttributes`

GetEvStation returns the EvStation field if non-nil, zero value otherwise.

### GetEvStationOk

`func (o *ExtendedAttribute) GetEvStationOk() (*EvChargingAttributes, bool)`

GetEvStationOk returns a tuple with the EvStation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvStation

`func (o *ExtendedAttribute) SetEvStation(v EvChargingAttributes)`

SetEvStation sets EvStation field to given value.

### HasEvStation

`func (o *ExtendedAttribute) HasEvStation() bool`

HasEvStation returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



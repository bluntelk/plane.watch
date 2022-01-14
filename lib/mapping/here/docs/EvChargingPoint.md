# EvChargingPoint

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**NumberOfConnectors** | Pointer to **int32** | Number of physical connectors in the connectors group | [optional] 
**ChargeMode** | Pointer to **string** | Charge mode of the connectors group. For more information, check the [IEC-61851-1](https://en.wikipedia.org/w/index.php?title&#x3D;Charging_station&amp;oldid&#x3D;1013010605#IEC-61851-1_Charging_Modes) standard. | [optional] 
**VoltsRange** | Pointer to **string** | Voltage provided by the connectors group | [optional] 
**Phases** | Pointer to **int32** | Number of phases provided by the connectors group | [optional] 
**AmpsRange** | Pointer to **string** | Amperage provided by the connectors group | [optional] 

## Methods

### NewEvChargingPoint

`func NewEvChargingPoint() *EvChargingPoint`

NewEvChargingPoint instantiates a new EvChargingPoint object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEvChargingPointWithDefaults

`func NewEvChargingPointWithDefaults() *EvChargingPoint`

NewEvChargingPointWithDefaults instantiates a new EvChargingPoint object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetNumberOfConnectors

`func (o *EvChargingPoint) GetNumberOfConnectors() int32`

GetNumberOfConnectors returns the NumberOfConnectors field if non-nil, zero value otherwise.

### GetNumberOfConnectorsOk

`func (o *EvChargingPoint) GetNumberOfConnectorsOk() (*int32, bool)`

GetNumberOfConnectorsOk returns a tuple with the NumberOfConnectors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfConnectors

`func (o *EvChargingPoint) SetNumberOfConnectors(v int32)`

SetNumberOfConnectors sets NumberOfConnectors field to given value.

### HasNumberOfConnectors

`func (o *EvChargingPoint) HasNumberOfConnectors() bool`

HasNumberOfConnectors returns a boolean if a field has been set.

### GetChargeMode

`func (o *EvChargingPoint) GetChargeMode() string`

GetChargeMode returns the ChargeMode field if non-nil, zero value otherwise.

### GetChargeModeOk

`func (o *EvChargingPoint) GetChargeModeOk() (*string, bool)`

GetChargeModeOk returns a tuple with the ChargeMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChargeMode

`func (o *EvChargingPoint) SetChargeMode(v string)`

SetChargeMode sets ChargeMode field to given value.

### HasChargeMode

`func (o *EvChargingPoint) HasChargeMode() bool`

HasChargeMode returns a boolean if a field has been set.

### GetVoltsRange

`func (o *EvChargingPoint) GetVoltsRange() string`

GetVoltsRange returns the VoltsRange field if non-nil, zero value otherwise.

### GetVoltsRangeOk

`func (o *EvChargingPoint) GetVoltsRangeOk() (*string, bool)`

GetVoltsRangeOk returns a tuple with the VoltsRange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoltsRange

`func (o *EvChargingPoint) SetVoltsRange(v string)`

SetVoltsRange sets VoltsRange field to given value.

### HasVoltsRange

`func (o *EvChargingPoint) HasVoltsRange() bool`

HasVoltsRange returns a boolean if a field has been set.

### GetPhases

`func (o *EvChargingPoint) GetPhases() int32`

GetPhases returns the Phases field if non-nil, zero value otherwise.

### GetPhasesOk

`func (o *EvChargingPoint) GetPhasesOk() (*int32, bool)`

GetPhasesOk returns a tuple with the Phases field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhases

`func (o *EvChargingPoint) SetPhases(v int32)`

SetPhases sets Phases field to given value.

### HasPhases

`func (o *EvChargingPoint) HasPhases() bool`

HasPhases returns a boolean if a field has been set.

### GetAmpsRange

`func (o *EvChargingPoint) GetAmpsRange() string`

GetAmpsRange returns the AmpsRange field if non-nil, zero value otherwise.

### GetAmpsRangeOk

`func (o *EvChargingPoint) GetAmpsRangeOk() (*string, bool)`

GetAmpsRangeOk returns a tuple with the AmpsRange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAmpsRange

`func (o *EvChargingPoint) SetAmpsRange(v string)`

SetAmpsRange sets AmpsRange field to given value.

### HasAmpsRange

`func (o *EvChargingPoint) HasAmpsRange() bool`

HasAmpsRange returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



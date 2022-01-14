# EvChargingAttributes

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Connectors** | Pointer to [**[]EvConnector**](EvConnector.md) | List of EV pool groups of connectors. Each group is defined by a common charging connector type and max power level. The numberOfConnectors field contains the number of connectors in the group. | [optional] 
**TotalNumberOfConnectors** | Pointer to **int32** | Total number of charging connectors in the EV charging pool | [optional] 

## Methods

### NewEvChargingAttributes

`func NewEvChargingAttributes() *EvChargingAttributes`

NewEvChargingAttributes instantiates a new EvChargingAttributes object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEvChargingAttributesWithDefaults

`func NewEvChargingAttributesWithDefaults() *EvChargingAttributes`

NewEvChargingAttributesWithDefaults instantiates a new EvChargingAttributes object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConnectors

`func (o *EvChargingAttributes) GetConnectors() []EvConnector`

GetConnectors returns the Connectors field if non-nil, zero value otherwise.

### GetConnectorsOk

`func (o *EvChargingAttributes) GetConnectorsOk() (*[]EvConnector, bool)`

GetConnectorsOk returns a tuple with the Connectors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectors

`func (o *EvChargingAttributes) SetConnectors(v []EvConnector)`

SetConnectors sets Connectors field to given value.

### HasConnectors

`func (o *EvChargingAttributes) HasConnectors() bool`

HasConnectors returns a boolean if a field has been set.

### GetTotalNumberOfConnectors

`func (o *EvChargingAttributes) GetTotalNumberOfConnectors() int32`

GetTotalNumberOfConnectors returns the TotalNumberOfConnectors field if non-nil, zero value otherwise.

### GetTotalNumberOfConnectorsOk

`func (o *EvChargingAttributes) GetTotalNumberOfConnectorsOk() (*int32, bool)`

GetTotalNumberOfConnectorsOk returns a tuple with the TotalNumberOfConnectors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTotalNumberOfConnectors

`func (o *EvChargingAttributes) SetTotalNumberOfConnectors(v int32)`

SetTotalNumberOfConnectors sets TotalNumberOfConnectors field to given value.

### HasTotalNumberOfConnectors

`func (o *EvChargingAttributes) HasTotalNumberOfConnectors() bool`

HasTotalNumberOfConnectors returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



# EvConnector

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SupplierName** | Pointer to **string** | The EV charge point operator | [optional] 
**ConnectorType** | Pointer to [**EvNameId**](EvNameId.md) | Id and name element pair representing the connector type in the EV pool group. For more information on the current connector types, see the [connectorTypes](https://developer.here.com/documentation/charging-stations/dev_guide/topics/resource-type-connector-types.html) values in the HERE EV Charge Points API. | [optional] 
**PowerFeedType** | Pointer to [**EvNameId**](EvNameId.md) | Details on type of power feed with respect to [SAE J1772](https://en.wikipedia.org/wiki/SAE_J1772#Charging) standard. | [optional] 
**MaxPowerLevel** | Pointer to **float64** | Maximum charge power (in kilowatt) of connectors in connectors group. | [optional] 
**ChargingPoint** | Pointer to [**EvChargingPoint**](EvChargingPoint.md) | Connectors group additional charging information | [optional] 

## Methods

### NewEvConnector

`func NewEvConnector() *EvConnector`

NewEvConnector instantiates a new EvConnector object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEvConnectorWithDefaults

`func NewEvConnectorWithDefaults() *EvConnector`

NewEvConnectorWithDefaults instantiates a new EvConnector object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSupplierName

`func (o *EvConnector) GetSupplierName() string`

GetSupplierName returns the SupplierName field if non-nil, zero value otherwise.

### GetSupplierNameOk

`func (o *EvConnector) GetSupplierNameOk() (*string, bool)`

GetSupplierNameOk returns a tuple with the SupplierName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupplierName

`func (o *EvConnector) SetSupplierName(v string)`

SetSupplierName sets SupplierName field to given value.

### HasSupplierName

`func (o *EvConnector) HasSupplierName() bool`

HasSupplierName returns a boolean if a field has been set.

### GetConnectorType

`func (o *EvConnector) GetConnectorType() EvNameId`

GetConnectorType returns the ConnectorType field if non-nil, zero value otherwise.

### GetConnectorTypeOk

`func (o *EvConnector) GetConnectorTypeOk() (*EvNameId, bool)`

GetConnectorTypeOk returns a tuple with the ConnectorType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectorType

`func (o *EvConnector) SetConnectorType(v EvNameId)`

SetConnectorType sets ConnectorType field to given value.

### HasConnectorType

`func (o *EvConnector) HasConnectorType() bool`

HasConnectorType returns a boolean if a field has been set.

### GetPowerFeedType

`func (o *EvConnector) GetPowerFeedType() EvNameId`

GetPowerFeedType returns the PowerFeedType field if non-nil, zero value otherwise.

### GetPowerFeedTypeOk

`func (o *EvConnector) GetPowerFeedTypeOk() (*EvNameId, bool)`

GetPowerFeedTypeOk returns a tuple with the PowerFeedType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPowerFeedType

`func (o *EvConnector) SetPowerFeedType(v EvNameId)`

SetPowerFeedType sets PowerFeedType field to given value.

### HasPowerFeedType

`func (o *EvConnector) HasPowerFeedType() bool`

HasPowerFeedType returns a boolean if a field has been set.

### GetMaxPowerLevel

`func (o *EvConnector) GetMaxPowerLevel() float64`

GetMaxPowerLevel returns the MaxPowerLevel field if non-nil, zero value otherwise.

### GetMaxPowerLevelOk

`func (o *EvConnector) GetMaxPowerLevelOk() (*float64, bool)`

GetMaxPowerLevelOk returns a tuple with the MaxPowerLevel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxPowerLevel

`func (o *EvConnector) SetMaxPowerLevel(v float64)`

SetMaxPowerLevel sets MaxPowerLevel field to given value.

### HasMaxPowerLevel

`func (o *EvConnector) HasMaxPowerLevel() bool`

HasMaxPowerLevel returns a boolean if a field has been set.

### GetChargingPoint

`func (o *EvConnector) GetChargingPoint() EvChargingPoint`

GetChargingPoint returns the ChargingPoint field if non-nil, zero value otherwise.

### GetChargingPointOk

`func (o *EvConnector) GetChargingPointOk() (*EvChargingPoint, bool)`

GetChargingPointOk returns a tuple with the ChargingPoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChargingPoint

`func (o *EvConnector) SetChargingPoint(v EvChargingPoint)`

SetChargingPoint sets ChargingPoint field to given value.

### HasChargingPoint

`func (o *EvConnector) HasChargingPoint() bool`

HasChargingPoint returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



# MapView

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**West** | **float64** | Longitude of the western-side of the box. For example: \&quot;8.80068\&quot; | 
**South** | **float64** | Latitude of the southern-side of the box. For example: \&quot;52.19333\&quot; | 
**East** | **float64** | Longitude of the eastern-side of the box. For example: \&quot;8.8167\&quot; | 
**North** | **float64** | Latitude of the northern-side of the box. For example: \&quot;52.19555\&quot; | 

## Methods

### NewMapView

`func NewMapView(west float64, south float64, east float64, north float64, ) *MapView`

NewMapView instantiates a new MapView object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMapViewWithDefaults

`func NewMapViewWithDefaults() *MapView`

NewMapViewWithDefaults instantiates a new MapView object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetWest

`func (o *MapView) GetWest() float64`

GetWest returns the West field if non-nil, zero value otherwise.

### GetWestOk

`func (o *MapView) GetWestOk() (*float64, bool)`

GetWestOk returns a tuple with the West field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWest

`func (o *MapView) SetWest(v float64)`

SetWest sets West field to given value.


### GetSouth

`func (o *MapView) GetSouth() float64`

GetSouth returns the South field if non-nil, zero value otherwise.

### GetSouthOk

`func (o *MapView) GetSouthOk() (*float64, bool)`

GetSouthOk returns a tuple with the South field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSouth

`func (o *MapView) SetSouth(v float64)`

SetSouth sets South field to given value.


### GetEast

`func (o *MapView) GetEast() float64`

GetEast returns the East field if non-nil, zero value otherwise.

### GetEastOk

`func (o *MapView) GetEastOk() (*float64, bool)`

GetEastOk returns a tuple with the East field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEast

`func (o *MapView) SetEast(v float64)`

SetEast sets East field to given value.


### GetNorth

`func (o *MapView) GetNorth() float64`

GetNorth returns the North field if non-nil, zero value otherwise.

### GetNorthOk

`func (o *MapView) GetNorthOk() (*float64, bool)`

GetNorthOk returns a tuple with the North field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNorth

`func (o *MapView) SetNorth(v float64)`

SetNorth sets North field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



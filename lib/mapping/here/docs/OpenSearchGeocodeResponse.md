# OpenSearchGeocodeResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Items** | [**[]GeocodeResultItem**](GeocodeResultItem.md) | The results are presented as a JSON list of candidates in ranked order (most-likely to least-likely) based on the matched location criteria. | 

## Methods

### NewOpenSearchGeocodeResponse

`func NewOpenSearchGeocodeResponse(items []GeocodeResultItem, ) *OpenSearchGeocodeResponse`

NewOpenSearchGeocodeResponse instantiates a new OpenSearchGeocodeResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpenSearchGeocodeResponseWithDefaults

`func NewOpenSearchGeocodeResponseWithDefaults() *OpenSearchGeocodeResponse`

NewOpenSearchGeocodeResponseWithDefaults instantiates a new OpenSearchGeocodeResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetItems

`func (o *OpenSearchGeocodeResponse) GetItems() []GeocodeResultItem`

GetItems returns the Items field if non-nil, zero value otherwise.

### GetItemsOk

`func (o *OpenSearchGeocodeResponse) GetItemsOk() (*[]GeocodeResultItem, bool)`

GetItemsOk returns a tuple with the Items field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetItems

`func (o *OpenSearchGeocodeResponse) SetItems(v []GeocodeResultItem)`

SetItems sets Items field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



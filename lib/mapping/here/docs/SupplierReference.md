# SupplierReference

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Supplier** | [**Supplier**](Supplier.md) | Information about the supplier of this reference. | 
**Id** | **string** | Identifier of the place as provided by the supplier. | 

## Methods

### NewSupplierReference

`func NewSupplierReference(supplier Supplier, id string, ) *SupplierReference`

NewSupplierReference instantiates a new SupplierReference object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSupplierReferenceWithDefaults

`func NewSupplierReferenceWithDefaults() *SupplierReference`

NewSupplierReferenceWithDefaults instantiates a new SupplierReference object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSupplier

`func (o *SupplierReference) GetSupplier() Supplier`

GetSupplier returns the Supplier field if non-nil, zero value otherwise.

### GetSupplierOk

`func (o *SupplierReference) GetSupplierOk() (*Supplier, bool)`

GetSupplierOk returns a tuple with the Supplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSupplier

`func (o *SupplierReference) SetSupplier(v Supplier)`

SetSupplier sets Supplier field to given value.


### GetId

`func (o *SupplierReference) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *SupplierReference) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *SupplierReference) SetId(v string)`

SetId sets Id field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



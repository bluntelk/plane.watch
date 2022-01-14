# OpeningHours

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Categories** | Pointer to [**[]Category**](Category.md) | The list of place categories, this set of opening hours refers to. | [optional] 
**Text** | **[]string** |  | 
**IsOpen** | Pointer to **bool** |  | [optional] 
**Structured** | [**[]StructuredOpeningHours**](StructuredOpeningHours.md) | List of iCalender-based structured representations of opening hours | 

## Methods

### NewOpeningHours

`func NewOpeningHours(text []string, structured []StructuredOpeningHours, ) *OpeningHours`

NewOpeningHours instantiates a new OpeningHours object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewOpeningHoursWithDefaults

`func NewOpeningHoursWithDefaults() *OpeningHours`

NewOpeningHoursWithDefaults instantiates a new OpeningHours object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCategories

`func (o *OpeningHours) GetCategories() []Category`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *OpeningHours) GetCategoriesOk() (*[]Category, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *OpeningHours) SetCategories(v []Category)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *OpeningHours) HasCategories() bool`

HasCategories returns a boolean if a field has been set.

### GetText

`func (o *OpeningHours) GetText() []string`

GetText returns the Text field if non-nil, zero value otherwise.

### GetTextOk

`func (o *OpeningHours) GetTextOk() (*[]string, bool)`

GetTextOk returns a tuple with the Text field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetText

`func (o *OpeningHours) SetText(v []string)`

SetText sets Text field to given value.


### GetIsOpen

`func (o *OpeningHours) GetIsOpen() bool`

GetIsOpen returns the IsOpen field if non-nil, zero value otherwise.

### GetIsOpenOk

`func (o *OpeningHours) GetIsOpenOk() (*bool, bool)`

GetIsOpenOk returns a tuple with the IsOpen field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsOpen

`func (o *OpeningHours) SetIsOpen(v bool)`

SetIsOpen sets IsOpen field to given value.

### HasIsOpen

`func (o *OpeningHours) HasIsOpen() bool`

HasIsOpen returns a boolean if a field has been set.

### GetStructured

`func (o *OpeningHours) GetStructured() []StructuredOpeningHours`

GetStructured returns the Structured field if non-nil, zero value otherwise.

### GetStructuredOk

`func (o *OpeningHours) GetStructuredOk() (*[]StructuredOpeningHours, bool)`

GetStructuredOk returns a tuple with the Structured field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStructured

`func (o *OpeningHours) SetStructured(v []StructuredOpeningHours)`

SetStructured sets Structured field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



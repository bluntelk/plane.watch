# ContactInformation

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Categories** | Pointer to [**[]Category**](Category.md) | The list of place categories, this set of contact details refers to. | [optional] 
**Phone** | Pointer to [**[]Contact**](Contact.md) |  | [optional] 
**Mobile** | Pointer to [**[]Contact**](Contact.md) |  | [optional] 
**TollFree** | Pointer to [**[]Contact**](Contact.md) |  | [optional] 
**Fax** | Pointer to [**[]Contact**](Contact.md) |  | [optional] 
**Www** | Pointer to [**[]Contact**](Contact.md) |  | [optional] 
**Email** | Pointer to [**[]Contact**](Contact.md) |  | [optional] 

## Methods

### NewContactInformation

`func NewContactInformation() *ContactInformation`

NewContactInformation instantiates a new ContactInformation object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewContactInformationWithDefaults

`func NewContactInformationWithDefaults() *ContactInformation`

NewContactInformationWithDefaults instantiates a new ContactInformation object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCategories

`func (o *ContactInformation) GetCategories() []Category`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *ContactInformation) GetCategoriesOk() (*[]Category, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *ContactInformation) SetCategories(v []Category)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *ContactInformation) HasCategories() bool`

HasCategories returns a boolean if a field has been set.

### GetPhone

`func (o *ContactInformation) GetPhone() []Contact`

GetPhone returns the Phone field if non-nil, zero value otherwise.

### GetPhoneOk

`func (o *ContactInformation) GetPhoneOk() (*[]Contact, bool)`

GetPhoneOk returns a tuple with the Phone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhone

`func (o *ContactInformation) SetPhone(v []Contact)`

SetPhone sets Phone field to given value.

### HasPhone

`func (o *ContactInformation) HasPhone() bool`

HasPhone returns a boolean if a field has been set.

### GetMobile

`func (o *ContactInformation) GetMobile() []Contact`

GetMobile returns the Mobile field if non-nil, zero value otherwise.

### GetMobileOk

`func (o *ContactInformation) GetMobileOk() (*[]Contact, bool)`

GetMobileOk returns a tuple with the Mobile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMobile

`func (o *ContactInformation) SetMobile(v []Contact)`

SetMobile sets Mobile field to given value.

### HasMobile

`func (o *ContactInformation) HasMobile() bool`

HasMobile returns a boolean if a field has been set.

### GetTollFree

`func (o *ContactInformation) GetTollFree() []Contact`

GetTollFree returns the TollFree field if non-nil, zero value otherwise.

### GetTollFreeOk

`func (o *ContactInformation) GetTollFreeOk() (*[]Contact, bool)`

GetTollFreeOk returns a tuple with the TollFree field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTollFree

`func (o *ContactInformation) SetTollFree(v []Contact)`

SetTollFree sets TollFree field to given value.

### HasTollFree

`func (o *ContactInformation) HasTollFree() bool`

HasTollFree returns a boolean if a field has been set.

### GetFax

`func (o *ContactInformation) GetFax() []Contact`

GetFax returns the Fax field if non-nil, zero value otherwise.

### GetFaxOk

`func (o *ContactInformation) GetFaxOk() (*[]Contact, bool)`

GetFaxOk returns a tuple with the Fax field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFax

`func (o *ContactInformation) SetFax(v []Contact)`

SetFax sets Fax field to given value.

### HasFax

`func (o *ContactInformation) HasFax() bool`

HasFax returns a boolean if a field has been set.

### GetWww

`func (o *ContactInformation) GetWww() []Contact`

GetWww returns the Www field if non-nil, zero value otherwise.

### GetWwwOk

`func (o *ContactInformation) GetWwwOk() (*[]Contact, bool)`

GetWwwOk returns a tuple with the Www field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWww

`func (o *ContactInformation) SetWww(v []Contact)`

SetWww sets Www field to given value.

### HasWww

`func (o *ContactInformation) HasWww() bool`

HasWww returns a boolean if a field has been set.

### GetEmail

`func (o *ContactInformation) GetEmail() []Contact`

GetEmail returns the Email field if non-nil, zero value otherwise.

### GetEmailOk

`func (o *ContactInformation) GetEmailOk() (*[]Contact, bool)`

GetEmailOk returns a tuple with the Email field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmail

`func (o *ContactInformation) SetEmail(v []Contact)`

SetEmail sets Email field to given value.

### HasEmail

`func (o *ContactInformation) HasEmail() bool`

HasEmail returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



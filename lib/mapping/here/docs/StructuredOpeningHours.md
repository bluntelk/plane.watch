# StructuredOpeningHours

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Start** | **string** | String with a modified [iCalendar DATE-TIME](https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.5) value. The date part is omitted, values starts with the time section maker \\\&quot;T\\\&quot;. Example: T132000 | 
**Duration** | **string** | String with an [iCalendar DURATION](https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.6) value. A closed day has the value PT00:00M | 
**Recurrence** | **string** | String with a [RECUR](https://datatracker.ietf.org/doc/html/rfc5545#section-3.3.10) rule | 

## Methods

### NewStructuredOpeningHours

`func NewStructuredOpeningHours(start string, duration string, recurrence string, ) *StructuredOpeningHours`

NewStructuredOpeningHours instantiates a new StructuredOpeningHours object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStructuredOpeningHoursWithDefaults

`func NewStructuredOpeningHoursWithDefaults() *StructuredOpeningHours`

NewStructuredOpeningHoursWithDefaults instantiates a new StructuredOpeningHours object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStart

`func (o *StructuredOpeningHours) GetStart() string`

GetStart returns the Start field if non-nil, zero value otherwise.

### GetStartOk

`func (o *StructuredOpeningHours) GetStartOk() (*string, bool)`

GetStartOk returns a tuple with the Start field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStart

`func (o *StructuredOpeningHours) SetStart(v string)`

SetStart sets Start field to given value.


### GetDuration

`func (o *StructuredOpeningHours) GetDuration() string`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *StructuredOpeningHours) GetDurationOk() (*string, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *StructuredOpeningHours) SetDuration(v string)`

SetDuration sets Duration field to given value.


### GetRecurrence

`func (o *StructuredOpeningHours) GetRecurrence() string`

GetRecurrence returns the Recurrence field if non-nil, zero value otherwise.

### GetRecurrenceOk

`func (o *StructuredOpeningHours) GetRecurrenceOk() (*string, bool)`

GetRecurrenceOk returns a tuple with the Recurrence field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRecurrence

`func (o *StructuredOpeningHours) SetRecurrence(v string)`

SetRecurrence sets Recurrence field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



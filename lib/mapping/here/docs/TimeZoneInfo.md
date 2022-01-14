# TimeZoneInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | The name of the time zone as defined in the [tz database](https://en.wikipedia.org/wiki/Tz_database). For example: \&quot;Europe/Berlin\&quot; | 
**UtcOffset** | **string** | The UTC offset for this time zone at request time. For example \&quot;+02:00\&quot; | 

## Methods

### NewTimeZoneInfo

`func NewTimeZoneInfo(name string, utcOffset string, ) *TimeZoneInfo`

NewTimeZoneInfo instantiates a new TimeZoneInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTimeZoneInfoWithDefaults

`func NewTimeZoneInfoWithDefaults() *TimeZoneInfo`

NewTimeZoneInfoWithDefaults instantiates a new TimeZoneInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *TimeZoneInfo) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *TimeZoneInfo) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *TimeZoneInfo) SetName(v string)`

SetName sets Name field to given value.


### GetUtcOffset

`func (o *TimeZoneInfo) GetUtcOffset() string`

GetUtcOffset returns the UtcOffset field if non-nil, zero value otherwise.

### GetUtcOffsetOk

`func (o *TimeZoneInfo) GetUtcOffsetOk() (*string, bool)`

GetUtcOffsetOk returns a tuple with the UtcOffset field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUtcOffset

`func (o *TimeZoneInfo) SetUtcOffset(v string)`

SetUtcOffset sets UtcOffset field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



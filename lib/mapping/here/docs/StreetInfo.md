# StreetInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BaseName** | Pointer to **string** | Base name part of the street name. | [optional] 
**StreetType** | Pointer to **string** | Street type part of the street name. | [optional] 
**StreetTypePrecedes** | Pointer to **bool** | Defines if the street type is before or after the base name. | [optional] 
**StreetTypeAttached** | Pointer to **bool** | Defines if the street type is attached or unattached to the base name. | [optional] 
**Prefix** | Pointer to **string** | A prefix is a directional identifier that precedes, but is not included in, the base name of a road. | [optional] 
**Suffix** | Pointer to **string** | A suffix is a directional identifier that follows, but is not included in, the base name of a road. | [optional] 
**Direction** | Pointer to **string** | Indicates the official directional identifiers assigned to highways, typically either \&quot;North/South\&quot; or \&quot;East/West\&quot; | [optional] 
**Language** | Pointer to **string** | BCP 47 compliant language code | [optional] 

## Methods

### NewStreetInfo

`func NewStreetInfo() *StreetInfo`

NewStreetInfo instantiates a new StreetInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewStreetInfoWithDefaults

`func NewStreetInfoWithDefaults() *StreetInfo`

NewStreetInfoWithDefaults instantiates a new StreetInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBaseName

`func (o *StreetInfo) GetBaseName() string`

GetBaseName returns the BaseName field if non-nil, zero value otherwise.

### GetBaseNameOk

`func (o *StreetInfo) GetBaseNameOk() (*string, bool)`

GetBaseNameOk returns a tuple with the BaseName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseName

`func (o *StreetInfo) SetBaseName(v string)`

SetBaseName sets BaseName field to given value.

### HasBaseName

`func (o *StreetInfo) HasBaseName() bool`

HasBaseName returns a boolean if a field has been set.

### GetStreetType

`func (o *StreetInfo) GetStreetType() string`

GetStreetType returns the StreetType field if non-nil, zero value otherwise.

### GetStreetTypeOk

`func (o *StreetInfo) GetStreetTypeOk() (*string, bool)`

GetStreetTypeOk returns a tuple with the StreetType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreetType

`func (o *StreetInfo) SetStreetType(v string)`

SetStreetType sets StreetType field to given value.

### HasStreetType

`func (o *StreetInfo) HasStreetType() bool`

HasStreetType returns a boolean if a field has been set.

### GetStreetTypePrecedes

`func (o *StreetInfo) GetStreetTypePrecedes() bool`

GetStreetTypePrecedes returns the StreetTypePrecedes field if non-nil, zero value otherwise.

### GetStreetTypePrecedesOk

`func (o *StreetInfo) GetStreetTypePrecedesOk() (*bool, bool)`

GetStreetTypePrecedesOk returns a tuple with the StreetTypePrecedes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreetTypePrecedes

`func (o *StreetInfo) SetStreetTypePrecedes(v bool)`

SetStreetTypePrecedes sets StreetTypePrecedes field to given value.

### HasStreetTypePrecedes

`func (o *StreetInfo) HasStreetTypePrecedes() bool`

HasStreetTypePrecedes returns a boolean if a field has been set.

### GetStreetTypeAttached

`func (o *StreetInfo) GetStreetTypeAttached() bool`

GetStreetTypeAttached returns the StreetTypeAttached field if non-nil, zero value otherwise.

### GetStreetTypeAttachedOk

`func (o *StreetInfo) GetStreetTypeAttachedOk() (*bool, bool)`

GetStreetTypeAttachedOk returns a tuple with the StreetTypeAttached field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreetTypeAttached

`func (o *StreetInfo) SetStreetTypeAttached(v bool)`

SetStreetTypeAttached sets StreetTypeAttached field to given value.

### HasStreetTypeAttached

`func (o *StreetInfo) HasStreetTypeAttached() bool`

HasStreetTypeAttached returns a boolean if a field has been set.

### GetPrefix

`func (o *StreetInfo) GetPrefix() string`

GetPrefix returns the Prefix field if non-nil, zero value otherwise.

### GetPrefixOk

`func (o *StreetInfo) GetPrefixOk() (*string, bool)`

GetPrefixOk returns a tuple with the Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrefix

`func (o *StreetInfo) SetPrefix(v string)`

SetPrefix sets Prefix field to given value.

### HasPrefix

`func (o *StreetInfo) HasPrefix() bool`

HasPrefix returns a boolean if a field has been set.

### GetSuffix

`func (o *StreetInfo) GetSuffix() string`

GetSuffix returns the Suffix field if non-nil, zero value otherwise.

### GetSuffixOk

`func (o *StreetInfo) GetSuffixOk() (*string, bool)`

GetSuffixOk returns a tuple with the Suffix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuffix

`func (o *StreetInfo) SetSuffix(v string)`

SetSuffix sets Suffix field to given value.

### HasSuffix

`func (o *StreetInfo) HasSuffix() bool`

HasSuffix returns a boolean if a field has been set.

### GetDirection

`func (o *StreetInfo) GetDirection() string`

GetDirection returns the Direction field if non-nil, zero value otherwise.

### GetDirectionOk

`func (o *StreetInfo) GetDirectionOk() (*string, bool)`

GetDirectionOk returns a tuple with the Direction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirection

`func (o *StreetInfo) SetDirection(v string)`

SetDirection sets Direction field to given value.

### HasDirection

`func (o *StreetInfo) HasDirection() bool`

HasDirection returns a boolean if a field has been set.

### GetLanguage

`func (o *StreetInfo) GetLanguage() string`

GetLanguage returns the Language field if non-nil, zero value otherwise.

### GetLanguageOk

`func (o *StreetInfo) GetLanguageOk() (*string, bool)`

GetLanguageOk returns a tuple with the Language field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguage

`func (o *StreetInfo) SetLanguage(v string)`

SetLanguage sets Language field to given value.

### HasLanguage

`func (o *StreetInfo) HasLanguage() bool`

HasLanguage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



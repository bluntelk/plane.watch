/*
 * Geocoding and Search API v7
 *
 * This document describes the Geocoding and Search API.
 *
 * API version: 7.78
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package here

import (
	"encoding/json"
)

// StreetInfo struct for StreetInfo
type StreetInfo struct {
	// Base name part of the street name.
	BaseName *string `json:"baseName,omitempty"`
	// Street type part of the street name.
	StreetType *string `json:"streetType,omitempty"`
	// Defines if the street type is before or after the base name.
	StreetTypePrecedes *bool `json:"streetTypePrecedes,omitempty"`
	// Defines if the street type is attached or unattached to the base name.
	StreetTypeAttached *bool `json:"streetTypeAttached,omitempty"`
	// A prefix is a directional identifier that precedes, but is not included in, the base name of a road.
	Prefix *string `json:"prefix,omitempty"`
	// A suffix is a directional identifier that follows, but is not included in, the base name of a road.
	Suffix *string `json:"suffix,omitempty"`
	// Indicates the official directional identifiers assigned to highways, typically either \"North/South\" or \"East/West\"
	Direction *string `json:"direction,omitempty"`
	// BCP 47 compliant language code
	Language *string `json:"language,omitempty"`
}

// NewStreetInfo instantiates a new StreetInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewStreetInfo() *StreetInfo {
	this := StreetInfo{}
	return &this
}

// NewStreetInfoWithDefaults instantiates a new StreetInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewStreetInfoWithDefaults() *StreetInfo {
	this := StreetInfo{}
	return &this
}

// GetBaseName returns the BaseName field value if set, zero value otherwise.
func (o *StreetInfo) GetBaseName() string {
	if o == nil || o.BaseName == nil {
		var ret string
		return ret
	}
	return *o.BaseName
}

// GetBaseNameOk returns a tuple with the BaseName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetBaseNameOk() (*string, bool) {
	if o == nil || o.BaseName == nil {
		return nil, false
	}
	return o.BaseName, true
}

// HasBaseName returns a boolean if a field has been set.
func (o *StreetInfo) HasBaseName() bool {
	if o != nil && o.BaseName != nil {
		return true
	}

	return false
}

// SetBaseName gets a reference to the given string and assigns it to the BaseName field.
func (o *StreetInfo) SetBaseName(v string) {
	o.BaseName = &v
}

// GetStreetType returns the StreetType field value if set, zero value otherwise.
func (o *StreetInfo) GetStreetType() string {
	if o == nil || o.StreetType == nil {
		var ret string
		return ret
	}
	return *o.StreetType
}

// GetStreetTypeOk returns a tuple with the StreetType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetStreetTypeOk() (*string, bool) {
	if o == nil || o.StreetType == nil {
		return nil, false
	}
	return o.StreetType, true
}

// HasStreetType returns a boolean if a field has been set.
func (o *StreetInfo) HasStreetType() bool {
	if o != nil && o.StreetType != nil {
		return true
	}

	return false
}

// SetStreetType gets a reference to the given string and assigns it to the StreetType field.
func (o *StreetInfo) SetStreetType(v string) {
	o.StreetType = &v
}

// GetStreetTypePrecedes returns the StreetTypePrecedes field value if set, zero value otherwise.
func (o *StreetInfo) GetStreetTypePrecedes() bool {
	if o == nil || o.StreetTypePrecedes == nil {
		var ret bool
		return ret
	}
	return *o.StreetTypePrecedes
}

// GetStreetTypePrecedesOk returns a tuple with the StreetTypePrecedes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetStreetTypePrecedesOk() (*bool, bool) {
	if o == nil || o.StreetTypePrecedes == nil {
		return nil, false
	}
	return o.StreetTypePrecedes, true
}

// HasStreetTypePrecedes returns a boolean if a field has been set.
func (o *StreetInfo) HasStreetTypePrecedes() bool {
	if o != nil && o.StreetTypePrecedes != nil {
		return true
	}

	return false
}

// SetStreetTypePrecedes gets a reference to the given bool and assigns it to the StreetTypePrecedes field.
func (o *StreetInfo) SetStreetTypePrecedes(v bool) {
	o.StreetTypePrecedes = &v
}

// GetStreetTypeAttached returns the StreetTypeAttached field value if set, zero value otherwise.
func (o *StreetInfo) GetStreetTypeAttached() bool {
	if o == nil || o.StreetTypeAttached == nil {
		var ret bool
		return ret
	}
	return *o.StreetTypeAttached
}

// GetStreetTypeAttachedOk returns a tuple with the StreetTypeAttached field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetStreetTypeAttachedOk() (*bool, bool) {
	if o == nil || o.StreetTypeAttached == nil {
		return nil, false
	}
	return o.StreetTypeAttached, true
}

// HasStreetTypeAttached returns a boolean if a field has been set.
func (o *StreetInfo) HasStreetTypeAttached() bool {
	if o != nil && o.StreetTypeAttached != nil {
		return true
	}

	return false
}

// SetStreetTypeAttached gets a reference to the given bool and assigns it to the StreetTypeAttached field.
func (o *StreetInfo) SetStreetTypeAttached(v bool) {
	o.StreetTypeAttached = &v
}

// GetPrefix returns the Prefix field value if set, zero value otherwise.
func (o *StreetInfo) GetPrefix() string {
	if o == nil || o.Prefix == nil {
		var ret string
		return ret
	}
	return *o.Prefix
}

// GetPrefixOk returns a tuple with the Prefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetPrefixOk() (*string, bool) {
	if o == nil || o.Prefix == nil {
		return nil, false
	}
	return o.Prefix, true
}

// HasPrefix returns a boolean if a field has been set.
func (o *StreetInfo) HasPrefix() bool {
	if o != nil && o.Prefix != nil {
		return true
	}

	return false
}

// SetPrefix gets a reference to the given string and assigns it to the Prefix field.
func (o *StreetInfo) SetPrefix(v string) {
	o.Prefix = &v
}

// GetSuffix returns the Suffix field value if set, zero value otherwise.
func (o *StreetInfo) GetSuffix() string {
	if o == nil || o.Suffix == nil {
		var ret string
		return ret
	}
	return *o.Suffix
}

// GetSuffixOk returns a tuple with the Suffix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetSuffixOk() (*string, bool) {
	if o == nil || o.Suffix == nil {
		return nil, false
	}
	return o.Suffix, true
}

// HasSuffix returns a boolean if a field has been set.
func (o *StreetInfo) HasSuffix() bool {
	if o != nil && o.Suffix != nil {
		return true
	}

	return false
}

// SetSuffix gets a reference to the given string and assigns it to the Suffix field.
func (o *StreetInfo) SetSuffix(v string) {
	o.Suffix = &v
}

// GetDirection returns the Direction field value if set, zero value otherwise.
func (o *StreetInfo) GetDirection() string {
	if o == nil || o.Direction == nil {
		var ret string
		return ret
	}
	return *o.Direction
}

// GetDirectionOk returns a tuple with the Direction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetDirectionOk() (*string, bool) {
	if o == nil || o.Direction == nil {
		return nil, false
	}
	return o.Direction, true
}

// HasDirection returns a boolean if a field has been set.
func (o *StreetInfo) HasDirection() bool {
	if o != nil && o.Direction != nil {
		return true
	}

	return false
}

// SetDirection gets a reference to the given string and assigns it to the Direction field.
func (o *StreetInfo) SetDirection(v string) {
	o.Direction = &v
}

// GetLanguage returns the Language field value if set, zero value otherwise.
func (o *StreetInfo) GetLanguage() string {
	if o == nil || o.Language == nil {
		var ret string
		return ret
	}
	return *o.Language
}

// GetLanguageOk returns a tuple with the Language field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *StreetInfo) GetLanguageOk() (*string, bool) {
	if o == nil || o.Language == nil {
		return nil, false
	}
	return o.Language, true
}

// HasLanguage returns a boolean if a field has been set.
func (o *StreetInfo) HasLanguage() bool {
	if o != nil && o.Language != nil {
		return true
	}

	return false
}

// SetLanguage gets a reference to the given string and assigns it to the Language field.
func (o *StreetInfo) SetLanguage(v string) {
	o.Language = &v
}

func (o StreetInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.BaseName != nil {
		toSerialize["baseName"] = o.BaseName
	}
	if o.StreetType != nil {
		toSerialize["streetType"] = o.StreetType
	}
	if o.StreetTypePrecedes != nil {
		toSerialize["streetTypePrecedes"] = o.StreetTypePrecedes
	}
	if o.StreetTypeAttached != nil {
		toSerialize["streetTypeAttached"] = o.StreetTypeAttached
	}
	if o.Prefix != nil {
		toSerialize["prefix"] = o.Prefix
	}
	if o.Suffix != nil {
		toSerialize["suffix"] = o.Suffix
	}
	if o.Direction != nil {
		toSerialize["direction"] = o.Direction
	}
	if o.Language != nil {
		toSerialize["language"] = o.Language
	}
	return json.Marshal(toSerialize)
}

type NullableStreetInfo struct {
	value *StreetInfo
	isSet bool
}

func (v NullableStreetInfo) Get() *StreetInfo {
	return v.value
}

func (v *NullableStreetInfo) Set(val *StreetInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableStreetInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableStreetInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableStreetInfo(val *StreetInfo) *NullableStreetInfo {
	return &NullableStreetInfo{value: val, isSet: true}
}

func (v NullableStreetInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableStreetInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}



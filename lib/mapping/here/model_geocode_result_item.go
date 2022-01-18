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

// GeocodeResultItem struct for GeocodeResultItem
type GeocodeResultItem struct {
	// The localized display name of this result item.
	Title string `json:"title"`
	// The unique identifier for the result item. This ID can be used for a Look Up by ID search as well.
	Id *string `json:"id,omitempty"`
	// ISO3 country code of the item political view (default for international). This response element is populated when the politicalView parameter is set in the query
	PoliticalView *string `json:"politicalView,omitempty"`
	// WARNING: The resultType values 'intersection' and 'postalCodePoint' are in BETA state
	ResultType *string `json:"resultType,omitempty"`
	// * PA - a Point Address represents an individual address as a point object. Point Addresses are coming from trusted sources.   We can say with high certainty that the address exists and at what position. A Point Address result contains two types of coordinates.   One is the access point (or navigation coordinates), which is the point to start or end a drive. The other point is the position or display point.   This point varies per source and country. The point can be the rooftop point, a point close to the building entry, or a point close to the building,   driveway or parking lot that belongs to the building. * interpolated - an interpolated address. These are approximate positions as a result of a linear interpolation based on address ranges.   Address ranges, especially in the USA, are typical per block. For interpolated addresses, we cannot say with confidence that the address exists in reality.   But the interpolation provides a good location approximation that brings people in most use cases close to the target location.   The access point of an interpolated address result is calculated based on the address range and the road geometry.   The position (display) point is pre-configured offset from the street geometry.   Compared to Point Addresses, interpolated addresses are less accurate.
	HouseNumberType *string `json:"houseNumberType,omitempty"`
	AddressBlockType *string `json:"addressBlockType,omitempty"`
	LocalityType *string `json:"localityType,omitempty"`
	AdministrativeAreaType *string `json:"administrativeAreaType,omitempty"`
	// Postal address of the result item.
	Address Address `json:"address"`
	// The coordinates (latitude, longitude) of a pin on a map corresponding to the searched place.
	Position *DisplayResponseCoordinate `json:"position,omitempty"`
	// Coordinates of the place you are navigating to (for example, driving or walking). This is a point on a road or in a parking lot.
	Access *[]AccessResponseCoordinate `json:"access,omitempty"`
	// The distance \\\"as the crow flies\\\" from the search center to this result item in meters. For example: \\\"172039\\\".  When searching along a route this is the distance\\nalong the route plus the distance from the route polyline to this result item.
	Distance *int64 `json:"distance,omitempty"`
	// The bounding box enclosing the geometric shape (area or line) that an individual result covers. `place` typed results have no `mapView`.
	MapView *MapView `json:"mapView,omitempty"`
	// The list of categories assigned to this place.
	Categories *[]Category `json:"categories,omitempty"`
	// The list of food types assigned to this place.
	FoodTypes *[]Category `json:"foodTypes,omitempty"`
	// If true, indicates that the requested house number was corrected to match the nearest known house number. This field is visible only when the value is true.
	HouseNumberFallback *bool `json:"houseNumberFallback,omitempty"`
	// BETA - Provides time zone information for this place. (rendered only if 'show=tz' is provided.)
	TimeZone *TimeZoneInfo `json:"timeZone,omitempty"`
	// Indicates for each result how good the result matches to the original query. This can be used by the customer application to accept or reject the results depending on how \"expensive\" is the mistake for their use case
	Scoring *Scoring `json:"scoring,omitempty"`
	// BETA - Parsed terms and their positions in the input query (only rendered if 'show=parsing' is provided.)
	Parsing *Parsing `json:"parsing,omitempty"`
	// Street Details (only rendered if 'show=streetInfo' is provided.)
	StreetInfo *[]StreetInfo `json:"streetInfo,omitempty"`
	// Country Details (only rendered if 'show=countryInfo' is provided.)
	CountryInfo *CountryInfo `json:"countryInfo,omitempty"`
}

// NewGeocodeResultItem instantiates a new GeocodeResultItem object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewGeocodeResultItem(title string, address Address, ) *GeocodeResultItem {
	this := GeocodeResultItem{}
	this.Title = title
	this.Address = address
	return &this
}

// NewGeocodeResultItemWithDefaults instantiates a new GeocodeResultItem object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewGeocodeResultItemWithDefaults() *GeocodeResultItem {
	this := GeocodeResultItem{}
	return &this
}

// GetTitle returns the Title field value
func (o *GeocodeResultItem) GetTitle() string {
	if o == nil  {
		var ret string
		return ret
	}

	return o.Title
}

// GetTitleOk returns a tuple with the Title field value
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetTitleOk() (*string, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Title, true
}

// SetTitle sets field value
func (o *GeocodeResultItem) SetTitle(v string) {
	o.Title = v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetId() string {
	if o == nil || o.Id == nil {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetIdOk() (*string, bool) {
	if o == nil || o.Id == nil {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasId() bool {
	if o != nil && o.Id != nil {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *GeocodeResultItem) SetId(v string) {
	o.Id = &v
}

// GetPoliticalView returns the PoliticalView field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetPoliticalView() string {
	if o == nil || o.PoliticalView == nil {
		var ret string
		return ret
	}
	return *o.PoliticalView
}

// GetPoliticalViewOk returns a tuple with the PoliticalView field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetPoliticalViewOk() (*string, bool) {
	if o == nil || o.PoliticalView == nil {
		return nil, false
	}
	return o.PoliticalView, true
}

// HasPoliticalView returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasPoliticalView() bool {
	if o != nil && o.PoliticalView != nil {
		return true
	}

	return false
}

// SetPoliticalView gets a reference to the given string and assigns it to the PoliticalView field.
func (o *GeocodeResultItem) SetPoliticalView(v string) {
	o.PoliticalView = &v
}

// GetResultType returns the ResultType field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetResultType() string {
	if o == nil || o.ResultType == nil {
		var ret string
		return ret
	}
	return *o.ResultType
}

// GetResultTypeOk returns a tuple with the ResultType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetResultTypeOk() (*string, bool) {
	if o == nil || o.ResultType == nil {
		return nil, false
	}
	return o.ResultType, true
}

// HasResultType returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasResultType() bool {
	if o != nil && o.ResultType != nil {
		return true
	}

	return false
}

// SetResultType gets a reference to the given string and assigns it to the ResultType field.
func (o *GeocodeResultItem) SetResultType(v string) {
	o.ResultType = &v
}

// GetHouseNumberType returns the HouseNumberType field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetHouseNumberType() string {
	if o == nil || o.HouseNumberType == nil {
		var ret string
		return ret
	}
	return *o.HouseNumberType
}

// GetHouseNumberTypeOk returns a tuple with the HouseNumberType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetHouseNumberTypeOk() (*string, bool) {
	if o == nil || o.HouseNumberType == nil {
		return nil, false
	}
	return o.HouseNumberType, true
}

// HasHouseNumberType returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasHouseNumberType() bool {
	if o != nil && o.HouseNumberType != nil {
		return true
	}

	return false
}

// SetHouseNumberType gets a reference to the given string and assigns it to the HouseNumberType field.
func (o *GeocodeResultItem) SetHouseNumberType(v string) {
	o.HouseNumberType = &v
}

// GetAddressBlockType returns the AddressBlockType field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetAddressBlockType() string {
	if o == nil || o.AddressBlockType == nil {
		var ret string
		return ret
	}
	return *o.AddressBlockType
}

// GetAddressBlockTypeOk returns a tuple with the AddressBlockType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetAddressBlockTypeOk() (*string, bool) {
	if o == nil || o.AddressBlockType == nil {
		return nil, false
	}
	return o.AddressBlockType, true
}

// HasAddressBlockType returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasAddressBlockType() bool {
	if o != nil && o.AddressBlockType != nil {
		return true
	}

	return false
}

// SetAddressBlockType gets a reference to the given string and assigns it to the AddressBlockType field.
func (o *GeocodeResultItem) SetAddressBlockType(v string) {
	o.AddressBlockType = &v
}

// GetLocalityType returns the LocalityType field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetLocalityType() string {
	if o == nil || o.LocalityType == nil {
		var ret string
		return ret
	}
	return *o.LocalityType
}

// GetLocalityTypeOk returns a tuple with the LocalityType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetLocalityTypeOk() (*string, bool) {
	if o == nil || o.LocalityType == nil {
		return nil, false
	}
	return o.LocalityType, true
}

// HasLocalityType returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasLocalityType() bool {
	if o != nil && o.LocalityType != nil {
		return true
	}

	return false
}

// SetLocalityType gets a reference to the given string and assigns it to the LocalityType field.
func (o *GeocodeResultItem) SetLocalityType(v string) {
	o.LocalityType = &v
}

// GetAdministrativeAreaType returns the AdministrativeAreaType field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetAdministrativeAreaType() string {
	if o == nil || o.AdministrativeAreaType == nil {
		var ret string
		return ret
	}
	return *o.AdministrativeAreaType
}

// GetAdministrativeAreaTypeOk returns a tuple with the AdministrativeAreaType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetAdministrativeAreaTypeOk() (*string, bool) {
	if o == nil || o.AdministrativeAreaType == nil {
		return nil, false
	}
	return o.AdministrativeAreaType, true
}

// HasAdministrativeAreaType returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasAdministrativeAreaType() bool {
	if o != nil && o.AdministrativeAreaType != nil {
		return true
	}

	return false
}

// SetAdministrativeAreaType gets a reference to the given string and assigns it to the AdministrativeAreaType field.
func (o *GeocodeResultItem) SetAdministrativeAreaType(v string) {
	o.AdministrativeAreaType = &v
}

// GetAddress returns the Address field value
func (o *GeocodeResultItem) GetAddress() Address {
	if o == nil  {
		var ret Address
		return ret
	}

	return o.Address
}

// GetAddressOk returns a tuple with the Address field value
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetAddressOk() (*Address, bool) {
	if o == nil  {
		return nil, false
	}
	return &o.Address, true
}

// SetAddress sets field value
func (o *GeocodeResultItem) SetAddress(v Address) {
	o.Address = v
}

// GetPosition returns the Position field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetPosition() DisplayResponseCoordinate {
	if o == nil || o.Position == nil {
		var ret DisplayResponseCoordinate
		return ret
	}
	return *o.Position
}

// GetPositionOk returns a tuple with the Position field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetPositionOk() (*DisplayResponseCoordinate, bool) {
	if o == nil || o.Position == nil {
		return nil, false
	}
	return o.Position, true
}

// HasPosition returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasPosition() bool {
	if o != nil && o.Position != nil {
		return true
	}

	return false
}

// SetPosition gets a reference to the given DisplayResponseCoordinate and assigns it to the Position field.
func (o *GeocodeResultItem) SetPosition(v DisplayResponseCoordinate) {
	o.Position = &v
}

// GetAccess returns the Access field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetAccess() []AccessResponseCoordinate {
	if o == nil || o.Access == nil {
		var ret []AccessResponseCoordinate
		return ret
	}
	return *o.Access
}

// GetAccessOk returns a tuple with the Access field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetAccessOk() (*[]AccessResponseCoordinate, bool) {
	if o == nil || o.Access == nil {
		return nil, false
	}
	return o.Access, true
}

// HasAccess returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasAccess() bool {
	if o != nil && o.Access != nil {
		return true
	}

	return false
}

// SetAccess gets a reference to the given []AccessResponseCoordinate and assigns it to the Access field.
func (o *GeocodeResultItem) SetAccess(v []AccessResponseCoordinate) {
	o.Access = &v
}

// GetDistance returns the Distance field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetDistance() int64 {
	if o == nil || o.Distance == nil {
		var ret int64
		return ret
	}
	return *o.Distance
}

// GetDistanceOk returns a tuple with the Distance field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetDistanceOk() (*int64, bool) {
	if o == nil || o.Distance == nil {
		return nil, false
	}
	return o.Distance, true
}

// HasDistance returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasDistance() bool {
	if o != nil && o.Distance != nil {
		return true
	}

	return false
}

// SetDistance gets a reference to the given int64 and assigns it to the Distance field.
func (o *GeocodeResultItem) SetDistance(v int64) {
	o.Distance = &v
}

// GetMapView returns the MapView field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetMapView() MapView {
	if o == nil || o.MapView == nil {
		var ret MapView
		return ret
	}
	return *o.MapView
}

// GetMapViewOk returns a tuple with the MapView field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetMapViewOk() (*MapView, bool) {
	if o == nil || o.MapView == nil {
		return nil, false
	}
	return o.MapView, true
}

// HasMapView returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasMapView() bool {
	if o != nil && o.MapView != nil {
		return true
	}

	return false
}

// SetMapView gets a reference to the given MapView and assigns it to the MapView field.
func (o *GeocodeResultItem) SetMapView(v MapView) {
	o.MapView = &v
}

// GetCategories returns the Categories field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetCategories() []Category {
	if o == nil || o.Categories == nil {
		var ret []Category
		return ret
	}
	return *o.Categories
}

// GetCategoriesOk returns a tuple with the Categories field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetCategoriesOk() (*[]Category, bool) {
	if o == nil || o.Categories == nil {
		return nil, false
	}
	return o.Categories, true
}

// HasCategories returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasCategories() bool {
	if o != nil && o.Categories != nil {
		return true
	}

	return false
}

// SetCategories gets a reference to the given []Category and assigns it to the Categories field.
func (o *GeocodeResultItem) SetCategories(v []Category) {
	o.Categories = &v
}

// GetFoodTypes returns the FoodTypes field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetFoodTypes() []Category {
	if o == nil || o.FoodTypes == nil {
		var ret []Category
		return ret
	}
	return *o.FoodTypes
}

// GetFoodTypesOk returns a tuple with the FoodTypes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetFoodTypesOk() (*[]Category, bool) {
	if o == nil || o.FoodTypes == nil {
		return nil, false
	}
	return o.FoodTypes, true
}

// HasFoodTypes returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasFoodTypes() bool {
	if o != nil && o.FoodTypes != nil {
		return true
	}

	return false
}

// SetFoodTypes gets a reference to the given []Category and assigns it to the FoodTypes field.
func (o *GeocodeResultItem) SetFoodTypes(v []Category) {
	o.FoodTypes = &v
}

// GetHouseNumberFallback returns the HouseNumberFallback field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetHouseNumberFallback() bool {
	if o == nil || o.HouseNumberFallback == nil {
		var ret bool
		return ret
	}
	return *o.HouseNumberFallback
}

// GetHouseNumberFallbackOk returns a tuple with the HouseNumberFallback field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetHouseNumberFallbackOk() (*bool, bool) {
	if o == nil || o.HouseNumberFallback == nil {
		return nil, false
	}
	return o.HouseNumberFallback, true
}

// HasHouseNumberFallback returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasHouseNumberFallback() bool {
	if o != nil && o.HouseNumberFallback != nil {
		return true
	}

	return false
}

// SetHouseNumberFallback gets a reference to the given bool and assigns it to the HouseNumberFallback field.
func (o *GeocodeResultItem) SetHouseNumberFallback(v bool) {
	o.HouseNumberFallback = &v
}

// GetTimeZone returns the TimeZone field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetTimeZone() TimeZoneInfo {
	if o == nil || o.TimeZone == nil {
		var ret TimeZoneInfo
		return ret
	}
	return *o.TimeZone
}

// GetTimeZoneOk returns a tuple with the TimeZone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetTimeZoneOk() (*TimeZoneInfo, bool) {
	if o == nil || o.TimeZone == nil {
		return nil, false
	}
	return o.TimeZone, true
}

// HasTimeZone returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasTimeZone() bool {
	if o != nil && o.TimeZone != nil {
		return true
	}

	return false
}

// SetTimeZone gets a reference to the given TimeZoneInfo and assigns it to the TimeZone field.
func (o *GeocodeResultItem) SetTimeZone(v TimeZoneInfo) {
	o.TimeZone = &v
}

// GetScoring returns the Scoring field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetScoring() Scoring {
	if o == nil || o.Scoring == nil {
		var ret Scoring
		return ret
	}
	return *o.Scoring
}

// GetScoringOk returns a tuple with the Scoring field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetScoringOk() (*Scoring, bool) {
	if o == nil || o.Scoring == nil {
		return nil, false
	}
	return o.Scoring, true
}

// HasScoring returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasScoring() bool {
	if o != nil && o.Scoring != nil {
		return true
	}

	return false
}

// SetScoring gets a reference to the given Scoring and assigns it to the Scoring field.
func (o *GeocodeResultItem) SetScoring(v Scoring) {
	o.Scoring = &v
}

// GetParsing returns the Parsing field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetParsing() Parsing {
	if o == nil || o.Parsing == nil {
		var ret Parsing
		return ret
	}
	return *o.Parsing
}

// GetParsingOk returns a tuple with the Parsing field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetParsingOk() (*Parsing, bool) {
	if o == nil || o.Parsing == nil {
		return nil, false
	}
	return o.Parsing, true
}

// HasParsing returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasParsing() bool {
	if o != nil && o.Parsing != nil {
		return true
	}

	return false
}

// SetParsing gets a reference to the given Parsing and assigns it to the Parsing field.
func (o *GeocodeResultItem) SetParsing(v Parsing) {
	o.Parsing = &v
}

// GetStreetInfo returns the StreetInfo field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetStreetInfo() []StreetInfo {
	if o == nil || o.StreetInfo == nil {
		var ret []StreetInfo
		return ret
	}
	return *o.StreetInfo
}

// GetStreetInfoOk returns a tuple with the StreetInfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetStreetInfoOk() (*[]StreetInfo, bool) {
	if o == nil || o.StreetInfo == nil {
		return nil, false
	}
	return o.StreetInfo, true
}

// HasStreetInfo returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasStreetInfo() bool {
	if o != nil && o.StreetInfo != nil {
		return true
	}

	return false
}

// SetStreetInfo gets a reference to the given []StreetInfo and assigns it to the StreetInfo field.
func (o *GeocodeResultItem) SetStreetInfo(v []StreetInfo) {
	o.StreetInfo = &v
}

// GetCountryInfo returns the CountryInfo field value if set, zero value otherwise.
func (o *GeocodeResultItem) GetCountryInfo() CountryInfo {
	if o == nil || o.CountryInfo == nil {
		var ret CountryInfo
		return ret
	}
	return *o.CountryInfo
}

// GetCountryInfoOk returns a tuple with the CountryInfo field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *GeocodeResultItem) GetCountryInfoOk() (*CountryInfo, bool) {
	if o == nil || o.CountryInfo == nil {
		return nil, false
	}
	return o.CountryInfo, true
}

// HasCountryInfo returns a boolean if a field has been set.
func (o *GeocodeResultItem) HasCountryInfo() bool {
	if o != nil && o.CountryInfo != nil {
		return true
	}

	return false
}

// SetCountryInfo gets a reference to the given CountryInfo and assigns it to the CountryInfo field.
func (o *GeocodeResultItem) SetCountryInfo(v CountryInfo) {
	o.CountryInfo = &v
}

func (o GeocodeResultItem) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if true {
		toSerialize["title"] = o.Title
	}
	if o.Id != nil {
		toSerialize["id"] = o.Id
	}
	if o.PoliticalView != nil {
		toSerialize["politicalView"] = o.PoliticalView
	}
	if o.ResultType != nil {
		toSerialize["resultType"] = o.ResultType
	}
	if o.HouseNumberType != nil {
		toSerialize["houseNumberType"] = o.HouseNumberType
	}
	if o.AddressBlockType != nil {
		toSerialize["addressBlockType"] = o.AddressBlockType
	}
	if o.LocalityType != nil {
		toSerialize["localityType"] = o.LocalityType
	}
	if o.AdministrativeAreaType != nil {
		toSerialize["administrativeAreaType"] = o.AdministrativeAreaType
	}
	if true {
		toSerialize["address"] = o.Address
	}
	if o.Position != nil {
		toSerialize["position"] = o.Position
	}
	if o.Access != nil {
		toSerialize["access"] = o.Access
	}
	if o.Distance != nil {
		toSerialize["distance"] = o.Distance
	}
	if o.MapView != nil {
		toSerialize["mapView"] = o.MapView
	}
	if o.Categories != nil {
		toSerialize["categories"] = o.Categories
	}
	if o.FoodTypes != nil {
		toSerialize["foodTypes"] = o.FoodTypes
	}
	if o.HouseNumberFallback != nil {
		toSerialize["houseNumberFallback"] = o.HouseNumberFallback
	}
	if o.TimeZone != nil {
		toSerialize["timeZone"] = o.TimeZone
	}
	if o.Scoring != nil {
		toSerialize["scoring"] = o.Scoring
	}
	if o.Parsing != nil {
		toSerialize["parsing"] = o.Parsing
	}
	if o.StreetInfo != nil {
		toSerialize["streetInfo"] = o.StreetInfo
	}
	if o.CountryInfo != nil {
		toSerialize["countryInfo"] = o.CountryInfo
	}
	return json.Marshal(toSerialize)
}

type NullableGeocodeResultItem struct {
	value *GeocodeResultItem
	isSet bool
}

func (v NullableGeocodeResultItem) Get() *GeocodeResultItem {
	return v.value
}

func (v *NullableGeocodeResultItem) Set(val *GeocodeResultItem) {
	v.value = val
	v.isSet = true
}

func (v NullableGeocodeResultItem) IsSet() bool {
	return v.isSet
}

func (v *NullableGeocodeResultItem) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableGeocodeResultItem(val *GeocodeResultItem) *NullableGeocodeResultItem {
	return &NullableGeocodeResultItem{value: val, isSet: true}
}

func (v NullableGeocodeResultItem) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableGeocodeResultItem) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


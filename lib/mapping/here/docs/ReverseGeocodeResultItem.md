# ReverseGeocodeResultItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | The localized display name of this result item. | 
**Id** | Pointer to **string** | The unique identifier for the result item. This ID can be used for a Look Up by ID search as well. | [optional] 
**PoliticalView** | Pointer to **string** | ISO3 country code of the item political view (default for international). This response element is populated when the politicalView parameter is set in the query | [optional] 
**ResultType** | Pointer to **string** | WARNING: The resultType values &#39;intersection&#39; and &#39;postalCodePoint&#39; are in BETA state | [optional] 
**HouseNumberType** | Pointer to **string** | * PA - a Point Address represents an individual address as a point object. Point Addresses are coming from trusted sources.   We can say with high certainty that the address exists and at what position. A Point Address result contains two types of coordinates.   One is the access point (or navigation coordinates), which is the point to start or end a drive. The other point is the position or display point.   This point varies per source and country. The point can be the rooftop point, a point close to the building entry, or a point close to the building,   driveway or parking lot that belongs to the building. * interpolated - an interpolated address. These are approximate positions as a result of a linear interpolation based on address ranges.   Address ranges, especially in the USA, are typical per block. For interpolated addresses, we cannot say with confidence that the address exists in reality.   But the interpolation provides a good location approximation that brings people in most use cases close to the target location.   The access point of an interpolated address result is calculated based on the address range and the road geometry.   The position (display) point is pre-configured offset from the street geometry.   Compared to Point Addresses, interpolated addresses are less accurate. | [optional] 
**AddressBlockType** | Pointer to **string** |  | [optional] 
**LocalityType** | Pointer to **string** |  | [optional] 
**AdministrativeAreaType** | Pointer to **string** |  | [optional] 
**Address** | [**Address**](Address.md) | Postal address of the result item. | 
**Position** | Pointer to [**DisplayResponseCoordinate**](DisplayResponseCoordinate.md) | The coordinates (latitude, longitude) of a pin on a map corresponding to the searched place. | [optional] 
**Access** | Pointer to [**[]AccessResponseCoordinate**](AccessResponseCoordinate.md) | Coordinates of the place you are navigating to (for example, driving or walking). This is a point on a road or in a parking lot. | [optional] 
**Distance** | Pointer to **int64** | The distance \\\&quot;as the crow flies\\\&quot; from the search center to this result item in meters. For example: \\\&quot;172039\\\&quot;.  When searching along a route this is the distance\\nalong the route plus the distance from the route polyline to this result item. | [optional] 
**MapView** | Pointer to [**MapView**](MapView.md) | The bounding box enclosing the geometric shape (area or line) that an individual result covers. &#x60;place&#x60; typed results have no &#x60;mapView&#x60;. | [optional] 
**Categories** | Pointer to [**[]Category**](Category.md) | The list of categories assigned to this place. | [optional] 
**FoodTypes** | Pointer to [**[]Category**](Category.md) | The list of food types assigned to this place. | [optional] 
**HouseNumberFallback** | Pointer to **bool** | If true, indicates that the requested house number was corrected to match the nearest known house number. This field is visible only when the value is true. | [optional] 
**TimeZone** | Pointer to [**TimeZoneInfo**](TimeZoneInfo.md) | BETA - Provides time zone information for this place. (rendered only if &#39;show&#x3D;tz&#39; is provided.) | [optional] 
**StreetInfo** | Pointer to [**[]StreetInfo**](StreetInfo.md) | Street Details (only rendered if &#39;show&#x3D;streetInfo&#39; is provided.) | [optional] 
**CountryInfo** | Pointer to [**CountryInfo**](CountryInfo.md) | Country Details (only rendered if &#39;show&#x3D;countryInfo&#39; is provided.) | [optional] 

## Methods

### NewReverseGeocodeResultItem

`func NewReverseGeocodeResultItem(title string, address Address, ) *ReverseGeocodeResultItem`

NewReverseGeocodeResultItem instantiates a new ReverseGeocodeResultItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReverseGeocodeResultItemWithDefaults

`func NewReverseGeocodeResultItemWithDefaults() *ReverseGeocodeResultItem`

NewReverseGeocodeResultItemWithDefaults instantiates a new ReverseGeocodeResultItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *ReverseGeocodeResultItem) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *ReverseGeocodeResultItem) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *ReverseGeocodeResultItem) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetId

`func (o *ReverseGeocodeResultItem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ReverseGeocodeResultItem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ReverseGeocodeResultItem) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *ReverseGeocodeResultItem) HasId() bool`

HasId returns a boolean if a field has been set.

### GetPoliticalView

`func (o *ReverseGeocodeResultItem) GetPoliticalView() string`

GetPoliticalView returns the PoliticalView field if non-nil, zero value otherwise.

### GetPoliticalViewOk

`func (o *ReverseGeocodeResultItem) GetPoliticalViewOk() (*string, bool)`

GetPoliticalViewOk returns a tuple with the PoliticalView field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoliticalView

`func (o *ReverseGeocodeResultItem) SetPoliticalView(v string)`

SetPoliticalView sets PoliticalView field to given value.

### HasPoliticalView

`func (o *ReverseGeocodeResultItem) HasPoliticalView() bool`

HasPoliticalView returns a boolean if a field has been set.

### GetResultType

`func (o *ReverseGeocodeResultItem) GetResultType() string`

GetResultType returns the ResultType field if non-nil, zero value otherwise.

### GetResultTypeOk

`func (o *ReverseGeocodeResultItem) GetResultTypeOk() (*string, bool)`

GetResultTypeOk returns a tuple with the ResultType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResultType

`func (o *ReverseGeocodeResultItem) SetResultType(v string)`

SetResultType sets ResultType field to given value.

### HasResultType

`func (o *ReverseGeocodeResultItem) HasResultType() bool`

HasResultType returns a boolean if a field has been set.

### GetHouseNumberType

`func (o *ReverseGeocodeResultItem) GetHouseNumberType() string`

GetHouseNumberType returns the HouseNumberType field if non-nil, zero value otherwise.

### GetHouseNumberTypeOk

`func (o *ReverseGeocodeResultItem) GetHouseNumberTypeOk() (*string, bool)`

GetHouseNumberTypeOk returns a tuple with the HouseNumberType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumberType

`func (o *ReverseGeocodeResultItem) SetHouseNumberType(v string)`

SetHouseNumberType sets HouseNumberType field to given value.

### HasHouseNumberType

`func (o *ReverseGeocodeResultItem) HasHouseNumberType() bool`

HasHouseNumberType returns a boolean if a field has been set.

### GetAddressBlockType

`func (o *ReverseGeocodeResultItem) GetAddressBlockType() string`

GetAddressBlockType returns the AddressBlockType field if non-nil, zero value otherwise.

### GetAddressBlockTypeOk

`func (o *ReverseGeocodeResultItem) GetAddressBlockTypeOk() (*string, bool)`

GetAddressBlockTypeOk returns a tuple with the AddressBlockType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddressBlockType

`func (o *ReverseGeocodeResultItem) SetAddressBlockType(v string)`

SetAddressBlockType sets AddressBlockType field to given value.

### HasAddressBlockType

`func (o *ReverseGeocodeResultItem) HasAddressBlockType() bool`

HasAddressBlockType returns a boolean if a field has been set.

### GetLocalityType

`func (o *ReverseGeocodeResultItem) GetLocalityType() string`

GetLocalityType returns the LocalityType field if non-nil, zero value otherwise.

### GetLocalityTypeOk

`func (o *ReverseGeocodeResultItem) GetLocalityTypeOk() (*string, bool)`

GetLocalityTypeOk returns a tuple with the LocalityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalityType

`func (o *ReverseGeocodeResultItem) SetLocalityType(v string)`

SetLocalityType sets LocalityType field to given value.

### HasLocalityType

`func (o *ReverseGeocodeResultItem) HasLocalityType() bool`

HasLocalityType returns a boolean if a field has been set.

### GetAdministrativeAreaType

`func (o *ReverseGeocodeResultItem) GetAdministrativeAreaType() string`

GetAdministrativeAreaType returns the AdministrativeAreaType field if non-nil, zero value otherwise.

### GetAdministrativeAreaTypeOk

`func (o *ReverseGeocodeResultItem) GetAdministrativeAreaTypeOk() (*string, bool)`

GetAdministrativeAreaTypeOk returns a tuple with the AdministrativeAreaType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdministrativeAreaType

`func (o *ReverseGeocodeResultItem) SetAdministrativeAreaType(v string)`

SetAdministrativeAreaType sets AdministrativeAreaType field to given value.

### HasAdministrativeAreaType

`func (o *ReverseGeocodeResultItem) HasAdministrativeAreaType() bool`

HasAdministrativeAreaType returns a boolean if a field has been set.

### GetAddress

`func (o *ReverseGeocodeResultItem) GetAddress() Address`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *ReverseGeocodeResultItem) GetAddressOk() (*Address, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *ReverseGeocodeResultItem) SetAddress(v Address)`

SetAddress sets Address field to given value.


### GetPosition

`func (o *ReverseGeocodeResultItem) GetPosition() DisplayResponseCoordinate`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *ReverseGeocodeResultItem) GetPositionOk() (*DisplayResponseCoordinate, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *ReverseGeocodeResultItem) SetPosition(v DisplayResponseCoordinate)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *ReverseGeocodeResultItem) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### GetAccess

`func (o *ReverseGeocodeResultItem) GetAccess() []AccessResponseCoordinate`

GetAccess returns the Access field if non-nil, zero value otherwise.

### GetAccessOk

`func (o *ReverseGeocodeResultItem) GetAccessOk() (*[]AccessResponseCoordinate, bool)`

GetAccessOk returns a tuple with the Access field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccess

`func (o *ReverseGeocodeResultItem) SetAccess(v []AccessResponseCoordinate)`

SetAccess sets Access field to given value.

### HasAccess

`func (o *ReverseGeocodeResultItem) HasAccess() bool`

HasAccess returns a boolean if a field has been set.

### GetDistance

`func (o *ReverseGeocodeResultItem) GetDistance() int64`

GetDistance returns the Distance field if non-nil, zero value otherwise.

### GetDistanceOk

`func (o *ReverseGeocodeResultItem) GetDistanceOk() (*int64, bool)`

GetDistanceOk returns a tuple with the Distance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistance

`func (o *ReverseGeocodeResultItem) SetDistance(v int64)`

SetDistance sets Distance field to given value.

### HasDistance

`func (o *ReverseGeocodeResultItem) HasDistance() bool`

HasDistance returns a boolean if a field has been set.

### GetMapView

`func (o *ReverseGeocodeResultItem) GetMapView() MapView`

GetMapView returns the MapView field if non-nil, zero value otherwise.

### GetMapViewOk

`func (o *ReverseGeocodeResultItem) GetMapViewOk() (*MapView, bool)`

GetMapViewOk returns a tuple with the MapView field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMapView

`func (o *ReverseGeocodeResultItem) SetMapView(v MapView)`

SetMapView sets MapView field to given value.

### HasMapView

`func (o *ReverseGeocodeResultItem) HasMapView() bool`

HasMapView returns a boolean if a field has been set.

### GetCategories

`func (o *ReverseGeocodeResultItem) GetCategories() []Category`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *ReverseGeocodeResultItem) GetCategoriesOk() (*[]Category, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *ReverseGeocodeResultItem) SetCategories(v []Category)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *ReverseGeocodeResultItem) HasCategories() bool`

HasCategories returns a boolean if a field has been set.

### GetFoodTypes

`func (o *ReverseGeocodeResultItem) GetFoodTypes() []Category`

GetFoodTypes returns the FoodTypes field if non-nil, zero value otherwise.

### GetFoodTypesOk

`func (o *ReverseGeocodeResultItem) GetFoodTypesOk() (*[]Category, bool)`

GetFoodTypesOk returns a tuple with the FoodTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFoodTypes

`func (o *ReverseGeocodeResultItem) SetFoodTypes(v []Category)`

SetFoodTypes sets FoodTypes field to given value.

### HasFoodTypes

`func (o *ReverseGeocodeResultItem) HasFoodTypes() bool`

HasFoodTypes returns a boolean if a field has been set.

### GetHouseNumberFallback

`func (o *ReverseGeocodeResultItem) GetHouseNumberFallback() bool`

GetHouseNumberFallback returns the HouseNumberFallback field if non-nil, zero value otherwise.

### GetHouseNumberFallbackOk

`func (o *ReverseGeocodeResultItem) GetHouseNumberFallbackOk() (*bool, bool)`

GetHouseNumberFallbackOk returns a tuple with the HouseNumberFallback field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumberFallback

`func (o *ReverseGeocodeResultItem) SetHouseNumberFallback(v bool)`

SetHouseNumberFallback sets HouseNumberFallback field to given value.

### HasHouseNumberFallback

`func (o *ReverseGeocodeResultItem) HasHouseNumberFallback() bool`

HasHouseNumberFallback returns a boolean if a field has been set.

### GetTimeZone

`func (o *ReverseGeocodeResultItem) GetTimeZone() TimeZoneInfo`

GetTimeZone returns the TimeZone field if non-nil, zero value otherwise.

### GetTimeZoneOk

`func (o *ReverseGeocodeResultItem) GetTimeZoneOk() (*TimeZoneInfo, bool)`

GetTimeZoneOk returns a tuple with the TimeZone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeZone

`func (o *ReverseGeocodeResultItem) SetTimeZone(v TimeZoneInfo)`

SetTimeZone sets TimeZone field to given value.

### HasTimeZone

`func (o *ReverseGeocodeResultItem) HasTimeZone() bool`

HasTimeZone returns a boolean if a field has been set.

### GetStreetInfo

`func (o *ReverseGeocodeResultItem) GetStreetInfo() []StreetInfo`

GetStreetInfo returns the StreetInfo field if non-nil, zero value otherwise.

### GetStreetInfoOk

`func (o *ReverseGeocodeResultItem) GetStreetInfoOk() (*[]StreetInfo, bool)`

GetStreetInfoOk returns a tuple with the StreetInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreetInfo

`func (o *ReverseGeocodeResultItem) SetStreetInfo(v []StreetInfo)`

SetStreetInfo sets StreetInfo field to given value.

### HasStreetInfo

`func (o *ReverseGeocodeResultItem) HasStreetInfo() bool`

HasStreetInfo returns a boolean if a field has been set.

### GetCountryInfo

`func (o *ReverseGeocodeResultItem) GetCountryInfo() CountryInfo`

GetCountryInfo returns the CountryInfo field if non-nil, zero value otherwise.

### GetCountryInfoOk

`func (o *ReverseGeocodeResultItem) GetCountryInfoOk() (*CountryInfo, bool)`

GetCountryInfoOk returns a tuple with the CountryInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCountryInfo

`func (o *ReverseGeocodeResultItem) SetCountryInfo(v CountryInfo)`

SetCountryInfo sets CountryInfo field to given value.

### HasCountryInfo

`func (o *ReverseGeocodeResultItem) HasCountryInfo() bool`

HasCountryInfo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



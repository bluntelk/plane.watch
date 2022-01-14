# AutosuggestEntityResultItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | The localized display name of this result item. | 
**Id** | Pointer to **string** | The unique identifier for the result item. This ID can be used for a Look Up by ID search as well. | [optional] 
**PoliticalView** | Pointer to **string** | ISO3 country code of the item political view (default for international). This response element is populated when the politicalView parameter is set in the query | [optional] 
**OntologyId** | Pointer to **string** | Related ontology ID | [optional] 
**ResultType** | Pointer to **string** | WARNING: The resultType values &#39;intersection&#39; and &#39;postalCodePoint&#39; are in BETA state | [optional] 
**HouseNumberType** | Pointer to **string** | * PA - a Point Address represents an individual address as a point object. Point Addresses are coming from trusted sources.   We can say with high certainty that the address exists and at what position. A Point Address result contains two types of coordinates.   One is the access point (or navigation coordinates), which is the point to start or end a drive. The other point is the position or display point.   This point varies per source and country. The point can be the rooftop point, a point close to the building entry, or a point close to the building,   driveway or parking lot that belongs to the building. * interpolated - an interpolated address. These are approximate positions as a result of a linear interpolation based on address ranges.   Address ranges, especially in the USA, are typical per block. For interpolated addresses, we cannot say with confidence that the address exists in reality.   But the interpolation provides a good location approximation that brings people in most use cases close to the target location.   The access point of an interpolated address result is calculated based on the address range and the road geometry.   The position (display) point is pre-configured offset from the street geometry.   Compared to Point Addresses, interpolated addresses are less accurate. | [optional] 
**AddressBlockType** | Pointer to **string** |  | [optional] 
**LocalityType** | Pointer to **string** |  | [optional] 
**AdministrativeAreaType** | Pointer to **string** |  | [optional] 
**Address** | [**Address**](Address.md) | Postal address of the result item. Detailed address fields are rendered only if &#39;show&#x3D;details&#39; is provided. | 
**Position** | Pointer to [**DisplayResponseCoordinate**](DisplayResponseCoordinate.md) | The coordinates (latitude, longitude) of a pin on a map corresponding to the searched place. | [optional] 
**Access** | Pointer to [**[]AccessResponseCoordinate**](AccessResponseCoordinate.md) | Coordinates of the place you are navigating to (for example, driving or walking). This is a point on a road or in a parking lot. | [optional] 
**Distance** | Pointer to **int64** | The distance \\\&quot;as the crow flies\\\&quot; from the search center to this result item in meters. For example: \\\&quot;172039\\\&quot;.  When searching along a route this is the distance\\nalong the route plus the distance from the route polyline to this result item. | [optional] 
**ExcursionDistance** | Pointer to **int64** | BETA - Two times the distance from the polyline to this result item in meters while searching along the route | [optional] 
**MapView** | Pointer to [**MapView**](MapView.md) | The bounding box enclosing the geometric shape (area or line) that an individual result covers. &#x60;place&#x60; typed results have no &#x60;mapView&#x60;. | [optional] 
**Categories** | Pointer to [**[]Category**](Category.md) | The list of categories assigned to this place. | [optional] 
**Chains** | Pointer to [**[]Chain**](Chain.md) | The list of chains assigned to this place. | [optional] 
**References** | Pointer to [**[]SupplierReference**](SupplierReference.md) | The list of supplier references available for this place. | [optional] 
**FoodTypes** | Pointer to [**[]Category**](Category.md) | The list of food types assigned to this place. | [optional] 
**Contacts** | Pointer to [**[]ContactInformation**](ContactInformation.md) | BETA - Contact information like phone, email, WWW. (rendered only if &#39;show&#x3D;details&#39; is provided.) | [optional] 
**OpeningHours** | Pointer to [**[]OpeningHours**](OpeningHours.md) | BETA - A list of hours during which the place is open for business. (rendered only if &#39;show&#x3D;details&#39; is provided.) | [optional] 
**TimeZone** | Pointer to [**TimeZoneInfo**](TimeZoneInfo.md) | BETA - Provides time zone information for this place. (rendered only if &#39;show&#x3D;tz&#39; is provided.) | [optional] 
**Highlights** | Pointer to [**TitleAndAddressHighlighting**](TitleAndAddressHighlighting.md) | Describes how the parts of the response element matched the input query | [optional] 
**Phonemes** | Pointer to [**PhonemesSection**](PhonemesSection.md) | Phonemes for address and place names. (rendered only if &#39;show&#x3D;phonemes&#39; is provided.) | [optional] 
**StreetInfo** | Pointer to [**[]StreetInfo**](StreetInfo.md) | Street Details (only rendered if &#39;show&#x3D;streetInfo&#39; is provided.) | [optional] 

## Methods

### NewAutosuggestEntityResultItem

`func NewAutosuggestEntityResultItem(title string, address Address, ) *AutosuggestEntityResultItem`

NewAutosuggestEntityResultItem instantiates a new AutosuggestEntityResultItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAutosuggestEntityResultItemWithDefaults

`func NewAutosuggestEntityResultItemWithDefaults() *AutosuggestEntityResultItem`

NewAutosuggestEntityResultItemWithDefaults instantiates a new AutosuggestEntityResultItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *AutosuggestEntityResultItem) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *AutosuggestEntityResultItem) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *AutosuggestEntityResultItem) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetId

`func (o *AutosuggestEntityResultItem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AutosuggestEntityResultItem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AutosuggestEntityResultItem) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AutosuggestEntityResultItem) HasId() bool`

HasId returns a boolean if a field has been set.

### GetPoliticalView

`func (o *AutosuggestEntityResultItem) GetPoliticalView() string`

GetPoliticalView returns the PoliticalView field if non-nil, zero value otherwise.

### GetPoliticalViewOk

`func (o *AutosuggestEntityResultItem) GetPoliticalViewOk() (*string, bool)`

GetPoliticalViewOk returns a tuple with the PoliticalView field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoliticalView

`func (o *AutosuggestEntityResultItem) SetPoliticalView(v string)`

SetPoliticalView sets PoliticalView field to given value.

### HasPoliticalView

`func (o *AutosuggestEntityResultItem) HasPoliticalView() bool`

HasPoliticalView returns a boolean if a field has been set.

### GetOntologyId

`func (o *AutosuggestEntityResultItem) GetOntologyId() string`

GetOntologyId returns the OntologyId field if non-nil, zero value otherwise.

### GetOntologyIdOk

`func (o *AutosuggestEntityResultItem) GetOntologyIdOk() (*string, bool)`

GetOntologyIdOk returns a tuple with the OntologyId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOntologyId

`func (o *AutosuggestEntityResultItem) SetOntologyId(v string)`

SetOntologyId sets OntologyId field to given value.

### HasOntologyId

`func (o *AutosuggestEntityResultItem) HasOntologyId() bool`

HasOntologyId returns a boolean if a field has been set.

### GetResultType

`func (o *AutosuggestEntityResultItem) GetResultType() string`

GetResultType returns the ResultType field if non-nil, zero value otherwise.

### GetResultTypeOk

`func (o *AutosuggestEntityResultItem) GetResultTypeOk() (*string, bool)`

GetResultTypeOk returns a tuple with the ResultType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResultType

`func (o *AutosuggestEntityResultItem) SetResultType(v string)`

SetResultType sets ResultType field to given value.

### HasResultType

`func (o *AutosuggestEntityResultItem) HasResultType() bool`

HasResultType returns a boolean if a field has been set.

### GetHouseNumberType

`func (o *AutosuggestEntityResultItem) GetHouseNumberType() string`

GetHouseNumberType returns the HouseNumberType field if non-nil, zero value otherwise.

### GetHouseNumberTypeOk

`func (o *AutosuggestEntityResultItem) GetHouseNumberTypeOk() (*string, bool)`

GetHouseNumberTypeOk returns a tuple with the HouseNumberType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumberType

`func (o *AutosuggestEntityResultItem) SetHouseNumberType(v string)`

SetHouseNumberType sets HouseNumberType field to given value.

### HasHouseNumberType

`func (o *AutosuggestEntityResultItem) HasHouseNumberType() bool`

HasHouseNumberType returns a boolean if a field has been set.

### GetAddressBlockType

`func (o *AutosuggestEntityResultItem) GetAddressBlockType() string`

GetAddressBlockType returns the AddressBlockType field if non-nil, zero value otherwise.

### GetAddressBlockTypeOk

`func (o *AutosuggestEntityResultItem) GetAddressBlockTypeOk() (*string, bool)`

GetAddressBlockTypeOk returns a tuple with the AddressBlockType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddressBlockType

`func (o *AutosuggestEntityResultItem) SetAddressBlockType(v string)`

SetAddressBlockType sets AddressBlockType field to given value.

### HasAddressBlockType

`func (o *AutosuggestEntityResultItem) HasAddressBlockType() bool`

HasAddressBlockType returns a boolean if a field has been set.

### GetLocalityType

`func (o *AutosuggestEntityResultItem) GetLocalityType() string`

GetLocalityType returns the LocalityType field if non-nil, zero value otherwise.

### GetLocalityTypeOk

`func (o *AutosuggestEntityResultItem) GetLocalityTypeOk() (*string, bool)`

GetLocalityTypeOk returns a tuple with the LocalityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalityType

`func (o *AutosuggestEntityResultItem) SetLocalityType(v string)`

SetLocalityType sets LocalityType field to given value.

### HasLocalityType

`func (o *AutosuggestEntityResultItem) HasLocalityType() bool`

HasLocalityType returns a boolean if a field has been set.

### GetAdministrativeAreaType

`func (o *AutosuggestEntityResultItem) GetAdministrativeAreaType() string`

GetAdministrativeAreaType returns the AdministrativeAreaType field if non-nil, zero value otherwise.

### GetAdministrativeAreaTypeOk

`func (o *AutosuggestEntityResultItem) GetAdministrativeAreaTypeOk() (*string, bool)`

GetAdministrativeAreaTypeOk returns a tuple with the AdministrativeAreaType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdministrativeAreaType

`func (o *AutosuggestEntityResultItem) SetAdministrativeAreaType(v string)`

SetAdministrativeAreaType sets AdministrativeAreaType field to given value.

### HasAdministrativeAreaType

`func (o *AutosuggestEntityResultItem) HasAdministrativeAreaType() bool`

HasAdministrativeAreaType returns a boolean if a field has been set.

### GetAddress

`func (o *AutosuggestEntityResultItem) GetAddress() Address`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *AutosuggestEntityResultItem) GetAddressOk() (*Address, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *AutosuggestEntityResultItem) SetAddress(v Address)`

SetAddress sets Address field to given value.


### GetPosition

`func (o *AutosuggestEntityResultItem) GetPosition() DisplayResponseCoordinate`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *AutosuggestEntityResultItem) GetPositionOk() (*DisplayResponseCoordinate, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *AutosuggestEntityResultItem) SetPosition(v DisplayResponseCoordinate)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *AutosuggestEntityResultItem) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### GetAccess

`func (o *AutosuggestEntityResultItem) GetAccess() []AccessResponseCoordinate`

GetAccess returns the Access field if non-nil, zero value otherwise.

### GetAccessOk

`func (o *AutosuggestEntityResultItem) GetAccessOk() (*[]AccessResponseCoordinate, bool)`

GetAccessOk returns a tuple with the Access field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccess

`func (o *AutosuggestEntityResultItem) SetAccess(v []AccessResponseCoordinate)`

SetAccess sets Access field to given value.

### HasAccess

`func (o *AutosuggestEntityResultItem) HasAccess() bool`

HasAccess returns a boolean if a field has been set.

### GetDistance

`func (o *AutosuggestEntityResultItem) GetDistance() int64`

GetDistance returns the Distance field if non-nil, zero value otherwise.

### GetDistanceOk

`func (o *AutosuggestEntityResultItem) GetDistanceOk() (*int64, bool)`

GetDistanceOk returns a tuple with the Distance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistance

`func (o *AutosuggestEntityResultItem) SetDistance(v int64)`

SetDistance sets Distance field to given value.

### HasDistance

`func (o *AutosuggestEntityResultItem) HasDistance() bool`

HasDistance returns a boolean if a field has been set.

### GetExcursionDistance

`func (o *AutosuggestEntityResultItem) GetExcursionDistance() int64`

GetExcursionDistance returns the ExcursionDistance field if non-nil, zero value otherwise.

### GetExcursionDistanceOk

`func (o *AutosuggestEntityResultItem) GetExcursionDistanceOk() (*int64, bool)`

GetExcursionDistanceOk returns a tuple with the ExcursionDistance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExcursionDistance

`func (o *AutosuggestEntityResultItem) SetExcursionDistance(v int64)`

SetExcursionDistance sets ExcursionDistance field to given value.

### HasExcursionDistance

`func (o *AutosuggestEntityResultItem) HasExcursionDistance() bool`

HasExcursionDistance returns a boolean if a field has been set.

### GetMapView

`func (o *AutosuggestEntityResultItem) GetMapView() MapView`

GetMapView returns the MapView field if non-nil, zero value otherwise.

### GetMapViewOk

`func (o *AutosuggestEntityResultItem) GetMapViewOk() (*MapView, bool)`

GetMapViewOk returns a tuple with the MapView field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMapView

`func (o *AutosuggestEntityResultItem) SetMapView(v MapView)`

SetMapView sets MapView field to given value.

### HasMapView

`func (o *AutosuggestEntityResultItem) HasMapView() bool`

HasMapView returns a boolean if a field has been set.

### GetCategories

`func (o *AutosuggestEntityResultItem) GetCategories() []Category`

GetCategories returns the Categories field if non-nil, zero value otherwise.

### GetCategoriesOk

`func (o *AutosuggestEntityResultItem) GetCategoriesOk() (*[]Category, bool)`

GetCategoriesOk returns a tuple with the Categories field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCategories

`func (o *AutosuggestEntityResultItem) SetCategories(v []Category)`

SetCategories sets Categories field to given value.

### HasCategories

`func (o *AutosuggestEntityResultItem) HasCategories() bool`

HasCategories returns a boolean if a field has been set.

### GetChains

`func (o *AutosuggestEntityResultItem) GetChains() []Chain`

GetChains returns the Chains field if non-nil, zero value otherwise.

### GetChainsOk

`func (o *AutosuggestEntityResultItem) GetChainsOk() (*[]Chain, bool)`

GetChainsOk returns a tuple with the Chains field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChains

`func (o *AutosuggestEntityResultItem) SetChains(v []Chain)`

SetChains sets Chains field to given value.

### HasChains

`func (o *AutosuggestEntityResultItem) HasChains() bool`

HasChains returns a boolean if a field has been set.

### GetReferences

`func (o *AutosuggestEntityResultItem) GetReferences() []SupplierReference`

GetReferences returns the References field if non-nil, zero value otherwise.

### GetReferencesOk

`func (o *AutosuggestEntityResultItem) GetReferencesOk() (*[]SupplierReference, bool)`

GetReferencesOk returns a tuple with the References field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReferences

`func (o *AutosuggestEntityResultItem) SetReferences(v []SupplierReference)`

SetReferences sets References field to given value.

### HasReferences

`func (o *AutosuggestEntityResultItem) HasReferences() bool`

HasReferences returns a boolean if a field has been set.

### GetFoodTypes

`func (o *AutosuggestEntityResultItem) GetFoodTypes() []Category`

GetFoodTypes returns the FoodTypes field if non-nil, zero value otherwise.

### GetFoodTypesOk

`func (o *AutosuggestEntityResultItem) GetFoodTypesOk() (*[]Category, bool)`

GetFoodTypesOk returns a tuple with the FoodTypes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFoodTypes

`func (o *AutosuggestEntityResultItem) SetFoodTypes(v []Category)`

SetFoodTypes sets FoodTypes field to given value.

### HasFoodTypes

`func (o *AutosuggestEntityResultItem) HasFoodTypes() bool`

HasFoodTypes returns a boolean if a field has been set.

### GetContacts

`func (o *AutosuggestEntityResultItem) GetContacts() []ContactInformation`

GetContacts returns the Contacts field if non-nil, zero value otherwise.

### GetContactsOk

`func (o *AutosuggestEntityResultItem) GetContactsOk() (*[]ContactInformation, bool)`

GetContactsOk returns a tuple with the Contacts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContacts

`func (o *AutosuggestEntityResultItem) SetContacts(v []ContactInformation)`

SetContacts sets Contacts field to given value.

### HasContacts

`func (o *AutosuggestEntityResultItem) HasContacts() bool`

HasContacts returns a boolean if a field has been set.

### GetOpeningHours

`func (o *AutosuggestEntityResultItem) GetOpeningHours() []OpeningHours`

GetOpeningHours returns the OpeningHours field if non-nil, zero value otherwise.

### GetOpeningHoursOk

`func (o *AutosuggestEntityResultItem) GetOpeningHoursOk() (*[]OpeningHours, bool)`

GetOpeningHoursOk returns a tuple with the OpeningHours field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOpeningHours

`func (o *AutosuggestEntityResultItem) SetOpeningHours(v []OpeningHours)`

SetOpeningHours sets OpeningHours field to given value.

### HasOpeningHours

`func (o *AutosuggestEntityResultItem) HasOpeningHours() bool`

HasOpeningHours returns a boolean if a field has been set.

### GetTimeZone

`func (o *AutosuggestEntityResultItem) GetTimeZone() TimeZoneInfo`

GetTimeZone returns the TimeZone field if non-nil, zero value otherwise.

### GetTimeZoneOk

`func (o *AutosuggestEntityResultItem) GetTimeZoneOk() (*TimeZoneInfo, bool)`

GetTimeZoneOk returns a tuple with the TimeZone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeZone

`func (o *AutosuggestEntityResultItem) SetTimeZone(v TimeZoneInfo)`

SetTimeZone sets TimeZone field to given value.

### HasTimeZone

`func (o *AutosuggestEntityResultItem) HasTimeZone() bool`

HasTimeZone returns a boolean if a field has been set.

### GetHighlights

`func (o *AutosuggestEntityResultItem) GetHighlights() TitleAndAddressHighlighting`

GetHighlights returns the Highlights field if non-nil, zero value otherwise.

### GetHighlightsOk

`func (o *AutosuggestEntityResultItem) GetHighlightsOk() (*TitleAndAddressHighlighting, bool)`

GetHighlightsOk returns a tuple with the Highlights field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHighlights

`func (o *AutosuggestEntityResultItem) SetHighlights(v TitleAndAddressHighlighting)`

SetHighlights sets Highlights field to given value.

### HasHighlights

`func (o *AutosuggestEntityResultItem) HasHighlights() bool`

HasHighlights returns a boolean if a field has been set.

### GetPhonemes

`func (o *AutosuggestEntityResultItem) GetPhonemes() PhonemesSection`

GetPhonemes returns the Phonemes field if non-nil, zero value otherwise.

### GetPhonemesOk

`func (o *AutosuggestEntityResultItem) GetPhonemesOk() (*PhonemesSection, bool)`

GetPhonemesOk returns a tuple with the Phonemes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPhonemes

`func (o *AutosuggestEntityResultItem) SetPhonemes(v PhonemesSection)`

SetPhonemes sets Phonemes field to given value.

### HasPhonemes

`func (o *AutosuggestEntityResultItem) HasPhonemes() bool`

HasPhonemes returns a boolean if a field has been set.

### GetStreetInfo

`func (o *AutosuggestEntityResultItem) GetStreetInfo() []StreetInfo`

GetStreetInfo returns the StreetInfo field if non-nil, zero value otherwise.

### GetStreetInfoOk

`func (o *AutosuggestEntityResultItem) GetStreetInfoOk() (*[]StreetInfo, bool)`

GetStreetInfoOk returns a tuple with the StreetInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreetInfo

`func (o *AutosuggestEntityResultItem) SetStreetInfo(v []StreetInfo)`

SetStreetInfo sets StreetInfo field to given value.

### HasStreetInfo

`func (o *AutosuggestEntityResultItem) HasStreetInfo() bool`

HasStreetInfo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



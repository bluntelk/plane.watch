# AutocompleteResultItem

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Title** | **string** | The unified display name of this result item. The result title is composed so that the customer application can use it to render the suggestions with highlighting. It is build in a unified way for all the countries starting from the country name and down to the address line. It is build out of the address components that are important for the end-user to recognize and eventually to choose a result and includes all the input terms. For example: \&quot;Germany, 32547, Bad Oeynhausen, Schulstraße 4\&quot; | 
**Id** | Pointer to **string** | The unique identifier for the result item. This ID can be used for a Look Up by ID search as well. | [optional] 
**Language** | Pointer to **string** | The preferred language of address elements in the result. | [optional] 
**PoliticalView** | Pointer to **string** | ISO3 country code of the item political view (default for international). This response element is populated when the politicalView parameter is set in the query | [optional] 
**ResultType** | Pointer to **string** | WARNING: The resultType values &#39;intersection&#39; and &#39;postalCodePoint&#39; are in BETA state | [optional] 
**HouseNumberType** | Pointer to **string** |  | [optional] 
**LocalityType** | Pointer to **string** |  | [optional] 
**AdministrativeAreaType** | Pointer to **string** |  | [optional] 
**Address** | [**Address**](Address.md) | Detailed address of the result item. | 
**Distance** | Pointer to **int64** | The distance \\\&quot;as the crow flies\\\&quot; from the search center to this result item in meters. For example: \\\&quot;172039\\\&quot;.  When searching along a route this is the distance\\nalong the route plus the distance from the route polyline to this result item. | [optional] 
**Highlights** | Pointer to [**TitleAndAddressHighlighting**](TitleAndAddressHighlighting.md) | Describes how the parts of the response element matched the input query | [optional] 
**StreetInfo** | Pointer to [**[]StreetInfo**](StreetInfo.md) | Street Details (only rendered if &#39;show&#x3D;streetInfo&#39; is provided.) | [optional] 

## Methods

### NewAutocompleteResultItem

`func NewAutocompleteResultItem(title string, address Address, ) *AutocompleteResultItem`

NewAutocompleteResultItem instantiates a new AutocompleteResultItem object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAutocompleteResultItemWithDefaults

`func NewAutocompleteResultItemWithDefaults() *AutocompleteResultItem`

NewAutocompleteResultItemWithDefaults instantiates a new AutocompleteResultItem object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTitle

`func (o *AutocompleteResultItem) GetTitle() string`

GetTitle returns the Title field if non-nil, zero value otherwise.

### GetTitleOk

`func (o *AutocompleteResultItem) GetTitleOk() (*string, bool)`

GetTitleOk returns a tuple with the Title field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTitle

`func (o *AutocompleteResultItem) SetTitle(v string)`

SetTitle sets Title field to given value.


### GetId

`func (o *AutocompleteResultItem) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *AutocompleteResultItem) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *AutocompleteResultItem) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *AutocompleteResultItem) HasId() bool`

HasId returns a boolean if a field has been set.

### GetLanguage

`func (o *AutocompleteResultItem) GetLanguage() string`

GetLanguage returns the Language field if non-nil, zero value otherwise.

### GetLanguageOk

`func (o *AutocompleteResultItem) GetLanguageOk() (*string, bool)`

GetLanguageOk returns a tuple with the Language field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLanguage

`func (o *AutocompleteResultItem) SetLanguage(v string)`

SetLanguage sets Language field to given value.

### HasLanguage

`func (o *AutocompleteResultItem) HasLanguage() bool`

HasLanguage returns a boolean if a field has been set.

### GetPoliticalView

`func (o *AutocompleteResultItem) GetPoliticalView() string`

GetPoliticalView returns the PoliticalView field if non-nil, zero value otherwise.

### GetPoliticalViewOk

`func (o *AutocompleteResultItem) GetPoliticalViewOk() (*string, bool)`

GetPoliticalViewOk returns a tuple with the PoliticalView field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoliticalView

`func (o *AutocompleteResultItem) SetPoliticalView(v string)`

SetPoliticalView sets PoliticalView field to given value.

### HasPoliticalView

`func (o *AutocompleteResultItem) HasPoliticalView() bool`

HasPoliticalView returns a boolean if a field has been set.

### GetResultType

`func (o *AutocompleteResultItem) GetResultType() string`

GetResultType returns the ResultType field if non-nil, zero value otherwise.

### GetResultTypeOk

`func (o *AutocompleteResultItem) GetResultTypeOk() (*string, bool)`

GetResultTypeOk returns a tuple with the ResultType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetResultType

`func (o *AutocompleteResultItem) SetResultType(v string)`

SetResultType sets ResultType field to given value.

### HasResultType

`func (o *AutocompleteResultItem) HasResultType() bool`

HasResultType returns a boolean if a field has been set.

### GetHouseNumberType

`func (o *AutocompleteResultItem) GetHouseNumberType() string`

GetHouseNumberType returns the HouseNumberType field if non-nil, zero value otherwise.

### GetHouseNumberTypeOk

`func (o *AutocompleteResultItem) GetHouseNumberTypeOk() (*string, bool)`

GetHouseNumberTypeOk returns a tuple with the HouseNumberType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHouseNumberType

`func (o *AutocompleteResultItem) SetHouseNumberType(v string)`

SetHouseNumberType sets HouseNumberType field to given value.

### HasHouseNumberType

`func (o *AutocompleteResultItem) HasHouseNumberType() bool`

HasHouseNumberType returns a boolean if a field has been set.

### GetLocalityType

`func (o *AutocompleteResultItem) GetLocalityType() string`

GetLocalityType returns the LocalityType field if non-nil, zero value otherwise.

### GetLocalityTypeOk

`func (o *AutocompleteResultItem) GetLocalityTypeOk() (*string, bool)`

GetLocalityTypeOk returns a tuple with the LocalityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalityType

`func (o *AutocompleteResultItem) SetLocalityType(v string)`

SetLocalityType sets LocalityType field to given value.

### HasLocalityType

`func (o *AutocompleteResultItem) HasLocalityType() bool`

HasLocalityType returns a boolean if a field has been set.

### GetAdministrativeAreaType

`func (o *AutocompleteResultItem) GetAdministrativeAreaType() string`

GetAdministrativeAreaType returns the AdministrativeAreaType field if non-nil, zero value otherwise.

### GetAdministrativeAreaTypeOk

`func (o *AutocompleteResultItem) GetAdministrativeAreaTypeOk() (*string, bool)`

GetAdministrativeAreaTypeOk returns a tuple with the AdministrativeAreaType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdministrativeAreaType

`func (o *AutocompleteResultItem) SetAdministrativeAreaType(v string)`

SetAdministrativeAreaType sets AdministrativeAreaType field to given value.

### HasAdministrativeAreaType

`func (o *AutocompleteResultItem) HasAdministrativeAreaType() bool`

HasAdministrativeAreaType returns a boolean if a field has been set.

### GetAddress

`func (o *AutocompleteResultItem) GetAddress() Address`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *AutocompleteResultItem) GetAddressOk() (*Address, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *AutocompleteResultItem) SetAddress(v Address)`

SetAddress sets Address field to given value.


### GetDistance

`func (o *AutocompleteResultItem) GetDistance() int64`

GetDistance returns the Distance field if non-nil, zero value otherwise.

### GetDistanceOk

`func (o *AutocompleteResultItem) GetDistanceOk() (*int64, bool)`

GetDistanceOk returns a tuple with the Distance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDistance

`func (o *AutocompleteResultItem) SetDistance(v int64)`

SetDistance sets Distance field to given value.

### HasDistance

`func (o *AutocompleteResultItem) HasDistance() bool`

HasDistance returns a boolean if a field has been set.

### GetHighlights

`func (o *AutocompleteResultItem) GetHighlights() TitleAndAddressHighlighting`

GetHighlights returns the Highlights field if non-nil, zero value otherwise.

### GetHighlightsOk

`func (o *AutocompleteResultItem) GetHighlightsOk() (*TitleAndAddressHighlighting, bool)`

GetHighlightsOk returns a tuple with the Highlights field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHighlights

`func (o *AutocompleteResultItem) SetHighlights(v TitleAndAddressHighlighting)`

SetHighlights sets Highlights field to given value.

### HasHighlights

`func (o *AutocompleteResultItem) HasHighlights() bool`

HasHighlights returns a boolean if a field has been set.

### GetStreetInfo

`func (o *AutocompleteResultItem) GetStreetInfo() []StreetInfo`

GetStreetInfo returns the StreetInfo field if non-nil, zero value otherwise.

### GetStreetInfoOk

`func (o *AutocompleteResultItem) GetStreetInfoOk() (*[]StreetInfo, bool)`

GetStreetInfoOk returns a tuple with the StreetInfo field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreetInfo

`func (o *AutocompleteResultItem) SetStreetInfo(v []StreetInfo)`

SetStreetInfo sets StreetInfo field to given value.

### HasStreetInfo

`func (o *AutocompleteResultItem) HasStreetInfo() bool`

HasStreetInfo returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



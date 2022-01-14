# Scoring

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**QueryScore** | Pointer to **float64** | Indicates how good the input matches the returned address. It is equal to 1 if all input tokens are recognized and matched. | [optional] 
**FieldScore** | Pointer to [**FieldScore**](FieldScore.md) | Indicates how good the individual result fields match to the corresponding part of the query. Is included only for the result fields that are actually matched to the query. | [optional] 

## Methods

### NewScoring

`func NewScoring() *Scoring`

NewScoring instantiates a new Scoring object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewScoringWithDefaults

`func NewScoringWithDefaults() *Scoring`

NewScoringWithDefaults instantiates a new Scoring object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetQueryScore

`func (o *Scoring) GetQueryScore() float64`

GetQueryScore returns the QueryScore field if non-nil, zero value otherwise.

### GetQueryScoreOk

`func (o *Scoring) GetQueryScoreOk() (*float64, bool)`

GetQueryScoreOk returns a tuple with the QueryScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueryScore

`func (o *Scoring) SetQueryScore(v float64)`

SetQueryScore sets QueryScore field to given value.

### HasQueryScore

`func (o *Scoring) HasQueryScore() bool`

HasQueryScore returns a boolean if a field has been set.

### GetFieldScore

`func (o *Scoring) GetFieldScore() FieldScore`

GetFieldScore returns the FieldScore field if non-nil, zero value otherwise.

### GetFieldScoreOk

`func (o *Scoring) GetFieldScoreOk() (*FieldScore, bool)`

GetFieldScoreOk returns a tuple with the FieldScore field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFieldScore

`func (o *Scoring) SetFieldScore(v FieldScore)`

SetFieldScore sets FieldScore field to given value.

### HasFieldScore

`func (o *Scoring) HasFieldScore() bool`

HasFieldScore returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)



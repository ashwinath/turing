/*
 * Turing Minimal Openapi Spec for SDK
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"fmt"
)

// Dataset - struct for Dataset
type Dataset struct {
	BigQueryDataset *BigQueryDataset
}

// BigQueryDatasetAsDataset is a convenience function that returns BigQueryDataset wrapped in Dataset
func BigQueryDatasetAsDataset(v *BigQueryDataset) Dataset {
	return Dataset{BigQueryDataset: v}
}

// Unmarshal JSON data into one of the pointers in the struct
func (dst *Dataset) UnmarshalJSON(data []byte) error {
	var err error
	match := 0
	// try to unmarshal data into BigQueryDataset
	err = json.Unmarshal(data, &dst.BigQueryDataset)
	if err == nil {
		jsonBigQueryDataset, _ := json.Marshal(dst.BigQueryDataset)
		if string(jsonBigQueryDataset) == "{}" { // empty struct
			dst.BigQueryDataset = nil
		} else {
			match++
		}
	} else {
		dst.BigQueryDataset = nil
	}

	if match > 1 { // more than 1 match
		// reset to nil
		dst.BigQueryDataset = nil

		return fmt.Errorf("Data matches more than one schema in oneOf(Dataset)")
	} else if match == 1 {
		return nil // exactly one match
	} else { // no match
		return fmt.Errorf("Data failed to match schemas in oneOf(Dataset)")
	}
}

// Marshal data from the first non-nil pointers in the struct to JSON
func (src Dataset) MarshalJSON() ([]byte, error) {
	if src.BigQueryDataset != nil {
		return json.Marshal(&src.BigQueryDataset)
	}

	return nil, nil // no data in oneOf schemas
}

// Get the actual instance
func (obj *Dataset) GetActualInstance() interface{} {
	if obj.BigQueryDataset != nil {
		return obj.BigQueryDataset
	}

	// all schemas are nil
	return nil
}

type NullableDataset struct {
	value *Dataset
	isSet bool
}

func (v NullableDataset) Get() *Dataset {
	return v.value
}

func (v *NullableDataset) Set(val *Dataset) {
	v.value = val
	v.isSet = true
}

func (v NullableDataset) IsSet() bool {
	return v.isSet
}

func (v *NullableDataset) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDataset(val *Dataset) *NullableDataset {
	return &NullableDataset{value: val, isSet: true}
}

func (v NullableDataset) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDataset) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
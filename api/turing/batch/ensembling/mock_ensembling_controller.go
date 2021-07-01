// Code generated by mockery v2.6.0. DO NOT EDIT.

package batchensembling

import (
	models "github.com/gojek/turing/api/turing/models"
	mock "github.com/stretchr/testify/mock"
)

// MockEnsemblingController is an autogenerated mock type for the EnsemblingController type
type MockEnsemblingController struct {
	mock.Mock
}

// Create provides a mock function with given fields: request
func (_m *MockEnsemblingController) Create(request *CreateEnsemblingJobRequest) error {
	ret := _m.Called(request)

	var r0 error
	if rf, ok := ret.Get(0).(func(*CreateEnsemblingJobRequest) error); ok {
		r0 = rf(request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: namespace, ensemblingJob
func (_m *MockEnsemblingController) Delete(namespace string, ensemblingJob *models.EnsemblingJob) error {
	ret := _m.Called(namespace, ensemblingJob)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *models.EnsemblingJob) error); ok {
		r0 = rf(namespace, ensemblingJob)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetStatus provides a mock function with given fields: namespace, ensemblingJob
func (_m *MockEnsemblingController) GetStatus(namespace string, ensemblingJob *models.EnsemblingJob) (SparkApplicationState, error) {
	ret := _m.Called(namespace, ensemblingJob)

	var r0 SparkApplicationState
	if rf, ok := ret.Get(0).(func(string, *models.EnsemblingJob) SparkApplicationState); ok {
		r0 = rf(namespace, ensemblingJob)
	} else {
		r0 = ret.Get(0).(SparkApplicationState)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *models.EnsemblingJob) error); ok {
		r1 = rf(namespace, ensemblingJob)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

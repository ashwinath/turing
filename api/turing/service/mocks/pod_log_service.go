// Code generated by mockery v2.6.0. DO NOT EDIT.

package mocks

import (
	client "github.com/gojek/mlp/api/client"
	mock "github.com/stretchr/testify/mock"

	models "github.com/gojek/turing/api/turing/models"

	service "github.com/gojek/turing/api/turing/service"
)

// PodLogService is an autogenerated mock type for the PodLogService type
type PodLogService struct {
	mock.Mock
}

// ListEnsemblingJobPodLogs provides a mock function with given fields: ensemblingJobName, project, componentType, opts
func (_m *PodLogService) ListEnsemblingJobPodLogs(ensemblingJobName string, project *client.Project, componentType string, opts *service.PodLogOptions) ([]*service.PodLog, error) {
	ret := _m.Called(ensemblingJobName, project, componentType, opts)

	var r0 []*service.PodLog
	if rf, ok := ret.Get(0).(func(string, *client.Project, string, *service.PodLogOptions) []*service.PodLog); ok {
		r0 = rf(ensemblingJobName, project, componentType, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*service.PodLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, *client.Project, string, *service.PodLogOptions) error); ok {
		r1 = rf(ensemblingJobName, project, componentType, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListRouterPodLogs provides a mock function with given fields: project, router, routerVersion, componentType, opts
func (_m *PodLogService) ListRouterPodLogs(project *client.Project, router *models.Router, routerVersion *models.RouterVersion, componentType string, opts *service.PodLogOptions) ([]*service.PodLog, error) {
	ret := _m.Called(project, router, routerVersion, componentType, opts)

	var r0 []*service.PodLog
	if rf, ok := ret.Get(0).(func(*client.Project, *models.Router, *models.RouterVersion, string, *service.PodLogOptions) []*service.PodLog); ok {
		r0 = rf(project, router, routerVersion, componentType, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*service.PodLog)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*client.Project, *models.Router, *models.RouterVersion, string, *service.PodLogOptions) error); ok {
		r1 = rf(project, router, routerVersion, componentType, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

package validation_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/gojek/turing/api/turing/api/request"
	"github.com/gojek/turing/api/turing/models"
	"github.com/gojek/turing/api/turing/service/mocks"
	"github.com/gojek/turing/api/turing/validation"
	"github.com/gojek/turing/engines/experiment/manager"
	expRequest "github.com/gojek/turing/engines/experiment/pkg/request"
	"github.com/gojek/turing/engines/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLogConfig(t *testing.T) {
	tt := map[string]struct {
		input  request.LogConfig
		hasErr bool
	}{
		"valid_nop": {
			input: request.LogConfig{
				ResultLoggerType: "nop",
			},
			hasErr: false,
		},
		"invalid_type": {
			input: request.LogConfig{
				ResultLoggerType: "nope",
			},
			hasErr: true,
		},
		"valid_bq": {
			input: request.LogConfig{
				ResultLoggerType: "bigquery",
				BigQueryConfig: &request.BigQueryConfig{
					Table:                "project.dataset.table",
					ServiceAccountSecret: "acc",
				},
			},
			hasErr: false,
		},
		"bq_missing_config": {
			input: request.LogConfig{
				ResultLoggerType: "bigquery",
			},
			hasErr: true,
		},
		"bq_invalid_table": {
			input: request.LogConfig{
				ResultLoggerType: "bigquery",
				BigQueryConfig: &request.BigQueryConfig{
					Table:                "project:dataset.table",
					ServiceAccountSecret: "acc",
				},
			},
			hasErr: true,
		},
		"bq_invalid_svc_account": {
			input: request.LogConfig{
				ResultLoggerType: "bigquery",
				BigQueryConfig: &request.BigQueryConfig{
					Table:                "project.dataset.table",
					ServiceAccountSecret: "",
				},
			},
			hasErr: true,
		},
		"kafka_valid_config": {
			input: request.LogConfig{
				ResultLoggerType: "kafka",
				KafkaConfig: &request.KafkaConfig{
					Brokers:             "broker1,broker2",
					Topic:               "topic",
					SerializationFormat: "json",
				},
			},
			hasErr: false,
		},
		"kafka_missing_config": {
			input: request.LogConfig{
				ResultLoggerType: "kafka",
			},
			hasErr: true,
		},
		"kafka_invalid_config_missing_brokers": {
			input: request.LogConfig{
				ResultLoggerType: "kafka",
				KafkaConfig: &request.KafkaConfig{
					Topic:               "topic",
					SerializationFormat: "json",
				},
			},
			hasErr: true,
		},
		"kafka_invalid_config_missing_topic": {
			input: request.LogConfig{
				ResultLoggerType: "kafka",
				KafkaConfig: &request.KafkaConfig{
					Brokers:             "broker1,broker2",
					SerializationFormat: "json",
				},
			},
			hasErr: true,
		},
		"kafka_invalid_config_invalid_serialization": {
			input: request.LogConfig{
				ResultLoggerType: "kafka",
				KafkaConfig: &request.KafkaConfig{
					Brokers:             "broker1,broker2",
					Topic:               "topic",
					SerializationFormat: "test",
				},
			},
			hasErr: true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			mockExperimentsService := &mocks.ExperimentsService{}
			mockExperimentsService.On("ListEngines").Return([]manager.Engine{})
			validate, err := validation.NewValidator(mockExperimentsService)
			assert.NoError(t, err)
			err = validate.Struct(&tc.input)
			assert.Equal(t, tc.hasErr, err != nil)
		})
	}
}

func TestValidateExperimentEngineConfig(t *testing.T) {
	validationErr := "test-error"
	// Create mock experiment service
	config1 := json.RawMessage(`{"client": {"id": "1", "username": "1"}}`)
	config2 := json.RawMessage(`{"client": {"id": "2", "username": "2"}}`)
	expSvc := &mocks.ExperimentsService{}
	expSvc.On("ListEngines").Return([]manager.Engine{{Name: "custom"}})
	expSvc.On("ValidateExperimentConfig", "custom", config1).
		Return(nil)
	expSvc.On("ValidateExperimentConfig", "custom", config2).
		Return(errors.New(validationErr))

	// Define tests
	tests := map[string]struct {
		input request.ExperimentEngineConfig
		err   string
	}{
		"success | valid nop": {
			input: request.ExperimentEngineConfig{
				Type: "nop",
			},
		},
		"success | valid experiment config": {
			input: request.ExperimentEngineConfig{
				Type:   "custom",
				Config: config1,
			},
		},
		"failure | unknown engine type": {
			input: request.ExperimentEngineConfig{
				Type:   "unknown",
				Config: config2,
			},
			err: "Key: 'ExperimentEngineConfig.type' " +
				"Error:Field validation for 'type' failed on the 'oneof' tag",
		},
		"failure | validation error": {
			input: request.ExperimentEngineConfig{
				Type:   "custom",
				Config: config2,
			},
			err: fmt.Sprintf(
				"Key: 'ExperimentEngineConfig.config' "+
					"Error:Field validation for 'config' failed on the '%s' tag", validationErr),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			validate, err := validation.NewValidator(expSvc)
			assert.NoError(t, err)
			err = validate.Struct(&tc.input)
			if tc.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.err)
			}
		})
	}
}

type routerConfigTestCase struct {
	routes         models.Routes
	defaultRouteID string
	trafficRules   models.TrafficRules
	expectedError  string
}

func (tt routerConfigTestCase) RouterConfig() *request.RouterConfig {
	return &request.RouterConfig{
		Routes:         tt.routes,
		DefaultRouteID: tt.defaultRouteID,
		TrafficRules:   tt.trafficRules,
		ExperimentEngine: &request.ExperimentEngineConfig{
			Type: "nop",
		},
		Timeout: "20s",
		LogConfig: &request.LogConfig{
			ResultLoggerType: "nop",
		},
	}
}

func TestValidateTrafficRules(t *testing.T) {
	suite := map[string]routerConfigTestCase{
		"success": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
				{
					ID:       "route-b",
					Type:     "PROXY",
					Endpoint: "http://example.com/b",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: expRequest.HeaderFieldSource,
							Field:       "X-Region",
							Operator:    router.InConditionOperator,
							Values:      []string{"region-a", "region-b"},
						},
					},
					Routes: []string{"route-b"},
				},
			},
		},
		"failure | empty conditions": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
				{
					ID:       "route-b",
					Type:     "PROXY",
					Endpoint: "http://example.com/b",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{},
					Routes:     []string{"route-b"},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Conditions' " +
				"Error:Field validation for 'Conditions' failed on the 'notBlank' tag",
		},
		"failure | empty routes": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: expRequest.HeaderFieldSource,
							Field:       "X-Region",
							Operator:    router.InConditionOperator,
							Values:      []string{"region-b"},
						},
					},
					Routes: []string{},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Routes' " +
				"Error:Field validation for 'Routes' failed on the 'notBlank' tag",
		},
		"failure | unsupported operator": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
				{
					ID:       "route-b",
					Type:     "PROXY",
					Endpoint: "http://example.com/b",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: expRequest.HeaderFieldSource,
							Field:       "X-Region",
							Operator:    router.RuleConditionOperator{},
							Values:      []string{"region-b"},
						},
					},
					Routes: []string{"route-b"},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Conditions[0].Operator' " +
				"Error:Field validation for 'Operator' failed on the 'required' tag",
		},
		"failure | unsupported field source": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
				{
					ID:       "route-b",
					Type:     "PROXY",
					Endpoint: "http://example.com/b",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: "unknown",
							Field:       "X-Region",
							Operator:    router.InConditionOperator,
							Values:      []string{"region-b"},
						},
					},
					Routes: []string{"route-b"},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Conditions[0].FieldSource' " +
				"Error:Field validation for 'FieldSource' failed on the 'oneof' tag",
		},
		"failure | incomplete condition": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
				{
					ID:       "route-b",
					Type:     "PROXY",
					Endpoint: "http://example.com/b",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: expRequest.HeaderFieldSource,
							Field:       "",
							Operator:    router.InConditionOperator,
							Values:      []string{},
						},
					},
					Routes: []string{"route-b"},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Conditions[0].Field' " +
				"Error:Field validation for 'Field' failed on the 'required' tag\n" +
				"Key: 'RouterConfig.TrafficRules[0].Conditions[0].Values' " +
				"Error:Field validation for 'Values' failed on the 'notBlank' tag",
		},
		"failure | unknown route": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: expRequest.PayloadFieldSource,
							Field:       "some_property",
							Operator:    router.InConditionOperator,
							Values:      []string{"some_value"},
						},
					},
					Routes: []string{"route-c"},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Routes[0]' " +
				"Error:Field validation for '' failed on the 'oneof' tag",
		},
		"failure | rule contains default route": {
			routes: models.Routes{
				{
					ID:       "route-a",
					Type:     "PROXY",
					Endpoint: "http://example.com/a",
					Timeout:  "10ms",
				},
			},
			defaultRouteID: "route-a",
			trafficRules: models.TrafficRules{
				{
					Conditions: []*router.TrafficRuleCondition{
						{
							FieldSource: expRequest.PayloadFieldSource,
							Field:       "some_property",
							Operator:    router.InConditionOperator,
							Values:      []string{"some_value"},
						},
					},
					Routes: []string{"route-a"},
				},
			},
			expectedError: "Key: 'RouterConfig.TrafficRules[0].Routes[0]' " +
				"Error:Field validation for '' failed on the 'necsfield' tag",
		},
	}

	for name, tt := range suite {
		t.Run(name, func(t *testing.T) {
			mockExperimentsService := &mocks.ExperimentsService{}
			mockExperimentsService.On("ListEngines").Return([]manager.Engine{})
			validate, err := validation.NewValidator(mockExperimentsService)
			require.NoError(t, err)

			err = validate.Struct(tt.RouterConfig())
			if tt.expectedError == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.expectedError)
			}
		})
	}
}

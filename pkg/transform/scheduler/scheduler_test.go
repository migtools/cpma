package scheduler_test

import (
	"io/ioutil"
	"testing"

	cpmatest "github.com/konveyor/cpma/pkg/transform/internal/test"
	"github.com/konveyor/cpma/pkg/transform/scheduler"
	configv1 "github.com/openshift/api/config/v1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	k8sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/kubernetes/scheme"
)

func loadExpectedScheduler(file string) (*configv1.Scheduler, error) {
	expectedContent, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	expectedSchedule := new(configv1.Scheduler)
	serializer := k8sjson.NewYAMLSerializer(k8sjson.DefaultMetaFactory, scheme.Scheme, scheme.Scheme)
	_, _, err = serializer.Decode(expectedContent, nil, expectedSchedule)
	if err != nil {
		return nil, err
	}

	return expectedSchedule, nil
}

func TestTransformScheduler(t *testing.T) {
	t.Parallel()
	masterConfig, err := cpmatest.LoadMasterConfig("testdata/master_config.yaml")
	require.NoError(t, err)

	expectedCrd, err := loadExpectedScheduler("testdata/expected-CR-scheduler.yaml")
	require.NoError(t, err)

	testCases := []struct {
		name        string
		expectedCrd *configv1.Scheduler
	}{
		{
			name:        "build basic scheduler",
			expectedCrd: expectedCrd,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			schedulerResources, err := scheduler.Translate(*masterConfig)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCrd, schedulerResources)
		})
	}
}

func TestBasicAuthValidation(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		requireError bool
		inputFile    string
		expectedErr  error
	}{
		{
			name:         "fail on invalid ",
			requireError: true,
			inputFile:    "testdata/master_config-invalid-defaultnodeselector.yaml",
			expectedErr:  errors.New("DefaultNodeSelector can't be empty"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			masterConfig, err := cpmatest.LoadMasterConfig(tc.inputFile)
			require.NoError(t, err)

			err = scheduler.Validate(*masterConfig)

			if tc.requireError {
				assert.Equal(t, tc.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

package registries

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		expectError int
		data        map[string]Registries
	}{
		{
			name:        "Unvalid emtpy registries ",
			expectError: 1,
			data:        map[string]Registries{},
		},
		{
			name:        "Undefined registries",
			expectError: 1,
			data: map[string]Registries{
				"test1": Registries{
					List: []string{"test1"},
				},
				"test2": Registries{
					List: []string{"test2"},
				},
				"test3": Registries{
					List: []string{"test3"},
				},
			},
		},
		{
			name:        "Missing registries",
			expectError: 1,
			data: map[string]Registries{
				"block": Registries{
					List: []string{},
				},
				"insecure": Registries{
					List: []string{},
				},
				"search": Registries{
					List: []string{},
				},
			},
		},
		{
			name:        "Valid registries",
			expectError: 0,
			data: map[string]Registries{
				"block": Registries{
					List: []string{"test1"},
				},
				"insecure": Registries{
					List: []string{"test2"},
				},
				"search": Registries{
					List: []string{"test3"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			err := Validate(tc.data)
			assert.Equal(t, err, tc.expectError)
		})
	}
}

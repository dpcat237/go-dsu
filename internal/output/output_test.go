package output_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dpcat237/go-dsu/internal/output"
)

// TestOutput tests the Output's method.
func TestOutput(t *testing.T) {
	testCases := []struct {
		// test details
		description string
		expected    string
		// struct data
		error    error
		method   string
		mode     output.Mode
		response string
	}{
		{
			description: "Successful response",
			expected:    "Updated successfully 3 dependencies",
			error:       nil,
			method:      "updater.updateDependencies",
			mode:        output.ModeProd,
			response:    "Updated successfully 3 dependencies",
		},
		{
			description: "Error response",
			expected:    "Error: downloading data",
			error:       errors.New("Error: downloading data"),
			method:      "updater.updateDependencies",
			mode:        output.ModeProd,
			response:    "Updated successfully 3 dependencies",
		},
		{
			description: "Error response in development mode",
			expected:    "[updater.updateDependencies] Error: downloading data",
			error:       errors.New("Error: downloading data"),
			method:      "updater.updateDependencies",
			mode:        output.ModeDev,
			response:    "Updated successfully 3 dependencies",
		},
	}

	t.Log("Test of Output")
	{
		for i, tc := range testCases {
			t.Logf("\tTest %d:\t%s", i, tc.description)
			{
				out := output.Create("updater.updateDependencies").WithResponse(tc.response).WithError(tc.error)
				assert.Equal(t, tc.expected, out.ToString(tc.mode))
			}
		}
	}
}

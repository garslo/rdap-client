package protocol

import (
	"encoding/json"
	"fmt"
	"reflect"

	"testing"
)

var jsonExample = []byte(`{
       "version": "1.0",
       "publication": "2015-04-17T16:00:00Z",
       "description": "Some text",
       "services": [
         [
           ["entry1", "entry2", "entry3"],
           [
             "https://registry.example.com/myrdap/",
             "http://registry.example.com/myrdap/"
           ]
         ],
         [
           ["entry4"],
           [
             "http://example.org/"
           ]
         ]
       ]
   }`)

func TestConformity(t *testing.T) {
	if err := json.Unmarshal(jsonExample, &ServiceRegistry{}); err != nil {
		t.Fatal(err)
	}
}

func TestMatchAS(t *testing.T) {
	tests := []struct {
		description   string
		registry      ServiceRegistry
		as            uint32
		expected      []string
		expectedError error
	}{
		{
			description: "it should match an as number",
			as:          65411,
			registry: ServiceRegistry{
				Services: ServicesList{
					{
						{"2045-2045"},
						{"https://rir3.example.com/myrdap/"},
					},
					{
						{"10000-12000", "300000-400000"},
						{"http://example.org/"},
					},
					{
						{"64512-65534"},
						{"http://example.net/rdaprir2/", "https://example.net/rdaprir2/"},
					},
				},
			},
			expected: []string{"http://example.net/rdaprir2/", "https://example.net/rdaprir2/"},
		},
		{
			description: "it should not match an as number due to invalid beginning of as range",
			as:          1,
			registry: ServiceRegistry{
				Services: ServicesList{
					{
						{"invalid-123"},
						{},
					},
				},
			},
			expectedError: fmt.Errorf("strconv.ParseInt: parsing \"invalid\": invalid syntax"),
		},
		{
			description: "it should not match an as number due to invalid end of as range",
			as:          1,
			registry: ServiceRegistry{
				Services: ServicesList{
					{
						{"123-invalid"},
						{},
					},
				},
			},
			expectedError: fmt.Errorf("strconv.ParseInt: parsing \"invalid\": invalid syntax"),
		},
	}

	for i, test := range tests {
		urls, err := test.registry.MatchAS(test.as)

		if test.expectedError != nil && err != nil {
			if test.expectedError.Error() != err.Error() {
				t.Fatalf("At index %d (%s): expected error %s, got %s", i, test.description, test.expectedError, err)
			}
		}

		if !reflect.DeepEqual(test.expected, urls) {
			t.Fatalf("At index %d (%s): expected %v, got %v", i, test.description, test.expected, urls)
		}
	}
}
package fns

import (
	"strings"
	"testing"

	"github.com/coredns/caddy"
)

func TestFNSParse(t *testing.T) {
	tests := []struct {
		input              string
		shouldErr          bool
		expectedErrContent string
		expectedRRsRepoNil bool
	}{
		// positive
		{
			"fns {\napi-url https://example.com/api/rrs\n}",
			false,
			"",
			false,
		},

		// negative
		{
			"fns {\napi-url example.com\n}",
			true,
			"api-url scheme not support",
			true,
		},
		{
			"fns",
			true,
			"missing",
			true,
		},
		{
			"fns {\napi-url https://example.com\n bad\n}",
			true,
			"unknown property",
			true,
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		f, err := fnsParse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error but found %s for input %s", i, err, test.input)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Expected no error but found one for input %s. Error was: %v", i, test.input, err)
				continue
			}

			if !strings.Contains(err.Error(), test.expectedErrContent) {
				t.Errorf("Test %d: Expected error to contain: %v, found error: %v, input: %s", i, test.expectedErrContent, err, test.input)
			}
		}

		if !test.shouldErr {
			if test.expectedRRsRepoNil {
				t.Errorf("Test %d, expected rrs-repo %v for input %s, got: %v", i, test.expectedRRsRepoNil, test.input, f.RRsRepo)
			}
		}
	}
}

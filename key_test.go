package kubelock

import (
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_isExpired(t *testing.T) {
	testCases := []struct {
		name            string
		inputLockData   lockData
		expectedExpired bool
		errorMatcher    func(err error) bool
	}{
		{
			name: "case 0: not expired",
			inputLockData: lockData{
				CreatedAt: time.Now().UTC().Add(2 * time.Minute),
				TTL:       30 * time.Second,
			},
			expectedExpired: false,
			errorMatcher:    nil,
		},
		{
			name: "case 1: expired",
			inputLockData: lockData{
				CreatedAt: time.Now().UTC().Add(-5 * time.Minute),
				TTL:       2 * time.Minute,
			},
			expectedExpired: true,
			errorMatcher:    nil,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			expired := isExpired(tc.inputLockData)
			if !cmp.Equal(expired, tc.expectedExpired) {
				t.Fatalf("\n\n%s\n", cmp.Diff(tc.expectedExpired, expired))
			}
		})
	}
}

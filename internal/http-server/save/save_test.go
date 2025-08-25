package save_test

import (
	"testing"
)

func TestSaveHandler(test *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "1_success",
			alias: "ggl",
			url:   "https://google.com",
		},
	}
	for _, testCase := range cases {
		test.Run(testCase.name, func(t *testing.T) {
		})
	}
}

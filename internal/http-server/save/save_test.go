package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"yaus/internal/http-server/save"
	save_test "yaus/internal/http-server/save/mocks"
	"yaus/internal/lib/logger/slogext"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
			urlSaverMock := save_test.NewMockURLSaver(test)

			if testCase.respError == "" || testCase.mockError != nil {
				urlSaverMock.On("SaveURL", testCase.url, mock.AnythingOfType("string")).
					Return(int64(1), testCase.mockError).
					Once()
			}

			handler := save.New(slogext.NewDiscardLogger(), urlSaverMock)

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, testCase.url, testCase.alias)

			request, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
			require.NoError(test, err)

			responseRecorder := httptest.NewRecorder()
			handler.ServeHTTP(responseRecorder, request)

			require.Equal(test, responseRecorder.Code, http.StatusOK)

			body := responseRecorder.Body.String()

			var resp save.Response

			require.NoError(test, json.Unmarshal([]byte(body), &resp))

			require.Equal(test, testCase.respError, resp.Error)
		})
	}
}

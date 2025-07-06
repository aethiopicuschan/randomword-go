package randomword_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aethiopicuschan/randomword-go"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	type testCase struct {
		name        string
		opts        []randomword.Option
		mockDo      func(*http.Request) (*http.Response, error)
		wantWords   []string
		wantErr     error
		paramChecks []string
	}

	tests := []testCase{
		{
			name: "successful response",
			mockDo: func(req *http.Request) (*http.Response, error) {
				words := []string{"foo", "bar", "baz"}
				body, _ := json.Marshal(words)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			wantWords: []string{"foo", "bar", "baz"},
		},
		{
			name: "do returns error",
			mockDo: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network failure")
			},
			wantErr: randomword.ErrInternal,
		},
		{
			name: "non-200 status code",
			mockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("ignored")),
				}, nil
			},
			wantErr: randomword.ErrUnexpectedResponse,
		},
		{
			name: "empty body",
			mockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
			wantErr: randomword.ErrUnexpectedResponse,
		},
		{
			name: "invalid JSON",
			mockDo: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("not a json")),
				}, nil
			},
			wantErr: randomword.ErrInternal,
		},
		{
			name: "URL parameters: number, length, language",
			opts: []randomword.Option{
				randomword.WithNumber(3),
				randomword.WithLength(5),
				randomword.WithLanguage(randomword.Spanish),
			},
			mockDo: func(req *http.Request) (*http.Response, error) {
				words := []string{"x"}
				body, _ := json.Marshal(words)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader(body)),
				}, nil
			},
			wantWords:   []string{"x"},
			paramChecks: []string{"number=3", "length=5", "lang=es"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var capturedURL string
			// wrap the mockDo to capture URL each call
			doFunc := func(req *http.Request) (*http.Response, error) {
				capturedURL = req.URL.String()
				return tc.mockDo(req)
			}

			// build Request with options + custom do
			opts := append(tc.opts, randomword.WithDoFunc(doFunc))
			req, err := randomword.NewRequest(opts...)
			assert.NoError(t, err)

			words, err := req.Fetch()
			if tc.wantErr != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr),
					"expected error %v, got %v", tc.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantWords, words)

			// verify query parameters if needed
			for _, substr := range tc.paramChecks {
				assert.Contains(t, capturedURL, substr,
					"URL %q should contain %q", capturedURL, substr)
			}
		})
	}
}

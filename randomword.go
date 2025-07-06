package randomword

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	BASE_URL = "https://random-word-api.herokuapp.com"
)

// Language represents the language for random words.
type Language string

// Supported languages for random words.
const (
	English             Language = ""
	Spanish             Language = "es"
	Italian             Language = "it"
	German              Language = "de"
	FRench              Language = "fr"
	Chinese             Language = "zh"
	BrazilianPortuguese Language = "pt-br"
)

// Request represents a request for random words.
type Request struct {
	number int
	length int
	lang   Language
	do     func(*http.Request) (*http.Response, error)
}

type Option func(*Request) error

// WithNumber sets the number of words to return.
func WithNumber(n int) Option {
	return func(r *Request) (err error) {
		if n < 1 {
			return ErrNumberLessThanOne
		}
		r.number = n
		return
	}
}

// WithLength sets the length of the words to return.
func WithLength(l int) Option {
	return func(r *Request) (err error) {
		if l < 1 {
			return ErrLengthLessThanOne
		}
		r.length = l
		return
	}
}

// WithLanguage sets the language of the words to return.
func WithLanguage(lang Language) Option {
	return func(r *Request) (err error) {
		switch lang {
		case English, Spanish, Italian, German, FRench, Chinese, BrazilianPortuguese:
			r.lang = lang
			return
		default:
			return ErrUnsupportedLanguage
		}
	}
}

// WithDoFunc sets a custom function to execute the HTTP request.
func WithDoFunc(do func(*http.Request) (*http.Response, error)) Option {
	return func(r *Request) (err error) {
		if do == nil {
			return ErrDoFuncCannotBeNil
		}
		r.do = do
		return
	}
}

// NewRequest creates a new Request with default values.
func NewRequest(opts ...Option) (r *Request, err error) {
	// Default values
	r = &Request{
		number: -1,
		length: -1,
		lang:   English,
		do:     http.DefaultClient.Do,
	}
	for _, opt := range opts {
		if err = opt(r); err != nil {
			return
		}
	}
	return
}

func (r *Request) Fetch() (words []string, err error) {
	endpoint, err := url.JoinPath(BASE_URL, "/word")
	if err != nil {
		return
	}
	params := url.Values{}
	if r.number > 0 {
		params.Add("number", strconv.Itoa(r.number))
	}
	if r.length > 0 {
		params.Add("length", strconv.Itoa(r.length))
	}
	if r.lang != English {
		params.Add("lang", string(r.lang))
	}
	fullURL := endpoint
	if len(params) > 0 {
		fullURL = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		err = errors.Join(ErrInternal, err)
		return
	}
	resp, err := r.do(req)
	if err != nil {
		err = errors.Join(ErrInternal, err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = ErrUnexpectedResponse
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.Join(ErrInternal, err)
		return
	}
	if len(body) == 0 {
		err = ErrUnexpectedResponse
		return
	}
	err = json.Unmarshal(body, &words)
	if err != nil {
		err = errors.Join(ErrInternal, err)
		return
	}
	return
}

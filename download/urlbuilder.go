package download

import (
	"net/url"
	"strconv"
)

type URLBuilder struct {
	baseURL string
	values  url.Values
}

func NewURL(baseURL string) *URLBuilder {
	return &URLBuilder{
		baseURL: baseURL,
		values:  make(url.Values),
	}
}

// Add string parameter (only if value isn't empty)
func (u *URLBuilder) AddString(key, value string) {
	if value != "" {
		u.values.Add(key, value)
	}
}

// Add int parameter (only if value isn't 0)
func (u *URLBuilder) AddInt(key string, value int) {
	if value != 0 {
		u.values.Add(key, strconv.Itoa(value))
	}
}

// Add bool parameter (only if true)
func (u *URLBuilder) AddBool(key string, value bool) {
	if value {
		u.values.Add(key, "true")
	}
}

// Build the final URL string
func (u *URLBuilder) Build() string {
	if len(u.values) == 0 {
		return u.baseURL
	}
	return u.baseURL + "?" + u.values.Encode()
}

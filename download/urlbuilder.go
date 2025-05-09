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

// SetString sets or updates a string parameter (only if value isn't empty)
func (u *URLBuilder) SetString(key, value string) {
	if value != "" {
		u.values.Set(key, value) // Uses Set instead of Add to replace existing
	} else {
		u.values.Del(key) // Remove if empty value
	}
}

// SetInt sets or updates an int parameter (only if value isn't 0)
func (u *URLBuilder) SetInt(key string, value int) {
	if value != 0 {
		u.values.Set(key, strconv.Itoa(value)) // Uses Set instead of Add
	} else {
		u.values.Del(key) // Remove if zero value
	}
}

// SetBool sets or updates a bool parameter (only if true)
func (u *URLBuilder) SetBool(key string, value bool) {
	if value {
		u.values.Set(key, "true") // Uses Set instead of Add
	} else {
		u.values.Del(key) // Remove if false
	}
}

// Build the final URL string
func (u *URLBuilder) Build() string {
	if len(u.values) == 0 {
		return u.baseURL
	}
	return u.baseURL + "?" + u.values.Encode()
}

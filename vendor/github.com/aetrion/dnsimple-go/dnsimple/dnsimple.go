// Package dnsimple provides a client for the DNSimple API.
// In order to use this package you will need a DNSimple account.
package dnsimple

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// libraryVersion identifies the current library version.
	// This is a pro-forma convention given that Go dependencies
	// tends to be fetched directly from the repo.
	// It is also used in the user-agent identify the client.
	libraryVersion = "0.10.0-beta"

	// defaultBaseURL to the DNSimple production API.
	defaultBaseURL = "https://api.dnsimple.com"

	// userAgent represents the default user agent used
	// when no other user agent is set.
	defaultUserAgent = "dnsimple-go/" + libraryVersion

	apiVersion = "v2"
)

// Client represents a client to the DNSimple API.
type Client struct {
	// HttpClient is the underlying HTTP client
	// used to communicate with the API.
	HttpClient *http.Client

	// Credentials used for accessing the DNSimple API
	Credentials Credentials

	// BaseURL for API requests.
	// Defaults to the public DNSimple API, but can be set to a different endpoint (e.g. the sandbox).
	BaseURL string

	// UserAgent used when communicating with the DNSimple API.
	UserAgent string

	// Services used for talking to different parts of the DNSimple API.
	Identity  *IdentityService
	Contacts  *ContactsService
	Domains   *DomainsService
	Oauth     *OauthService
	Registrar *RegistrarService
	Tlds      *TldsService
	Webhooks  *WebhooksService
	Zones     *ZonesService

	// Set to true to output debugging logs during API calls
	Debug bool
}

// ListOptions contains the common options you can pass to a List method
// in order to control parameters such as paginations and page number.
type ListOptions struct {
	// The page to return
	Page int `url:"page,omitempty"`

	// The number of entries to return per page
	PerPage int `url:"per_page,omitempty"`
}

// NewClient returns a new DNSimple API client using the given credentials.
func NewClient(credentials Credentials) *Client {
	c := &Client{Credentials: credentials, HttpClient: &http.Client{}, BaseURL: defaultBaseURL, UserAgent: defaultUserAgent}
	c.Identity = &IdentityService{client: c}
	c.Contacts = &ContactsService{client: c}
	c.Domains = &DomainsService{client: c}
	c.Oauth = &OauthService{client: c}
	c.Registrar = &RegistrarService{client: c}
	c.Tlds = &TldsService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Zones = &ZonesService{client: c}
	return c
}

// NewRequest creates an API request.
// The path is expected to be a relative path and will be resolved
// according to the BaseURL of the Client. Paths should always be specified without a preceding slash.
func (c *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	url := c.BaseURL + path

	body := new(bytes.Buffer)
	if payload != nil {
		err := json.NewEncoder(body).Encode(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", c.UserAgent)
	for key, value := range c.Credentials.Headers() {
		req.Header.Add(key, value)
	}

	return req, nil
}

func versioned(path string) string {
	return fmt.Sprintf("/%s/%s", apiVersion, strings.Trim(path, "/"))
}

func (c *Client) get(path string, obj interface{}) (*http.Response, error) {
	req, err := c.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil, obj)
}

func (c *Client) post(path string, payload, obj interface{}) (*http.Response, error) {
	req, err := c.NewRequest("POST", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

func (c *Client) put(path string, payload, obj interface{}) (*http.Response, error) {
	req, err := c.NewRequest("PUT", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

func (c *Client) patch(path string, payload, obj interface{}) (*http.Response, error) {
	req, err := c.NewRequest("PATCH", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

func (c *Client) delete(path string, payload interface{}, obj interface{}) (*http.Response, error) {
	req, err := c.NewRequest("DELETE", path, payload)
	if err != nil {
		return nil, err
	}

	return c.Do(req, payload, obj)
}

// Do sends an API request and returns the API response.
//
// The API response is JSON decoded and stored in the value pointed by obj,
// or returned as an error if an API error has occurred.
// If obj implements the io.Writer interface, the raw response body will be written to obj,
// without attempting to decode it.
func (c *Client) Do(req *http.Request, payload, obj interface{}) (*http.Response, error) {
	if c.Debug {
		log.Printf("Executing request (%v): %#v", req.URL, req)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if c.Debug {
		log.Printf("Response received: %#v", resp)
	}

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	// If obj implements the io.Writer,
	// the response body is decoded into v.
	if obj != nil {
		if w, ok := obj.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(obj)
		}
	}

	return resp, err
}

// A Response represents an API response.
type Response struct {
	// HTTP response
	HttpResponse *http.Response

	// If the response is paginated, the Pagination will store them.
	Pagination *Pagination `json:"pagination"`
}

// RateLimit returns the maximum amount of requests this account can send in an hour.
func (r *Response) RateLimit() int {
	value, _ := strconv.Atoi(r.HttpResponse.Header.Get("X-RateLimit-Limit"))
	return value
}

// RateLimitRemaining returns the remaining amount of requests this account can send within this hour window.
func (r *Response) RateLimitRemaining() int {
	value, _ := strconv.Atoi(r.HttpResponse.Header.Get("X-RateLimit-Remaining"))
	return value
}

// RateLimitReset returns when the throttling window will be reset for this account.
func (r *Response) RateLimitReset() time.Time {
	value, _ := strconv.ParseInt(r.HttpResponse.Header.Get("X-RateLimit-Reset"), 10, 64)
	return time.Unix(value, 0)
}

// If the response is paginated, Pagination represents the pagination information.
type Pagination struct {
	CurrentPage  int `json:"current_page"`
	PerPage      int `json:"per_page"`
	TotalPages   int `json:"total_pages"`
	TotalEntries int `json:"total_entries"`
}

// An ErrorResponse represents an API response that generated an error.
type ErrorResponse struct {
	Response

	// human-readable message
	Message string `json:"message"`
}

// Error implements the error interface.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %v %v",
		r.HttpResponse.Request.Method, r.HttpResponse.Request.URL,
		r.HttpResponse.StatusCode, r.Message)
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if the status code is different than 2xx. Specific requests
// may have additional requirements, but this is sufficient in most of the cases.
func CheckResponse(resp *http.Response) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{}
	errorResponse.HttpResponse = resp

	err := json.NewDecoder(resp.Body).Decode(errorResponse)
	if err != nil {
		return err
	}

	return errorResponse
}

// see encoding/json
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addURLQueryOptions(path string, options interface{}) (string, error) {
	val := reflect.ValueOf(options)
	qso := map[string]string{}

	// options is a pointer
	// return if the value of the pointer is nil,
	// otherwise replace the pointer with the value.
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return path, nil
		}
		val = val.Elem()
	}

	// extract all the options from the struct
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		sf := typ.Field(i)
		sv := val.Field(i)

		tag := sf.Tag.Get("url")

		// The field has a different tag
		if tag == "" {
			continue
		}

		// The field is ignored with `url:"-"`
		if tag == "-" {
			continue
		}

		splits := strings.Split(tag, ",")
		name, opts := splits[0], splits[1:]

		if optionsContains(opts, "omitempty") && isEmptyValue(sv) {
			continue
		}

		qso[name] = fmt.Sprint(sv.Interface())
	}

	// append the options to the URL
	u, err := url.Parse(path)
	if err != nil {
		return path, err
	}
	qs := u.Query()
	for k, v := range qso {
		qs.Add(k, v)
	}
	u.RawQuery = qs.Encode()

	return u.String(), nil
}

func optionsContains(options []string, option string) bool {
	for _, s := range options {
		if s == option {
			return true
		}
	}
	return false
}

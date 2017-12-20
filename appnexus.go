package appnexus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

// Credentials required to login
type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Client used to make HTTP requests
type Client struct {
	client      *http.Client
	EndPoint    *url.URL
	Rate        Rate
	UserAgent   string
	token       string
	credentials credentials
	MemberID    int

	Members    *MemberService
	Segments   *SegmentService
	Publishers *PublisherService
	Sites      *SiteService
	Placements *PlacementService
	Deals      *DealService
}

// Rate contains information on the current rate limit in operation
type Rate struct {
	Reads             int `json:"reads"`
	ReadLimit         int `json:"read_limit"`
	ReadLimitSeconds  int `json:"read_limit_seconds"`
	Writes            int `json:"writes"`
	WriteLimit        int `json:"write_limit"`
	WriteLimitSeconds int `json:"write_limit_seconds"`
}

// Response is a AppNexus API response object
type Response struct {
	*http.Response
	Obj struct {
		Status           string      `json:"status"`
		ID               json.Number `json:"id,omitempty,Number"`
		ErrorID          string      `json:"error_id,omitempty"`
		Error            string      `json:"error,omitempty"`
		ErrorDescription string      `json:"error_description,omitempty"`
		ErrorCode        string      `json:"error_code,omitempty"`
		Token            string      `json:"token,omitempty"`
		Service          string      `json:"service,omitempty"`
		Method           string      `json:"method,omitempty"`
		Count            int         `json:"count,omitempty"`
		StartElement     int         `json:"start_element,omitempty"`
		NumElements      int         `json:"num_elements,omitempty"`
		Deal             Deal        `json:"deal,omitempty"`
		Placement        Placement   `json:"placement,omitempty"`
		Site             Site        `json:"site,omitempty"`
		Publisher        Publisher   `json:"publisher,omitempty"`
		Member           Member      `json:"member,omitempty"`
		Segments         []Segment   `json:"segments,omitempty"`
		Rate             Rate        `json:"dbg_info"`
	} `json:"response"`
}

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	StartElement int  `url:"start_element,omitempty"`
	NumElements  int  `url:"num_elements,omitempty"`
	Active       bool `url:"active,omitempty"`
}

// NewClient returns a new AppNexus API client
func NewClient(endPointURL string) (*Client, error) {

	httpClient := http.DefaultClient
	baseURL, err := url.Parse(endPointURL)
	if err != nil {
		return nil, errors.New("Invalid AppNexus API endpoint")
	}

	c := &Client{
		client:    httpClient,
		EndPoint:  baseURL,
		UserAgent: "github.com/tnako/appnexus go-appnexus-client",
	}

	c.Members = &MemberService{client: c}
	c.Segments = &SegmentService{client: c}
	c.Publishers = &PublisherService{client: c}
	c.Sites = &SiteService{client: c}
	c.Placements = &PlacementService{client: c}
	c.Deals = &DealService{client: c}

	return c, nil
}

// NewRequest creates an API request using a relative URL
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	u := c.EndPoint.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)

	if c.token != "" {
		req.Header.Add("Authorization", c.token)
	}

	return req, nil
}

// Do sends an API request and returns the API response.  The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred.  If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) do(req *http.Request, v interface{}) (*Response, error) {

	_ = c.waitForRateLimit(req.Method)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.New("client.do.do: " + err.Error())
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("client.do.readall: " + err.Error())
	}

	response, err := c.checkResponse(resp, data)
	if err != nil {

		// If the call failed with a NOAUTH error, attempt to reauthenticate
		// and try the request again:
		if response != nil && req != nil && response.Obj.ErrorID == "NOAUTH" && req.URL.Path != "auth" {
			c.token = ""
			err = c.Login(c.credentials.Username, c.credentials.Password)
			if err != nil {
				return nil, errors.New("Could not reauthenticate:\n" + err.Error())
			}

			return c.do(req, v)
		}

		return nil, errors.New("client.do.checkResponse: " + err.Error())
	}

	if v != nil {
		err := json.Unmarshal(data, v)
		if err != nil {
			return nil, errors.New("client.do.unmarshal: " + err.Error())
		}
	}

	return response, nil
}

// Wait for the Write or Read rate limit timeout
func (c *Client) waitForRateLimit(method string) time.Duration {

	var duration time.Duration

	// Write limit for POST, PUT, DELETE:
	limit := c.Rate.WriteLimit
	actions := c.Rate.Writes
	period := c.Rate.WriteLimitSeconds

	// Read limit for GET:
	if method == "GET" {
		limit = c.Rate.ReadLimit
		actions = c.Rate.Reads
		period = c.Rate.ReadLimitSeconds
	}

	limit--

	// More actions than the limit on the requested operation:
	if actions >= limit {
		duration = time.Duration(period) * time.Second
		time.Sleep(duration)
	}

	return duration
}

// CheckResponse checks the API response for errors, and returns them if
// present.
func (c *Client) checkResponse(r *http.Response, data []byte) (*Response, error) {
	var resp *Response
	if r.StatusCode < 200 || r.StatusCode > 299 {
		// ToDo: add case for code=429 to wait before next request
		return nil, fmt.Errorf("%s | %#v", r.Status, r.Header)
	}

	if len(data) > 0 {

		resp = &Response{Response: r}
		err := json.Unmarshal(data, resp)
		if err != nil {
			return nil, err
		}

		c.Rate = resp.Obj.Rate

		if resp.Obj.ErrorID != "" || resp.Obj.Error != "" {
			str := fmt.Sprintf("AppNexus:checkResponse [%s]: %s", resp.Obj.ErrorID, resp.Obj.Error)
			return resp, errors.New(str)
		}
	}

	return resp, nil
}

// Login to the AppNexus API and get an authentication token
func (c *Client) Login(username string, password string) error {

	c.credentials = credentials{
		Username: username,
		Password: password,
	}

	auth := struct {
		credentials `json:"auth"`
	}{c.credentials}

	req, err := c.newRequest("POST", "auth", auth)
	if err != nil {
		return err
	}

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}

	c.token = resp.Cookies()[0].Value
	return nil
}

// addOptions adds the parameters in opt as URL query parameters to s.  opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

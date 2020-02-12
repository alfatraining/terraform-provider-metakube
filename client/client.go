package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

const (
	defaultBaseURL = "https://metakube.syseleven.de"
)

// Client is a metkube api client.
type Client struct {
	client  *http.Client
	BaseURL *url.URL

	// Services
	Datacenters *DatacentersService
	Projects    *ProjectsService
}

// CreateOpt represent api clients construction option.
type CreateOpt func() *http.Client

type tokenSource struct {
	accessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.accessToken,
	}
	return token, nil
}

// WithBearerToken used for api client with Bearer Authentication.
func WithBearerToken(token string) CreateOpt {
	return func() *http.Client {
		return oauth2.NewClient(context.Background(), &tokenSource{token})
	}
}

// WithDefault used to create api client with default http client.
func WithDefault() CreateOpt {
	return func() *http.Client {
		return http.DefaultClient
	}
}

// New returns new default metakube api client.
func New() *Client {
	return NewClient(WithDefault())
}

// NewClient returns new metakube api client.
func NewClient(opt CreateOpt) *Client {
	client := &Client{
		client: opt(),
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	client.BaseURL = baseURL

	client.Datacenters = &DatacentersService{client}
	client.Projects = &ProjectsService{client}

	return client
}

// NewRequest returns new request to api configured in client.
func (c *Client) NewRequest(method, path string, payload interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	if payload != nil {
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, err
		}

	}
	return http.NewRequest(method, u.String(), buf)
}

// Do performs a request.
func (c *Client) Do(ctx context.Context, req *http.Request, out interface{}) error {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		// TODO: error cases
		return fmt.Errorf("non-ok status returned: %v", resp.Status)
	}

	if out != nil {
		return json.NewDecoder(resp.Body).Decode(out)
	}
	return nil
}

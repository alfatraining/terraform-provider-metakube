package gometakube

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

const (
	defaultBaseURL = "https://metakube.syseleven.de"
)

// Client is a metkube api client.
type Client struct {
	client  *http.Client
	BaseURL *url.URL

	// retry patch request on conflict status node 409.
	retriesOnConflict     uint
	retryOnConflictPeriod time.Duration

	// Services
	Datacenters     *DatacentersService
	Projects        *ProjectsService
	Clusters        *ClustersService
	NodeDeployments *NodeDeploymentsService
	Openstack       *OpenstackService
	SSHKeys         *SSHKeysService
}

// An ErrorMessage details the error caused by an API request.
type ErrorMessage struct {
	Code    int      `json:"code"`
	Details []string `json:"details"`
	Message string   `json:"message"`
}

// An ErrorResponse reports the error caused by an API request.
type ErrorResponse struct {
	// HTTP response of the error request.
	Response *http.Response `json:"-"`

	// API response message explaining the error.
	ErrorMessage *ErrorMessage `json:"error"`

	// Raw message that failed to be parsed into ErrorMessage.
	RawErrorMessage string `json:"-"`
}

func (r *ErrorResponse) Error() string {
	msg := r.RawErrorMessage
	if r.ErrorMessage != nil {
		msg = fmt.Sprintf("%+v", r.ErrorMessage)
	}
	return fmt.Sprintf("%v %v: %d %s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, msg)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, errorResponse)
		if err != nil {
			errorResponse.RawErrorMessage = string(data)
		}
	}
	return errorResponse
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
		client:                opt(),
		retriesOnConflict:     3,
		retryOnConflictPeriod: 5 * time.Second,
	}
	baseURL, _ := url.Parse(defaultBaseURL)
	client.BaseURL = baseURL

	client.Datacenters = &DatacentersService{client}
	client.Projects = &ProjectsService{client}
	client.Clusters = &ClustersService{client}
	client.NodeDeployments = &NodeDeploymentsService{client}
	client.Openstack = &OpenstackService{client}
	client.SSHKeys = &SSHKeysService{client}

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
func (c *Client) Do(ctx context.Context, req *http.Request, out interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if verr := resp.Body.Close(); err == nil {
			err = verr
		}
	}()

	err = checkResponse(resp)
	if err != nil {
		return resp, err
	}

	if out != nil {
		err = json.NewDecoder(resp.Body).Decode(out)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}

func (c *Client) resourceList(ctx context.Context, path string, ret interface{}) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, ret)
}

func (c *Client) resourceDelete(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, nil)
}

func (c *Client) resourceCreate(ctx context.Context, path string, v, ret interface{}) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPost, path, v)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, ret)
}

func (c *Client) resourceGet(ctx context.Context, path string, ret interface{}) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, &ret)
}

func (c *Client) resourcePut(ctx context.Context, path string, put, ret interface{}) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPut, path, put)
	if err != nil {
		return nil, err
	}
	return c.Do(ctx, req, ret)
}

func (c *Client) resourcePatch(ctx context.Context, path string, patch, ret interface{}) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPatch, path, patch)
	if err != nil {
		return nil, err
	}
	// TODO(furkhat): move retries out.
	ticker := time.NewTicker(c.retryOnConflictPeriod)
	defer ticker.Stop()
	resp, err := c.Do(ctx, req, &ret)
	for i := uint(0); i < c.retriesOnConflict; i++ {
		select {
		case <-ticker.C:
			if resp != nil && (resp.StatusCode < 400 || resp.StatusCode > 499) {
				break
			}
			// TODO(furkhat): add warning log.
			resp, err = c.Do(ctx, req, &ret)
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return resp, err
}

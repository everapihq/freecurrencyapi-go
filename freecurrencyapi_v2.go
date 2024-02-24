package freecurrencyapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

var (
	ErrInvalidStatusCode = errors.New("invalid status code")
	ErrUnauthorized      = errors.New("unauthorized")
)

type Client interface {
	Latest(ctx context.Context, request LatestRequest) (LatestResponse, error)
	Currencies(ctx context.Context, request CurrenciesRequest) (CurrenciesResponse, error)
	Historical(ctx context.Context, request HistoricalRequest) (HistoricalResponse, error)
	Status(ctx context.Context) (StatusResponse, error)
}

type v2Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

func (v v2Client) Latest(ctx context.Context, request LatestRequest) (LatestResponse, error) {
	res, err := v.doRequest(ctx, http.MethodGet, "latest", nil, request.toParams())
	if err != nil {
		return LatestResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return LatestResponse{}, ErrUnauthorized
		}
		return LatestResponse{}, ErrInvalidStatusCode
	}

	var decodedResponse latestResponse
	err = json.NewDecoder(res.Body).Decode(&decodedResponse)
	if err != nil {
		return LatestResponse{}, err
	}

	return decodedResponse.toLatestResponse(), nil
}

func (v v2Client) Currencies(ctx context.Context, request CurrenciesRequest) (CurrenciesResponse, error) {
	res, err := v.doRequest(ctx, http.MethodGet, "currencies", nil, request.toParams())
	if err != nil {
		return CurrenciesResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return CurrenciesResponse{}, ErrUnauthorized
		}
		return CurrenciesResponse{}, ErrInvalidStatusCode
	}

	var decodedResponse currenciesResponse
	err = json.NewDecoder(res.Body).Decode(&decodedResponse)
	if err != nil {
		return CurrenciesResponse{}, err
	}

	return decodedResponse.toCurrenciesResponse(), nil
}

func (v v2Client) Historical(ctx context.Context, request HistoricalRequest) (HistoricalResponse, error) {
	res, err := v.doRequest(ctx, http.MethodGet, "historical", nil, request.toParams())
	if err != nil {
		return HistoricalResponse{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return HistoricalResponse{}, ErrUnauthorized
		}
		return HistoricalResponse{}, ErrInvalidStatusCode
	}

	var decodedResponse historicalResponse
	err = json.NewDecoder(res.Body).Decode(&decodedResponse)
	if err != nil {
		return HistoricalResponse{}, err
	}

	return decodedResponse.toHistoricalResponse(), nil
}

func (v v2Client) Status(ctx context.Context) (StatusResponse, error) {
	res, err := v.doRequest(ctx, http.MethodGet, "status", nil, nil)
	if err != nil {
		return StatusResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusUnauthorized {
			return StatusResponse{}, ErrUnauthorized
		}
		return StatusResponse{}, ErrInvalidStatusCode
	}

	var decodedResponse statusResponse
	err = json.NewDecoder(res.Body).Decode(&decodedResponse)
	if err != nil {
		return StatusResponse{}, err
	}

	return decodedResponse.toStatusResponse(), nil
}

func (v v2Client) doRequest(ctx context.Context, method string, path string, body io.Reader, query map[string]string) (*http.Response, error) {
	url := v.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", v.apiKey)
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	resp, err := v.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ClientOption struct {
	HTTPClient *http.Client
	BaseURL    string
}

func (c *ClientOption) WithHTTPClient(httpClient *http.Client) *ClientOption {
	c.HTTPClient = httpClient
	return c
}

func (c *ClientOption) WithBaseURL(baseURL string) *ClientOption {
	c.BaseURL = baseURL
	return c
}

func Options() *ClientOption {
	return &ClientOption{}
}

func NewClient(apiKey string, opts ...*ClientOption) Client {
	opt := Options()
	if len(opts) > 0 && opts[0] != nil {
		opt = opts[0]
	}

	httpClient := opt.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL := opt.BaseURL
	if baseURL == "" {
		baseURL = BaseUrl
	}

	return &v2Client{
		apiKey:     apiKey,
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

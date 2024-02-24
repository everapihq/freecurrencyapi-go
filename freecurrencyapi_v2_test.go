package freecurrencyapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var sampleApiKey = "sample-api-key"

func baseServerHandler(method string, path string, apiKey string, response string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != path {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if r.Header.Get("apikey") != apiKey {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
}

func Test_FreeCurrencyAPIv2_Status(t *testing.T) {
	mockResponse := `{"quotas":{"month":{"total":300,"used":71,"remaining":229}}}`
	defaultHandler := baseServerHandler(http.MethodGet, "/v1/status", sampleApiKey, mockResponse)

	cases := []struct {
		name          string
		ctx           context.Context
		apiKey        string
		assert        func(t *testing.T, response StatusResponse, err error)
		serverHandler http.HandlerFunc
		serverKill    bool
	}{
		{
			name:   "should fail if no api key is provided",
			apiKey: "",
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, StatusResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if unauthorized",
			apiKey: "invalid-api-key",
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, StatusResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if invalid status code",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrInvalidStatusCode, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, StatusResponse{}, response, "expected response to be empty")
			},
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
		},
		{
			name:   "should return error if client.Do fails",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "connect: connection refused", "expected error to contain connect: connection refused")

				assert.Equal(t, StatusResponse{}, response, "expected response to be empty")
			},
			serverKill: true,
		},
		{
			name:   "should return error if cannot decode response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, StatusResponse{}, response, "expected response to be empty")
			},
			serverHandler: baseServerHandler(http.MethodGet, "/v1/status", sampleApiKey, `{"quotas":{"month":{"total":300,"used":71,"remaining":229}`),
		},
		{
			name:   "should return error if context is canceled",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "context canceled", "expected error to contain context canceled")

				assert.Equal(t, StatusResponse{}, response, "expected response to be empty")
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
		},
		{
			name:   "should return status response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response StatusResponse, err error) {
				assert.Nil(t, err, "expected error to be nil")
				assert.Equal(t, StatusResponse{
					Month: StatusItem{
						Total:     300,
						Used:      71,
						Remaining: 229,
					},
				}, response, "expected response to be equal")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := defaultHandler
			if tc.serverHandler != nil {
				handler = tc.serverHandler
			}
			testServer := httptest.NewServer(handler)
			if tc.serverKill {
				testServer.Close()
			}
			cli := NewClient(tc.apiKey, Options().WithHTTPClient(testServer.Client()).WithBaseURL(testServer.URL+"/v1/"))
			ctx := context.Background()
			if tc.ctx != nil {
				ctx = tc.ctx
			}
			response, err := cli.Status(ctx)
			tc.assert(t, response, err)
			t.Cleanup(testServer.Close)
		})
	}
}

func Test_FreeCurrencyAPIv2_Currencies(t *testing.T) {
	mockResponse := `{"data":{"AED":{"symbol":"AED","name":"United Arab Emirates Dirham","symbol_native":"د.إ","decimal_digits":2,"rounding":0,"code":"AED","name_plural":"UAE dirhams"},"AFN":{"symbol":"Af","name":"Afghan Afghani","symbol_native":"؋","decimal_digits":0,"rounding":0,"code":"AFN","name_plural":"Afghan Afghanis"}}}`

	defaultHandler := baseServerHandler(http.MethodGet, "/v1/currencies", sampleApiKey, mockResponse)

	cases := []struct {
		name          string
		ctx           context.Context
		apiKey        string
		assert        func(t *testing.T, response CurrenciesResponse, err error)
		serverHandler http.HandlerFunc
		serverKill    bool
	}{
		{
			name:   "should fail if no api key is provided",
			apiKey: "",
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, CurrenciesResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if unauthorized",
			apiKey: "invalid-api-key",
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, CurrenciesResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if invalid status code",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrInvalidStatusCode, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, CurrenciesResponse{}, response, "expected response to be empty")
			},
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
		},
		{
			name:   "should return error if client.Do fails",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "connect: connection refused", "expected error to contain connect: connection refused")

				assert.Equal(t, CurrenciesResponse{}, response, "expected response to be empty")
			},
			serverKill: true,
		},
		{
			name:   "should return error if cannot decode response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, CurrenciesResponse{}, response, "expected response to be empty")
			},
			serverHandler: baseServerHandler(http.MethodGet, "/v1/currencies", sampleApiKey, `{"data":{"AED":{"symbol":"AED","name":"United Arab Emirates Dirham","symbol_native":"د.إ","decimal_digits":2,"rounding":0,"code":"AED","name_plural":"UAE dirhams"}`),
		},
		{
			name:   "should return error if context is canceled",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "context canceled", "expected error to contain context canceled")

				assert.Equal(t, CurrenciesResponse{}, response, "expected response to be empty")
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
		},
		{
			name:   "should return currencies response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response CurrenciesResponse, err error) {
				assert.Nil(t, err, "expected error to be nil")
				assert.Equal(t, CurrenciesResponse{
					Currencies: map[string]CurrencyItem{
						"AED": {
							Symbol:        "AED",
							Name:          "United Arab Emirates Dirham",
							SymbolNative:  "د.إ",
							DecimalDigits: 2,
							Rounding:      0,
							Code:          "AED",
							NamePlural:    "UAE dirhams",
						},
						"AFN": {
							Symbol:        "Af",
							Name:          "Afghan Afghani",
							SymbolNative:  "؋",
							DecimalDigits: 0,
							Rounding:      0,
							Code:          "AFN",
							NamePlural:    "Afghan Afghanis",
						},
					},
				}, response, "expected response to be equal")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := defaultHandler
			if tc.serverHandler != nil {
				handler = tc.serverHandler
			}
			testServer := httptest.NewServer(handler)
			if tc.serverKill {
				testServer.Close()
			}
			cli := NewClient(tc.apiKey, Options().WithHTTPClient(testServer.Client()).WithBaseURL(testServer.URL+"/v1/"))
			ctx := context.Background()
			if tc.ctx != nil {
				ctx = tc.ctx
			}
			response, err := cli.Currencies(ctx, CurrenciesRequest{})
			tc.assert(t, response, err)
			t.Cleanup(testServer.Close)
		})
	}
}

func Test_FreeCurrencyAPIv2_Latest(t *testing.T) {
	mockResponse := `{"data":{"AED":3.67306,"AFN":91.80254,"ALL":108.22904,"AMD":480.41659}}`

	defaultHandler := baseServerHandler(http.MethodGet, "/v1/latest", sampleApiKey, mockResponse)
	cases := []struct {
		name          string
		ctx           context.Context
		apiKey        string
		assert        func(t *testing.T, response LatestResponse, err error)
		serverHandler http.HandlerFunc
		serverKill    bool
	}{
		{
			name:   "should fail if no api key is provided",
			apiKey: "",
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, LatestResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if unauthorized",
			apiKey: "invalid-api-key",
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, LatestResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if invalid status code",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrInvalidStatusCode, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, LatestResponse{}, response, "expected response to be empty")
			},
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
		},
		{
			name:   "should return error if client.Do fails",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "connect: connection refused", "expected error to contain connect: connection refused")

				assert.Equal(t, LatestResponse{}, response, "expected response to be empty")
			},
			serverKill: true,
		},
		{
			name:   "should return error if cannot decode response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, LatestResponse{}, response, "expected response to be empty")
			},
			serverHandler: baseServerHandler(http.MethodGet, "/v1/latest", sampleApiKey, `{"data":{"AED":3.67306`),
		},
		{
			name:   "should return error if context is canceled",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "context canceled", "expected error to contain context canceled")

				assert.Equal(t, LatestResponse{}, response, "expected response to be empty")
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
		},
		{
			name:   "should return latest response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response LatestResponse, err error) {
				assert.Nil(t, err, "expected error to be nil")
				assert.Equal(t, LatestResponse{
					Rates: map[string]float64{
						"AED": 3.67306,
						"AFN": 91.80254,
						"ALL": 108.22904,
						"AMD": 480.41659,
					},
				}, response, "expected response to be equal")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := defaultHandler
			if tc.serverHandler != nil {
				handler = tc.serverHandler
			}
			testServer := httptest.NewServer(handler)
			if tc.serverKill {
				testServer.Close()
			}
			cli := NewClient(tc.apiKey, Options().WithHTTPClient(testServer.Client()).WithBaseURL(testServer.URL+"/v1/"))
			ctx := context.Background()
			if tc.ctx != nil {
				ctx = tc.ctx
			}
			response, err := cli.Latest(ctx, LatestRequest{})
			tc.assert(t, response, err)
			t.Cleanup(testServer.Close)
		})
	}
}

func Test_FreeCurrencyAPIv2_Historical(t *testing.T) {
	mockResponse := `{"data":{"2022-01-01":{"AED":3.67306,"AFN":91.80254,"ALL":108.22904,"AMD":480.41659}}}`

	defaultHandler := baseServerHandler(http.MethodGet, "/v1/historical", sampleApiKey, mockResponse)
	cases := []struct {
		name          string
		ctx           context.Context
		apiKey        string
		assert        func(t *testing.T, response HistoricalResponse, err error)
		serverHandler http.HandlerFunc
		serverKill    bool
	}{
		{
			name:   "should fail if no api key is provided",
			apiKey: "",
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, HistoricalResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if unauthorized",
			apiKey: "invalid-api-key",
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrUnauthorized, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, HistoricalResponse{}, response, "expected response to be empty")
			},
		},
		{
			name:   "should fail if invalid status code",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, ErrInvalidStatusCode, err, "expected error to be ErrInvalidStatusCode")

				assert.Equal(t, HistoricalResponse{}, response, "expected response to be empty")
			},
			serverHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			}),
		},
		{
			name:   "should return error if client.Do fails",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "connect: connection refused", "expected error to contain connect: connection refused")

				assert.Equal(t, HistoricalResponse{}, response, "expected response to be empty")
			},
			serverKill: true,
		},
		{
			name:   "should return error if cannot decode response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Equal(t, HistoricalResponse{}, response, "expected response to be empty")
			},
			serverHandler: baseServerHandler(http.MethodGet, "/v1/historical", sampleApiKey, `{"data":{"2022-01-01":{"AED":3.67306`),
		},
		{
			name:   "should return error if context is canceled",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.NotNil(t, err, "expected error to be not nil")
				assert.Contains(t, err.Error(), "context canceled", "expected error to contain context canceled")

				assert.Equal(t, HistoricalResponse{}, response, "expected response to be empty")
			},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
		},
		{
			name:   "should return historical response",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, response HistoricalResponse, err error) {
				assert.Nil(t, err, "expected error to be nil")
				assert.Equal(t, HistoricalResponse{
					Rates: map[string]float64{
						"AED": 3.67306,
						"AFN": 91.80254,
						"ALL": 108.22904,
						"AMD": 480.41659,
					},
				}, response, "expected response to be equal")
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			handler := defaultHandler
			if tc.serverHandler != nil {
				handler = tc.serverHandler
			}
			testServer := httptest.NewServer(handler)
			if tc.serverKill {
				testServer.Close()
			}
			cli := NewClient(tc.apiKey, Options().WithHTTPClient(testServer.Client()).WithBaseURL(testServer.URL+"/v1/"))
			ctx := context.Background()
			if tc.ctx != nil {
				ctx = tc.ctx
			}
			response, err := cli.Historical(ctx, HistoricalRequest{})
			tc.assert(t, response, err)
			t.Cleanup(testServer.Close)
		})
	}
}

func Test_FreeCurrencyAPIv2_NewClient(t *testing.T) {
	cases := []struct {
		name    string
		apiKey  string
		assert  func(t *testing.T, client Client)
		options *ClientOption
	}{
		{
			name:   "should return client with default options",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, client Client) {
				assert.Equal(t, sampleApiKey, client.(*v2Client).apiKey, "expected api key to be equal")
				assert.Equal(t, BaseUrl, client.(*v2Client).baseURL, "expected base url to be equal")
				assert.Equal(t, http.DefaultClient, client.(*v2Client).httpClient, "expected http client to be equal")
			},
		},
		{
			name:   "should return client with custom options",
			apiKey: sampleApiKey,
			assert: func(t *testing.T, client Client) {
				assert.Equal(t, sampleApiKey, client.(*v2Client).apiKey, "expected api key to be equal")
				assert.Equal(t, "https://api.example.com/v1/", client.(*v2Client).baseURL, "expected base url to be equal")
				assert.Equal(t, &http.Client{Timeout: 10}, client.(*v2Client).httpClient, "expected http client to be equal")
			},
			options: Options().WithBaseURL("https://api.example.com/v1/").WithHTTPClient(&http.Client{Timeout: 10}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := NewClient(tc.apiKey, tc.options)
			tc.assert(t, client)
		})
	}
}

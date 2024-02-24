package freecurrencyapi

type latestResponse struct {
	Data map[string]float64 `json:"data"`
}

type LatestResponse struct {
	Rates map[string]float64
}

func (r latestResponse) toLatestResponse() LatestResponse {
	return LatestResponse{
		Rates: r.Data,
	}
}

type currencyResponseItem struct {
	Symbol        string `json:"symbol"`
	Name          string `json:"name"`
	SymbolNative  string `json:"symbol_native"`
	DecimalDigits int    `json:"decimal_digits"`
	Rounding      int    `json:"rounding"`
	Code          string `json:"code"`
	NamePlural    string `json:"name_plural"`
}

type currenciesResponse struct {
	Data map[string]currencyResponseItem `json:"data"`
}

type CurrencyItem struct {
	Symbol        string
	Name          string
	SymbolNative  string
	DecimalDigits int
	Rounding      int
	Code          string
	NamePlural    string
}

type CurrenciesResponse struct {
	Currencies map[string]CurrencyItem
}

func (r currenciesResponse) toCurrenciesResponse() CurrenciesResponse {
	currencies := make(map[string]CurrencyItem, len(r.Data))
	for key, value := range r.Data {
		currencies[key] = CurrencyItem{
			Symbol:        value.Symbol,
			Name:          value.Name,
			SymbolNative:  value.SymbolNative,
			DecimalDigits: value.DecimalDigits,
			Rounding:      value.Rounding,
			Code:          value.Code,
			NamePlural:    value.NamePlural,
		}
	}

	return CurrenciesResponse{
		Currencies: currencies,
	}
}

type historicalResponse struct {
	Data map[string]map[string]float64 `json:"data"`
}

type HistoricalResponse struct {
	Rates map[string]float64
}

func (r historicalResponse) toHistoricalResponse() HistoricalResponse {
	firstKey := ""
	for key := range r.Data {
		firstKey = key
		break
	}

	return HistoricalResponse{
		Rates: r.Data[firstKey],
	}
}

type statusResponse struct {
	Quotas struct {
		Month struct {
			Total     int `json:"total"`
			Used      int `json:"used"`
			Remaining int `json:"remaining"`
		} `json:"month"`
	} `json:"quotas"`
}

type StatusItem struct {
	Total     int
	Used      int
	Remaining int
}

type StatusResponse struct {
	Month StatusItem
}

func (r statusResponse) toStatusResponse() StatusResponse {
	return StatusResponse{
		Month: StatusItem{
			Total:     r.Quotas.Month.Total,
			Used:      r.Quotas.Month.Used,
			Remaining: r.Quotas.Month.Remaining,
		},
	}
}

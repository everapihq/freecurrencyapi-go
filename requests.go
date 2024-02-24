package freecurrencyapi

import (
	"strings"
	"time"
)

type LatestRequest struct {
	BaseCurrency string
	Currencies   []string
}

func (r LatestRequest) toParams() map[string]string {
	p := map[string]string{}

	if r.BaseCurrency != "" {
		p["base_currency"] = r.BaseCurrency
	}

	if len(r.Currencies) > 0 {
		p["currencies"] = strings.Join(r.Currencies, ",")
	}

	return p
}

type CurrenciesRequest struct {
	Currencies []string
}

func (r CurrenciesRequest) toParams() map[string]string {
	p := map[string]string{}

	if len(r.Currencies) > 0 {
		p["currencies"] = strings.Join(r.Currencies, ",")
	}

	return p
}

type HistoricalRequest struct {
	Date         time.Time
	BaseCurrency string
	Currencies   []string
}

func (r HistoricalRequest) toParams() map[string]string {
	p := map[string]string{
		"date": r.Date.Format("2006-01-02"),
	}

	if r.BaseCurrency != "" {
		p["base_currency"] = r.BaseCurrency
	}

	if len(r.Currencies) > 0 {
		p["currencies"] = strings.Join(r.Currencies, ",")
	}

	return p
}

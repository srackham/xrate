package xrates

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/srackham/xrate/internal/cache"
	"github.com/srackham/xrate/internal/config"
	"github.com/srackham/xrate/internal/helpers"
)

// Cache data types.
type Rates map[string]float64        // Key = currency symbol; value = value in USD.
type RatesCacheData map[string]Rates // Key = date string "YYYY-MM-DD".

// Define a custom type for the HTTP Get function
type httpGet func(url string) (*http.Response, error)

type ExchangeRates struct {
	ConfigFile string
	cache.Cache[RatesCacheData]
	HttpGet httpGet // Set to mock http.Get function when testing
}

func New() ExchangeRates {
	data := make(RatesCacheData)
	result := ExchangeRates{
		path.Join(helpers.GetXDGConfigDir(), "xrate", "config.yaml"),
		*cache.New(&data),
		http.Get,
	}
	result.CacheFile = path.Join(helpers.GetXDGCacheDir(), "xrate", "exchange-rates.json")
	return result
}

// getRates fetches a list of currency exchange rates against the USD
// TODO getRates should be an IXRatesAPI interface cf. prices.IPriceAPI.
func (x *ExchangeRates) getRates() (Rates, error) {
	if helpers.GithubActions() {
		// getRates() requires HTTP access and should never execute from Github Actions.
		mockRates := Rates{"usd": 1.0, "nzd": 1.5, "aud": 1.6}
		return mockRates, nil
	}
	rates := make(Rates)
	conf, err := config.LoadConfig(x.ConfigFile)
	if err != nil {
		return rates, err
	}

	url := "https://openexchangerates.org/api/latest.json?app_id=" + conf.XratesAppId
	resp, err := x.HttpGet(url)
	if err != nil {
		return rates, fmt.Errorf("exchange rate request: %s: %s", url, err.Error())
	}
	defer resp.Body.Close()

	// See https://www.sohamkamani.com/golang/json/#decoding-json-to-maps---unstructured-data
	var m map[string]any
	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return rates, fmt.Errorf("exchange rate decode: %s", err.Error())
	}
	_, exists := m["rates"]
	if !exists {
		jsonData, err := json.Marshal(m)
		if err != nil {
			return rates, fmt.Errorf("invalid exchange rate response: %s: %v", url, m)
		}
		return rates, fmt.Errorf("invalid exchange rate response: %s: %s", url, string(jsonData))
	}
	m = m["rates"].(map[string]any)
	for k, v := range m {
		rates[strings.ToUpper(k)] = v.(float64)
	}
	return rates, nil
}

// GetCachedRate returns the amount of `currency` that $1 USD would buy at today's rates.
// `symbol` is a currency symbol.
// If `force` is `true` then then today's rates are unconditionally fetched and the cache updated.
// TODO tests
func (x *ExchangeRates) GetCachedRate(currency string, force bool) (float64, error) {
	if currency == "" {
		return 0.0, fmt.Errorf("no currency specified")
	}
	if currency == "USD" {
		return 1.00, nil
	}
	var rate float64
	var ok bool
	today := helpers.TodaysDate()
	if rate, ok = (*x.CacheData)[today][strings.ToUpper(currency)]; !ok || force {
		rates, err := x.getRates()
		if err != nil {
			return 0.0, err
		}
		x.CacheData = &RatesCacheData{today: rates}
		if rate, ok = (*x.CacheData)[today][strings.ToUpper(currency)]; !ok {
			return 0.0, fmt.Errorf("unknown currency: %s", currency)
		}
	}
	return rate, nil
}

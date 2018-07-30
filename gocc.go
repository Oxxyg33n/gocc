package gocc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/juju/errors"
)

const baseURL = "https://openexchangerates.org/api/"

type exchangeRatesResponse struct {
	Timestamp   int                `json:"timestamp"`
	Rates       map[string]float64 `json:"rates"`
	Error       bool               `json:"error"`
	Description string             `json:"description"`
}

func (e exchangeRatesResponse) getTimeStamp() int {
	return e.Timestamp
}

func (e exchangeRatesResponse) getRates() map[string]float64 {
	return e.Rates
}

func (e exchangeRatesResponse) getError() bool {
	return e.Error
}

func (e exchangeRatesResponse) getDescription() string {
	return e.Description
}

// AvailableCurrencies returns list of available currencies
func AvailableCurrencies() (map[string]string, error) {
	url := baseURL + "currencies.json"
	resp, err := newRequest().Get(url)
	if err != nil {
		return nil, errors.Annotatef(err, "getting response from %s url failed", url)
	}

	cMap := map[string]string{}
	err = json.Unmarshal(resp, &cMap)
	if err != nil {
		log.Println(err)
		return nil, errors.Annotatef(err, "unmarshalling json response %s failed", resp)
	}

	return cMap, nil
}

// ShowExchangeRate returns a map with exchange rates for currency baseCurrency
// ex: ShowExchangeRate("USD", false)
// Additional param showAlt can be either true or false (bool)
// If it's true then it will retrieve the list with alternative currencies
func ShowExchangeRate(baseCurrency string, showAlt bool) (int, map[string]float64, error) {
	appID := os.Getenv("OER_APP_ID")

	var url string
	if !showAlt {
		url = fmt.Sprintf(baseURL+"latest.json?app_id=%s&base=%s", appID, baseCurrency)
	} else {
		url = fmt.Sprintf(baseURL+"latest.json?app_id=%s&base=%s&show_alternative=1", appID, baseCurrency)
	}
	resp, err := newRequest().Get(url)
	if err != nil {
		return 0, map[string]float64{}, errors.Annotatef(err, "getting response from %s url failed", url)
	}

	rates := exchangeRatesResponse{}
	err = json.Unmarshal(resp, &rates)
	if err != nil {
		log.Println(err)
		return 0, map[string]float64{}, errors.Annotatef(err, "unmarshalling json response %s failed", resp)
	}

	if rates.getError() {
		return 0, nil, errors.New(rates.getDescription())
	}

	return rates.getTimeStamp(), rates.getRates(), nil
}

// ConvertCurrency converts <amount> USD dollars to EURO/GBP/CAD/YEN/BTC
// ex: ConvertCurrency(100)
// nolint: gocyclo
func ConvertCurrency(baseCurrency string, amount float64) (string, error) {
	// Show alternative currencies?
	showAlt := true

	_, rates, err := ShowExchangeRate(baseCurrency, showAlt)
	if err != nil {
		return "", err
	}

	var currencies [6]float64
	for k := range rates {
		switch k {
		case "EUR":
			currencies[0] = amount * rates["EUR"]
		case "GBP":
			currencies[1] = amount * rates["GBP"]
		case "CAD":
			currencies[2] = amount * rates["CAD"]
		case "JPY":
			currencies[3] = amount * rates["JPY"]
		case "BTC":
			currencies[4] = amount * rates["BTC"]
		case "LTC":
			currencies[5] = amount * rates["LTC"]
		}
	}

	var curString string
	curSymbols := []string{"EUR", "GBP", "CAD", "YEN", "BTC", "LTC"}
	for z := range curSymbols {
		if curSymbols[z] == "LTC" && !showAlt {
			break
		}
		curString += curSymbols[z] + ": " + strconv.FormatFloat(currencies[z], 'f', 6, 64) + "\n"
	}

	return curString, nil
}

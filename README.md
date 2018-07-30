# GOCC - GO Currency Converter & Wrapper for OXR

Wrapper & Currency Converter written in Golang for [Open Exchange Rates](https://openexchangerates.org/).
It provides simple to use methods to make calls to the API and get responses

## Installation
Use go get:

`$ go get github.com/oxxyg33n/gocc`

or just clone the repository in desired location with HTTPS:

`git clone https://github.com/Oxxyg33n/gocc.git`

## How to use
First, import the package:

`import "github.com/oxxyg33n/gocc"`

Now, to show a list of all available currencies

`
		cMap, _ := gocc.AvailableCurrencies()
`

you will have a `map` of type `[string]string` that you can do anything with, for example iterate through it and
print all the key-value pairs


    for k := range cMap {
        fmt.Println(k + ": " + cMap[k]) // Print out key-value pairs
    }
    
Now, to show exchange rates for specific currency (available only to developer, enterprise and unlimited plans) or
for USD (available for free).

Let's say we have variable `baseCur` which is string that contains currency symbols(3-letter ISO code, USD/EUR, etc..)

Then we can do something like this:


    cMap, _ := gocc.AvailableCurrencies()   // Get list of all vailable currencies

		for k := range cMap {   // Iterate through the map
			if baseCur == k {   // If the currency specified is found then show exchange rates
				timestamp, rates, _ := gocc.ShowExchangeRate(baseCur, false) // Timestamp is UNIX timestamp
				
				time := time.Unix(int64(timestamp), 0)
				fmt.Println(time.String())
				
				var floatToString string
				for rate := range rates {
				    floatToString = strconv.FormatFloat(rates[rate], 'f', 6, 64)
				    fmt.Println(rate + ": " + floatToString) // Print out exchange rate for currency relative to baseCur
				}
			} else { // If baseCur is not the code from list then continue
				continue
			}
		}

Now, to convert USD dollars or other base currency (available only to developer, enterprise and unlimited plans)
to EUR/GBP/CAD/YET/BTC/LTC and print out results

        convertedValue, _ := gocc.ConvertCurrency("USD", 100.50)
        
## Credits
Thanks to [DiSiquera and his GoCurrency](https://github.com/DiSiqueira/GoCurrency) for request.go :)
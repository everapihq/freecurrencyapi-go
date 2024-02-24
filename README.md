<p align="center">
<img src="https://app.freecurrencyapi.com/img/logo/freecurrencyapi.png" width="300"/>
</p>

# freecurrencyapi-go: Golang Currency Converter

This package is a Golang wrapper for [freecurrencyapi.com](https://freecurrencyapi.com) that aims to make the usage of freecurrencyapi.com API as easy as possible in your project.

## Usage

Initialize the API with your API Key (get one for free at [freecurrencyapi.com](https://freecurrencyapi.com)):

```go
client := freecurrencyapi.NewClient("YOUR-API-KEY")
```

Afterwards you can make calls to the API like this:

### Status Endpoint

```go
client.Status(context.TODO())
```

### Currencies Endpoint

```go
response, err := client.Currencies(context.TODO(), freeCurrencyApi.CurrenciesRequest{
	Currencies: []string{"EUR", "GBP"},
})
```

### Latest Endpoint

```go
response, err := client.Latest(context.TODO(), freeCurrencyApi.LatestRequest{
    BaseCurrency: "USD",
    Currencies: []string{"EUR", "GBP"},
})
```

### Historical Endpoint

```go
response, err := client.Historical(context.TODO(), freeCurrencyApi.HistoricalRequest{
    BaseCurrency: "USD",
    Currencies: []string{"EUR", "GBP"},
    Date: "2021-01-01",
})
```

Find out more about our endpoints, parameters and response data structure in the [docs](https://freecurrencyapi.com/docs)

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.

[docs]: https://freecurrencyapi.com/docs
[freecurrencyapi.com]: https://freecurrencyapi.com

# Fiat Currency Exchange Rates

A simple CLI command to fetch the amount of a fiat currency that would be exchanged for one US dollar.

    Usage: xrate SYMBOL

`SYMBOL` is the currency's ticker symbol e.g. `NZD`, `AUD`, `EUR`:

```
$ xrate NZD
1.76
```

## Implementation

- Fetches exchange rates using the [Open Exchange Rates](https://openexchangerates.org/) Web API.
- The exchange rates are cached and only fetched once per day.

## Installation

Clone the [xrate Github repo](https://github.com/srackham/xrate) and compile and install using the `go` command:

```
$ git clone https://github.com/srackham/xrate.git
$ cd xrate
$ go install
```

## Configuration
You will need to obtain an App ID from [Open Exchange Rates](https://openexchangerates.org/) and create the `$HOME/.config/xrate/config.yaml` configuration file and put the App ID in it. Here's the file format:

```yaml
xrates-appid: <your open exchange rates App ID goes here>
```

The exchange rates are cached in `$HOME/.cache/xrate/exchange-rates.json`.

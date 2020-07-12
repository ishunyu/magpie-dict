# magpie-dict
An online Chinese-English dictionary based on the subtitle translations of shows by the The Magpie Bridge Brigade.

### Requirements
- [Go](https://golang.org/)

## Getting Started
### Configure
Copy default config file and then replace properties suited to your environment.
```
cp config/config.json config/local.json
```

### Build
```
go build -o bin/ ./pkg/server
```

### Run
```
bin/server config/local.json
```

Starting the server for the first time may take a while as it's indexing everything. Subsequent restart should be faster.

If reindexing is necessary, simply delete `indexPath`.

## Development
```
go run pkg/server/*.go config/local.json
```
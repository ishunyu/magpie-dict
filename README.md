# Magpie Dictionary
An online dictionary based on the subtitle translations from the The Magpie Bridge Brigade. Built in Go.

### Requirements
- [Go](https://golang.org/)
- Data files

## Getting Started
### Configure
Copy the configuration template file and populate the properties
```
mkdir tmp && cp config/template.json tmp/config.json
```

### Build
```
go build -o bin/ ./pkg/server
```

### Run
```
bin/server tmp/config.json
```
Or in the background
```
nohup bin/server tmp/config.json > tmp/server.log &
```

- Starting the server for the first time may take a while as it's indexing everything. Subsequent restart should be faster.

- If reindexing is necessary, simply delete directory referenced by the `indexPath` configuration property.
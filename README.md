# Magpie Dictionary
An online dictionary based on the subtitle translations from the The Magpie Bridge Brigade. Built in Go.

The goal of this project is to provide subtitle translations by comparing translated text side-by-side. It supports a proprietary format that allows subtitles from different time slices to be compared, please see [DATA_FORMAT.md](DATA_FORMAT.md). This data file can be produced by scripts from the [CompareSBV](https://github.com/BrigadeMagpie/CompareSBV) project.

## Getting Started
### Requirements
- [Go](https://golang.org/)
- Data files

### Configuration
Copy the configuration template file and populate the properties appropriately.
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

- Starting the server for the first time may take a while as it needs to index everything.

- If reindexing is necessary, simply delete the `indexPath` directory.
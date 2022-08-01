# Deployment
## Navigate
```
cd <project_root>
```
## Build
```
go build -o bin/ ./pkg/server
```
## Start
```
bin/server config/local.json
```
or headless
```
nohup bin/server config/local.json > somefile.log &
```
## Stop
```
kill `ps -ef | grep bin/server | grep -v grep | awk '{print $2}'`
```
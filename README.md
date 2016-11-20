# makina
unknown 

## Run
```
go get -u ./...
go build

# without log file (console log)
./makina

# with log file
./makina -log=makina.log
```

## Daemon
```
use makina.service

$ sudo systemctl daemon-reload
$ sudo systemctl stop makina
$ sudo systemctl start makina
```

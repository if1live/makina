# makina
unknown 

## Install
```
# dependencies
go get github.com/ChimeraCoder/anaconda
go get github.com/tj/go-dropy
go get github.com/kardianos/osext
go get github.com/stretchr/testify/assert
go get github.com/xconstruct/go-pushbullet
```

## Run
```
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

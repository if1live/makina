package main

import (
	raven "github.com/getsentry/raven-go"
	"github.com/kardianos/osext"
)

func check(e error) {
	if e != nil {
		raven.CaptureErrorAndWait(e, nil)
		panic(e)
	}
}

func GetExecutablePath() string {
	path, err := osext.ExecutableFolder()
	check(err)
	return path
}

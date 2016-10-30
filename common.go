package main

import (
	"github.com/kardianos/osext"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetExecutablePath() string {
	path, err := osext.ExecutableFolder()
	check(err)
	return path
}

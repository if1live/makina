package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func SaveTweetJsonFile(t *anaconda.Tweet, filepath string) {
	b, err := json.Marshal(t)
	check(err)

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	f, err := os.Create(filepath)
	check(err)

	w := bufio.NewWriter(f)
	out.WriteTo(w)

	w.Flush()
}

func LoadTweetJsonFile(filepath string) anaconda.Tweet {
	// https://gist.github.com/border/775526
	file, e := ioutil.ReadFile(filepath)
	check(e)

	tweet := anaconda.Tweet{}
	json.Unmarshal(file, &tweet)
	return tweet
}

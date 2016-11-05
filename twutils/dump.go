package twutils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/ChimeraCoder/anaconda"
)

func LoadTweetJsonFile(filepath string) anaconda.Tweet {
	// https://gist.github.com/border/775526
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		panic(e)
	}

	tweet := anaconda.Tweet{}
	json.Unmarshal(file, &tweet)
	return tweet
}

func SaveTweetJsonFile(t *anaconda.Tweet, filepath string) {
	b, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	f, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	out.WriteTo(w)

	w.Flush()
}

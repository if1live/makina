package hitomi

import "github.com/if1live/makina/storages"

type Config struct {
	MyName   string
	Accessor storages.Accessor

	HaruExecutable string
	HaruHostName   string
	ShowLog        bool
}

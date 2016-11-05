package hitomi

import "github.com/if1live/makina/storages"
import "github.com/if1live/makina/senders"

type Config struct {
	MyName       string
	Accessor     storages.Accessor
	StatusSender *senders.Sender

	HaruExecutable string
	HaruHostName   string
	ShowLog        bool
}

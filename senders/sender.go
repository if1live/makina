package senders

type Message struct {
	Title string
	Body  string
}

type SendStrategy interface {
	Send(title, body string)
}

type Sender struct {
	strategy SendStrategy
	last     Message
}

func New(s SendStrategy) *Sender {
	return &Sender{
		strategy: s,
	}
}

func (s *Sender) Send(title, body string) {
	s.last = Message{title, body}
	s.strategy.Send(title, body)
}

func (s *Sender) SendTitleOnly(title string) {
	s.last = Message{title, ""}
	s.strategy.Send(title, "")
}

func (s *Sender) GetLastMessage() Message {
	return s.last
}

package senders

import pushbullet "github.com/xconstruct/go-pushbullet"

type PushbulletSendStrategy struct {
	pb  *pushbullet.Client
	dev *pushbullet.Device
}

func NewPushbullet(token string) SendStrategy {
	pb := pushbullet.New(token)
	devs, err := pb.Devices()
	if err != nil {
		panic(err)
	}

	dev := selectDevice(devs)
	return &PushbulletSendStrategy{
		pb:  pb,
		dev: dev,
	}
}

func selectDevice(devs []*pushbullet.Device) *pushbullet.Device {
	for _, dev := range devs {
		return dev
	}
	return nil
}

func (s *PushbulletSendStrategy) Send(title, body string) {
	err := s.pb.PushNote(s.dev.Iden, title, body)
	if err != nil {
		panic(err)
	}
}

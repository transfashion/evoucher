package qiscus

type Qiscus struct {
	BaseUrl string
	AppCode string
	Secret  string
	Sender  string
}

type Chat struct {
	RoomId string
	Number string
	Name   string
}

package qiscus

type Qiscus struct {
	Config *QiscusConfig
}

type QiscusConfig struct {
	BaseUrl string
	AppCode string
	Secret  string
	Sender  string
}

func NewQiscus(cfg *QiscusConfig) *Qiscus {
	return &Qiscus{
		Config: cfg,
	}
}

func (q *Qiscus) InternalHitTest() {

}

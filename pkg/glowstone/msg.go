package glowstone

type Msg struct {
	Src     string
	Dest    string
	Payload []byte
}

func NewMsg(payload []byte, src string, dest string) Msg {
	return Msg{
		Src:     src,
		Dest:    dest,
		Payload: payload,
	}
}

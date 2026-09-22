package typesafe

type Request struct {
	State     EntryType           `json:"state"`
	Model     string              `json:"model"`
	Questions map[string]Question `json:"questions"`
}

func NewRequest(state EntryType) Request {
	return Request{
		State:     state,
		Model:     ModelLatest,
		Questions: make(map[string]Question),
	}
}
func (r *Request) AddQuestion(name string, q Question) {
	r.Questions[name] = q
}

package typesafe

// This is the standard Request type per the Typesafe.ai documentation
// found at https://docs.typesafe.ai/api
type Request struct {
	// The data you are asking questions about
	State EntryType `json:"state"`
	// The identifier of the model you wish to use, defaults to "jev-latest"
	Model string `json:"model"`
	// The questions you wish to have jev provide answers for.  Must be non-empty before
	// the request is sent
	Questions map[string]Question `json:"questions"`
}

// Create a new [Request], defaulting to using the 'jev-latest' model
// The model can be updated directly.  Please use [*Request.AddQuestion] to add
// new questions to the request
func NewRequest(state EntryType) Request {
	return Request{
		State:     state,
		Model:     ModelLatest,
		Questions: make(map[string]Question),
	}
}

// Adds a new [Question] with the given `name`
func (r *Request) AddQuestion(name string, q Question) {
	r.Questions[name] = q
}

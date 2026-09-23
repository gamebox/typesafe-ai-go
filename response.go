package typesafe

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
)

// Details token usage for a request
type Usage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
}

// This is the standard Response type as documented in Typesafe.ai's
// own documentation at https://docs.typesafe.ai/api
type Response struct {
	// The model used
	Model string `json:"model"`
	// The answers to the questions asked.  Should never be nil or zero-length
	Answers map[string]jsontext.Value `json:"answers"`
	// Token usage information for this request
	Usage Usage `json:"usage"`
}

// Looks up a answer with `name` and if present attempts to decode
// it as a [NoulAnswer]. Will return false if an answer does not exist,
// or less likely if it can not be decoded as a NoulAnswer
func (r Response) DecodeNoul(name string) (NoulAnswer, bool) {
	answer := NoulAnswer{}

	a, ok := r.Answers[name]
	if !ok {
		return answer, false
	}

	reader := bytes.NewReader(a)
	dec := jsontext.NewDecoder(reader)
	err := json.UnmarshalDecode(dec, &answer)
	return answer, err == nil
}

// Looks up a answer with `name` and if present attempts to decode
// it as a [ChoiceAnswer]. Will return false if an answer does not exist,
// or less likely if it can not be decoded as a ChoiceAnswer
func (r Response) DecodeChoice(name string) (ChoiceAnswer, bool) {
	answer := ChoiceAnswer{}

	a, ok := r.Answers[name]
	if !ok {
		return answer, false
	}

	reader := bytes.NewReader(a)
	dec := jsontext.NewDecoder(reader)
	err := json.UnmarshalDecode(dec, &answer)
	return answer, err == nil
}

// Looks up a answer with `name` and if present attempts to decode
// it as a [ScoreAnswer]. Will return false if an answer does not exist,
// or less likely if it can not be decoded as a ScoreAnswer
func (r Response) DecodeScore(name string) (ScoreAnswer, bool) {
	answer := ScoreAnswer{}

	a, ok := r.Answers[name]
	if !ok {
		return answer, false
	}

	reader := bytes.NewReader(a)
	dec := jsontext.NewDecoder(reader)
	err := json.UnmarshalDecode(dec, &answer)
	return answer, err == nil
}

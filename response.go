package typesafe

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
)

type Usage struct {
	InputTokens  int64 `json:"input_tokens"`
	OutputTokens int64 `json:"output_tokens"`
}

type Response struct {
	Model   string                    `json:"model"`
	Answers map[string]jsontext.Value `json:"answers"`
	Usage   Usage                     `json:"usage"`
}

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

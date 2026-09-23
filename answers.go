package typesafe

// Any of the answer types allowed
type Answer any

// A [Noul answer]
//
// [Noul answer]: https://docs.typesafe.ai/api#noul-answer
type NoulAnswer struct {
	// This will always be `noul`
	Type string `json:"type"`
	// The yes/no answer on a scale from 0 (no) to 1 (yes).
	Noul float64 `json:"noul"`
}

// A [Choice answer]
//
// [Choice answer]: https://docs.typesafe.ai/api#choice-answer
type ChoiceAnswer struct {
	// This wil always be `choice`
	Type string `json:"type"`
	// The highest-probability option.
	Choice string `json:"choice"`
	// Every option mapped to its probability (floats that sum to 1).
	Probabilities map[string]float64 `json:"probabilities"`
	// How certain the model is, derived from probabilities.
	Confidence float64 `json:"confidence"`
}

// A [Score answer]
//
// [Score answer]: https://docs.typesafe.ai/api#score-answer
type ScoreAnswer struct {
	// This wil always be `score`
	Type string `json:"type"`
	// The probability-weighted answer across the levels; can land between levels.
	Score float64 `json:"score"`
	// Each level number mapped back to its description.
	Legend map[string]string `json:"legend"`
	// Each level (string key) mapped to its probability (floats that sum to 1).
	Probabilities map[string]float64 `json:"probabilities"`
	// How certain the model is, derived from probabilities.
	Confidence float64 `json:"confidence"`
}

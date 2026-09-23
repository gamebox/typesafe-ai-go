package typesafe

// A marker interface for all question types
type Question interface {
	Ask()
}

// A [Noul question type] following [Typesafe.ai's documentation]
//
// [Noul question type]: https://docs.typesafe.ai/api#noul
// [Typesafe.ai's documentation]: https://docs.typesafe.ai/api
type NoulQuestion struct {
	Type         string               `json:"type"`
	Instructions EntryType            `json:"instructions"`
	Criteria     map[string]EntryType `json:"criteria"`
}

// Create a new [NoulQuestion].  Use this over creating a struct directly.
func Noul(instructions EntryType) NoulQuestion {
	return NoulQuestion{Instructions: instructions, Type: "noul"}
}

// Sets the criteria for the "true" condition
func (q *NoulQuestion) TrueCriteria(desc string) {
	q.Criteria["true"] = desc
}

// Sets the criteria for the "true" condition
func (q *NoulQuestion) FalseCriteria(desc string) {
	q.Criteria["false"] = desc
}

// Marker method for [Question] interface. Noop
func (q NoulQuestion) Ask() {}

// A [Choice question type] following [Typesafe.ai's documentation]
//
// [Choice question type]: https://docs.typesafe.ai/api#choice
// [Typesafe.ai's documentation]: https://docs.typesafe.ai/api
type ChoiceQuestion struct {
	Type         string               `json:"type"`
	Instructions EntryType            `json:"instructions"`
	Criteria     map[string]EntryType `json:"criteria"`
}

// Create a new [ChoiceQuestion].  Use this over creating a struct directly.
func Choice(instructions EntryType) ChoiceQuestion {
	return ChoiceQuestion{
		Instructions: instructions,
		Type:         "choice",
		Criteria:     make(map[string]EntryType),
	}
}

// Add criteria for choice `name` using the given description `desc`
func (q *ChoiceQuestion) AddCriteria(name string, desc any) {
	if q.Criteria == nil {
		q.Criteria = make(map[string]EntryType)
	}
	q.Criteria[name] = desc
}

// Marker method for [Question] interface. Noop
func (q ChoiceQuestion) Ask() {}

// A [Score question type] following [Typesafe.ai's documentation]
//
// [Score question type]: https://docs.typesafe.ai/api#score
// [Typesafe.ai's documentation]: https://docs.typesafe.ai/api
type ScoreQuestion struct {
	Type         string    `json:"type"`
	Instructions EntryType `json:"instructions"`
	Criteria     []any     `json:"criteria"`
}

// Create a new [ScoreQuestion].  Use this over creating a struct directly.
func Score(instructions EntryType) ScoreQuestion {
	return ScoreQuestion{
		Instructions: instructions,
		Type:         "score",
		Criteria:     make([]any, 0),
	}
}

// Adds a new score criteria using the description `desc`
func (q *ScoreQuestion) AddCriteria(desc any) {
	q.Criteria = append(q.Criteria, desc)
}

// Marker method for [Question] interface. Noop
func (q ScoreQuestion) Ask() {}

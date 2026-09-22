package typesafe

type Question interface {
	Ask()
}

type NoulQuestion struct {
	Type         string               `json:"type"`
	Instructions EntryType            `json:"instructions"`
	Criteria     map[string]EntryType `json:"criteria"`
}

func Noul(instructions EntryType) NoulQuestion {
	return NoulQuestion{Instructions: instructions, Type: "noul"}
}

func (q *NoulQuestion) TrueCriteria(desc string) {
	q.Criteria["true"] = desc
}
func (q *NoulQuestion) FalseCriteria(desc string) {
	q.Criteria["false"] = desc
}
func (q NoulQuestion) Ask() {}

type ChoiceQuestion struct {
	Type         string               `json:"type"`
	Instructions EntryType            `json:"instructions"`
	Criteria     map[string]EntryType `json:"criteria"`
}

func Choice(instructions EntryType) ChoiceQuestion {
	return ChoiceQuestion{
		Instructions: instructions,
		Type:         "choice",
		Criteria:     make(map[string]EntryType),
	}
}

func (q *ChoiceQuestion) AddCriteria(name string, desc any) {
	if q.Criteria == nil {
		q.Criteria = make(map[string]EntryType)
	}
	q.Criteria[name] = desc
}
func (q ChoiceQuestion) Ask() {}

type ScoreQuestion struct {
	Type         string    `json:"type"`
	Instructions EntryType `json:"instructions"`
	Criteria     []any     `json:"criteria"`
}

func Score(instructions EntryType) ScoreQuestion {
	return ScoreQuestion{
		Instructions: instructions,
		Type:         "score",
		Criteria:     make([]any, 0),
	}
}

func (q *ScoreQuestion) AddCriteria(desc any) {
	q.Criteria = append(q.Criteria, desc)
}
func (q ScoreQuestion) Ask() {}

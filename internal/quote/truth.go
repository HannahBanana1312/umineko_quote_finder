package quote

type Truth string

const (
	TruthAll  Truth = ""
	TruthRed  Truth = "red"
	TruthBlue Truth = "blue"
)

func (Truth) Parse(s string) Truth {
	switch s {
	case "red":
		return TruthRed
	case "blue":
		return TruthBlue
	default:
		return TruthAll
	}
}

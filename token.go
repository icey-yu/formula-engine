package formula_engine

type token struct {
	Type  TT
	Value string
	Start int
	End   int
}

func newToken(type_ TT, value string, start int, end int) *token {
	return &token{
		Type:  type_,
		Value: value,
		Start: start,
		End:   end,
	}
}

package formula_engine

import "github.com/shopspring/decimal"

func GetAstTreeByString(str string) (AstNode, error) {
	tokens, err := newLexer(str).MakeTokens()
	if err != nil {
		return nil, err
	}
	return newParser(tokens).Parse()
}

func CalByAstTree(node AstNode, identifierMap map[string]string) (*decimal.Decimal, error) {
	return newInterpreter(node, identifierMap).Interpret()
}

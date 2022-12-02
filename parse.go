package formula_engine

import (
	"fmt"
	"github.com/pkg/errors"
)

// parser 语法解析器
type parser struct {
	Tokens       []*token
	CurrentToken *token
	LastIdx      int
	Idx          int
}

func newParser(t []*token) *parser {
	p := &parser{
		Tokens:  t,
		LastIdx: -1,
		Idx:     -1,
	}
	p.advance()
	return p
}

// Parse 解析
// <expr> ::= <and_term> { OR <expr> }
// <and_term> ::= <not_expr> { AND <and_expr>}
// <not_term> ::= { NOT } <com_expr>
// <com_term> ::= <pri_pre> { GT|LT|EQ|NEQ|GTE|LTE <pri_pre> }                     // compare term
// <pri_ope> ::= <pri_ope> { +|- <sec_ope> }                                   // Primary operation
// <sec_ope> ::= <sec_ope> { *|/ <ter_ope> }                                   // Secondary operation
// <ter_ope> ::= <factor> { ^ <ter_ope> }                                      // Tertiary operation
// <factor> ::= NUM| FUNCTION LPAREN [ expr { COMMA expr }] RPAREN| IDENTIFIER| { PLUS | MINUS } factor| LPAREN expr RPAREN
func (p *parser) Parse() (AstNode, error) {
	res, err := p.expr()
	if err != nil {
		return nil, err
	} else if p.CurrentToken.Type != TTEof {
		return nil, p.makeErr(illegalSyntaxErrMsg, fmt.Sprintf("Unable to parse completely.Idx: %d", p.Idx))
	}
	return res, nil
}

// expr <expr> ::= <and_term> { OR <expr> }
func (p *parser) expr() (AstNode, error) {
	return p.binOpLeft(p.andTerm, p.expr, []TT{TTOr})
}

func (p *parser) andTerm() (AstNode, error) {
	return p.binOpLeft(p.notTerm, p.andTerm, []TT{TTAnd})
}

func (p *parser) notTerm() (AstNode, error) {
	if p.CurrentToken.Type == TTNot {
		tok := p.CurrentToken
		p.advance()
		node, err := p.notTerm()
		if err != nil {
			return nil, err
		}
		return newAstUnNode(tok, node), nil
	} else {
		return p.comTerm()
	}
}

func (p *parser) comTerm() (AstNode, error) {
	return p.binOpLeft(p.priOpe, p.priOpe, []TT{TTGt, TTGte, TTEq, TTNeq, TTLt, TTLte})
}

func (p *parser) priOpe() (AstNode, error) {
	return p.binOpRight(p.secOpe, []TT{TTPlus, TTMinus})
}

func (p *parser) secOpe() (AstNode, error) {
	return p.binOpRight(p.terOpe, []TT{TTMul, TTDiv})
}

func (p *parser) terOpe() (AstNode, error) {
	return p.binOpLeft(p.factor, p.terOpe, []TT{TTPow})
}

// factor <factor> ::= NUM | FUNCTION LPAREN [ expr { COMMA IDENTIFIER }] RPAREN | IDENTIFIER | { PLUS | MINUS } factor | LPAREN expr RPAREN
func (p *parser) factor() (AstNode, error) {
	tok := p.CurrentToken
	switch {
	case InSlice([]TT{TTNum, TTIdentifier}, tok.Type):
		// NUM | IDENTIFIER
		p.advance()
		return newAstSinNode(tok), nil
	case InSlice([]TT{TTPlus, TTMinus}, tok.Type):
		// { PLUS | MINUS } factor
		p.advance()
		fac, err := p.factor()
		if err != nil {
			return nil, err
		}
		return newAstUnNode(tok, fac), nil
	case tok.Type == TTFunction:
		// FUNCTION LPAREN [ expr { COMMA IDENTIFIER }] RPAREN
		p.advance()
		if p.CurrentToken.Type != TTLparen {
			return nil, p.makeErr(illegalSyntaxErrMsg, fmt.Sprintf("UnExpected tokType:'%s', expected '(' after function name", p.CurrentToken.Type))
		}
		params := make([]AstNode, 0)
		p.advance()
		if p.CurrentToken.Type != TTRparen {
			node, err := p.expr()
			if err != nil {
				return nil, err
			}
			params = append(params, node)
			for p.CurrentToken.Type == TTComma {
				p.advance()
				node, err := p.expr()
				if err != nil {
					return nil, err
				}
				params = append(params, node)
			}
		}

		if p.CurrentToken.Type != TTRparen {
			return nil, p.makeErr(illegalSyntaxErrMsg, fmt.Sprintf("UnExpected tokType:'%s', expected ')' when there is '(' before", p.CurrentToken.Type))
		}
		num, ok := FuncParNumMap[tok.Value]
		if !ok {
			return nil, p.makeErr(systemErrMsg, fmt.Sprintf("Can not found function name %s in FuncParNumMap, please plus it", tok.Value))
		}
		err := checkParNum(num, len(params))
		if err != nil {
			return nil, errors.Wrapf(err, getTokPos(tok))
		}

		p.advance()
		return newAstGeneralNode(tok, params...), nil

	case tok.Type == TTLparen:
		// LPAREN expr RPAREN
		p.advance()
		expr, err := p.expr()
		if err != nil {
			return nil, err
		}
		if p.CurrentToken.Type == TTRparen {
			p.advance()
			return expr, nil
		} else {
			return nil, p.makeErr(illegalSyntaxErrMsg, fmt.Sprintf("UnExpected tokType:'%s', expected ')' when there is '(' before", tok.Type))
		}
	default:
		return nil, p.makeErr(illegalSyntaxErrMsg, fmt.Sprintf("UnExpected tokType:'%s'", tok.Type))
	}
}

// binOpRight 二元操作生成默认右枝存在二叉树（如：+，-，*，/）
func (p *parser) binOpRight(f func() (AstNode, error), ops []TT) (AstNode, error) {
	p.LastIdx = p.Idx
	left, err := f()
	if err != nil {
		return nil, err
	}
	for InSlice(ops, p.CurrentToken.Type) {

		tok := p.CurrentToken
		p.advance()
		right, err := f()

		if err != nil {
			return nil, err
		}
		// 最后构建的时候，如果只有left，则返回left。如果有tok和right，组装构建后返回
		left = newAstBinNode(tok, left, right)
	}
	return left, nil
}

// binOpLeft 二元操作生成默认左枝存在二叉树（如：^）
func (p *parser) binOpLeft(f1 func() (AstNode, error), f2 func() (AstNode, error), ops []TT) (AstNode, error) {
	left, err := f1()
	if err != nil {
		return nil, err
	}
	if InSlice(ops, p.CurrentToken.Type) {
		tok := p.CurrentToken
		p.advance()
		right, err := f2()
		if err != nil {
			return nil, err
		}
		// 最后构建的时候，如果只有left，则返回left。如果有tok和right，组装构建后返回
		left = newAstBinNode(tok, left, right)
	}
	return left, nil
}

// advance 下一个
func (p *parser) advance() {
	p.Idx += 1
	if p.Idx < len(p.Tokens) {
		p.CurrentToken = p.Tokens[p.Idx]
	} else {
		p.CurrentToken = nil
	}
}

// resetLast 重设记录点
func (p *parser) resetLast() {
	p.LastIdx = -1
}

// rollBackToLast 回滚到记录点
func (p *parser) rollBackToLast() {
	p.Idx = p.LastIdx
	p.CurrentToken = p.Tokens[p.Idx]
}

// makeErrWithIdx 组装错误
func (p *parser) makeErr(errName, details string) error {
	return makeErrWithIdx(p.Idx, errName, details)
}

package formula_engine

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// interpreter 解释器
type interpreter struct {
	Root          AstNode
	IdentifierMap map[string]string
	CurrentToken  *token
	// 该map能够根据节点类型决定访问哪个visit方法
	visitMap map[string]func(node AstNode) (*decimal.Decimal, error)
	// 该map能够通过TT类型决定访问哪个一元计算方法
	unVisMap map[TT]func(p *decimal.Decimal) (*decimal.Decimal, error)
	// 该map能够通过TT类型决定访问哪个二元计算方法
	binVisMap map[TT]func(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error)
}

func newInterpreter(root AstNode, identifierMap map[string]string) *interpreter {
	i := &interpreter{
		Root:          root,
		IdentifierMap: identifierMap,
	}
	i.visitMap = map[string]func(node AstNode) (*decimal.Decimal, error){
		astSinNodeName:     i.visitAstSinNode,
		astUnNodeName:      i.visitAstUnNode,
		astBinNodeName:     i.visitAstBinNode,
		astGeneralNodeName: i.visitAstGeneralNode,
	}
	i.unVisMap = map[TT]func(p *decimal.Decimal) (*decimal.Decimal, error){
		TTPlus:  unPlus,
		TTMinus: unMinus,
		TTNot:   not,
	}
	i.binVisMap = map[TT]func(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error){
		TTPlus:  plus,
		TTMinus: minus,
		TTMul:   mul,
		TTDiv:   div,
		TTPow:   pow,
		TTAnd:   and,
		TTOr:    or,
		TTEq:    eq,
		TTNeq:   neq,
		TTGt:    gt,
		TTGte:   gte,
		TTLt:    lt,
		TTLte:   lte,
	}
	return i
}

func (i *interpreter) Interpret() (*decimal.Decimal, error) {
	return i.visit(i.Root)
}

// visit 通用访问入口
func (i *interpreter) visit(node AstNode) (*decimal.Decimal, error) {
	i.CurrentToken = node.GetTok()
	return i.visitMap[node.GetName()](node)
}

// visitAstSinNode 访问单节点
func (i *interpreter) visitAstSinNode(node AstNode) (*decimal.Decimal, error) {
	tok := node.GetTok()
	// 如果该token为变量，通过IdentifierMap获取其值。
	if tok.Type == TTIdentifier {
		val, ok := i.IdentifierMap[tok.Value]
		if !ok {
			return nil, makeErrWithToken(tok, illegalCalErrMsg, fmt.Sprintf("Cannot found a value by key %s in IdentifierMap, please plus it.", tok.Value))
		}
		tok.Value = val
	}
	dec, err := decimal.NewFromString(tok.Value)
	if err != nil {
		return nil, makeErrWithToken(tok, systemErrMsg, err.Error())
	}
	return &dec, nil
}

// visitAstUnNode 访问单支节点
func (i *interpreter) visitAstUnNode(node AstNode) (*decimal.Decimal, error) {
	tok := i.CurrentToken
	binNode, ok := node.(*astUnNode)
	if !ok {
		return nil, makeErrWithToken(node.GetTok(), systemErrMsg, "Is not astUnNode type,please check method GetName().")
	}
	child, err := i.visit(binNode.Node)
	if err != nil {
		return nil, err
	}
	fun, ok := i.unVisMap[tok.Type]
	if !ok {
		return nil, makeErrWithToken(tok, systemErrMsg, fmt.Sprintf("UnKnow Unary type %s", i.CurrentToken.Type))
	}
	return fun(child)
}

// visitAstBinNode 访问二叉节点
func (i *interpreter) visitAstBinNode(node AstNode) (*decimal.Decimal, error) {
	tok := i.CurrentToken
	binNode, ok := node.(*astBinNode)
	if !ok {
		return nil, makeErrWithToken(node.GetTok(), systemErrMsg, "Is not astBinNode type,please check method GetName().")
	}
	left, err := i.visit(binNode.LNode)
	if err != nil {
		return nil, err
	}
	right, err := i.visit(binNode.RNode)
	if err != nil {
		return nil, err
	}
	fun, ok := i.binVisMap[tok.Type]
	if !ok {
		return nil, makeErrWithToken(i.CurrentToken, systemErrMsg, fmt.Sprintf("UnKnow Binary type %s", i.CurrentToken.Type))
	}
	res, err := fun(left, right)
	if err != nil {
		return nil, errors.Wrapf(err, getTokPos(tok))
	}
	return res, nil
}

// visitAstGeneralNode 访问一般节点
func (i *interpreter) visitAstGeneralNode(node AstNode) (*decimal.Decimal, error) {
	tok := i.CurrentToken
	binNode, ok := node.(*astGeneralNode)
	if !ok {
		return nil, makeErrWithToken(node.GetTok(), systemErrMsg, "Is not astGeneralNode type,please check method GetName().")
	}
	params := make([]*decimal.Decimal, 0)
	for _, n := range binNode.Nodes {
		param, err := i.visit(n)
		if err != nil {
			return nil, err
		}
		params = append(params, param)
	}

	res, err := function(tok.Value, params...)
	if err != nil {
		return nil, errors.Wrapf(err, getTokPos(tok))
	}
	return res, nil
}

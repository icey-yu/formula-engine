package formula_engine

type AstNode interface {
	GetName() string
	GetTok() *token
}

// astGeneralNode 一般树节点
type astGeneralNode struct {
	Tok   *token
	Nodes []AstNode
}

func newAstGeneralNode(t *token, node ...AstNode) *astGeneralNode {
	return &astGeneralNode{
		Tok:   t,
		Nodes: node,
	}
}

func (a *astGeneralNode) GetName() string {
	return "astGeneralNode"
}

func (a *astGeneralNode) GetTok() *token {
	return a.Tok
}

//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------

// astBinNode BinaryNode 二叉树
type astBinNode struct {
	Tok   *token
	LNode AstNode
	RNode AstNode
}

func newAstBinNode(t *token, lNode AstNode, rNode AstNode) *astBinNode {
	return &astBinNode{
		Tok:   t,
		LNode: lNode,
		RNode: rNode,
	}
}

func (a *astBinNode) GetName() string {
	return "astBinNode"
}

func (a *astBinNode) GetTok() *token {
	return a.Tok
}

//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------

// astUnNode UnaryNode,单叉树
type astUnNode struct {
	Tok  *token
	Node AstNode
}

func newAstUnNode(t *token, node AstNode) *astUnNode {
	return &astUnNode{
		Tok:  t,
		Node: node,
	}
}

func (a *astUnNode) GetName() string {
	return "astUnNode"
}

func (a *astUnNode) GetTok() *token {
	return a.Tok
}

//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------

// astSinNode single,单个节点
type astSinNode struct {
	Tok *token
}

func newAstSinNode(t *token) *astSinNode {
	return &astSinNode{
		Tok: t,
	}
}

func (a *astSinNode) GetName() string {
	return "astSinNode"
}

func (a *astSinNode) GetTok() *token {
	return a.Tok
}

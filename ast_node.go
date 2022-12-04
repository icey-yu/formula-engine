package formula_engine

const (
	astGeneralNodeName = "astGeneralNode"
	astBinNodeName     = "astBinNode"
	astUnNodeName      = "astUnNode"
	astSinNodeName     = "astSinNode"
)

type AstNode interface {
	GetName() string
	GetTok() *token
}

func DeepCopyAstNode(node AstNode) AstNode {
	if node == nil {
		return nil
	}
	switch node.GetName() {
	case astSinNodeName:
		sNode := node.(*astSinNode)
		t := &token{
			Type:  sNode.Tok.Type,
			Value: sNode.Tok.Value,
			Start: sNode.Tok.Start,
			End:   sNode.Tok.End,
		}
		return newAstSinNode(t)
	case astUnNodeName:
		sNode := node.(*astUnNode)
		t := &token{
			Type:  sNode.Tok.Type,
			Value: sNode.Tok.Value,
			Start: sNode.Tok.Start,
			End:   sNode.Tok.End,
		}
		child := DeepCopyAstNode(sNode.Node)
		return newAstUnNode(t, child)
	case astBinNodeName:
		sNode := node.(*astBinNode)
		t := &token{
			Type:  sNode.Tok.Type,
			Value: sNode.Tok.Value,
			Start: sNode.Tok.Start,
			End:   sNode.Tok.End,
		}
		lNode := DeepCopyAstNode(sNode.LNode)
		RNode := DeepCopyAstNode(sNode.RNode)
		return newAstBinNode(t, lNode, RNode)
	case astGeneralNodeName:
		sNode := node.(*astGeneralNode)
		t := &token{
			Type:  sNode.Tok.Type,
			Value: sNode.Tok.Value,
			Start: sNode.Tok.Start,
			End:   sNode.Tok.End,
		}
		children := make([]AstNode, 0)
		for _, c := range sNode.Nodes {
			children = append(children, DeepCopyAstNode(c))
		}
		return newAstGeneralNode(t, children...)
	}
	return nil
}

//--------------------------------------------------------------------------------
//--------------------------------------------------------------------------------

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
	return astGeneralNodeName
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
	return astBinNodeName
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
	return astUnNodeName
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
	return astSinNodeName
}

func (a *astSinNode) GetTok() *token {
	return a.Tok
}

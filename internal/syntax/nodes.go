package syntax

// Clause is the clause.
type Clause struct {
	ExprList []Expr
}

// Expr is the expression.
type Expr interface {
	Node
	aExpr()
}

// Node is the node.
type Node interface {
	Pos() Pos
	aNode()
}

type expr struct {
	node
}

func (*expr) aExpr() {}

type node struct {
	pos Pos
}

func (n *node) Pos() Pos { return n.pos }
func (*node) aNode()     {}

// BindVarExpr is the bind var expression.
type BindVarExpr struct {
	Type  bindVarStyle
	Index int
	expr
}

// bindVarStyle is the type of bind vars.
type bindVarStyle int

const (
	// bindVarDollar is the style of bind vars like $1, $2, $3
	bindVarDollar bindVarStyle = iota
	// bindVarQuestion is the style of bind vars like ?, ?, ?
	bindVarQuestion
)

// PlainExpr is the plain text expression.
type PlainExpr struct {
	Text string
	expr
}

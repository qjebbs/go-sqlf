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
	Type  BindVarStyle
	Index int
	expr
}

// BindVarStyle is the type of bind vars.
type BindVarStyle int

const (
	// Dollar is the style of bind vars like $1, $2, $3
	Dollar BindVarStyle = iota
	// Question is the style of bind vars like ?, ?, ?
	Question
)

// PlainExpr is the plain text expression.
type PlainExpr struct {
	Text string
	expr
}

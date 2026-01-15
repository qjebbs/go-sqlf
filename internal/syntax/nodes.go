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
	typ   bindStyle
	Index int
	Name  string
	expr
}

// bindStyle is the type of bind vars.
type bindStyle int

const (
	// bindStyleQuestion is the style of bind vars like ?, ?, ?
	bindStyleQuestion bindStyle = iota
	// bindStyleDollarNumbered is the style of bind vars like $1, $2, $3
	bindStyleDollarNumbered
	// bindStyleQuestionNumbered is the style of bind vars like ?1, ?2, ?3
	bindStyleQuestionNumbered
	// bindStyleColonNamed is the style of bind vars like :name, :other
	bindStyleColonNamed
	// bindStyleColonNumbered is the style of bind vars like :1, :2, :3
	bindStyleColonNumbered
	// bindStyleAtNamed is the style of bind vars like @name, @other
	bindStyleAtNamed
)

// PlainExpr is the plain text expression.
type PlainExpr struct {
	Text string
	expr
}

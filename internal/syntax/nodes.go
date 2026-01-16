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
	typ   BindStyle
	Index int
	Name  string
	expr
}

// BindStyle is the type of bind vars.
type BindStyle int

const (
	// BindStyleUnknown is the unknown style of bind vars.
	BindStyleUnknown BindStyle = iota
	// BindStyleQuestion is the style of bind vars like ?, ?, ?
	BindStyleQuestion
	// BindStyleDollarNumbered is the style of bind vars like $1, $2, $3
	BindStyleDollarNumbered
	// BindStyleQuestionNumbered is the style of bind vars like ?1, ?2, ?3
	BindStyleQuestionNumbered
	// BindStyleColonNamed is the style of bind vars like :name, :other
	BindStyleColonNamed
	// BindStyleColonNumbered is the style of bind vars like :1, :2, :3
	BindStyleColonNumbered
	// BindStyleAtNamed is the style of bind vars like @name, @other
	BindStyleAtNamed
)

// PlainExpr is the plain text expression.
type PlainExpr struct {
	Text string
	expr
}

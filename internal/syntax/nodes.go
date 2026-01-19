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
	typ   BindVarStyle
	Index int
	Name  string
	expr
}

// BindVarStyle is the type of bind vars.
type BindVarStyle int

const (
	// BindVarStyleUnknown is the unknown style of bind vars.
	BindVarStyleUnknown BindVarStyle = iota
	// BindVarStyleQuestion is the style of bind vars like ?, ?, ?
	BindVarStyleQuestion
	// BindVarStyleDollarNumbered is the style of bind vars like $1, $2, $3
	BindVarStyleDollarNumbered
	// BindVarStyleQuestionNumbered is the style of bind vars like ?1, ?2, ?3
	BindVarStyleQuestionNumbered
	// BindVarStyleColonNamed is the style of bind vars like :name, :other
	BindVarStyleColonNamed
	// BindVarStyleColonNumbered is the style of bind vars like :1, :2, :3
	BindVarStyleColonNumbered
	// BindVarStyleAtNamed is the style of bind vars like @name, @other
	BindVarStyleAtNamed
)

// IdentityExpr is the identity expression.
type IdentityExpr struct {
	Name string
	expr
}

// PlainExpr is the plain text expression.
type PlainExpr struct {
	Text string
	expr
}

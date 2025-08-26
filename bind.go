package sqlf

// BindStyle is the type of bind vars.
type BindStyle int

const (
	// BindStyleDollar is the style of bind vars like $1, $2, $3
	BindStyleDollar BindStyle = iota
	// BindStyleQuestion is the style of bind vars like ?, ?, ?
	BindStyleQuestion
)

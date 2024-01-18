package sqlb

import (
	"fmt"
	"log"
	"strings"

	"github.com/qjebbs/go-sqlf"
	"github.com/qjebbs/go-sqlf/util"
)

// Build builds the query.
func (b *QueryBuilder) Build() (query string, args []any, err error) {
	ctx := sqlf.NewContext()
	ctx.WithBindVarStyle(b.bindVarStyle)
	query, err = b.buildInternal(ctx)
	if err != nil {
		return "", nil, err
	}
	args = ctx.Args()
	return query, args, nil
}

// BuildContext builds the query with the context.
func (b *QueryBuilder) BuildContext(ctx *sqlf.Context) (query string, err error) {
	return b.buildInternal(ctx)
}

// Debug enables debug mode.
func (b *QueryBuilder) Debug() {
	b.debug = true
}

// buildInternal builds the query with the selects.
func (b *QueryBuilder) buildInternal(ctx *sqlf.Context) (string, error) {
	if b == nil {
		return "", nil
	}
	if err := b.anyError(); err != nil {
		return "", err
	}
	clauses := make([]string, 0)

	dep, err := b.calcDependency()
	if err != nil {
		return "", err
	}

	sq, err := b.buildCTEs(ctx, dep)
	if err != nil {
		return "", err
	}
	if sq != "" {
		clauses = append(clauses, sq)
	}

	sel, err := b.buildSelects(ctx)
	if err != nil {
		return "", err
	}
	clauses = append(clauses, sel)
	from, err := b.buildFrom(ctx, dep)
	if err != nil {
		return "", err
	}
	if from != "" {
		clauses = append(clauses, from)
	}
	where, err := b.conditions.BuildContext(ctx)
	if err != nil {
		return "", err
	}
	if where != "" {
		clauses = append(clauses, where)
	}
	groupby, err := b.groupbys.BuildContext(ctx)
	if err != nil {
		return "", err
	}
	if groupby != "" {
		clauses = append(clauses, groupby)
	}
	order, err := b.orders.BuildContext(ctx)
	if err != nil {
		return "", err
	}
	if order != "" {
		clauses = append(clauses, order)
	}
	if b.limit > 0 {
		clauses = append(clauses, fmt.Sprintf(`LIMIT %d`, b.limit))
	}
	if b.offset > 0 {
		clauses = append(clauses, fmt.Sprintf(`OFFSET %d`, b.offset))
	}
	query := strings.TrimSpace(strings.Join(clauses, " "))
	if len(b.unions) > 0 {
		union, err := b.buildUnion(ctx)
		if err != nil {
			return "", err
		}
		query = strings.TrimSpace(query + " " + union)
	}
	if b.debug {
		interpolated, err := util.Interpolate(query, ctx.Args())
		if err != nil {
			log.Printf("debug: interpolated query: %s\n", err)
		}
		log.Println(interpolated)
	}
	return query, nil
}

func (b *QueryBuilder) buildCTEs(ctx *sqlf.Context, dep map[Table]bool) (string, error) {
	if len(b.ctes) == 0 {
		return "", nil
	}
	clauses := make([]string, 0, len(b.ctes))
	for _, cte := range b.ctes {
		if !dep[cte.table] {
			continue
		}
		query, err := cte.BuildContext(ctx)
		if err != nil {
			return "", fmt.Errorf("build CTE '%s': %w", cte.table, err)
		}
		if query == "" {
			continue
		}
		clauses = append(clauses, fmt.Sprintf(
			"%s AS (%s)",
			cte.table.Name, query,
		))
	}
	if len(clauses) == 0 {
		return "", nil
	}
	return "With " + strings.Join(clauses, ", "), nil
}

func (b *QueryBuilder) buildSelects(ctx *sqlf.Context) (string, error) {
	if b.distinct {
		b.selects.Prefix = "SELECT DISTINCT"
	} else {
		b.selects.Prefix = "SELECT"
	}
	sel, err := b.selects.BuildContext(ctx)
	if err != nil {
		return "", err
	}
	touches, err := b.touches.BuildContext(ctx)
	if err != nil {
		return "", err
	}
	if sel == "" {
		return "", fmt.Errorf("no columns selected")
	}
	if touches == "" {
		return sel, nil
	}
	return sel + ", " + touches, nil
}

func (b *QueryBuilder) buildFrom(ctx *sqlf.Context, dep map[Table]bool) (string, error) {
	tables := make([]string, 0, len(b.tables))
	for _, t := range b.tables {
		ft, ok := b.froms[t]
		if !ok {
			// should not happen
			return "", fmt.Errorf("table '%s' not found", t)
		}
		if b.distinct && ft.Optional && !dep[t] {
			continue
		}
		from := b.froms[t]
		c, err := from.Fragment.BuildContext(ctx)
		if err != nil {
			return "", fmt.Errorf("build FROM '%s': %w", from.Fragment.Raw, err)
		}
		tables = append(tables, c)
	}
	return "FROM " + strings.Join(tables, " "), nil
}

func (b *QueryBuilder) buildUnion(ctx *sqlf.Context) (string, error) {
	clauses := make([]string, 0, len(b.unions))
	for _, union := range b.unions {
		query, err := union.BuildContext(ctx)
		if err != nil {
			return "", err
		}
		if query == "" {
			continue
		}
		clauses = append(clauses, query)
	}
	return "UNION (" + strings.Join(clauses, ") UNION (") + ")", nil
}

package sqlb

import (
	"fmt"
	"log"
	"strings"

	"github.com/qjebbs/go-sqlf/v3"
	"github.com/qjebbs/go-sqlf/v3/syntax"
	"github.com/qjebbs/go-sqlf/v3/util"
)

var _ Builder = (*QueryBuilder)(nil)
var _ sqlf.Builder = (*QueryBuilder)(nil)

// BuildQuery builds the query.
func (b *QueryBuilder) BuildQuery(bindVarStyle syntax.BindVarStyle) (query string, args []any, err error) {
	ctx := sqlf.NewContext(bindVarStyle)
	query, err = b.buildInternal(ctx)
	if err != nil {
		return "", nil, err
	}
	args = ctx.Args()
	return query, args, nil
}

// Build implements sqlf.Builder
func (b *QueryBuilder) Build(ctx *sqlf.Context) (query string, err error) {
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

	dep, err := b.collectDependencies()
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

	// reserve a position for select
	selectAt := len(clauses)
	clauses = append(clauses, "")
	from, err := b.buildFrom(ctx, dep)
	if err != nil {
		return "", err
	}
	if from != "" {
		clauses = append(clauses, from)
	}
	where, err := sqlf.Prefix(
		"WHERE",
		sqlf.Join(" AND ", util.Ttoa(b.conditions)...),
	).Build(ctx)
	if err != nil {
		return "", err
	}
	if where != "" {
		clauses = append(clauses, where)
	}
	groupby, err := sqlf.Prefix(
		"GROUP BY",
		sqlf.Join(", ", util.Ttoa(b.groupbys)...),
	).Build(ctx)
	if err != nil {
		return "", err
	}
	if groupby != "" {
		clauses = append(clauses, groupby)
		having, err := sqlf.Prefix(
			"HAVING",
			sqlf.Join(" AND ", util.Ttoa(b.havings)...),
		).Build(ctx)
		if err != nil {
			return "", err
		}
		if having != "" {
			clauses = append(clauses, having)
		}
	}
	order, err := b.buildOrders(ctx)
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
	// select must build order, because buildOrders may add columns to touches
	sel, err := b.buildSelects(ctx)
	if err != nil {
		return "", err
	}
	clauses[selectAt] = sel
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

func (b *QueryBuilder) buildCTEs(ctx *sqlf.Context, dep map[TableAliased]bool) (string, error) {
	if len(b.ctes) == 0 {
		return "", nil
	}
	clauses := make([]string, 0, len(b.ctes))
	for _, cte := range b.ctes {
		if !dep[NewTableAliased(cte.name, "")] {
			continue
		}
		query, err := cte.Build(ctx)
		if err != nil {
			return "", fmt.Errorf("build CTE '%s': %w", cte.name, err)
		}
		if query == "" {
			continue
		}
		clauses = append(clauses, fmt.Sprintf(
			"%s AS (%s)",
			cte.name, query,
		))
	}
	if len(clauses) == 0 {
		return "", nil
	}
	return "With " + strings.Join(clauses, ", "), nil
}

func (b *QueryBuilder) buildSelects(ctx *sqlf.Context) (string, error) {
	prefix := "SELECT"
	if b.distinct {
		prefix = "SELECT DISTINCT"
	}
	sel, err := sqlf.Prefix(
		prefix,
		sqlf.Join(", ", util.Ttoa(b.selects)...),
	).Build(ctx)
	if err != nil {
		return "", err
	}
	touches, err := sqlf.Join(", ", util.Ttoa(b.touches)...).Build(ctx)
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

func (b *QueryBuilder) buildFrom(ctx *sqlf.Context, dep map[TableAliased]bool) (string, error) {
	tables := make([]string, 0, len(b.tables))
	for _, t := range b.tables {
		if (b.distinct || len(b.groupbys) > 0) && t.Optional && !dep[t.Names] {
			continue
		}
		c, err := t.Fragment.Build(ctx)
		if err != nil {
			return "", fmt.Errorf("build FROM '%s': %w", t.Names, err)
		}
		tables = append(tables, c)
	}
	return "FROM " + strings.Join(tables, " "), nil
}

func (b *QueryBuilder) buildUnion(ctx *sqlf.Context) (string, error) {
	clauses := make([]string, 0, len(b.unions))
	for _, union := range b.unions {
		query, err := union.Build(ctx)
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

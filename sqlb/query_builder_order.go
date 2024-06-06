package sqlb

import (
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

// Order is the sorting order.
type Order uint

// orders
const (
	Asc Order = iota
	AscNullsFirst
	AscNullsLast
	Desc
	DescNullsFirst
	DescNullsLast
)

var orders = []string{
	"ASC",
	"ASC NULLS FIRST",
	"ASC NULLS LAST",
	"DESC",
	"DESC NULLS FIRST",
	"DESC NULLS LAST",
}

type orderItem struct {
	column *Column
	order  Order
}

// OrderBy set the sorting order. the order can be "ASC", "DESC", "ASC NULLS FIRST" or "DESC NULLS LAST"
func (b *QueryBuilder) OrderBy(column *Column, order Order) *QueryBuilder {
	b.orders = append(b.orders, &orderItem{column: column, order: order})
	return b
}

func (b *QueryBuilder) buildOrders(ctx *sqlf.Context) (string, error) {
	f := sqlf.F("#join('#fragment', ', ')").WithPrefix("ORDER BY")
	for i, item := range b.orders {
		if item.order > DescNullsLast {
			b.pushError(fmt.Errorf("invalid order: %d", item.order))
			continue
		}
		if !b.distinct {
			f.AppendFragments(sqlf.Ff(
				"#f1 "+orders[item.order],
				item.column,
			))
			continue
		}
		// pq: for SELECT DISTINCT, ORDER BY expressions must appear in select list
		alias := fmt.Sprintf("_order_%d", i+1)
		orderStr := orders[item.order]
		b.touches.AppendFragments(sqlf.Ff("#f1 AS "+alias, item.column))
		f.AppendFragments(sqlf.F(
			fmt.Sprintf("%s %s", alias, orderStr),
		))
	}
	return f.BuildFragment(ctx)
}

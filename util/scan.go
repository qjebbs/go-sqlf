package util

import (
	"database/sql"
	"fmt"

	"github.com/qjebbs/go-sqlf/v2"
)

// QueryAble is the interface for query-able *sql.DB, *sql.Tx, etc.
type QueryAble interface {
	Exec(query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
}

// NewScanDestFunc is the function to create a new scan destination,
// returning the destination and its fields to scan.
type NewScanDestFunc[T any] func() (T, []any)

// ScanBuilder is like Scan, but it builds query from sqlf.Builder
func ScanBuilder[T any](db QueryAble, b sqlf.QueryBuilder, fn NewScanDestFunc[T]) ([]T, error) {
	query, args, err := b.BuildQuery()
	if err != nil {
		return nil, err
	}
	return Scan[T](db, query, args, fn)
}

// Scan scans query rows with scanner
func Scan[T any](db QueryAble, query string, args []any, fn NewScanDestFunc[T]) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		dest, fields := fn()
		err = ScanRow(rows, fields...)
		if err != nil {
			return nil, err
		}
		results = append(results, dest)
	}
	return results, nil
}

// ScanRow scans a single row to dest, unlike rows.Scan(), it drops the extra columns.
// It's useful when *sqlb.QueryBuilder.OrderBy() add extra column to the query.
func ScanRow(rows *sql.Rows, dest ...any) error {
	cols, err := rows.Columns()
	if err != nil {
		return err
	}
	nBlacholes := len(cols) - len(dest)
	bh := &blackhole{}
	for i := 0; i < nBlacholes; i++ {
		dest = append(dest, &bh)
	}
	return rows.Scan(dest...)
}

// CountBuilder is like Count, but it builds query from sqlf.Builder.
func CountBuilder(db QueryAble, b sqlf.QueryBuilder) (count int64, err error) {
	query, args, err := b.BuildQuery()
	if err != nil {
		return 0, err
	}
	return Count(db, query, args)
}

// Count count the number of rows of the query.
func Count(db QueryAble, query string, args []any) (count int64, err error) {
	query = fmt.Sprintf(`SELECT COUNT(1) FROM (%s) list`, query)
	err = db.QueryRow(query, args...).Scan(&count)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		query, _ := Interpolate(query, args)
		return 0, fmt.Errorf("%w: %s", err, query)
	}
	return count, nil
}

type blackhole struct{}

func (b *blackhole) Scan(_ any) error { return nil }

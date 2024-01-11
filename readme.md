Package sqls focuses only on building SQL queries by free combination
of fragments. Thus, it works naturally with all sql dialects without
having to deal with the differences between them. Unlike any other
sql builder or ORMs, Fragment is the only concept you need to learn.

## Fragment

Fragment is the builder for a part of or even the full query, it allows you
to write and combine fragments with freedom.

With the help of Fragment, we pay attention only to the reference relationships
inside the fragment, for example, use "$1" to refer the first element of s.Args.

The syntax of the fragment is exactly the same as the syntax of the "database/sql",
plus preprocessing functions support:

	SELECT * FROM foo WHERE id IN ($1, $2, $3) AND #fragment(1)
	SELECT * FROM foo WHERE id IN (?, ?, ?) AND #fragment(1)
	SELECT * FROM foo WHERE #join('#fragment', ' AND ')

## Preprocessing Functions

| name           | description                        | example                     |
| -------------- | ---------------------------------- | --------------------------- |
| c, col, column | Column by index                    | #c1, #c(1)                  |
| t, table       | Table name / alias by index        | #t1, #t(1)                  |
| f, fragment    | Fragment by index                  | #f1, #f(1)                  |
| b, builder     | Builder by index                   | #b1, #b(1)                  |
| join           | Join the template by the separator | #join('#fragment', ' AND ') |
| join           | Join from index 3 to end           | #join('#?', ',', 3)         |
| join           | Join from index 3 to 6             | #join('#?', ',', 3, 6)      |
| $              | Bindvar, usually used in #join()   | #join('#$', ', ')           |
| ?              | Bindvar, usually used in #join()   | #join('#?', ', ')           |
| global$        | global context  var by index       | #global$1                   |
| global?        | global context  var by index       | #global?1                   |

Note:
  - References in the #join template are functions, not function calls.
  - #c1 is equivalent to #c(1), which is a special syntax to call preprocessing functions when a number is the only argument.

## Examples

> See [example_test.go](./example_test.go) for more examples.

In most cases, it's easy and flexible to create your own builder  for simple queries, with a few lines of code.

<details>

```go
func Example_update() {
	update := &sqlf.Fragment{
		Prefix: "",
		Raw:    "UPDATE #t1 SET #join('#c=#$', ', ')",
	}
	where := &sqlf.Fragment{
		Prefix: "WHERE",
		Raw:    "#join('#fragment', ' AND ')",
	}
	// consider wrapping it with your own builder 
	// to provide a more friendly APIs
	builder := &sqlf.Fragment{
		Raw: "#join('#fragment', ' ')",
		Fragments: []*sqlf.Fragment{
			update,
			where,
		},
	}

	var users sqlf.Table = "users"
	update.WithTables(users)
	update.WithColumns(users.Expressions("name", "email")...)
	update.WithArgs("jebbs", "qjebbs@gmail.com")
	// append as many conditions as you want
	where.AppendFragments(&sqlf.Fragment{
		Raw:     "#c1=$1",
		Columns: users.Expressions("id"),
		Args:    []any{1},
	})

	bulit, args, err := builder.Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(bulit)
	fmt.Println(args)
	// Output:
	// UPDATE users SET name=$1, email=$2 WHERE id=$3
	// [jebbs qjebbs@gmail.com 1]
}
```
</details>

The repo also provides `*sqlb.QueryBuilder` for building complex queries.

<details>

```go
func ExampleQueryBuilder_Build() {
	var (
		foo = sqlb.NewTable("foo", "f")
		bar = sqlb.NewTable("bar", "b")
	)
	b := sqlb.NewQueryBuilder().
		Select(foo.Column("*")).
		From(foo).
		InnerJoin(bar, &sqlf.Fragment{
			Raw: "#c1=#c2",
			Columns: []*sqlf.TableColumn{
				bar.Column("foo_id"),
				foo.Column("id"),
			},
		}).
		Where(&sqlf.Fragment{
			Raw:     "(#c1=$1 OR #c2=$1)",
			Columns: foo.Columns("a", "b"),
			Args:    []any{1},
		}).
		Where2(bar.Column("c"), "=", 2)

	query, args, err := b.BindVar(syntax.Dollar).Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	fmt.Println(args)
	query, args, err = b.BindVar(syntax.Question).Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(query)
	fmt.Println(args)
	// Output:
	// SELECT f.* FROM foo AS f INNER JOIN bar AS b ON b.foo_id=f.id WHERE (f.a=$1 OR f.b=$1) AND b.c=$2
	// [1 2]
	// SELECT f.* FROM foo AS f INNER JOIN bar AS b ON b.foo_id=f.id WHERE (f.a=? OR f.b=?) AND b.c=?
	// [1 1 2]
}
```
</details>
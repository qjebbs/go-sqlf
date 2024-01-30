package sqlf

// Fa is a shortcut for creating a new Fragment with args.
func Fa(raw string, args ...any) *Fragment {
	return &Fragment{
		Raw:  raw,
		Args: args,
	}
}

// For code readability reasons, we don't provide shortcuts other than FArgs().
// `sqlf.Fa("id = $1", 1)` is familiar to most people, it's just like writing
// a regular SQL statement.
//
// Properties involved in other shortcuts, are usually referred by preprocessing
// functions, which are not familiar to most people. It's better to write following
// instead,
//
// ```go
// &sqlf.Fragment{
//     Raw: "#builder1 UNION #builder2",
//     Builders: []Builder{
//         b1,
// 	       b2,
//     },
// }
// ```
//
// so that people can easily find the definition of `builder1` in the code.

// // Fc is a shortcut for creating a new Fragment with columns.
// func Fc(raw string, columns ...*Column) *Fragment {
// 	return &Fragment{
// 		Raw:     raw,
// 		Columns: columns,
// 	}
// }

// // Ft is a shortcut for creating a new Fragment with tables.
// func Ft(raw string, tables ...Table) *Fragment {
// 	return &Fragment{
// 		Raw:    raw,
// 		Tables: tables,
// 	}
// }

// // Ff is a shortcut for creating a new Fragment with fragments.
// func Ff(raw string, fragments ...*Fragment) *Fragment {
// 	return &Fragment{
// 		Raw:       raw,
// 		Fragments: fragments,
// 	}
// }

// // Fb is a shortcut for creating a new Fragment with builders.
// func Fb(raw string, builders ...Builder) *Fragment {
// 	return &Fragment{
// 		Raw:      raw,
// 		Builders: builders,
// 	}
// }

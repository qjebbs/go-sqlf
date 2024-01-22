package sqlf

// FArgs is a shortcut for creating a new fragment with args.
func FArgs(raw string, args ...any) *Fragment {
	return &Fragment{
		Raw:  raw,
		Args: args,
	}
}

// For code readability reasons, we don't provide shortcuts other than FArgs().
// `sqlf.FArgs("id = $1", 1)` is familiar to most people, it's just like writing
// a regular SQL statement.

// Properties involved in other shortcuts, are usually referred by preprocessing
// functions, which are not familiar to most people. It's better to write instead,
//
// ```go
// &sqlf.Fragment{
//     Raw: "#builder1",
//     Builders: []Builder{
//         b,
//     },
// }
// ```
//
// so that people can easily find the definition of `builder1` in the code.

// // FColumns is a shortcut for creating a new fragment with columns.
// func FColumns(raw string, columns ...*Column) *Fragment {
// 	return &Fragment{
// 		Raw:     raw,
// 		Columns: columns,
// 	}
// }

// // FTables is a shortcut for creating a new fragment with tables.
// func FTables(raw string, tables ...Table) *Fragment {
// 	return &Fragment{
// 		Raw:    raw,
// 		Tables: tables,
// 	}
// }

// // FFragments is a shortcut for creating a new fragment with fragments.
// func FFragments(raw string, fragments ...*Fragment) *Fragment {
// 	return &Fragment{
// 		Raw:       raw,
// 		Fragments: fragments,
// 	}
// }

// // FBuilers is a shortcut for creating a new fragment with builders.
// func FBuilers(raw string, builders ...Builder) *Fragment {
// 	return &Fragment{
// 		Raw:      raw,
// 		Builders: builders,
// 	}
// }

package argstore

// Store defines argument storage for SQL building.
type Store interface {
	// Args returns the built args of the store.
	Args() []any
	// CommitArg commits an built arg to the store and returns the built bindvar.
	CommitArg(arg any) string
}

package storage

var global Store

// Init sets the global storage backend. Must be called before any
// storage operations (typically in cmd/root.go PersistentPreRunE).
func Init(s Store) { global = s }

// Get returns the global storage backend.
func Get() Store { return global }

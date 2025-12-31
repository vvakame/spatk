package testdata

// This file intentionally has compile errors for testing -allow-errors flag

// UndefinedType is not defined anywhere - this causes a compile error
// +sig
type BrokenModel struct {
	ID          string
	Name        string
	BrokenField UndefinedType // This type doesn't exist
}

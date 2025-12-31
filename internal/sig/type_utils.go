package sig

import (
	"go/types"
)

// isInvalidOrNil returns true if the type is nil or invalid (e.g., undefined type due to compile errors)
func isInvalidOrNil(typ types.Type) bool {
	if typ == nil {
		return true
	}
	if basic, ok := typ.(*types.Basic); ok && basic.Kind() == types.Invalid {
		return true
	}
	return false
}

func isBasicKindOrUnderlying(typ types.Type, kind types.BasicKind) bool {
	if basic, ok := typ.(*types.Basic); ok && basic.Kind() == kind {
		return true
	}

	if named, ok := typ.(*types.Named); ok {
		if basic, ok := named.Underlying().(*types.Basic); ok && basic.Kind() == kind {
			return true
		}
	}

	return false
}

func isTimeTime(typ types.Type) bool {
	if named, ok := typ.(*types.Named); ok {
		if named.Obj().Pkg().Path() == "time" && named.Obj().Name() == "Time" {
			return true
		}
	}

	return false
}

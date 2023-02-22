package sqb

func NoIndent() Option {
	return func(b *builder) {
		b.noIndent = true
	}
}

func InitialIndentLevel(indentLevel int) Option {
	return func(b *builder) {
		b.indentLevel = indentLevel
	}
}

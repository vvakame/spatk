package scur

func (cursor *Cursor) ClearValue() {
	for _, p := range cursor.params {
		p.Value = nil
	}
}

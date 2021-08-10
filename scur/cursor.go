package scur

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

func New() *Cursor {
	return &Cursor{}
}

type Cursor struct {
	params []*CursorParameter
}

func (cursor *Cursor) Unshift(p *CursorParameter) {
	cursor.params = append([]*CursorParameter{p}, cursor.params...)
}

func (cursor *Cursor) Push(p *CursorParameter) {
	cursor.params = append(cursor.params, p)
}

func (cursor *Cursor) Clone() *Cursor {
	newCursor := &Cursor{}
	for _, p := range cursor.params {
		newP := *p
		newCursor.params = append(newCursor.params, &newP)
	}

	return newCursor
}

func (cursor *Cursor) SetValue(obj interface{}) {
	for _, p := range cursor.params {
		p.Value = p.ToValue(obj)
	}
}

// CursorParameter is single column info for spanner cursor.
type CursorParameter struct {
	Name    string
	Order   Order
	ToValue func(obj interface{}) interface{}
	Value   interface{}
}

// EncodeParameters is encode cursor parameters value to string format.
func (cursor *Cursor) EncodeParameters() (string, error) {
 	var buf bytes.Buffer
	for idx, p := range cursor.params {
		err := encodeParameter(&buf, p.Value)
		if err != nil {
			return "", err
		}

		if idx != len(cursor.params)-1 {
			buf.WriteString(":")
		}
	}

	return buf.String(), nil
}

// DecodeCursorParameters into Cursor.
func DecodeCursorParameters(cursor *Cursor, s string) error {
	ss := strings.SplitN(s, ":", -1)
	if len(cursor.params) != len(ss) {
		return errors.New("cursor length mismatch")
	}

	for idx, p := range cursor.params {
		v, err := decodeParameter(ss[idx])
		if err != nil {
			return err
		}
		p.Value = v
	}

	return nil
}

// WhereExpression returns part of where expression about cursor.
func (cursor *Cursor) WhereExpression() (string, map[string]interface{}, error) {
	if len(cursor.params) == 0 {
		return "", nil, errors.New("cursor length is zero")
	}

	params := make(map[string]interface{})
	var buf bytes.Buffer
	buf.WriteString("(\n")
	for idx1, p1 := range cursor.params {
		paramName := fmt.Sprintf("cursor%d", idx1+1)
		params[paramName] = p1.Value

		buf.WriteString(" ")
		if len(cursor.params) != 1 {
			buf.WriteString(" (")
		}

		for idx2, p2 := range cursor.params {
			if idx1 <= idx2 {
				break
			}
			buf.WriteString(" ")
			buf.WriteString(p2.Name)
			buf.WriteString(" = ")
			buf.WriteString("@")
			buf.WriteString(fmt.Sprintf("cursor%d", idx2+1))
			buf.WriteString(" AND")
		}

		buf.WriteString(" ")
		buf.WriteString(p1.Name)
		if p1.Order == OrderAsc {
			buf.WriteString(" > ")
		} else {
			buf.WriteString(" < ")
		}
		buf.WriteString("@")
		buf.WriteString(paramName)

		if len(cursor.params) != 1 {
			buf.WriteString(" )")
		}
		buf.WriteString("\n")

		if idx1 != len(cursor.params)-1 {
			buf.WriteString("  OR\n")
		}
	}
	buf.WriteString(")\n")

	return buf.String(), params, nil
}

// OrderByExpression returns part of order by expression likes "A ASC, B DESC".
func (cursor *Cursor) OrderByExpression() (string, error) {
	if len(cursor.params) == 0 {
		return "", errors.New("cursor length is zero")
	}

	var buf bytes.Buffer
	for idx, p := range cursor.params {
		buf.WriteString(" ")
		buf.WriteString(p.Name)
		buf.WriteString(" ")
		buf.WriteString(p.Order.String())
		if idx != len(cursor.params)-1 {
			buf.WriteString(",")
		}
	}

	return buf.String(), nil
}

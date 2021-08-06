package scur

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/spanner"
)

// Order means sort order in spanner query.
type Order int

const (
	// OrderAsc represent of ASC order.
	OrderAsc Order = iota
	// OrderDesc represent of DESC order.
	OrderDesc
)

// Cursor provides query cursor structure for spanner query.
type Cursor []CursorParameter

// CursorParameter is single column info for spanner cursor.
type CursorParameter struct {
	Name  string
	Order Order
	Value interface{}
}

type cursorParameterType int

const (
	cursorParameterTypeUndefined cursorParameterType = iota
	cursorParameterTypeString
	cursorParameterTypeInt
	cursorParameterTypeInt64
	cursorParameterTypeBool
	cursorParameterTypeFloat64
	cursorParameterTypeTime
)

func (order Order) String() string {
	switch order {
	case OrderAsc:
		return "ASC"
	case OrderDesc:
		return "DESC"
	default:
		return strconv.Itoa(int(order))
	}
}

// DecodeCursorParameters into Cursor.
func DecodeCursorParameters(cursor Cursor, s string) error {
	ss := strings.SplitN(s, ":", -1)
	if len(cursor) != len(ss) {
		return errors.New("cursor length mismatch")
	}

	for idx := range cursor {
		err := cursor.elemDecode(idx, ss[idx])
		if err != nil {
			return err
		}
	}

	return nil
}

func (cursor Cursor) elemDecode(idx int, s string) error {
	if len(s) < 2 {
		return fmt.Errorf("unexpected encoded value: %s", s)
	}
	n, err := strconv.ParseInt(s[0:2], 16, 64)
	if err != nil {
		return err
	}
	sv := s[2:]
	switch cursorParameterType(n) {
	case cursorParameterTypeString:
		b, err := base64.RawURLEncoding.DecodeString(sv)
		if err != nil {
			return err
		}
		cursor[idx].Value = string(b)
	case cursorParameterTypeInt:
		v, err := strconv.ParseInt(sv, 36, 64)
		if err != nil {
			return fmt.Errorf("int value decode failed: %w", err)
		}
		cursor[idx].Value = int(v)
	case cursorParameterTypeInt64:
		v, err := strconv.ParseInt(sv, 36, 64)
		if err != nil {
			return fmt.Errorf("int64 value decode failed: %w", err)
		}
		cursor[idx].Value = v
	case cursorParameterTypeBool:
		v, err := strconv.ParseInt(sv, 10, 64)
		if err != nil {
			return fmt.Errorf("bool value decode failed: %w", err)
		}
		cursor[idx].Value = v == 1
	case cursorParameterTypeFloat64:
		v, err := strconv.ParseUint(sv, 36, 64)
		if err != nil {
			return fmt.Errorf("float64 value decode failed: %w", err)
		}
		cursor[idx].Value = math.Float64frombits(v)
	case cursorParameterTypeTime:
		v, err := strconv.ParseInt(sv, 36, 64)
		if err != nil {
			return fmt.Errorf("time.Time value decode failed: %w", err)
		}
		cursor[idx].Value = time.Unix(0, v)
	default:
		return fmt.Errorf("unsupported cursor value type: %d", n)
	}

	return nil
}

// EncodeParameters is encode cursor parameters value to string format.
func (cursor Cursor) EncodeParameters() (string, error) {
	var buf bytes.Buffer
	for idx, cc := range cursor {
		err := cursor.elemEncode(&buf, cc.Value)
		if err != nil {
			return "", err
		}

		if idx != len(cursor)-1 {
			buf.WriteString(":")
		}
	}

	return buf.String(), nil
}

func (cursor Cursor) elemEncode(buf *bytes.Buffer, v interface{}) error {
	// NOTE : をdelimiterに使うのでこれが出現しないようにしないようにやる必要がある
	switch v := v.(type) {
	case string:
		buf.WriteString(fmt.Sprintf("%02x", cursorParameterTypeString))
		buf.WriteString(base64.RawURLEncoding.EncodeToString([]byte(v)))
	case int:
		buf.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeInt, strconv.FormatInt(int64(v), 36)))
	case int64:
		buf.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeInt64, strconv.FormatInt(v, 36)))
	case bool:
		var nv int
		if v {
			nv = 1
		} else {
			nv = 0
		}
		buf.WriteString(fmt.Sprintf("%02x%d", cursorParameterTypeBool, nv))
	case float64:
		buf.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeFloat64, strconv.FormatUint(math.Float64bits(v), 36)))
	case time.Time:
		buf.WriteString(fmt.Sprintf("%02x%s", cursorParameterTypeTime, strconv.FormatInt(v.UnixNano(), 36)))
	case spanner.Encoder:
		v2, err := v.EncodeSpanner()
		if err != nil {
			return err
		}
		err = cursor.elemEncode(buf, v2)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported cursor value type: %T", v)
	}

	return nil
}

// WhereExpression returns part of where expression about cursor.
func (cursor Cursor) WhereExpression() (string, map[string]interface{}, error) {
	if len(cursor) == 0 {
		return "", nil, errors.New("cursor length is zero")
	}

	params := make(map[string]interface{})
	var buf bytes.Buffer
	buf.WriteString("(\n")
	for idx1, cc1 := range cursor {
		paramName := fmt.Sprintf("cursor%d", idx1+1)
		params[paramName] = cc1.Value

		buf.WriteString(" ")
		if len(cursor) != 1 {
			buf.WriteString(" (")
		}

		for idx2, cc2 := range cursor {
			if idx1 <= idx2 {
				break
			}
			buf.WriteString(" ")
			buf.WriteString(cc2.Name)
			buf.WriteString(" = ")
			buf.WriteString("@")
			buf.WriteString(fmt.Sprintf("cursor%d", idx2+1))
			buf.WriteString(" AND")
		}

		buf.WriteString(" ")
		buf.WriteString(cc1.Name)
		if cc1.Order == OrderAsc {
			buf.WriteString(" > ")
		} else {
			buf.WriteString(" < ")
		}
		buf.WriteString("@")
		buf.WriteString(paramName)

		if len(cursor) != 1 {
			buf.WriteString(" )")
		}
		buf.WriteString("\n")

		if idx1 != len(cursor)-1 {
			buf.WriteString("  OR\n")
		}
	}
	buf.WriteString(")\n")

	return buf.String(), params, nil
}

// OrderByExpression returns part of order by expression likes "A ASC, B DESC".
func (cursor Cursor) OrderByExpression() (string, error) {
	if len(cursor) == 0 {
		return "", errors.New("cursor length is zero")
	}

	var buf bytes.Buffer
	for idx, cc := range cursor {
		buf.WriteString(" ")
		buf.WriteString(cc.Name)
		buf.WriteString(" ")
		buf.WriteString(cc.Order.String())
		if idx != len(cursor)-1 {
			buf.WriteString(",")
		}
	}

	return buf.String(), nil
}

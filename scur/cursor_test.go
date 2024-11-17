package scur_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/MakeNowJust/heredoc/v2"
	"github.com/vvakame/spatk/scur"
)

var _ spanner.Encoder = FooID("")
var _ spanner.Decoder = (*FooID)(nil)

type FooID string

func (id FooID) EncodeSpanner() (any, error) {
	return string(id), nil
}

func (id *FooID) DecodeSpanner(input any) error {
	s, ok := input.(string)
	if !ok {
		return fmt.Errorf("unexpected id type: %T", input)
	}

	*id = FooID(s)
	return nil
}

func TestCursor_EncodeDecodeParameters(t *testing.T) {
	tests := []struct {
		name     string
		cursor   *scur.Cursor
		expected *scur.Cursor // nil の場合 cursor と同じものにする
	}{
		{
			name: "1 column",
			cursor: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "String",
					Order: scur.OrderAsc,
					Value: "foobar",
				})
				return cursor
			}(),
		},
		{
			name: "3 columns",
			cursor: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "String",
					Order: scur.OrderAsc,
					Value: "foobar",
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Int",
					Order: scur.OrderAsc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Time",
					Order: scur.OrderAsc,
					Value: time.Now().Round(0),
				})
				return cursor
			}(),
		},
		{
			name: "5 columns",
			cursor: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "String",
					Order: scur.OrderAsc,
					Value: "foobar",
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Int",
					Order: scur.OrderAsc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Time",
					Order: scur.OrderAsc,
					Value: time.Now().Round(0),
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Bool",
					Order: scur.OrderAsc,
					Value: true,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Float",
					Order: scur.OrderAsc,
					Value: 1.25,
				})
				return cursor
			}(),
		},
		{
			name: "all types",
			cursor: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "String",
					Order: scur.OrderAsc,
					Value: "foobar",
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Int",
					Order: scur.OrderAsc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Int64",
					Order: scur.OrderAsc,
					Value: int64(1),
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Bool",
					Order: scur.OrderAsc,
					Value: true,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Float",
					Order: scur.OrderAsc,
					Value: 1.25,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "Time",
					Order: scur.OrderAsc,
					Value: time.Now().Round(0),
				})
				return cursor
			}(),
		},
		{
			name: "complex string",
			cursor: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "String",
					Order: scur.OrderAsc,
					Value: "https://godoc.org/cloud.google.com/go/spanner :::",
				})
				return cursor
			}(),
		},
		{
			name: "custom type value",
			cursor: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "CustomType",
					Order: scur.OrderAsc,
					Value: FooID("aaa"),
				})
				return cursor
			}(),
			expected: func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "CustomType",
					Order: scur.OrderAsc,
					Value: "aaa",
				})
				return cursor
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := tt.cursor.EncodeParameters()
			if err != nil {
				t.Fatal(err)
			}

			t.Logf("encoded: %d %s", len(s), s)

			actual := tt.cursor.Clone()
			actual.ClearValue()

			err = scur.DecodeCursorParameters(actual, s)
			if err != nil {
				t.Fatal(err)
			}

			expected := tt.expected
			if expected == nil {
				expected = tt.cursor
			}

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("unexpected expected = %v, decoded %v", expected, actual)
			}
		})
	}
}

func TestCursor_WhereExpression(t *testing.T) {
	tests := []struct {
		name   string
		cursor *scur.Cursor
		sql    string
		params map[string]any
	}{
		{
			"1 column",
			func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "A",
					Order: scur.OrderAsc,
					Value: 1,
				})
				return cursor
			}(),
			heredoc.Doc(`
				(
				  A > @cursor1
				)
			`),
			map[string]any{
				"cursor1": 1,
			},
		},
		{
			"2 columns",
			func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "A",
					Order: scur.OrderAsc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "B",
					Order: scur.OrderDesc,
					Value: "b",
				})
				return cursor
			}(),
			heredoc.Doc(`
				(
				  ( A > @cursor1 )
				  OR
				  ( A = @cursor1 AND B < @cursor2 )
				)
			`),
			map[string]any{
				"cursor1": 1,
				"cursor2": "b",
			},
		},
		{
			"3 columns",
			func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "A",
					Order: scur.OrderDesc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "B",
					Order: scur.OrderAsc,
					Value: "b",
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "C",
					Order: scur.OrderAsc,
					Value: 1.25,
				})
				return cursor
			}(),
			heredoc.Doc(`
				(
				  ( A < @cursor1 )
				  OR
				  ( A = @cursor1 AND B > @cursor2 )
				  OR
				  ( A = @cursor1 AND B = @cursor2 AND C > @cursor3 )
				)
			`),
			map[string]any{
				"cursor1": 1,
				"cursor2": "b",
				"cursor3": 1.25,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, params, err := tt.cursor.WhereExpression()
			if err != nil {
				t.Fatal(err)
			}
			if sql != tt.sql {
				t.Errorf("WhereExpression() sql = %v, want %v", sql, tt.sql)
			}
			if !reflect.DeepEqual(params, tt.params) {
				t.Errorf("WhereExpression() params = %v, want %v", params, tt.params)
			}
		})
	}
}

func TestCursor_OrderByExpression(t *testing.T) {
	tests := []struct {
		name   string
		cursor *scur.Cursor
		sql    string
	}{
		{
			"1 column",
			func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "A",
					Order: scur.OrderAsc,
					Value: 1,
				})
				return cursor
			}(),
			" A ASC",
		},
		{
			"2 columns",
			func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "A",
					Order: scur.OrderAsc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "B",
					Order: scur.OrderDesc,
					Value: "b",
				})
				return cursor
			}(),
			" A ASC, B DESC",
		},
		{
			"3 columns",
			func() *scur.Cursor {
				cursor := scur.New()
				cursor.Push(&scur.CursorParameter{
					Name:  "A",
					Order: scur.OrderDesc,
					Value: 1,
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "B",
					Order: scur.OrderAsc,
					Value: "b",
				})
				cursor.Push(&scur.CursorParameter{
					Name:  "C",
					Order: scur.OrderAsc,
					Value: 1.25,
				})
				return cursor
			}(),
			" A DESC, B ASC, C ASC",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, err := tt.cursor.OrderByExpression()
			if err != nil {
				t.Fatal(err)
			}
			if sql != tt.sql {
				t.Errorf("OrderByExpression() sql = %v, want %v", sql, tt.sql)
			}
		})
	}
}

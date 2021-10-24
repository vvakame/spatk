package sigtest

import (
	"reflect"
	"testing"
)

func TestGeneratedModel(t *testing.T) {
	t.Run("ModelA", func(t *testing.T) {
		t.Run("table name", func(t *testing.T) {
			if v := spannerInfoModelA.TableName(); v != "ModelA" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.As("A").TableName(); v != "ModelA AS A" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.As("A").ForceIndex("IndexA").TableName(); v != "ModelA@{FORCE_INDEX=IndexA} AS A" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.TableName(); v != "ModelA" {
				t.Errorf("unexpected: %v", v)
			}
		})
		t.Run("column name", func(t *testing.T) {
			if v := spannerInfoModelA.ID(); v != "ModelAID" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.IDAs("ID").ID(); v != "ModelAID AS ID" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.ID(); v != "ModelAID" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.IDAs("A").NameAs("B").ID(); v != "ModelAID AS A" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.IDAs("A").NameAs("B").UpdatedAtAs("C").CreatedAtAs("D").ColumnNames(); !reflect.DeepEqual(v, []string{"ModelAID AS A", "Name AS B", "UpdatedAt AS C", "CreatedAt AS D"}) {
				t.Errorf("unexpected: %v", v)
			}
		})
		t.Run("complex pattern", func(t *testing.T) {
			if v := spannerInfoModelA.IDAs("ID").As("TABLE").ID(); v != "TABLE.ModelAID AS ID" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.As("TABLE").IDAs("I").TableName(); v != "ModelA AS TABLE" {
				t.Errorf("unexpected: %v", v)
			}
			if v := spannerInfoModelA.As("TABLE").IDAs("I").ColumnNames()[0]; v != "TABLE.ModelAID AS I" {
				t.Errorf("unexpected: %v", v)
			}
		})
	})
	t.Run("ModelB", func(t *testing.T) {
		t.Run("table name", func(t *testing.T) {
			if v := spannerInfoModelB.TableName(); v != "ModelB" {
				t.Errorf("unexpected: %v", v)
			}
		})
	})
}

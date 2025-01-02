//go:generate go run github.com/vvakame/spatk/cmd/sig -private -output model_gen.go .

package sigtest

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/vvakame/spatk/scur"
)

type ModelAID string
type ModelBID string

// +sig
type ModelA struct {
	ID        ModelAID `spanner:"ModelAID"`
	Name      string
	UpdatedAt time.Time
	CreatedAt time.Time
}

// +sig
type ModelBar struct {
	TableName string   `spanner:"-" sig:"table=ModelB"`
	ID        ModelBID `spanner:"ModelBID"`
	Name      string
	UpdatedAt time.Time `sig:"min=TimestampMinValue,max=TimestampMaxValue"`
	CreatedAt time.Time `sig:"minmax=Timestamp"`
}

func TimestampMinValue() time.Time {
	return scur.TimestampMinValue()
}

func TimestampMaxValue() time.Time {
	return scur.TimestampMaxValue()
}

type Time time.Time

// +sig
type ModelC struct {
	TableName    string     `spanner:"-" sig:"table=ModelC"`
	ID           string     `spanner:"ModelCID"`
	OwnTimeType  Time       ``
	UUID         uuid.UUID  ``
	LocalType1   localType  ``
	LocalType2   *localType ``
	IgnoreColumn string     `spanner:"-"`
}

// TimeSpannerMinValue is a helper function to get min value of Time type.
// "Time" + "SpannerMinValue" will be used for min value.
func TimeSpannerMinValue() Time {
	return Time(scur.TimestampMinValue())
}

// TimeSpannerMaxValue is a helper function to get max value of Time type.
// "Time" + "SpannerMaxValue" will be used for max value.
func TimeSpannerMaxValue() Time {
	return Time(scur.TimestampMaxValue())
}

// UUIDSpannerMinValue is a helper function to get min value of uuid.UUID type.
// If type is coming from external package, drop that package name and use only type name.
func UUIDSpannerMinValue() uuid.UUID {
	return uuid.Nil
}

// UUIDSpannerMaxValue is a helper function to get max value of uuid.UUID type.
// If type is coming from external package, drop that package name and use only type name.
func UUIDSpannerMaxValue() uuid.UUID {
	return uuid.Max
}

type localType struct {
	Int int
}

func localTypeSpannerMinValue() localType {
	return localType{
		Int: math.MinInt,
	}
}

func localTypeSpannerMaxValue() localType {
	return localType{
		Int: math.MaxInt,
	}
}

func localTypePointerSpannerMinValue() *localType {
	return &localType{
		Int: math.MinInt,
	}
}

func localTypePointerSpannerMaxValue() *localType {
	return &localType{
		Int: math.MaxInt,
	}
}

// +sig
type ModelD struct {
	ID        string `spanner:"ModelDID"`
	IntValue  int    ``
	BoolValue bool   ``
}

//go:generate go run github.com/vvakame/spatk/cmd/sig -private -output model_gen.go .

package sigtest

import "time"

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
	UpdatedAt time.Time
	CreatedAt time.Time
}

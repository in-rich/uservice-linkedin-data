package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type Company struct {
	bun.BaseModel `bun:"table:companies"`

	ID               *uuid.UUID `bun:"id,pk,type:uuid"`
	PublicIdentifier string     `bun:"public_identifier,unique,notnull"`

	Name string `bun:"name"`

	UpdatedAt *time.Time `bun:"updated_at"`
}

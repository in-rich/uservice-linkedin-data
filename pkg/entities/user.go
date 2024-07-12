package entities

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID               *uuid.UUID `bun:"id,pk,type:uuid"`
	PublicIdentifier string     `bun:"public_identifier,unique,notnull"`

	FirstName string `bun:"first_name"`
	LastName  string `bun:"last_name"`

	UpdatedAt *time.Time `bun:"updated_at"`
}

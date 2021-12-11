package model

import "github.com/google/uuid"

type firstShared struct {
	FirstShared bool `db:"firstshared"`
}

type heldBy struct {
	HeldBy uuid.UUID `db:"heldBy"`
}

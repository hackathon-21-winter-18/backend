package model

import "github.com/google/uuid"

type firstShared struct {
	FirstShared bool `db:"firstshared"`
}

type heldBy struct {
	heldBy uuid.UUID `db:"heldBy"`
}

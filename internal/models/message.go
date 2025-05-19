package models

import (
	"time"

	"github.com/google/uuid"
)




type Message struct{
	Message string
	ID uuid.UUID
	Sender string
	Receiver string
	IsRead bool
	CreatedAt time.Time
}
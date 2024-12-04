// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Problem struct {
	ID              uuid.UUID
	Title           string
	DescriptionPath string
	TestcasesPath   string
	Tags            []string
	TimeLimit       float64
	MemoryLimit     float64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

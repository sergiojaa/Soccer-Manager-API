package domain

import "time"

type Player struct {
	ID        int64
	TeamID    int64
	FirstName string
	LastName  string
	Country   string
	Age       int
	Position  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PlayerInput struct {
	TeamID    int64
	FirstName string
	LastName  string
	Country   string
	Age       int
	Position  string
}

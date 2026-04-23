package domain

import "time"

type Team struct {
	ID        int64
	UserID    int64
	Name      string
	Country   string
	Budget    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

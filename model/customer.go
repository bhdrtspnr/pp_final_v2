package model

import (
	"time"
)

type Customer struct {
	Id                            int       `json:"id"`
	Name                          string    `json:"name"`
	Balance                       float64   `json:"balance"`
	Consecutive_discount          int       `json:"consecutive_discount"`
	Has_subsequent_discount_until time.Time `json:"has_subsequent_discount_until"`
}

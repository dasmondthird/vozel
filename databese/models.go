package database

import (
	"time"
)

type User struct {
	ID           int64 `gorm:"primaryKey"`
	Username     string
	Registered   bool
	VPNKey       string
	Location     string
	Expiry       time.Time
	Balance      int
	RegisteredAt time.Time
	Config       string
}

type Tariff struct {
	ID       int64 `gorm:"primaryKey"`
	Name     string
	Price    int
	Duration time.Duration
}

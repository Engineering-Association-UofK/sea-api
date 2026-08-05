package models

import "time"

type RateLimitEndpoints string

const (
	LimitSendCode       RateLimitEndpoints = "send_code"
	LimitUpdateUsername RateLimitEndpoints = "update_username"
	LimitUpdateEmail    RateLimitEndpoints = "update_email"
)

var LimitMultipliers = map[RateLimitEndpoints]int{
	LimitSendCode:       1,
	LimitUpdateUsername: 10080,
	LimitUpdateEmail:    10,
}

type RateLimitModel struct {
	ID           string             `db:"id"`
	IpAddress    string             `db:"ip_address"`
	RequestCount int                `db:"request_count"`
	Endpoint     RateLimitEndpoints `db:"endpoint"`
	LastRequest  time.Time          `db:"last_request"`
}

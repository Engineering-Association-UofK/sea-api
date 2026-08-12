package models

import "time"

type Log struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	Action    string    `json:"action" db:"action"`
	TableName string    `json:"table_name" db:"table_name"`
	ObjectID  int64     `json:"object_id" db:"object_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

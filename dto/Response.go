package dto

import "time"

type Response struct {
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
	Reason    string    `json:"reason,omitempty"`
}

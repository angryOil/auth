package req

import "time"

type Create struct {
	Email     string
	Password  string
	Role      []string
	CreatedAt time.Time
}

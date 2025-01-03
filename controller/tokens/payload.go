package tokens

import "time"

type Payload struct {
	Username  string        `json:"username"`
	ExpiredAt time.Duration `json:"expired_at"`
}

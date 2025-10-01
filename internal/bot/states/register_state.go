package states

import "time"

type RegisterState struct {
	Name      string
	Birthdate time.Time
}

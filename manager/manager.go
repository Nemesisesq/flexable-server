package manager

import "github.com/nemesisesq/flexable/user"

type Manager struct {
	User user.User `json:"user"`
}

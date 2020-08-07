package user

import (
	"errors"
	"time"
)

type User struct {
	ID        int64 `json:"user,string"`
	Username  string
	CreatedAt time.Time
}

type Users interface {
	Find(id int64) (*User, error)
	Create(u *User) (int64, error)
}

func (u *User) ValidCreate() error {
	if u.Username == "" {
		return errors.New("bad username")
	}

	return nil
}

func (u *User) ValidGet() error {
	if u.ID == 0 {
		return errors.New("bad user id")
	}

	return nil
}

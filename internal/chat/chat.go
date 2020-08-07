package chat

import (
	"time"

	"github.com/pkg/errors"
)

// Chat ...
type Chat struct {
	ID        int64
	Name      string
	Users     []string
	CreatedAt time.Time
}

type Chats interface {
	Find(userID int64) ([]*Chat, error)
	Create(m *Chat) (int64, error)
}

func (c *Chat) ValidCreate() error {
	switch {
	case c.Name == "":
		return errors.New("bad chat name")
	case c.Users == nil:
		return errors.New("bad members")
	}

	return nil
}

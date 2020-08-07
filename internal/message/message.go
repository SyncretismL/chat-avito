package message

import (
	"errors"
	"time"
)

// Message ...
type Message struct {
	ID        int64 `json:"message_id,string"`
	Chat      int64 `json:"chat,string"`
	Author    int64 `json:"author,string"`
	Text      string
	CreatedAt time.Time
}

type Messages interface {
	Find(chatID int64) ([]*Message, error)
	Create(m *Message) (int64, error)
}

func (m *Message) ValidCreate() error {
	switch {
	case m.Chat == 0:
		return errors.New("bad chat id")
	case m.Author == 0:
		return errors.New("bad author id")
	case m.Text == "":
		return errors.New("bad text")
	}

	return nil
}

func (m *Message) ValidGet() error {
	if m.Chat == 0 {
		return errors.New("bad chat id")
	}

	return nil
}

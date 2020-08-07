package postgres

import (
	"chat/internal/message"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

var _ message.Messages = &MsgStorage{}

// RobotStorage ...
type MsgStorage struct {
	statementStorage

	createStmt *sql.Stmt
	findStmt   *sql.Stmt
}

// NewRobotStorage ...
func NewMsgStorage(db *DB) (*MsgStorage, error) {
	s := &MsgStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: createMsgQuery, Dst: &s.createStmt},
		{Query: findMsgQuery, Dst: &s.findStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const msgFields = "chat, author, text, created_at"

const createMsgQuery = "INSERT INTO public.messages (" + msgFields + ") VALUES ((SELECT chat_id FROM public.users_chats WHERE chat_id = $1 AND user_id = $2), $2, $3, now()) RETURNING id"

// Create ...
func (s *MsgStorage) Create(m *message.Message) (int64, error) {
	if err := s.createStmt.QueryRow(&m.Chat, &m.Author, &m.Text).Scan(&m.ID); err != nil {
		msg := fmt.Sprintf("can not exec query with msgID %v", m.ID)
		return 0, errors.WithMessage(err, msg)
	}

	return m.ID, nil
}

const findMsgQuery = "SELECT id, " + msgFields + " FROM public.messages WHERE chat=$1 ORDER BY created_at ASC"

// Find ...
func (s *MsgStorage) Find(id int64) ([]*message.Message, error) {
	var ms []*message.Message

	rows, err := s.findStmt.Query(id)
	if err != nil {
		msg := fmt.Sprintf("can't scan msgs with chat id %v", id)
		return nil, errors.WithMessage(err, msg)
	}

	defer rows.Close()

	for rows.Next() {
		var m message.Message

		if err := scanMsg(rows, &m); err != nil {
			msg := fmt.Sprintf("failed to scan msgs with chat id %v", id)
			return nil, errors.WithMessage(err, msg)
		}

		ms = append(ms, &m)
	}

	return ms, nil
}

func scanMsg(scanner sqlScanner, m *message.Message) error {
	return scanner.Scan(&m.ID, &m.Chat, &m.Author, &m.Text, &m.CreatedAt)
}

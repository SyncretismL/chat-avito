package postgres

import (
	"chat/internal/chat"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lib/pq"

	"github.com/pkg/errors"
)

var _ chat.Chats = &ChatStorage{}

// RobotStorage ...
type ChatStorage struct {
	statementStorage

	createStmt   *sql.Stmt
	createUCStmt *sql.Stmt
	findStmt     *sql.Stmt
}

// NewRobotStorage ...
func NewChatStorage(db *DB) (*ChatStorage, error) {
	s := &ChatStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: createChatQuery, Dst: &s.createStmt},
		{Query: findChatQuery, Dst: &s.findStmt},
		{Query: createChatUserQuery, Dst: &s.createUCStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const chatFields = "users, name, created_at"

const createChatQuery = "INSERT INTO public.chats (name, created_at) VALUES ($1, now()) RETURNING id"
const createChatUserQuery = "INSERT INTO public.users_chats(chat_id, user_id) SELECT $1, unnest($2::integer[])"

func (s *ChatStorage) Create(c *chat.Chat) (int64, error) {
	tx, err := s.db.Session.Begin()
	if err != nil {
		return 0, err
	}

	{
		stmt := tx.Stmt(s.createStmt)

		defer stmt.Close()

		if err := stmt.QueryRow(&c.Name).Scan(&c.ID); err != nil {
			tx.Rollback()
			msg := fmt.Sprintf("can not exec query with chatID %v", c.ID)
			return 0, errors.WithMessage(err, msg)
		}
	}

	{
		stmt := tx.Stmt(s.createUCStmt)
		valueArgs := []int64{}

		defer stmt.Close()

		for _, v := range c.Users {
			num, _ := strconv.ParseInt(v, 10, 64)
			valueArgs = append(valueArgs, num)
		}

		if _, err := stmt.Exec(c.ID, pq.Array(valueArgs)); err != nil {
			tx.Rollback()
			msg := fmt.Sprintf("failed to exec users with chat id %v", c.ID)
			return 0, errors.WithMessage(err, msg)
		}
	}

	return c.ID, tx.Commit()
}

const findChatQuery = "SELECT chats.id, array_agg(users_chats.user_id), chats.name, chats.created_at FROM public.chats INNER JOIN users_chats ON (users_chats.chat_id = chats.id)" +
	"GROUP BY chats.id HAVING $1=ANY(array_agg(users_chats.user_id)) ORDER BY (SELECT MAX(created_at) FROM public.messages WHERE messages.chat = chats.id) DESC"

func (s *ChatStorage) Find(id int64) ([]*chat.Chat, error) {
	var chs []*chat.Chat

	rows, err := s.findStmt.Query(id)
	if err != nil {
		msg := fmt.Sprintf("can't scan chats with user id %v", id)
		return nil, errors.WithMessage(err, msg)
	}

	defer rows.Close()

	for rows.Next() {
		var c chat.Chat

		if err := scanChat(rows, &c); err != nil {
			msg := fmt.Sprintf("failed to scan msgs with chat id %v", id)
			return nil, errors.WithMessage(err, msg)
		}

		chs = append(chs, &c)
	}

	return chs, nil
}

func scanChat(scanner sqlScanner, c *chat.Chat) error {
	return scanner.Scan(&c.ID, pq.Array(&c.Users), &c.Name, &c.CreatedAt)
}

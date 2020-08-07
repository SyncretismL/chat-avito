package postgres

import (
	"chat/internal/user"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

var _ user.Users = &UserStorage{}

// RobotStorage ...
type UserStorage struct {
	statementStorage

	createStmt *sql.Stmt
	findStmt   *sql.Stmt
}

// NewRobotStorage ...
func NewUserStorage(db *DB) (*UserStorage, error) {
	s := &UserStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: createUserQuery, Dst: &s.createStmt},
		{Query: findUserQuery, Dst: &s.findStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const userFields = "username, created_at"

const createUserQuery = "INSERT INTO public.users (" + userFields + ") VALUES ($1, now()) RETURNING id"

// Create ...
func (s *UserStorage) Create(u *user.User) (int64, error) {
	if err := s.createStmt.QueryRow(&u.Username).Scan(&u.ID); err != nil {
		msg := fmt.Sprintf("can not exec query with userID %v", u.ID)
		return 0, errors.WithMessage(err, msg)
	}

	return u.ID, nil
}

const findUserQuery = "SELECT id, " + userFields + " FROM public.users WHERE id=$1"

// Find ...
func (s *UserStorage) Find(id int64) (*user.User, error) {
	var u user.User

	row := s.findStmt.QueryRow(id)
	if err := scanUser(row, &u); err != nil {
		msg := fmt.Sprintf("can't scan user with id %v", id)
		return nil, errors.WithMessage(err, msg)
	}

	return &u, nil
}

func scanUser(scanner sqlScanner, u *user.User) error {
	return scanner.Scan(&u.ID, &u.Username, &u.CreatedAt)
}

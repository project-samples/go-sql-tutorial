package services

import (
	"context"
	"database/sql"
	"fmt"
	s "github.com/core-go/sql"
	"reflect"
	"strings"

	. "go-service/internal/models"
)

type SqlUserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *SqlUserService {
	return &SqlUserService{DB: db}
}


func (m *SqlUserService) GetAll(ctx context.Context) (*[]User, error) {
	query := "select id, username, email, phone, date_of_birth from users"
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var result []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth)
		result = append(result, user)
	}
	return &result, nil
}

func (m *SqlUserService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	query := "select id, username, email, phone, date_of_birth from users where id = ?"
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.DateOfBirth)
	if err != nil {
		errMsg := err.Error()
		if strings.Compare(fmt.Sprintf(errMsg), "0 row(s) returned") == 0 {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (m *SqlUserService) Insert(ctx context.Context, user *User) (int64, error) {
	query := "insert into users (id, username, email, phone, date_of_birth) values (?, ?, ?, ?, ?)"
	stmt, er0 := m.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	result, er1 := stmt.ExecContext(ctx, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth)
	if er1 != nil {
		return -1, nil
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Update(ctx context.Context, user *User) (int64, error) {
	query := "update users set username = ?, email = ?, phone = ?, date_of_birth = ? where id = ?"
	stmt, er0 := m.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	result, er1 := stmt.ExecContext(ctx, user.Username, user.Email, user.Phone, user.DateOfBirth, user.Id)
	if er1 != nil {
		return -1, er1
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	userType := reflect.TypeOf(User{})
	jsonColumnMap := s.MakeJsonColumnMap(userType)
	colMap := s.JSONToColumns(user, jsonColumnMap)
	keys, _ := s.FindPrimaryKeys(userType)
	query, args := s.BuildToPatch("users", colMap, keys, s.BuildParam)
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = ?"
	stmt, er0 := m.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	result, er1 := stmt.ExecContext(ctx, id)
	if er1 != nil {
		return -1, er1
	}
	rowAffect, er2 := result.RowsAffected()
	if er2 != nil {
		return 0, er2
	}
	return rowAffect, nil
}

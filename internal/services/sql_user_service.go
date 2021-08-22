package services

import (
	"context"
	"database/sql"
	"fmt"
	s "github.com/core-go/sql"
	"log"
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
	var result []User
	query := "select id, username, phone, email, url, locked, date_of_birth from users"
	err := s.Query(ctx, m.DB, &result, query)
	return &result, err
	/*

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.Locked, &user.DateOfBirth)
		result = append(result, user)
	}
	return &result, nil
	 */
}

func (m *SqlUserService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	query := "select id, username, email, phone, url, active, locked, date_of_birth from users where id = ?"
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
	// stm, args := s.BuildInsertSql("users", user, 0, s.BuildParam)
	stm, args, _ := s.BuildToSave(m.DB,"users", user)
	log.Print(fmt.Sprintf(stm, args))
	result, er1 := m.DB.ExecContext(ctx, stm, args...)
	/*
	query := "insert into users (id, username, email, phone, date_of_birth) values (?, ?, ?, ?, ?)"
	stmt, er0 := m.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	result, er1 := stmt.ExecContext(ctx, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth)
	*/
	if er1 != nil {
		return -1, nil
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Update(ctx context.Context, user *User) (int64, error) {
	stm, args := s.BuildToUpdate("users", user, 0, s.BuildParam)
	result, er1 := m.DB.ExecContext(ctx, stm, args...)
	/*
	query := "update users set username = ?, email = ?, phone = ?, date_of_birth = ? where id = ?"
	stmt, er0 := m.DB.Prepare(query)
	if er0 != nil {
		return -1, nil
	}
	result, er1 := stmt.ExecContext(ctx, user.Username, user.Email, user.Phone, user.DateOfBirth, user.Id)
	 */
	if er1 != nil {
		return -1, er1
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Patch(ctx context.Context, user map[string]interface{}) (int64, error) {
	userType := reflect.TypeOf(User{})
	result, err := s.Patch(ctx, m.DB, "users", user, userType)
	if err != nil {
		return result, err
	}
	return result, nil
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
func (m *SqlUserService) Batch(ctx context.Context, user *[]User) (int64, error) {
	return s.SaveBatch(ctx, m.DB, "users", user)
}

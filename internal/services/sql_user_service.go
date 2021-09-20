package services

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	. "go-service/internal/models"
	sq "go-service/sql"
	"log"
	"reflect"
	"strings"
)

type SqlUserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *SqlUserService {
	return &SqlUserService{DB: db}
}

func (m *SqlUserService) GetAll(ctx context.Context) (*[]User, error) {
	query := "select id, username, email, phone, date_of_birth, interests, skills, achievements, settings from users"
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	var res []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth)
		if err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	return &res, nil
}

func (m *SqlUserService) Load(ctx context.Context, id string) (*User, error) {
	var user User
	query := "select id, username, email, phone, date_of_birth, interests, skills, achievements, settings from users WHERE ID = $1"
	row,err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	for row.Next(){
		err = row.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth)
	}
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
	test:= *user
	query, values, err := BuildToSaveOracle("users", test, buildParam)
	if err != nil {
		return 0, err
	}
	log.Println(query)
	log.Println(values)
	//query := "insert into users (id, username, email, phone, date_of_birth, interests, skills, achievements, settings) values ($1, $2, $3, $4, $5, $6, $7, $8, $9)"
	//result,err := m.DB.Exec(query, user.Id, user.Version, user.Username, user.Email, user.Phone, user.DateOfBirth, pq.Array(user.Interests), pq.Array(user.Skills), pq.Array(user.Achievements), user.Settings)
	result,err := m.DB.Exec(query, values...)
	if err != nil {
		log.Println(err)
		return -1, nil
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Update(ctx context.Context, user *User) (int64, error) {
	test:= *user
	query, values := sq.BuildToUpdateWithVersion("users", test, 1, buildParam, pq.Array, true)
	log.Println(query)
	log.Println(values)
	//query := "update users set username = $2, email = $3, phone = $4, date_of_birth = $5, interests = $6, skills = $7, achievements = $8, settings = $9 where id = $1"
	//result, err := m.DB.Exec(query, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth, pq.Array(user.Interests), pq.Array(user.Skills), pq.Array(user.Achievements), user.Settings)
	result,err := m.DB.Exec(query, values...)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (m *SqlUserService) Delete(ctx context.Context, id string) (int64, error) {
	query := "delete from users where id = :1"
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return -1, err
	}
	rowAffect, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	fmt.Printf("Total rows/records has been deleted %v", rowAffect)
	return rowAffect, nil
}

func buildParam(s int) string {
	return fmt.Sprintf(`:%d`,s)
}

func BuildToSaveOracle(table string, model interface{}, buildParam func(int) string, options...bool) (string, []interface{}, error) {
	modelType := reflect.Indirect(reflect.ValueOf(model)).Type()
	mv := reflect.Indirect(reflect.ValueOf(model))
	cols, keys, schema := sq.MakeSchema(modelType)
	variables := make([]string, 0)
	var setColumns []string
	uniqueCols := make([]string, 0)
	inColumns := make([]string, 0)
	values := make([]interface{}, 0)
	insertCols := make([]string, 0)
	i := 0
	for _, key := range cols {
		fdb := schema[key]
		f := mv.Field(fdb.Index)
		fieldValue := f.Interface()
		tkey := `"` + strings.Replace(key, `"`, `""`, -1) + `"`
		tkey = strings.ToUpper(tkey)
		inColumns = append(inColumns, "temp."+key)
		for _, k := range keys {
			if key == k {
				onDupe := "a." + tkey + "=" + "temp." + tkey
				uniqueCols = append(uniqueCols, onDupe)
			}else{
				setColumns = append(setColumns, "a."+tkey+" = temp."+tkey)
			}
		}
		isNil := false
		if f.Kind() == reflect.Ptr {
			if reflect.ValueOf(fieldValue).IsNil() {
				isNil = true
			} else {
				fieldValue = reflect.Indirect(reflect.ValueOf(fieldValue)).Interface()
			}
		}
		if isNil {
			variables = append(variables,"null "+tkey)
		}else {
			v, ok := sq.GetDBValue(fieldValue)
			if ok {
				variables = append(variables, v+" "+tkey)
			} else {
				if boolValue, ok := fieldValue.(bool); ok {
					if boolValue {
						if fdb.True != nil {
							variables = append(variables, buildParam(i)+" "+tkey)
							values = append(values, *fdb.True)
							i++
						} else {
							variables = append(variables,"1 "+tkey)
						}
					}else {
						if fdb.False != nil {
							variables = append(variables, buildParam(i)+" "+tkey)
							values = append(values, *fdb.False)
							i++
						} else {
							variables = append(variables,"0 "+tkey)
						}
					}
				}else {
					variables = append(variables, buildParam(i)+" "+tkey)
					values = append(values, fieldValue)
					i++
				}
			}
		}
		insertCols = append(insertCols, tkey)
	}

	query := fmt.Sprintf("MERGE INTO %s a USING (SELECT %s FROM dual) temp ON  (%s) WHEN MATCHED THEN UPDATE SET %s WHEN NOT MATCHED THEN INSERT (%s) VALUES (%s)",
		table,
		strings.Join(variables, ", "),
		strings.Join(uniqueCols, " AND "),
		strings.Join(setColumns, ", "),
		strings.Join(insertCols, ", "),
		strings.Join(inColumns, ", "),
	)
	return query, values, nil
}

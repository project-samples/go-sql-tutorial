package app

import (
	"context"
	"database/sql"
	"github.com/core-go/health"
	s "github.com/core-go/health/sql"
	//_ "github.com/lib/pq"
	_ "github.com/godror/godror"

	"go-service/internal/handlers"
	"go-service/internal/services"
)

const (
	CreateDatabase = `CREATE if not exists DATABASE mydb
WITH
OWNER = postgres
ENCODING = 'UTF8'
LC_COLLATE = 'English_United States.1252'
LC_CTYPE = 'English_United States.1252'
TABLESPACE = pg_default
CONNECTION LIMIT = -1`
)

const (
	CreateTable = `
create table if not exists users (
  id varchar(40) not null,
    username varchar(120),
    email varchar(120),
    phone varchar(45),
    date_of_birth timestamp with time zone,
    interests varchar[],
    skills json[],
    settings json,
    primary key (id)
)`
)

type ApplicationContext struct {
	HealthHandler *health.Handler
	UserHandler   *handlers.UserHandler
}

func NewApp(context context.Context, conf DatabaseConfig) (*ApplicationContext, error) {
	db, err := sql.Open(conf.Driver, conf.DataSourceName)
	if err != nil {
		return nil, err
	}

	/*fmt.Sprintf("%s", "create database if not exists mydb")
	_, err = db.ExecContext(context, CreateDatabase)
	if err != nil {
		return nil, err
	}

	stmtUseDB := fmt.Sprintf("%s", "use masterdata")
	_, err = db.ExecContext(context, stmtUseDB)
	if err != nil {
		return nil, err
	}*/

	//_, err = db.ExecContext(context, CreateTable)
	//if err != nil {
	//	return nil, err
	//}

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	sqlChecker := s.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}

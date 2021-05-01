package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/core-go/health"
	s "github.com/core-go/health/sql"
	_ "github.com/go-sql-driver/mysql"

	"go-service/internal/handlers"
	"go-service/internal/services"
)

const (
	CreateTable = `
create table if not exists users (
  id varchar(40) not null,
  username varchar(120),
  email varchar(120),
  phone varchar(45),
  date_of_birth date,
  primary key (id)
)`
)

type ApplicationContext struct {
	HealthHandler *health.HealthHandler
	UserHandler   *handlers.UserHandler
}

func NewApp(context context.Context, conf DatabaseConfig) (*ApplicationContext, error) {
	db, err := sql.Open(conf.Driver, conf.DataSourceName)
	if err != nil {
		return nil, err
	}

	stmtCreate := fmt.Sprintf("%s", "create database if not exists masterdata")
	_, err = db.ExecContext(context, stmtCreate)
	if err != nil {
		return nil, err
	}

	stmtUseDB := fmt.Sprintf("%s", "use masterdata")
	_, err = db.ExecContext(context, stmtUseDB)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(context, CreateTable)
	if err != nil {
		return nil, err
	}

	userService := services.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	sqlChecker := s.NewHealthChecker(db)
	healthHandler := health.NewHealthHandler(sqlChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}

package app

import (
	"context"
	"github.com/core-go/health"
	s "github.com/core-go/health/sql"
	"github.com/core-go/sql"
	_ "github.com/go-sql-driver/mysql"

	"go-service/internal/handler"
	"go-service/internal/service"
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
	Health *health.Handler
	User   *handler.UserHandler
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	db, err := sql.OpenByConfig(cfg.Sql)
	if err != nil {
		return nil, err
	}

	stmtCreate := "create database if not exists masterdata"
	_, err = db.ExecContext(ctx, stmtCreate)
	if err != nil {
		return nil, err
	}

	stmtUseDB := "use masterdata"
	_, err = db.ExecContext(ctx, stmtUseDB)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, CreateTable)
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	sqlChecker := s.NewHealthChecker(db)
	healthHandler := health.NewHandler(sqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}

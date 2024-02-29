package controllers

import (
	"context"

	"github.com/guard-ai/guard-server/pkg"
	"github.com/guard-ai/guard-server/pkg/notifications"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgxuuid "github.com/vgarvardt/pgx-google-uuid/v5"
)

type Controller struct {
	db *pgxpool.Pool
	notifier *notifications.Notifier
}

func New() *Controller {
	ctx := context.Background()
	conf, err := pgxpool.ParseConfig(pkg.Env().PostgresConnectionString)
	conf.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		pgxuuid.Register(c.TypeMap())
		return nil
	}

	db, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		panic(err)
	}

	notifier, err := notifications.NewNotifier()
	if err != nil {
		panic(err)
	}

	return &Controller{
		db,
		notifier,
	} 
}

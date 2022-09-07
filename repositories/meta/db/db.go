package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	pkgdb "github.com/JanCalebManzano/tag-microservices/pkg/db"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/model"
)

type MetaDB interface {
	GetAllSystems(ctx context.Context) ([]*model.System, error)
	MonitorSystems(interval time.Duration) chan struct{}
}

func NewMetaDB(ctx context.Context, log *zap.Logger) (MetaDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_META_USERNAME"),
		os.Getenv("DB_META_PASSWORD"),
		os.Getenv("DB_META_HOST"),
		os.Getenv("DB_META_PORT"),
		os.Getenv("DB_META_DB"),
	)

	conn, err := pkgdb.SQLConnect("mysql", dsn)
	if err != nil {
		return nil, err
	}

	d := &db{
		SqlDB: pkgdb.NewSqlDB(
			conn,
			"2006-01-02 15:04:05",
			"2006-01-02 15:04:05 UTC -0700",
		),
		ctx: ctx,
		log: log.Named("db"),
	}

	d.systems, err = d.fetchAllSystems()
	if err != nil {
		return nil, err
	}

	return d, nil
}

type db struct {
	*pkgdb.SqlDB

	ctx     context.Context
	log     *zap.Logger
	systems []*model.System
}

func (d *db) fetchAllSystems() ([]*model.System, error) {
	rows, cancel, err := d.Query(d.ctx, "SELECT * FROM t_system")
	if err != nil {
		return nil, err
	}
	defer cancel()
	defer rows.Close()

	systems := make([]*model.System, 0)
	for rows.Next() {
		system := &model.System{}
		if err := rows.Scan(
			&system.SystemNo,
			&system.SystemName,
			&system.SystemShortName,
			&system.SetUser,
			&system.SetTimestamp,
		); err != nil {
			return nil, err
		}

		systems = append(systems, system)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return systems, nil
}

func (d *db) GetAllSystems(_ context.Context) ([]*model.System, error) {
	return d.systems, nil
}

func (d *db) MonitorSystems(interval time.Duration) chan struct{} {
	var err error
	ch := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)

		for {
			select {
			case <-ticker.C:
				d.systems, err = d.fetchAllSystems()
				if err != nil {
					d.log.Error("server: Unable to get updated systems", zap.Error(err))
					return
				}

				// notify updates, this will block unless there is a listener on the other end
				ch <- struct{}{}
			}
		}
	}()

	return ch
}

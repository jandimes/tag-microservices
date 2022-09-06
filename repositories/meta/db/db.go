package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	pkgdb "github.com/JanCalebManzano/tag-microservices/pkg/db"

	"github.com/JanCalebManzano/tag-microservices/repositories/meta/model"
)

type MetaDB interface {
	GetAllSystems(ctx context.Context) ([]*model.System, error)
}

func NewMetaDB() (MetaDB, error) {
	var (
		conn *sql.DB
		err  error
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_META_USERNAME"),
		os.Getenv("DB_META_PASSWORD"),
		os.Getenv("DB_META_HOST"),
		os.Getenv("DB_META_PORT"),
		os.Getenv("DB_META_DB"),
	)

	if conn, err = pkgdb.SQLConnect("mysql", dsn); err != nil {
		return nil, err
	}

	return &db{
		SqlDB: pkgdb.NewSqlDB(
			conn,
			"2006-01-02 15:04:05",
			"2006-01-02 15:04:05 UTC -0700",
		),
	}, nil
}

type db struct {
	*pkgdb.SqlDB
}

func (d *db) GetAllSystems(ctx context.Context) ([]*model.System, error) {
	rows, cancel, err := d.Query(ctx, "SELECT * FROM t_system")
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

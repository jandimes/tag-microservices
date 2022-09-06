package db

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	pkgerrors "github.com/JanCalebManzano/tag-microservices/pkg/errors"

	"github.com/JanCalebManzano/tag-microservices/pkg/masker"
)

// SQLConnect is used to connect to a sql database
func SQLConnect(driver, dsn string) (conn *sql.DB, err error) {
	conn, err = sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(100)
	conn.SetMaxIdleConns(100)
	conn.SetConnMaxLifetime(5 * time.Minute)

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, err
}

type SqlDB struct {
	pool            *sql.DB
	timestampInFmt  string
	timestampOutFmt string
}

func NewSqlDB(pool *sql.DB, timestampInFmt string, timestampOutFmt string) *SqlDB {
	return &SqlDB{pool: pool, timestampInFmt: timestampInFmt, timestampOutFmt: timestampOutFmt}
}

func (db *SqlDB) Close() error {
	return db.pool.Close()
}

func (db *SqlDB) processRow(
	msk *masker.Masker, cols []sql.RawBytes, colTypes []*sql.ColumnType, colNames []string,
) (row map[string]interface{}, err error) {
	row = map[string]interface{}{}
	for i, col := range cols {
		switch colTypes[i].DatabaseTypeName() {
		case "INT", "INT4", "INT8", "BIGINT":
			if string(col) == "" {
				row[colNames[i]] = 0
				continue
			}
			if row[colNames[i]], err = strconv.Atoi(string(col)); err != nil {
				return nil, err
			}
		case "FLOAT":
			if row[colNames[i]], err = strconv.ParseFloat(string(col), 32); err != nil {
				return nil, err
			}
		case "TIMESTAMP":
			t, err := time.Parse(db.timestampInFmt, string(col))
			if err != nil {
				row[colNames[i]] = string(col)
			} else {
				row[colNames[i]] = t.Format(db.timestampOutFmt)
			}
		default:
			s := string(col)
			switch colNames[i] {
			case "member_name":
				s = msk.Name(s)
			case "email_address":
				s = msk.Email(s)
			case "phone_number":
				s = msk.Mobile(s)
			}
			row[colNames[i]] = s
		}
	}
	return row, nil
}

func (db *SqlDB) Query(ctx context.Context, stmt string, args ...interface{}) (
	res *sql.Rows, cancel context.CancelFunc, err error) {
	c, cancel := context.WithTimeout(ctx, time.Second*10)

	conn, err := db.pool.Conn(c)
	if err != nil {
		return nil, cancel, pkgerrors.NewQueryError("DB", "", err)
	}

	// Execute the query
	rows, err := conn.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, cancel, pkgerrors.NewQueryError("DB", "", err)
	}

	return rows, cancel, nil

	/**
	// Get column names
	colNames, err := rows.Columns()
	if err != nil {
		return nil, pkgerrors.NewQueryError("DB", "", err)
	}

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, pkgerrors.NewQueryError("DB", "", err)
	}

	cols := make([]sql.RawBytes, len(colNames))
	colPtrs := make([]interface{}, len(colNames))
	for i := range cols {
		colPtrs[i] = &cols[i]
	}

	res = make([]map[string]interface{}, 0)

	msk := masker.New()
	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		if err = rows.Scan(colPtrs...); err != nil {
			return nil, err
		}

		row, err := db.processRow(msk, cols, colTypes, colNames)
		if err != nil {
			return nil, err
		}

		res = append(res, row)
	}

	if err = rows.Err(); err != nil {
		return nil, pkgerrors.NewQueryError("DB", "", err)
	}

	return res, nil

	*/
}

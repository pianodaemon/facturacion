package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

type (
	PgSqlSettings struct {
		Host     string `default:"rdbms_postgres"`
		Db       string `default:"interview"`
		User     string `default:"postgres"`
		Password string `default:"postgres"`
		Port     int    `default:"5432"`
	}
)

var pgSettings PgSqlSettings

func init() {

	envconfig.Process("postgres", &pgSettings)
}

func shapeConnStr() string {

	// SSL mode disable to use in containers
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		pgSettings.User,
		pgSettings.Password,
		pgSettings.Host,
		pgSettings.Port,
		pgSettings.Db)
}

func ProbeConn() error {

	dbinfo := shapeConnStr()

	db, err := sql.Open("postgres", dbinfo)

	if err != nil {

		return err
	}

	defer db.Close()

	if err = db.Ping(); err != nil {

		return err
	}

	return err
}

func execEntitySp(dbPtr *sql.DB, sph EntitySpHandler) (int, error) {

	var lastID int
	var rmsg string

	ctx := context.Background()
	tx, err := dbPtr.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	err = sph(&ctx, dbPtr, &lastID, &rmsg)

	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if lastID == -1 {
		return -1, fmt.Errorf(rmsg)
	}

	err = tx.Commit()

	if err != nil {
		return -1, err
	}

	return lastID, err
}

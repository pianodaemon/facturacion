package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type (
	Entity struct {
		ID   int
		code string
	}

	EntitySpHandler func(ctxPtr *context.Context, dbPtr *sql.DB, lastIdPtr *int, rmsgPtr *string) error
)

func alterEntity(dbPtr *sql.DB, ID int, code string) (int, error) {

	sph := func(ctxPtr *context.Context, dbPtr *sql.DB, lastIdPtr *int, rmsgPtr *string) error {

		fmt.Printf("%d %s\n", ID, code)
		return dbPtr.QueryRowContext(*ctxPtr, "SELECT * FROM alter_entity($1, $2) AS (rc integer, msg text)", ID, code).Scan(lastIdPtr, rmsgPtr)
	}

	return execEntitySp(dbPtr, sph)
}

// Logical deletion for the entity data
func deleteEntity(dbPtr *sql.DB, ID int) error {

	updateSQL := "UPDATE entitys SET blocked = true, last_touch_time = now() WHERE NOT blocked AND id = $1"
	res, err := dbPtr.Exec(updateSQL, ID)

	if err != nil {
		return err
	}

	number, err := res.RowsAffected()
	if number == 0 {

		return fmt.Errorf("No rows affected during logical deletion")
	}

	return err
}

func (selfPtr *Entity) CreateEntity(dbPtr *sql.DB) (int, error) {

	return alterEntity(dbPtr, selfPtr.ID, selfPtr.code)
}

func (selfPtr *Entity) EditEntity(dbPtr *sql.DB) (int, error) {

	return alterEntity(dbPtr, selfPtr.ID, selfPtr.code)
}

func (selfPtr *Entity) DelEntity(dbPtr *sql.DB) error {

	return deleteEntity(dbPtr, selfPtr.ID)
}

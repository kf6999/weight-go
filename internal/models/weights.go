package models

import (
	"database/sql"
	"errors"
	"time"
)

type Weight struct {
	ID     int
	Weight int
	Notes  string
	Date   time.Time
}

//  WeightModel type which wraps a sql.DB pool

type WeightModel struct {
	DB *sql.DB
}

func (m *WeightModel) Insert(weight int, notes string) (int, error) {
	stmt := `insert into weights (weight, notes) values (?,?)`

	result, err := m.DB.Exec(stmt, weight, notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *WeightModel) Get(id int) (*Weight, error) {
	stmt := `select weight, coalesce(notes,'') from weights where id = ?`

	row := m.DB.QueryRow(stmt, id)

	w := &Weight{}
	err := row.Scan(&w.ID, &w.Notes)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return w, nil
}

func (m *WeightModel) Latest() ([]*Weight, error) {
	return nil, nil
}

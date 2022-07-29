package models

import (
	"database/sql"
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

func (m *WeightModel) Insert(weight int, notes string, date int) (int, error) {
	return 0, nil
}

func (m *WeightModel) Get(id int) (*Weight, error) {
	return nil, nil
}

func (m *WeightModel) Latest() ([]*Weight, error) {
	return nil, nil
}

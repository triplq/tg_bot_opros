package database

import (
	"database/sql"
	"time"
)

type Model struct {
	DB *sql.DB
}

type Form struct {
	ID        int
	Msg       string
	Posted_at time.Time
	Channel   string
	PostData  string
}

func (m *Model) PasteForms(forms []*Form) error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`insert into list (msg, posted_at, post_data, channel) values ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, f := range forms {
		_, err := stmt.Exec(f.Msg, f.Posted_at, f.PostData, f.Channel)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) LastHash(channel_name string) (string, error) {
	stmt := `select post_data from list
	where channel=$1
	order by posted_at desc
	limit 1`

	var hash string

	res := m.DB.QueryRow(stmt, channel_name)
	err := res.Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (m *Model) Show() ([]*Form, error) {
	forms, err := m.list()
	if err != nil {
		return nil, err
	}
	err = m.setSeen()
	if err != nil {
		return nil, err
	}

	return forms, nil
}

func (m *Model) setSeen() error {
	stmt := `update list set seen=true where seen=false`

	_, err := m.DB.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) list() ([]*Form, error) {
	stmt := `select msg from list where seen=false`
	forms := []*Form{}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := &Form{}
		err := rows.Scan(&f.Msg)
		if err != nil {
			return nil, err
		}
		forms = append(forms, f)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return forms, err
}

package database

import "database/sql"

type Model struct {
	DB *sql.DB
}

type Form struct {
	ID  int
	msg string
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
	stmt := `select ID, msg from list where seen=false`
	forms := []*Form{}

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		f := &Form{}
		err := rows.Scan(&f.ID, &f.msg)
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

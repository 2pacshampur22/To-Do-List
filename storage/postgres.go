package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	DB *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{DB: db}, nil

}

func (p *PostgresStorage) Add(t Task) error {
	_, err := p.DB.Exec(`
	insert into tasks ( name, description, is_done)
	values ($1, $2, $3);`, t.Name, t.Description, t.IsDone)
	return err
}

func (p *PostgresStorage) List() ([]Task, error) {
	rows, err := p.DB.Query("select id, name, description, is_done from tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []Task

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Name, &t.Description, &t.IsDone)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (p *PostgresStorage) Delete(id int) error {
	_, err := p.DB.Exec("delete from tasks where id = $1", id)
	return err
}

func (p *PostgresStorage) Done(id int) error {
	_, err := p.DB.Exec("update tasks set is_done = true where id = $1", id)
	return err
}

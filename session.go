package toygo

import (
	"database/sql"
)

type Session struct {
	db *sql.DB
}

func Init(driverName string, dataSourceName string) (*Session, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err.Error())
	}

	return &Session{db}, err
}

func (s *Session) Close() {
	s.db.Close()
}

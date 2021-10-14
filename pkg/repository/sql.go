package repository

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type PostgreSql struct {
	Host     string
	Port     int
	Db       string
	User     string
	Password string
}

func LoadPostgreConf(fp string) (PostgreSql, error) {
	content, err := ioutil.ReadFile(fp)
	if err != nil {
		return PostgreSql{}, err
	}
	c := strings.Split(string(content), ":")
	port, err := strconv.Atoi(c[1])
	if err != nil {
		return PostgreSql{}, err
	}
	return PostgreSql{
		Host:     c[0],
		Port:     port,
		Db:       c[2],
		User:     c[3],
		Password: c[4],
	}, nil
}

func NewPostgreConnection(p PostgreSql) (*sql.DB, error) {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Db,
	)
	return sql.Open("postgres", psqlconn)
}

func GetLastId(db *sql.DB, table string) (int64, error) {
	rows, err := db.Query(
		fmt.Sprintf("SELECT id FROM %s ORDER BY id DESC LIMIT 1",
			table,
		),
	)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	lastId := int64(0)
	for rows.Next() {
		if err = rows.Scan(&lastId); err != nil {
			return 0, err
		}
	}

	return lastId, nil
}

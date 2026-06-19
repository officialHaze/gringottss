package sqliteclient

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	db "github.com/officialhaze/gringottss/api-server/db/sqlc"
	"github.com/officialhaze/gringottss/api-server/logger"
	"github.com/pressly/goose/v3"
)

var (
	Queries *db.Queries
)

type Client struct {
	DriverName   string
	DBName       string
	MaxOpenConns int8
	DB           *sql.DB
}

func Init(drivername, dbname string, maxopenconns int8) *Client {
	return &Client{
		DriverName:   drivername,
		DBName:       dbname,
		MaxOpenConns: maxopenconns,
	}
}

func (c *Client) OpenDB() (*Client, error) {
	logger.INFO().Println("Opening DB....")
	query := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)", c.DBName)
	db, err := sql.Open(c.DriverName, query)
	if err != nil {
		return nil, err
	}
	logger.SUCCESS().Println("DB opened successfully.")

	// Set maximum open connections
	logger.INFO().Printf("Setting maximum db connections - %d", c.MaxOpenConns)
	db.SetMaxOpenConns(int(c.MaxOpenConns))

	// Set db for global use in Client
	logger.INFO().Println("Setting DB for global use.")
	c.DB = db

	return c, nil
}

func (c *Client) RunMigrations() error {
	logger.INFO().Println("Running migrations...")
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("Migrations failed!\n%s\n", err.Error())
	}
	if err := goose.Up(c.DB, "db/migrations"); err != nil {
		return fmt.Errorf("Migrations failed!\n%s\n", err.Error())
	}

	return nil
}

func (c *Client) LoadQueries() error {
	logger.INFO().Println("Loading queries...")
	if c.DB == nil {
		return fmt.Errorf("No DB found to load queries!")
	}

	logger.SUCCESS().Println("Queries loaded successfully!")

	logger.INFO().Println("Setting queries for global use.")
	Queries = db.New(c.DB)
	return nil
}

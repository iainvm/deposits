package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const driver = "postgres"

type DataSource struct {
	Host         string
	Port         string
	User         string
	Password     string
	DatabaseName string
	SSLMode      bool
}

func NewDataSource(host string, port string, user string, password string, databaseName string, sslmode bool) DataSource {
	return DataSource{
		Host:         host,
		Port:         port,
		User:         user,
		Password:     password,
		DatabaseName: databaseName,
		SSLMode:      false,
	}
}

func (dataSource DataSource) String() string {
	sslMode := "disable"
	if dataSource.SSLMode {
		sslMode = "enable"
	}
	return fmt.Sprintf(
		"host=%s port= %s user=%s password=%s dbname=%s sslmode=%s",
		dataSource.Host,
		dataSource.Port,
		dataSource.User,
		dataSource.Password,
		dataSource.DatabaseName,
		sslMode,
	)
}

func Connect(dataSource DataSource) (*sqlx.DB, error) {
	return sqlx.Connect(
		driver,
		dataSource.String(),
	)
}

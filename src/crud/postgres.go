package crud

import (
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/geometry-labs/go-service-template/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type PostgresConn struct {
	conn *gorm.DB
}

var postgresInstance *PostgresConn
var postgresConnOnce sync.Once

func NewPostgresConn(dsn string) (*PostgresConn, error) { // Only for testing
	session, err := createSession(dsn)
	if err != nil {
		zap.S().Info("Cannot create a connection to postgres", err)
	}
	postgresInstance = &PostgresConn{
		conn: session,
	}
	return postgresInstance, err
}

func GetPostgresConn() *PostgresConn {
	postgresConnOnce.Do(func() {
		// TODO: create dsn string from env variables
		//dsn := NewDsn("postgres", "5432", "postgres", "changeme", "postgres", "disable", "UTC")
		dsn := NewDsn(config.Config.DbHost, config.Config.DbPort, config.Config.DbUser,
			config.Config.DbPassword, config.Config.DbName, config.Config.DbSslmode,
			config.Config.DbTimezone)
		//session, err := createSession(dsn)
		session, err := retryGetPostgresSession(dsn)
		if err != nil {
			zap.S().Fatal("Finally Cannot create a connection to postgres", err)
		} else {
			zap.S().Info("Finally successful connection to postgres")
		}
		postgresInstance = &PostgresConn{
			conn: session,
		}
	})
	return postgresInstance
}

func (p *PostgresConn) GetConn() *gorm.DB {
	return p.conn
}

func createSession(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.S().Info("err:", err)
	}

	return db, err
}

func NewDsn(host string, port string, user string, password string, dbname string, sslmode string, timezone string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)
}

func retryGetPostgresSession(dsn string) (*gorm.DB, error) {
	var session *gorm.DB
	operation := func() error {
		sess, err := createSession(dsn)
		if err != nil {
			zap.S().Info("POSTGRES SESSION Error : ", err.Error())
		} else {
			session = sess
		}
		return err
	}
	neb := backoff.NewExponentialBackOff()
	//neb.MaxElapsedTime = time.Minute
	err := backoff.Retry(operation, neb)
	return session, err
}

package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"selfstudy/crawl/product/util"
	"sync"
	"time"
)

type PostgresDataSource struct {
	connection *pgxpool.Pool
}

var (
	postgresDataSourcePool *pgxpool.Pool
	once                   sync.Once
)

// config: https://github.com/jackc/pgx/discussions/1989
func NewPostgresDataSource() PostgresDataSource {

	once.Do(func() {
		config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
		if err != nil {
			util.LogError(err.Error())
		}
		config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			util.LogDebug("Postgres Connection Established")
			return nil
		}
		// Customize pool settings
		config.MaxConns = 10
		config.MinConns = 2
		config.MaxConnLifetime = time.Hour
		config.MaxConnIdleTime = time.Minute * 30
		config.HealthCheckPeriod = time.Minute

		config.PrepareConn = func(ctx context.Context, c *pgx.Conn) (bool, error) {
			util.LogDebug("Before acquiring the connection pool to the database!!")
			return true, nil
		}

		config.AfterRelease = func(c *pgx.Conn) bool {
			util.LogDebug("After releasing the connection pool to the database!!")
			return true
		}

		config.BeforeClose = func(c *pgx.Conn) {
			util.LogDebug("Closed the connection pool to the database!!")
		}

		config.ShouldPing = func(ctx context.Context, params pgxpool.ShouldPingParams) bool {
			util.LogDebug("ShouldPing is called!")
			return true
		}

		connPool, err := pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			util.LogError(err.Error())
			panic(err)
		}
		postgresDataSourcePool = connPool

	})
	return PostgresDataSource{connection: postgresDataSourcePool}
}

func getOPostgresDataSourceInstance() PostgresDataSource {
	return NewPostgresDataSource()
}

var GetOPostgresDataSourceInstance = getOPostgresDataSourceInstance

func (p PostgresDataSource) insert() {

}

func (p PostgresDataSource) insertBatch() {

}

func (p PostgresDataSource) update() {

}

func (p PostgresDataSource) delete() {

}

func (p PostgresDataSource) close() {
	p.connection.Close()
}

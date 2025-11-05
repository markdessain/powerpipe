package db_client

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/apache/arrow/go/v17/arrow/flight/flightsql/driver"
	"github.com/turbot/pipe-fittings/v2/backend"
	"github.com/turbot/steampipe-plugin-sdk/v5/sperr"
)

type FlightBackend struct {
	connectionString string
	rowReader        backend.RowReader
}

func NewFlightBackend(connString string) (*FlightBackend, error) {
	connString = strings.TrimSpace(connString) // remove any leading or trailing whitespace
	// connString = strings.TrimPrefix(connString, "flight://")

	fmt.Println(connString)

	return &FlightBackend{
		connectionString: connString,
		rowReader:        backend.NewBasicRowReader(),
	}, nil
}

// Connect implements Backend.
func (b *FlightBackend) Connect(ctx context.Context, options ...backend.BackendOption) (*sql.DB, error) {

	_ = &driver.Driver{}
	config := backend.NewBackendConfig(options)
	db, err := sql.Open("flightsql", b.connectionString)
	if err != nil {
		return nil, sperr.WrapWithMessage(err, "could not connect to flightsql backend")
	}
	db.SetConnMaxIdleTime(config.MaxConnIdleTime)
	db.SetConnMaxLifetime(config.MaxConnLifeTime)
	db.SetMaxOpenConns(config.MaxOpenConns)

	return db, nil
}

func (b *FlightBackend) ConnectionString() string {
	return b.connectionString
}

func (b *FlightBackend) Name() string {
	return "flightsql"
}

// RowReader implements Backend.
func (b *FlightBackend) RowReader() backend.RowReader {
	return b.rowReader
}

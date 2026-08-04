package main

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()

	var (
		database = "testdb"
		user     = "testuser"
		password = "testpassword"
	)

	pgContainer, err := postgres.Run(ctx,
		"postgis/postgis:12-3.0",
		postgres.WithDatabase(database),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
	)
	t.Cleanup(func() {
		assert.NoError(t, testcontainers.TerminateContainer(pgContainer))
	})
	if err != nil {
		t.Fatal(err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	db, err := sql.Open("postgres", connStr)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, db.Close())
	}()

	assert.NoError(t, createDB(ctx, db))

	assert.NoError(t, populateDB(ctx, db))

	r := bytes.NewBufferString(`{"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}`)
	assert.NoError(t, readGeoJSON(ctx, db, r))

	w := &strings.Builder{}
	assert.NoError(t, writeGeoJSON(ctx, db, w))
	assert.Equal(t, strings.Join([]string{
		`{"id":1,"name":"London","geometry":{"type":"Point","coordinates":[0.1275,51.50722]}}`,
		`{"id":2,"name":"Berlin","geometry":{"type":"Point","coordinates":[13.405,52.52]}}`,
		`{"id":3,"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}`,
	}, "\n")+"\n", w.String())
}

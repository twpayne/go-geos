// Demonstrate integration with PostgreSQL/PostGIS.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	_ "github.com/lib/pq" // Add PostgreSQL to database/sql.

	"github.com/twpayne/go-geos"
	"github.com/twpayne/go-geos/geometry"
)

var (
	dsn = flag.String("dsn", "postgres://localhost/geomtest?binary_parameters=yes&sslmode=disable", "data source name")

	create   = flag.Bool("create", false, "create database schema")
	populate = flag.Bool("populate", false, "populate waypoints")
	read     = flag.Bool("read", false, "import waypoint from stdin in GeoJSON format")
	write    = flag.Bool("write", false, "write waypoints to stdout in GeoJSON format")
)

// A Waypoint is a location with an identifier and a name.
type Waypoint struct {
	ID       int                `json:"id"`
	Name     string             `json:"name"`
	Geometry *geometry.Geometry `json:"geometry"`
}

// createDB demonstrates create a PostgreSQL/PostGIS database with a table with
// a geometry column.
func createDB(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE EXTENSION IF NOT EXISTS postgis;
		CREATE TABLE IF NOT EXISTS waypoints (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			geom geometry(POINT, 4326) NOT NULL
		);
	`)
	return err
}

// populateDB demonstrates populating a PostgreSQL/PostGIS database using
// pq.CopyIn for fast imports.
func populateDB(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	stmt, err := tx.PrepareContext(ctx, "COPY waypoints (name, geom) FROM stdin")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, waypoint := range []Waypoint{
		{
			Name:     "London",
			Geometry: geometry.NewGeometry(geos.NewPoint([]float64{0.1275, 51.50722}).SetSRID(4326)),
		},
		{
			Name:     "Berlin",
			Geometry: geometry.NewGeometry(geos.NewPoint([]float64{13.405, 52.52}).SetSRID(4326)),
		},
	} {
		if _, err := stmt.ExecContext(ctx, waypoint.Name, waypoint.Geometry); err != nil {
			return err
		}
	}
	if _, err := stmt.ExecContext(ctx); err != nil {
		return err
	}
	return tx.Commit()
}

// readGeoJSON demonstrates reading GeoJSON data and inserting it into a
// database with INSERT.
func readGeoJSON(ctx context.Context, db *sql.DB, r io.Reader) error {
	var waypoint Waypoint
	if err := json.NewDecoder(r).Decode(&waypoint); err != nil {
		return err
	}
	waypoint.Geometry.SetSRID(4326)
	_, err := db.ExecContext(ctx, `
		INSERT INTO waypoints(name, geom) VALUES ($1, $2);
	`, waypoint.Name, waypoint.Geometry)
	return err
}

// writeGeoJSON demonstrates reading data from a database with SELECT and
// writing it as GeoJSON.
func writeGeoJSON(ctx context.Context, db *sql.DB, w io.Writer) error {
	rows, err := db.QueryContext(ctx, `
		SELECT id, name, ST_AsEWKB(geom) FROM waypoints ORDER BY id ASC;
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var waypoint Waypoint
		if err := rows.Scan(&waypoint.ID, &waypoint.Name, &waypoint.Geometry); err != nil {
			return err
		}
		if err := json.NewEncoder(w).Encode(&waypoint); err != nil {
			return err
		}
	}
	return rows.Err()
}

func run() error {
	ctx := context.Background()

	flag.Parse()
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	if *create {
		if err := createDB(ctx, db); err != nil {
			return err
		}
	}
	if *populate {
		if err := populateDB(ctx, db); err != nil {
			return err
		}
	}
	if *read {
		if err := readGeoJSON(ctx, db, os.Stdin); err != nil {
			return err
		}
	}
	if *write {
		if err := writeGeoJSON(ctx, db, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Println(err)
		os.Exit(1)
	}
}

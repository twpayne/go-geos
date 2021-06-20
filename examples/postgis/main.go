package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/lib/pq"

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
func createDB(db *sql.DB) error {
	_, err := db.Exec(`
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
func populateDB(db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(pq.CopyIn("waypoints", "name", "geom"))
	if err != nil {
		return err
	}
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
		if _, err := stmt.Exec(waypoint.Name, waypoint.Geometry); err != nil {
			return err
		}
	}
	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return tx.Commit()
}

// readGeoJSON demonstrates reading GeoJSON data and inserting it into a
// database with INSERT.
func readGeoJSON(db *sql.DB, r io.Reader) error {
	var waypoint Waypoint
	if err := json.NewDecoder(r).Decode(&waypoint); err != nil {
		return err
	}
	_, err := db.Exec(`
		INSERT INTO waypoints(name, geom) VALUES ($1, $2);
	`, waypoint.Name, waypoint.Geometry)
	return err
}

// writeGeoJSON demonstrates reading data from a database with SELECT and
// writing it as GeoJSON.
func writeGeoJSON(db *sql.DB, w io.Writer) error {
	rows, err := db.Query(`
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
	flag.Parse()
	db, err := sql.Open("postgres", *dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return err
	}
	if *create {
		if err := createDB(db); err != nil {
			return err
		}
	}
	if *populate {
		if err := populateDB(db); err != nil {
			return err
		}
	}
	if *read {
		if err := readGeoJSON(db, os.Stdin); err != nil {
			return err
		}
	}
	if *write {
		if err := writeGeoJSON(db, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

# PostGIS example

This example demonstrates:

 * Connecting to a PostgreSQL/PostGIS database.
 * Importing data in GeoJSON format and storing it in the database.
 * Exporting data from the database and converting it to GeoJSON.


## Quick start

Change to this directory:

```console
$ cd ${GOPATH}/src/github.com/twpayne/go-geos/examples/postgis
```

Create a database called `geomtest`:

```console
$ createdb geomtest
```

Save the data source name in an environment variable, for example:

```
$ DSN="postgres://username:password@localhost/geomtest?binary_parameters=yes&sslmode=disable"
```

Create the database schema, including the PostGIS extension and a table with a
geometry column:

```console
$ go run . -dsn $DSN -create
```

Populate the database using [`pq.CopyIn`](https://pkg.go.dev/github.com/lib/pq#CopyIn):

```console
$ go run . -dsn $DSN -populate
```

Write data from the database in GeoJSON format:

```console
$ go run . -dsn $DSN -write
{"id":1,"name":"London","geometry":{"type":"Point","coordinates":[0.1275,51.50722]}}
{"id":2,"name":"Berlin","geometry":{"type":"Point","coordinates":[13.405,52.52]}}
```

Import new data into the database in GeoJSON format:

```console
$ echo '{"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}' | go run . -dsn $DSN -read
```

Verify that the data was imported:

```console
$ go run . -dsn $DSN -write
{"id":1,"name":"London","geometry":{"type":"Point","coordinates":[0.1275,51.50722]}}
{"id":2,"name":"Berlin","geometry":{"type":"Point","coordinates":[13.405,52.52]}}
{"id":3,"name":"Paris","geometry":{"type":"Point","coordinates":[2.3508,48.8567]}}
```

Delete the database:

```console
$ dropdb geomtest
```

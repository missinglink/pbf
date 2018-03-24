package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // required database driver
)

// schema copied from https://github.com/mfn/osmlib-sqlite

// Connection - Connection
type Connection struct {
	db   *sql.DB
	Stmt *Statements
}

// GetDB - expose the underlying db object
func (c *Connection) GetDB() *sql.DB {
	return c.db
}

// Open - open connection and set up
func (c *Connection) Open(path string) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	c.db = db

	c.tables()
	c.prepare()

	// https://github.com/mattn/go-sqlite3/issues/274
	db.SetMaxOpenConns(1)

	// start transaction
	_, err = c.db.Exec("BEGIN TRANSACTION")
	if err != nil {
		panic(err)
	}
}

// Close - close connection and clean up
func (c *Connection) Close() {

	defer c.Stmt.Close()
	defer c.db.Close()

	// commit transaction
	_, err := c.db.Exec("END TRANSACTION")
	if err != nil {
		panic(err)
	}
}

// created tables
func (c *Connection) tables() {
	_, err := c.db.Exec(`

		PRAGMA main.foreign_keys=OFF;
    PRAGMA main.page_size=4096;
    PRAGMA main.cache_size=-2000;
    PRAGMA main.synchronous=OFF;
    PRAGMA main.journal_mode=OFF;
    PRAGMA main.temp_store=MEMORY;

		BEGIN TRANSACTION;
		CREATE TABLE IF NOT EXISTS nodes (
		    id INTEGER NOT NULL PRIMARY KEY,
		    lon REAL NOT NULL,
		    lat REAL NOT NULL
		);
		CREATE TABLE IF NOT EXISTS node_tags (
		    ref INTEGER NOT NULL,
		    key TEXT,
		    value TEXT,
				UNIQUE( ref, key ) ON CONFLICT REPLACE
		);
		CREATE TABLE IF NOT EXISTS ways (
		    id INTEGER NOT NULL PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS way_tags (
		    ref INTEGER NOT NULL,
		    key TEXT,
		    value TEXT,
				UNIQUE( ref, key ) ON CONFLICT REPLACE
		);
		CREATE TABLE IF NOT EXISTS way_nodes (
		    way INTEGER NOT NULL,
		    num INTEGER NOT NULL,
		    node INTEGER NOT NULL,
				UNIQUE( way, num ) ON CONFLICT REPLACE
		);
		CREATE TABLE IF NOT EXISTS relations (
		    id INTEGER NOT NULL PRIMARY KEY
		);
		CREATE TABLE IF NOT EXISTS relation_tags (
		    ref INTEGER NOT NULL,
		    key TEXT,
		    value TEXT,
				UNIQUE( ref, key ) ON CONFLICT REPLACE
		);
		CREATE TABLE IF NOT EXISTS members (
		    relation INTEGER NOT NULL,
		    type TEXT,
		    ref INTEGER NOT NULL,
		    role TEXT
		);
		COMMIT TRANSACTION;`)
	if err != nil {
		panic(err)
	}
}

// created prepared statement
func (c *Connection) prepare() {

	node, err := c.db.Prepare("INSERT OR REPLACE INTO nodes (id, lon, lat) VALUES (:id, :lon, :lat)")
	if err != nil {
		panic(err)
	}

	nodeTags, err := c.db.Prepare("INSERT OR REPLACE INTO node_tags (ref, key, value) VALUES (:ref, :key, :value)")
	if err != nil {
		panic(err)
	}

	way, err := c.db.Prepare("INSERT OR REPLACE INTO ways (id) VALUES (:id)")
	if err != nil {
		panic(err)
	}

	wayTags, err := c.db.Prepare("INSERT OR REPLACE INTO way_tags (ref, key, value) VALUES (:ref, :key, :value)")
	if err != nil {
		panic(err)
	}

	wayNodes, err := c.db.Prepare("INSERT OR REPLACE INTO way_nodes (way, num, node) VALUES (:way, :num, :node)")
	if err != nil {
		panic(err)
	}

	relation, err := c.db.Prepare("INSERT OR REPLACE INTO relations (id) VALUES (:id)")
	if err != nil {
		panic(err)
	}

	relationTags, err := c.db.Prepare("INSERT OR REPLACE INTO relation_tags (ref, key, value) VALUES (:ref, :key, :value)")
	if err != nil {
		panic(err)
	}

	member, err := c.db.Prepare("INSERT OR REPLACE INTO members (relation, type, ref, role) VALUES (:relation, :type, :ref, :role)")
	if err != nil {
		panic(err)
	}

	c.Stmt = &Statements{
		Node:         node,
		NodeTags:     nodeTags,
		Way:          way,
		WayTags:      wayTags,
		WayNodes:     wayNodes,
		Relation:     relation,
		RelationTags: relationTags,
		Member:       member,
	}
}

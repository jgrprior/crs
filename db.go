package crs

import "gopkg.in/mgo.v2"

// Persister is the interface that supports writing Entry data to persistent storage.
type Persister interface {
	Save(*Entry) error
	Close()
}

// DB maintains the original MongoDB session.
type DB struct {
	url        string
	collection string
	session    *mgo.Session
	err        error
}

// Connect establishes a database connection given a database URL.
func Connect(url, collection string) (Persister, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return &DB{url: url, err: err}, err
	}
	if err := session.Ping(); err != nil {
		return &DB{url: url, err: err}, err
	}
	return &DB{url: url, collection: collection, session: session}, err
}

// Save writes an entry using the existing database connection.
func (db *DB) Save(e *Entry) error {
	s := db.session.Copy()
	defer s.Close()
	c := s.DB("").C(db.collection)
	return c.Insert(e)
}

// Close closes the original session created at dial time.
func (db *DB) Close() {
	db.session.Close()
}

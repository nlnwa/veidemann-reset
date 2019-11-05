package rethinkdb

import (
	"fmt"

	"github.com/pkg/errors"
	r "gopkg.in/rethinkdb/rethinkdb-go.v5"
)

type Options struct {
	Host     string
	Port     int
	Name     string
	Username string
	Password string
}

type Client struct {
	connectOpts r.ConnectOpts
	*r.Session
}

type ConnectOption func(*r.ConnectOpts)

func WithName(name string) ConnectOption {
	return func(opts *r.ConnectOpts) {
		opts.Database = name
	}
}

func WithAddress(host string, port int) ConnectOption {
	return func(opts *r.ConnectOpts) {
		opts.Address = fmt.Sprintf("%s:%d", host, port)
	}
}

func WithCredentials(username string, password string) ConnectOption {
	return func(opts *r.ConnectOpts) {
		opts.Username = username
		opts.Password = password
	}
}

func NewClient(options Options) *Client {
	return newClient(
		WithName(options.Name),
		WithAddress(options.Host, options.Port),
		WithCredentials(options.Username, options.Password),
	)
}

func newClient(options ...ConnectOption) *Client {
	var db Client
	for _, option := range options {
		option(&db.connectOpts)
	}
	return &db
}

func (db *Client) Connect() error {
	var err error
	if db.Session, err = r.Connect(db.connectOpts); err != nil {
		return errors.Wrap(err, "failed to connect to database")
	} else {
		return nil
	}
}

func (db *Client) Disconnect() error {
	if err := db.Session.Close(); err != nil {
		return errors.Wrap(err, "failed to disconnect from database")
	} else {
		return nil
	}
}

func (db *Client) Get(table string, id string, value interface{}) error {
	if cursor, err := r.Table(table).Get(id).Run(db.Session); err != nil {
		return errors.Wrapf(err, "failed to get document (id: %s) in table: %s", id, table)
	} else {
		if err = cursor.One(&value); err != nil {
			return err
		} else {
			return nil
		}
	}
}

// DeleteTable deletes all documents in table using soft durability (no write sync)
func (db *Client) DeleteTable(database string, table string) error {
	_, err := r.DB(database).Table(table).Delete(r.DeleteOpts{
		Durability:    "soft",
		ReturnChanges: false,
	}).Run(db.Session)

	if err != nil {
		return errors.Wrapf(err, "failed to delete documents in table: %s", table)
	}

	return nil
}

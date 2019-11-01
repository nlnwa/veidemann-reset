package rethinkdb

import "log"

type Query interface {
	Clean(map[string][]string) error
}

func (db *Client) Clean(laundry map[string][]string) error {
	err := db.Connect()
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Disconnect()
	}()

	for database, tables := range laundry {
		for _, table := range tables {
			log.Printf("removing all documents from table: %s.%s\n", database, table)
			err := db.DeleteTable(database, table)
			if err != nil {
				log.Printf(err.Error())
				continue
			}
		}
	}

	return nil
}

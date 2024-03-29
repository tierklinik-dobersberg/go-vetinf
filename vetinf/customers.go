package vetinf

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

// CustomerDB wraps the vetkldat.dbf file
type CustomerDB struct {
	infdat *Infdat

	l          sync.Mutex
	table      *godbf.DbfTable
	idToRowIdx map[int]int
}

// All loads all customers stored in vetkldat.dbf
func (db *CustomerDB) All() ([]Customer, error) {
	db.l.Lock()
	defer db.l.Unlock()

	len := db.table.NumberOfRecords()
	all := make([]Customer, len)
	for i := 0; i < len; i++ {
		var c Customer

		if err := db.table.DecodeRow(i, &c, true); err != nil {
			return nil, err
		}

		all[i] = c
	}

	return all, nil
}

// StreamAll streams all customers found in db.
func (db *CustomerDB) StreamAll(ctx context.Context) (<-chan Customer, <-chan error, int) {
	db.l.Lock()

	total := db.table.NumberOfRecords()
	customers := make(chan Customer, 10)
	errors := make(chan error, 10)

	go func() {
		defer db.l.Unlock()
		defer close(customers)
		defer close(errors)

		for i := 0; i < total; i++ {
			var c Customer

			if err := db.table.DecodeRow(i, &c, true); err != nil {
				select {
				case errors <- fmt.Errorf("row #%d: %w", i, err):
				case <-ctx.Done():
					return
				}
				continue
			}

			if db.table.RowIsDeleted(i) {
				c.Meta.Deleted = true
			}

			select {
			case customers <- c:
			case <-ctx.Done():
				return
			}
		}
	}()

	return customers, errors, total
}

// ByID loads a customer entry by it's ID
func (db *CustomerDB) ByID(id int) (*Customer, error) {
	db.l.Lock()
	defer db.l.Unlock()

	recordIdx, ok := db.idToRowIdx[id]
	if !ok {
		return nil, errors.New("ID not found")
	}

	var c Customer
	if err := db.table.DecodeRow(recordIdx, &c, true); err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *CustomerDB) buildIndex() error {
	db.l.Lock()
	defer db.l.Unlock()

	recordCount := db.table.NumberOfRecords()
	if db.idToRowIdx == nil {
		db.idToRowIdx = make(map[int]int, recordCount)
	}

	for i := 0; i < recordCount; i++ {
		id, err := db.table.Int64FieldValueByName(i, "KNR")
		if err != nil {
			if _, ok := err.(*strconv.NumError); !ok {
				return err
			}

			continue
		}

		db.idToRowIdx[int(id)] = i
	}

	return nil
}

// PrintFields prints all supported fields to stdout
func (db *CustomerDB) PrintFields() {
	for _, f := range db.table.Fields() {
		fmt.Printf("%s length=%d type=%c\n", f.Name(), f.Length(), f.FieldType())
	}
}

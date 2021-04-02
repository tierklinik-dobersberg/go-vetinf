package vetinf

import (
	"context"
	"fmt"
	"sync"

	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

type AnimalDB struct {
	infdat *Infdat

	l     sync.Mutex
	table *godbf.DbfTable
}

func (db *AnimalDB) StreamAll(ctx context.Context) (<-chan SmallAnimalRecord, <-chan error, int) {
	records := make(chan SmallAnimalRecord, 10)
	errors := make(chan error, 10)

	db.l.Lock()
	total := db.table.NumberOfRecords()

	go func() {
		defer db.l.Unlock()
		defer close(records)
		defer close(errors)

		for i := 0; i < total; i++ {
			var r SmallAnimalRecord

			if err := db.table.DecodeRow(i, &r, true); err != nil {
				select {
				case errors <- fmt.Errorf("row #%d: %w", i, err):
				case <-ctx.Done():
					return
				}
				continue
			}

			if db.table.RowIsDeleted(i) {
				r.Meta.Deleted = true
			}

			select {
			case records <- r:
			case <-ctx.Done():
				return
			}
		}
	}()

	return records, errors, total
}

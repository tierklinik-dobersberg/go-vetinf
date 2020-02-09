package vetinf

import (
	"fmt"

	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

// CustomerDB wraps the vetkldat.dbf file
type CustomerDB struct {
	infdat *Infdat
	table  *godbf.DbfTable
}

// PrintFields prints all supported fields to stdout
func (db *CustomerDB) PrintFields() {
	for _, f := range db.table.Fields() {
		fmt.Printf("%s length=%d type=%c\n", f.Name(), f.Length(), f.FieldType())
	}
}

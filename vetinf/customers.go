package vetinf

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

// Customer represents customer data stored in the vetkldat.dbf file
type Customer struct {
	ID        int    `dbf:"knr" json:"id"`
	Group     string `dbf:"gruppe" json:"group"`
	Name      string `dbf:"name" json:"name"`
	Firstname string `dbf:"vorname" json:"firstname"`
	Titel     string `dbf:"titel" json:"title"`
	Street    string `dbf:"strasse" json:"street"`
	CityCode  int    `dbf:"plz" json:"cityCode"`
	City      string `dbf:"ort" json:"city"`
	Phone     string `dbf:"telefon" json:"phone"`
}

func (c Customer) String() string {
	return fmt.Sprintf("Customer{id:%d name:%q}", c.ID, c.Name)
}

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

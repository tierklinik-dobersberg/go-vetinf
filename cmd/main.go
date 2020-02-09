package main

import (
	"log"

	"github.com/spf13/afero"
	"github.com/tierklinik-dobersberg/go-vetinf/vetinf"
)

func main() {
	infdat := vetinf.OpenReadonlyFs("/tmp", afero.NewOsFs())

	db, err := infdat.CustomerDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	db.PrintFields()
}

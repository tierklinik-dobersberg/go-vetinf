package main

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
	"github.com/tierklinik-dobersberg/go-vetinf/vetinf"
)

func main() {
	infdat := vetinf.OpenReadonlyFs("/tmp", afero.NewOsFs())

	db, err := infdat.CustomerDB("IBM852")
	if err != nil {
		log.Fatal(err.Error())
	}

	//db.PrintFields()
	all, err := db.All()
	if err != nil {
		log.Fatal(err.Error())
	}

	for _, a := range all {
		fmt.Println(a)
	}

	fmt.Println(db.ByID(13455))
}

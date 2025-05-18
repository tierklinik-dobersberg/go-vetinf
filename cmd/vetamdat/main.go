package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
)

type record struct {
	client [6]byte
	animal [6]byte
	index  [3]byte
	data   [(89 - 15)]byte
}

type analysis struct {
	highest int
	indexes map[int]struct{}
	records []record
}

func main() {
	reader, err := os.Open("./vetamdat")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	var idx int

	an := make(map[string]analysis)

	//encoder := mahonia.NewEncoder("IBM852")
	decoder := mahonia.NewDecoder("IBM852")

	for {
		idx++
		var r record

		if _, err := io.ReadFull(reader, r.client[:]); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			log.Fatal(err)
		}
		if _, err := io.ReadFull(reader, r.animal[:]); err != nil {
			log.Fatal(err)
		}
		if _, err := io.ReadFull(reader, r.index[:]); err != nil {
			log.Fatal(err)
		}
		if _, err := io.ReadFull(reader, r.data[:]); err != nil {
			log.Fatal(err)
		}

		key := string(r.client[:]) + "-" + string(r.animal[:])
		a := an[key]

		a.records = append(a.records, r)

		if strings.TrimSpace(string(r.index[:])) == "" {
			log.Println("skipping empty row")
			continue
		}

		i, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(string(r.index[:]), "0")), 10, 0)
		if err != nil {
			log.Fatalf("failed to parse index: %x", r.index[:])
		}

		if i > int64(a.highest) {
			a.highest = int(i)
		}

		if a.indexes == nil {
			a.indexes = make(map[int]struct{})
		}

		_, ok := a.indexes[int(i)]
		if ok {
			log.Printf("found duplicate index")
		} else {
			a.indexes[int(i)] = struct{}{}
		}

		an[key] = a

		text := decoder.ConvertString(string(r.data[:]))

		fmt.Printf("block-%10d | %6s | %6s | %3s | %s\n", idx, string(r.client[:]), string(r.animal[:]), string(r.index[:]), text)

		/*
			var input [1]byte
			if _, err := io.ReadFull(os.Stdin, input[:]); err != nil {
				log.Fatal(err)
			}
		*/
	}

	fmt.Printf("\n========\nFound %d records\n", idx)
}

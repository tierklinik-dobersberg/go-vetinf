package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

func getToJSONCOmmand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "to-json",
		Short: "Dump DBF files as JSON",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
		Files:
			for _, file := range args {
				content, err := afero.Afero{Fs: infdat}.ReadFile(file)
				if err != nil {
					log.Printf("%s: %s", file, err)
					continue Files
				}

				db, err := godbf.NewFromByteArray(content, encoding)
				if err != nil {
					log.Printf("%s: %s", file, err)
					continue Files
				}

				base := filepath.Base(file)
				ext := filepath.Ext(file)

				fileName := fmt.Sprintf(
					"%s.json",
					strings.TrimSuffix(base, ext),
				)
				target := filepath.Join(
					filepath.Dir(file),
					fileName,
				)

				f, err := os.Create(target)
				if err != nil {
					log.Printf("%s: %s", file, err)
					continue Files
				}

				encoder := json.NewEncoder(f)
				encoder.SetIndent("", "  ")
				encoder.SetEscapeHTML(false)

				for rowIdx := 0; rowIdx < db.NumberOfRecords(); rowIdx++ {
					m := make(map[string]interface{})

					m["_rowIndex"] = rowIdx
					m["_file"] = base

					/*
						if db.RowIsDeleted(rowIdx) {
							m["_deleted"] = true
						}
					*/

					for fieldIdx, field := range db.Fields() {
						var x interface{}
						var err error

						fieldValue := db.FieldValue(rowIdx, fieldIdx)

						if fieldValue != "" {
							switch field.FieldType() {
							case godbf.Character:
								x = db.FieldValue(rowIdx, fieldIdx)
							case godbf.Logical:
								v := db.FieldValue(rowIdx, fieldIdx)
								switch v {
								case "y", "Y", "j", "J", "1":
									x = true
								default:
									x = false
								}
							case godbf.Numeric:
								x, err = db.Float64FieldValueByName(rowIdx, field.Name())
							case godbf.Float:
								x, err = db.Float64FieldValueByName(rowIdx, field.Name())
							case godbf.Date:
								val := db.FieldValue(rowIdx, fieldIdx)
								x, err = toDate(val)
							}

							if err != nil {
								log.Printf("%s: failed to get field %s for row %d (value: %q): %s", file, field.Name(), rowIdx, db.FieldValue(rowIdx, fieldIdx), err)
								x = db.FieldValue(rowIdx, fieldIdx)
							}

							m[strings.ToLower(field.Name())] = x
						}
					}

					if err := encoder.Encode(m); err != nil {
						log.Printf("%s: failed to write row %d: %s", file, rowIdx, err)
					}
				}

				f.Close()
			}
		},
	}

	return cmd
}

func toDate(val string) (time.Time, error) {
	if len(val) != 8 {
		return time.Time{}, fmt.Errorf("invalid date")
	}
	year := val[0:4]
	month := val[4:6]
	day := val[6:8]

	y, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse '%s' as year: %w", year, err)
	}

	m, err := strconv.ParseInt(month, 10, 32)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse '%s' as month: %w", month, err)
	}

	d, err := strconv.ParseInt(day, 10, 32)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse '%s' as d: %w", day, err)
	}

	t := time.Date(int(y), time.Month(m), int(d), 0, 0, 0, 0, time.UTC)
	return t, nil
}

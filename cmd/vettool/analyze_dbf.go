package main

import (
	"fmt"
	"log"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

func getAnalyzeCommand() *cobra.Command {
	var (
		skipUnused bool
	)

	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze a VetInf dbf file",
		Args:  cobra.ExactArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			content, err := afero.Afero{Fs: infdat}.ReadFile(args[0])
			if err != nil {
				log.Fatal(err)
			}

			db, err := godbf.NewFromByteArray(content, encoding)
			if err != nil {
				log.Fatal(err)
			}

			usage := make(map[int]int, len(db.Fields()))

			for i := 0; i < db.NumberOfRecords(); i++ {
				slice := db.GetRowAsSlice(i)
				for idx, value := range slice {
					if value != "" {
						usage[idx]++
					}
				}
			}

			fmt.Printf("DBF: %s\nTotal Record: %d\n=========================\n\n", args[0], db.NumberOfRecords())

			for idx, field := range db.Fields() {
				if skipUnused && usage[idx] == 0 {
					continue
				}

				fmt.Printf("%s: %s (%d), usage %.1f%% (%d)\n",
					field.Name(),
					field.FieldType().String(),
					field.Length(),
					float64(100)/float64(db.NumberOfRecords())*float64(usage[idx]),
					usage[idx])
			}
		},
	}

	flags := cmd.Flags()
	{
		flags.BoolVar(&skipUnused, "skip-unused", true, "Skip unused fields")
	}

	return cmd
}

package main

import (
	"log"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/tierklinik-dobersberg/go-vetinf/vetinf"
)

var (
	encoding  string
	infdatDir string

	infdat *vetinf.Infdat
)

func main() {
	if err := getRootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func getRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "vettool",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if infdatDir == "" {
				infdatDir = "."
			}

			infdat = vetinf.OpenReadonlyFs(infdatDir, afero.NewOsFs())
			return nil
		},
	}

	flags := cmd.PersistentFlags()
	{
		flags.StringVarP(&encoding, "encoding", "e", "IBM852", "Database encoding")
		flags.StringVarP(&infdatDir, "infdat", "i", os.Getenv("INFDAT_DIR"), "Path to VetInf Infdat")
	}

	cmd.AddCommand(
		getAnalyzeCommand(),
	)

	return cmd
}

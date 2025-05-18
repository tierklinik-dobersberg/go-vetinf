package vetinf

import (
	"errors"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/axgle/mahonia"
	"github.com/spf13/afero"
	"github.com/tierklinik-dobersberg/go-dbf/godbf"
)

// Infdat represents the Infdat folder of a VetInf installation.
// If VetInf is used with multiple hosts, make sure to use the shared
// Infdat directory instead of per-host one.
// Infdat also implements the afero.Fs interface but bound to
// the actual installation directory
type Infdat struct {
	// Fs provides access to the file-system on which the Infdat
	// directory is stored
	afero.Fs

	// RootPath is the path to the Infdat directory stored on Fs
	RootPath string
}

// OpenFs opens an VetInf installation directory located at root
// on fs
func OpenFs(root string, fs afero.Fs) *Infdat {
	return &Infdat{
		RootPath: root,
		Fs:       afero.NewBasePathFs(fs, root),
	}
}

// OpenReadonlyFs is like OpenFs but denies any write operations
func OpenReadonlyFs(root string, fs afero.Fs) *Infdat {
	return OpenFs(root, afero.NewReadOnlyFs(fs))
}

// OpenCachedFs is like OpenReadonlyFs but provides a file cache where each file
// accessed is pulled into an overlay FS (memory mapped) and subsequent reads are performed
// on the cached file (as long as the file does not exceed cacheTime). See
// afero.NewCacheOnReadFs for more information
func OpenCachedFs(root string, cacheTime time.Duration, base afero.Fs) *Infdat {
	layer := afero.NewMemMapFs()
	ufs := afero.NewCacheOnReadFs(base, layer, cacheTime)
	return OpenReadonlyFs(root, ufs)
}

// CustomerDB opens the customer DBase file (vetkldat.dbf)
func (inf *Infdat) CustomerDB(encoding string) (*CustomerDB, error) {
	content, err := afero.Afero{Fs: inf}.ReadFile("vetkldat.dbf")
	if err != nil {
		return nil, err
	}

	db, err := godbf.NewFromByteArray(content, encoding)
	if err != nil {
		return nil, err
	}

	customers := &CustomerDB{
		infdat: inf,
		table:  db,
	}

	if err := customers.buildIndex(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (inf *Infdat) AnimalDB(encoding string) (*AnimalDB, error) {
	content, err := afero.Afero{Fs: inf}.ReadFile("vetktdat.dbf")
	if err != nil {
		return nil, err
	}

	db, err := godbf.NewFromByteArray(content, encoding)
	if err != nil {
		return nil, err
	}

	animals := &AnimalDB{
		infdat: inf,
		table:  db,
	}
	return animals, nil
}

func (inf *Infdat) Vetamdat() (<-chan VetamdatRecord, error) {
	reader, err := afero.Afero{Fs: inf}.Open("./vetamdat")
	if err != nil {
	}

	//encoder := mahonia.NewEncoder("IBM852")
	decoder := mahonia.NewDecoder("IBM852")

	ch := make(chan VetamdatRecord, 100)

	go func() {
		for {
			var r record

			if _, err := io.ReadFull(reader, r.client[:]); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}

				log.Println(err)
				return
			}
			if _, err := io.ReadFull(reader, r.animal[:]); err != nil {
				log.Println(err)
				return
			}
			if _, err := io.ReadFull(reader, r.index[:]); err != nil {
				log.Println(err)
				return
			}
			if _, err := io.ReadFull(reader, r.data[:]); err != nil {
				log.Println(err)
				return
			}

			if strings.TrimSpace(string(r.index[:])) == "" {
				log.Println("skipping empty row")
				continue
			}

			i, err := strconv.ParseInt(strings.TrimSpace(strings.TrimPrefix(string(r.index[:]), "0")), 10, 0)
			if err != nil {
				log.Fatalf("failed to parse index: %x", r.index[:])
			}

			text := decoder.ConvertString(string(r.data[:]))

			amr := VetamdatRecord{
				ClientID: strings.TrimPrefix(string(r.client[:]), "0"),
				AnimalID: strings.TrimPrefix(string(r.animal[:]), "0"),
				Index:    int(i),
				Data:     text,
			}

			ch <- amr
		}
	}()

	return ch, nil
}

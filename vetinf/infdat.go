package vetinf

import (
	"time"

	"github.com/spf13/afero"
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

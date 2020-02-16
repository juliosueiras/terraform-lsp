package memfs

import "github.com/spf13/afero"

var base = afero.NewOsFs()
var roBase = afero.NewReadOnlyFs(base)
var MemFs = afero.NewCopyOnWriteFs(roBase, afero.NewMemMapFs())

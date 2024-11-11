package renamer

import (
	"fmt"
	"strings"
)

type PadRenamer struct {
	pad string
}

func NewPadRenamer(fileCount int) *PadRenamer {
	count := 999
	if fileCount > 0 {
		count = fileCount
	}

	return &PadRenamer{pad: strings.Repeat("0", len(fmt.Sprint(count)))}
}

func (r *PadRenamer) NewName(index int, extension string) string {
	str := fmt.Sprint(index)
	return r.pad[:len(r.pad)-len(str)] + str + extension
}

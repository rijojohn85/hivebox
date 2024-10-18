package hivebox

import (
	"fmt"
	"io"
	"os"
)

const version = "0.0.1"

func PrintVersion(fd ...io.Writer) {
	var out io.Writer
	switch len(fd) {
	case 0:
		out = os.Stdout
	case 1:
		out = fd[0]
	default:
		panic("too many params")
	}
	fmt.Fprint(out, version)
}

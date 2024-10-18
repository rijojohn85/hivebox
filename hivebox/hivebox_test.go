package hivebox

import (
	"bytes"
	"testing"
)

func TestVersion(t *testing.T) {
	buffer := bytes.Buffer{}
	PrintVersion(&buffer)
	got := buffer.String()
	want := version

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

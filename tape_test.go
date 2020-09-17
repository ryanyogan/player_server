package poker

import (
	"io/ioutil"
	"testing"
)

func TestTape_Write(t *testing.T) {
	file, clean := CreateTempFile(t, "123456")
	defer clean()

	tape := &tape{file}
	tape.Write([]byte("abc"))

	file.Seek(0, 0)
	newFileContent, _ := ioutil.ReadAll(file)

	got := string(newFileContent)
	want := "abc"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

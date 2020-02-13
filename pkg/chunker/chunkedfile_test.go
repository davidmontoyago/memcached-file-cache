package chunker

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestSplitsFileIntoRandomSizedChunks(t *testing.T) {
	f, err := os.Open("fixture/file.dat")
	if err != nil {
		t.Error(err)
	}
	// 10MB
	expectedFileSize := 10485760
	expectedChecksum := "91388263e7c545ebea3952fb2637dffa"

	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	chunkedFile := NewFromFile(file)

	totalParts := len(chunkedFile.Chunks())
	if totalParts < 10 {
		t.Errorf("got %d parts but expected at least 10 parts", totalParts)
	}
	var chunksSizeSum int
	for _, chunk := range chunkedFile.Chunks() {
		chunksSizeSum += len(chunk.Bytes())
	}
	if chunksSizeSum != expectedFileSize {
		t.Errorf("got %d from adding up chunks but expected %d", chunksSizeSum, expectedFileSize)
	}
	if err := chunkedFile.Validate(expectedChecksum); err != nil {
		t.Errorf("got checksum %s but expected %s", chunkedFile.Checksum(), expectedChecksum)
	}
}

func TestAssemblesFileFromChunks(t *testing.T) {
	f, err := os.Open("fixture/file.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	chunkedFile := NewFromFile(file)
	chunks := chunkedFile.Chunks()

	testChunkedFile := NewFromChunks(chunks)
	assembledFile := testChunkedFile.File()
	if !bytes.Equal(file, assembledFile) {
		t.Errorf("expected assembled file to contain same bytes as fixture but they differ")
	}
	if err := testChunkedFile.Validate("91388263e7c545ebea3952fb2637dffa"); err != nil {
		t.Error(err)
	}
}

func TestHandlesAnEmptyFile(t *testing.T) {
	f, err := os.Open("fixture/empty.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	chunkedFile := NewFromFile(file)
	chunks := chunkedFile.Chunks()

	if len(chunks) > 0 {
		t.Errorf("expected file to have no chunks but has %d", len(chunks))
	}
}

package chunker

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestSplitsFileIntoRandomSizedChunks(t *testing.T) {
	chunker := &Chunker{}
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
	chunkedFile := chunker.Split(file)

	totalParts := len(chunkedFile.Parts)
	if totalParts < 10 {
		t.Errorf("got %d parts but expected at least 10 parts", totalParts)
	}
	var chunksSizeSum int
	for _, chunk := range chunkedFile.Parts {
		chunksSizeSum += len(chunk.Bytes)
	}
	if chunksSizeSum != expectedFileSize {
		t.Errorf("got %d from adding up chunks but expected %d", chunksSizeSum, expectedFileSize)
	}
	if chunkedFile.Checksum != expectedChecksum {
		t.Errorf("got checksum %s but expected %s", chunkedFile.Checksum, expectedChecksum)
	}
}

func TestAssemblesFileFromChunks(t *testing.T) {
	chunker := &Chunker{}
	f, err := os.Open("fixture/file.dat")
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	expectedChecksum := "91388263e7c545ebea3952fb2637dffa"

	chunkedFile := chunker.Split(file)

	assembledFile := chunker.Assemble(chunkedFile)
	if !bytes.Equal(file, assembledFile) {
		t.Errorf("expected assembled file to contain same bytes as fixture but they differ")
	}
	checksum := fmt.Sprintf("%x", md5.Sum(assembledFile))
	if chunkedFile.Checksum != expectedChecksum {
		t.Errorf("got checksum %s but expected %s", checksum, expectedChecksum)
	}
}

package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testDumpExport(t *testing.T) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", serverUser, serverPassword, serverHost, serverPort, serverDatabase)

	savePath := "/tmp/dump.sql.gz"
	
	if err := os.Remove(savePath); err != nil {
		panic(err)
	}

	saveFile, err := os.Create(savePath)
	if err != nil {
		t.Fatal(err.Error())
	}

	defer func() {

		if err := saveFile.Close(); err != nil {
			panic(err)
		}
		if err := os.Remove(savePath); err != nil {
			panic(err)
    }
	}()

	dump := Dump{}

	// Test for pg_dump presence
	assert.True(t, dump.CanExport())

	// Test full db dump
	err = dump.Export(url, saveFile)
	assert.NoError(t, err)

	// Test nonexistent database
	invalidURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", serverUser, serverPassword, serverHost, serverPort, "foobar")
	err = dump.Export(invalidURL, saveFile)
	assert.Contains(t, err.Error(), `database "foobar" does not exist`)

	// Test dump of non existent db
	dump = Dump{Table: "foobar"}
	err = dump.Export(url, saveFile)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no matching tables were found")

	// Should drop "search_path" param from URI
	dump = Dump{}
	searchPathURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=private", serverUser, serverPassword, serverHost, serverPort, serverDatabase)
	err = dump.Export(searchPathURL, saveFile)
	assert.NoError(t, err)
}

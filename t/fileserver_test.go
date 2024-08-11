// As this is a side project, these tests aren't super extensive.
// Just checking that there aren't horrible errors.

package fileserver_test

import (
	"net/rpc"
	"os"
	"path/filepath"
	"testing"

	"github.com/jonasiwnl/distributed-fileserver/v2/fileserver"
)

// Share global client for all tests.
var client *rpc.Client

func TestMain(m *testing.M) {
	quit := make(chan bool, 1)
	go fileserver.Start(quit)

	var err error
	client, err = rpc.Dial("tcp", "localhost"+fileserver.PORT)
	if err == nil {
		m.Run()
	}
	client.Close()
	quit <- true
}

func TestDir(t *testing.T) {
	testDir := "testdir"
	directoryPath := filepath.Join(fileserver.DIRECTORY, testDir)

	args := &fileserver.DirArgs{Path: testDir, Mode: 0755}
	var reply bool
	err := client.Call("FileServer.MakeDirectory", args, &reply)
	if err != nil {
		t.Fatal("making directory: ", err)
	}
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) || !reply {
		t.Fatal("directory not created")
	}

	err = client.Call("FileServer.RemoveDirectory", args, &reply)
	if err != nil {
		t.Fatal("removing directory: ", err)
	}
	if _, err := os.Stat(directoryPath); err == nil || !reply {
		t.Fatal("directory not removed")
	}

	err = client.Call("FileServer.RemoveDirectory", args, &reply)
	if err != nil {
		t.Fatal("error removing non-existent directory: ", err)
	}
}

func TestFile(t *testing.T) {
	testFile := "testfile"
	filePath := filepath.Join(fileserver.DIRECTORY, testFile)

	args := &fileserver.FileArgs{Path: testFile, Data: []byte("test"), Mode: 0644}
	var reply bool
	err := client.Call("FileServer.WriteFile", args, &reply)
	if err != nil {
		t.Fatal("writing file: ", err)
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) || !reply {
		t.Fatal("file not created")
	}

	err = client.Call("FileServer.RemoveFile", args, &reply)
	if err != nil {
		t.Fatal("removing file: ", err)
	}
	if _, err := os.Stat(filePath); err == nil || !reply {
		t.Fatal("file not removed")
	}
}

func TestListDir(t *testing.T) {
}

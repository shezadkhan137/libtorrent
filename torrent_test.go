package libtorrent

import (
	//"bytes"
	"github.com/op/go-logging"
	"github.com/torrance/libtorrent/metainfo"
	//"io"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"
)

func loadTorrentFile(t *testing.T) *os.File {
	//torrentFile, err := os.Open("testData/test.txt.torrent")
	//tstring := "testData/test.txt.torrent"
	tstring := "testData/multitest.torrent"
	//tstring := "/Users/shaz/Downloads/ubuntu-12.04.3-server-amd64.iso.torrent"
	torrentFile, err := os.Open(tstring)
	if err != nil {
		t.Fatal("Could not open torrent file: ", err)
	}
	return torrentFile
}

// func TestNewTorrent(t *testing.T) {
// 	torrentFile := loadTorrentFile(t)

// 	minfo, err := metainfo.ParseMetainfo(torrentFile)

// 	if err != nil {
// 		t.Fatal("Could not parse torrent into metainfo", err)
// 	}

// 	tmpDir, err := ioutil.TempDir("", "libtorrentTesting")
// 	if err != nil {
// 		t.Fatal("Could not create temporary directory to run tests: ", err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	config := &Config{RootDirectory: tmpDir}
// 	_, err = NewTorrent(minfo, config)
// 	if err != nil {
// 		t.Error("Could not create torrent from valid metainfo file: ", err)
// 	}
// }

// func TestAgainstTransmission(t *testing.T) {
// 	logging.SetLevel(logging.DEBUG, "libtorrent.torrent")

// 	torrentFile := loadTorrentFile(t)

// 	tmpDir, err := ioutil.TempDir("", "libtorrentTesting")
// 	if err != nil {
// 		t.Fatal("Could not create temporary directory to run tests: ", err)
// 	}
// 	defer os.RemoveAll(tmpDir)

// 	testFile, _ := os.Create(tmpDir + string(os.PathSeparator) + "test.txt")
// 	originalFile, _ := os.Open("testData/test.txt")
// 	io.Copy(testFile, originalFile)

// 	minfo, err := metainfo.ParseMetainfo(torrentFile)
// 	if err != nil {
// 		t.Fatal("Could not parse torrent into metainfo", err)
// 	}

// 	config := &Config{RootDirectory: tmpDir}
// 	tor, err := NewTorrent(minfo, config)
// 	if err != nil {
// 		t.Fatal("Could not create torrent from valid metainfo file: ", err)
// 	}

// 	if !bytes.Equal(tor.bitf.Bytes(), []byte{0xC0}) {
// 		t.Fatalf("Torrent data was not loaded and validated, got: %#x", tor.bitf)
// 	}

// 	conn, err := net.Dial("tcp", "localhost:51413")
// 	if err != nil {
// 		t.Fatal("Failed to create connection: ", err)
// 	}

// 	tor.Start()
// 	tor.AddPeer(conn, nil)
// 	time.Sleep(60 * time.Second)
// 	panic("Hello World")
// }

func TestDownloadAgainstTransmission(t *testing.T) {
	logging.SetLevel(logging.DEBUG, "libtorrent.torrent")

	torrentFile := loadTorrentFile(t)

	tmpDir, err := ioutil.TempDir("", "libtorrentTesting")
	if err != nil {
		t.Fatal("Could not create temporary directory to run tests: ", err)
	}
	logger.Debug(tmpDir)
	defer os.RemoveAll(tmpDir)

	//testFile, _ := os.Create(tmpDir + string(os.PathSeparator) + "test.txt")
	//originalFile, _ := os.Open("testData/test.txt")
	//io.Copy(testFile, originalFile)

	minfo, err := metainfo.ParseMetainfo(torrentFile)
	if err != nil {
		t.Fatal("Could not parse torrent into metainfo", err)
	}

	config := &Config{RootDirectory: tmpDir}
	tor, err := NewTorrent(minfo, config)
	if err != nil {
		t.Fatal("Could not create torrent from valid metainfo file: ", err)
	}

	// if !bytes.Equal(tor.bitf.Bytes(), []byte{0xC0}) {
	// 	t.Fatalf("Torrent data was not loaded and validated, got: %#x", tor.bitf)
	// }

	conn, err := net.Dial("tcp", "localhost:51413")
	if err != nil {
		t.Fatal("Failed to create connection: ", err)
	}

	tor.Start()
	tor.AddPeer(conn, nil)
	time.Sleep(120 * time.Second)
	//panic("Hello World")
}

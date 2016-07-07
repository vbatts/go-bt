package torrent

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/vbatts/go-bt/bencode"
)

func TestDecode(t *testing.T) {
	data := new(map[string]interface{})
	tData, err := ioutil.ReadFile("./testdata/farts.torrent")
	if err != nil {
		t.Fatal(err)
	}

	// currently failing in ./bencode/struct.go:134
	if err := bencode.Unmarshal(tData, &data); err != nil {
		t.Error(err)
	}

	f, err := DecocdeTorrentData(data)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", f)

}

func TestFileMarshal(t *testing.T) {
	f1 := File{
		Announce: "http://foo.bar.com:9090/announce",
		AnnounceList: []string{"http://foo.bar.com:9091/announce",
			"http://foo.bar.com:9092/announce",
			"http://foo.bar.com:9093/announce",
		},
	}

	buf, err := bencode.Marshal(f1)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Contains(buf, []byte("omitempty")) || bytes.Contains(buf, []byte("created by")) {
		t.Errorf("should not have the string 'omitempty' or 'created by' in %q", buf)
	}

	f2 := File{}
	err = bencode.Unmarshal(buf, &f2)
	if err != nil {
		t.Fatal(err)
	}

	if f1.Announce != f2.Announce {
		t.Errorf("expected %q, got %q", f1.Announce, f2.Announce)
	}
	if len(f1.AnnounceList) != len(f2.AnnounceList) {
		t.Errorf("expected %q, got %q", len(f1.AnnounceList), len(f2.AnnounceList))
	}
}

func TestTime(t *testing.T) {
	f1 := File{}
	if f1.CreationDateTime().Unix() != 0 {
		t.Errorf("%s -- %d", f1.CreationDateTime(), f1.CreationDateTime().Unix())
	}
}

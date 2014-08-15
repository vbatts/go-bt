package torrent

import (
	"bytes"
	"github.com/vbatts/go-bt/bencode"
	"testing"
)

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

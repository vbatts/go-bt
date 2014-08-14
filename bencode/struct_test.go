package bencode

import (
	"testing"
)

type testStruct struct {
	Field1 string `bencode:"my field1"`
	Field2 int64  `bencode:"my field2"`
}

type testOldTag struct {
	Field1 string "my field1"
	Field2 int64  "my field2"
}

func TestMarshalling(t *testing.T) {
	ts1 := testStruct{"foo", 123456}
	buf, err := Marshal(ts1)
	if err != nil {
		t.Fatal(err)
	}

	ts2 := testStruct{}
	err = Unmarshal(buf, &ts2)
	if err != nil {
		t.Fatal(err)
	}
	if ts1.Field1 != ts2.Field1 {
		t.Errorf("expected %q, got %q", ts1.Field1, ts2.Field1)
	}
	if ts1.Field2 != ts2.Field2 {
		t.Errorf("expected %q, got %q", ts1.Field2, ts2.Field2)
	}
}

func TestOldMarshalling(t *testing.T) {
	ts1 := testOldTag{"foo", 123456}
	buf, err := Marshal(ts1)
	if err != nil {
		t.Fatal(err)
	}

	ts2 := testStruct{}
	err = Unmarshal(buf, &ts2)
	if err != nil {
		t.Fatal(err)
	}
	if ts1.Field1 != ts2.Field1 {
		t.Errorf("expected %q, got %q", ts1.Field1, ts2.Field1)
	}
	if ts1.Field2 != ts2.Field2 {
		t.Errorf("expected %q, got %q", ts1.Field2, ts2.Field2)
	}
}

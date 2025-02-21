package eagle

import (
	"testing"
)

func TestMtimeCount(t *testing.T) {
	e, err := New()
	if err != nil {
		t.Fatalf("err %s", err)
	}

	cnt, err := e.MtimeCount()
	t.Log(cnt)
	if err != nil {
		t.Fatalf("err %s", err)
	}
}

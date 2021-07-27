package sdk

import "testing"

func TestModelCommand(t *testing.T) {
	var r responsePacket

	if want, res := false, r.validPrefix(); want != res {
		t.Fatalf("want %t, got %t", want, res)
	}

	if want, res := false, r.validSize(); want != res {
		t.Fatalf("want %t, got %t", want, res)
	}

	if want, res := false, r.validCmdCode(); want != res {
		t.Fatalf("want %t, got %t", want, res)
	}

	if want, res := false, r.validResCode(); want != res {
		t.Fatalf("want %t, got %t", want, res)
	}

	if want, res := 0, r.size(); want != res {
		t.Fatalf("want %d, got %d", want, res)
	}

	if want, res := false, r.matchWith(nil); want != res {
		t.Fatalf("want %t, got %t", want, res)
	}

	cmd, _ := getCmdByCode(0, 0)
	if want, res := false, r.matchWith(cmd); want != res {
		t.Fatalf("want %t, got %t", want, res)
	}
}

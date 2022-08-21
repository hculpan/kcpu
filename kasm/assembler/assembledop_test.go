package assembler

import "testing"

func TestNewAssembledOp(t *testing.T) {
	op := NewAssembledOp(0, 1, 34, 82, "")
	if op.ToString() != "00:01:22:52" {
		t.Fatalf("Expected '00:01:22:52', got '%s'", op.ToString())
	}
}

func TestNewAssembledOpWithAddress(t *testing.T) {
	op := NewAssembledOpWithAddress(0, 1, 0x0D9A, "")
	if op.ToString() != "00:01:0D:9A" {
		t.Fatalf("Expected '00:01:OD:9A', got '%s'", op.ToString())
	}
}

func TestGetAddressFromData(t *testing.T) {
	op := NewAssembledOpWithAddress(0, 1, 3482, "")
	if op.ToString() != "00:01:0D:9A" {
		t.Fatalf("Expected '00:01:OD:9A', got '%s'", op.ToString())
	}
	if op.GetDataAsAddress() != 3482 {
		t.Fatalf("Expected address '0D9A', got '%04X'", op.GetDataAsAddress())
	}
}

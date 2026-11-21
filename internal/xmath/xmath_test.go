package xmath

import "testing"

func TestXmath(t *testing.T) {
	if Mod(-1, 5) != 4 {
		t.Error("Mod")
	}
	if GCD(12, 18) != 6 || LCM(4, 6) != 12 {
		t.Error("GCD/LCM")
	}
	if Abs(-3) != 3 || Sign(-3) != -1 {
		t.Error("Abs/Sign")
	}
}

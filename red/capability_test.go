package red

import (
	"testing"
)

func TestCapabilityNew(t *testing.T) {
	var v Capability
	if v.Test(0) {
		t.Errorf("Unable to make a capability and read its 0th value.")
	}
}

func TestCapabilityIsClear(t *testing.T) {
	var v Capability
	for i := uint32(0); i < 32; i++ {
		if v.Test(i) {
			t.Errorf("Bit %d is set, and it shouldn't be.", i)
		}
	}
}

func TestCapabilitySetTo(t *testing.T) {
	var v Capability
	v.SetTo(10, true)
	if !v.Test(10) {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 10)
	}
	v.SetTo(10, false)
	if v.Test(10) {
		t.Errorf("Bit %d is set, and it shouldn't be.", 10)
	}
}

func TestCapabilitySetAndGet(t *testing.T) {
	var v Capability
	v.Set(10)
	if !v.Test(10) {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 10)
	}
}

func TestCapabilityChain(t *testing.T) {
	var v Capability
	if !v.Set(10).Set(9).Clear(9).Test(10) {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 10)
	}
}

func TestCapabilityFlip(t *testing.T) {
	var v Capability
	v.SetTo(10, false)
	v.Flip(10)
	if !v.Test(10) {
		t.Errorf("Bit %d is clear, and it shouldn't be.", 10)
	}
	v.SetTo(10, true)
	v.Flip(10)
	if v.Test(10) {
		t.Errorf("Bit %d is set, and it shouldn't be.", 10)
	}
}

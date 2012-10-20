package pastry

import (
	"math/big"
	"testing"
)

// Make sure the NodeIDDigits returned by NodeIDDigitsFromByte actually add up to equal the original byte.
func TestNodeIDDigitsFromByteEqualsByte(t *testing.T) {
	bytes := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-_ ")
	for _, b := range bytes {
		d1, d2 := NodeIDDigitsFromByte(b)
		if uint8(d1)+uint8(d2) != uint8(b) {
			t.Errorf("%v + %v should equal %v, but instead equals %v.", uint8(d1), uint8(d2), uint8(b), uint8(d1)+uint8(d2))
		}
	}
}

// Make sure an error is thrown if a NodeID is created from less than 32 bytes
func TestNodeIDFromBytesWithInsufficientBytes(t *testing.T) {
	bytes := []byte("123456789012345")
	id, err := NodeIDFromBytes(bytes)
	if err == nil {
		t.Errorf("Source length of %v bytes, but no error thrown. Instead returned NodeID of %v", len(bytes), id)
	}
}

// Make sure NodeIDDigits discard their insignificant bits when comparing for equality
func TestNodeIDDigitEqualsDiscardsInsignificantBits(t *testing.T) {
	d1 := NodeIDDigit(0xf)
	d2 := NodeIDDigit(0xf0)
	if !d1.Equals(d2) {
		t.Errorf("%s should equal %s, but it doesn't.", d1, d2)
	}
}

// Make sure an error is *not* thrown if enough bytes are passed in.
func TestNodeIDFromBytesWithSufficientBytes(t *testing.T) {
	bytes := []byte("1234567890123456")
	_, err := NodeIDFromBytes(bytes)
	if err != nil {
		t.Errorf("Source length of %v bytes threw an error when no error should have been thrown.", len(bytes))
		t.Logf(err.Error())
	}
}

// Make sure the correct common prefix length is reported for two NodeIDs
func TestNodeIDCommonPrefixLen(t *testing.T) {
	n1 := NodeID([]NodeIDDigit{0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0})
	n2 := NodeID([]NodeIDDigit{0xf, 0xd0, 0xf, 0xd0, 0xd0, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0})
	diff1 := 4

	n3 := NodeID([]NodeIDDigit{0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf})
	n4 := NodeID([]NodeIDDigit{0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xa, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf})
	diff2 := 6

	if n1.CommonPrefixLen(n2) != diff1 {
		t.Errorf("Common prefix length should be %v, is %v instead.", diff1, n1.CommonPrefixLen(n2))
		t.Log(n1)
		t.Log(n2)
		if len(n1) > n1.CommonPrefixLen(n2) && len(n2) > n1.CommonPrefixLen(n2) {
			t.Logf("First significant digit: %v vs. %v", n1[n1.CommonPrefixLen(n2)], n2[n1.CommonPrefixLen(n2)])
		}
	}
	if n2.CommonPrefixLen(n3) != 0 {
		t.Errorf("Common prefix length should be %v, is %v instead.", 0, n2.CommonPrefixLen(n3))
		t.Log(n2)
		t.Log(n3)
		if len(n2) > n2.CommonPrefixLen(n3) && len(n3) > n2.CommonPrefixLen(n3) {
			t.Logf("First significant digit: %v vs. %v", n2[n2.CommonPrefixLen(n3)], n3[n2.CommonPrefixLen(n3)])
		}
	}
	if n3.CommonPrefixLen(n4) != diff2 {
		t.Errorf("Common prefix length should be %v, is %v instead.", diff2, n3.CommonPrefixLen(n4))
		t.Log(n3)
		t.Log(n4)
		if len(n3) > n3.CommonPrefixLen(n4) && len(n4) > n3.CommonPrefixLen(n4) {
			t.Logf("First significant digit: %v vs. %v", n3[n3.CommonPrefixLen(n4)], n4[n3.CommonPrefixLen(n4)])
		}
	}
	if n4.CommonPrefixLen(n4) != len(n4) {
		t.Errorf("Common prefix length should be %v, is %v instead.", len(n4), n4)
		if n4.CommonPrefixLen(n4) < len(n4) {
			t.Logf("First significant digit: %v vs. %v", n4[n4.CommonPrefixLen(n4)], n4[n4.CommonPrefixLen(n4)])
		}
	}
}

// Make sure the correct difference is reported between NodeIDs
func TestNodeIDDiff(t *testing.T) {
	n1 := NodeID([]NodeIDDigit{0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0})
	n2 := NodeID([]NodeIDDigit{0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xd0, 0xf, 0xb0})
	diff1 := n1.Diff(n2)
	if diff1.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("Difference should be 2, was %v instead", diff1)
	}
	diff2 := n2.Diff(n1)
	if diff2.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("Difference should be 2, was %v instead", diff2)
	}
	diff3 := n2.Diff(n2)
	if diff3.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Difference should be 0, was %v instead", diff3)
	}
}

// Make sure NodeID comparisons wrap around the circle
func TestNodeIDDiffWrap(t *testing.T) {
	n1, err := NodeIDFromBytes(make([]byte, 16))
	if err != nil {
		t.Fatalf(err.Error())
	}
	n2, err := NodeIDFromBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255})
	if err != nil {
		t.Fatalf(err.Error())
	}
	diff1 := n1.Diff(n2)
	if diff1.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Difference should be 1, was %v instead", diff1)
	}
	diff2 := n2.Diff(n1)
	if diff2.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("Difference should be 1, was %v instead", diff2)
	}
	diff3 := n2.Diff(n2)
	if diff3.Cmp(big.NewInt(0)) != 0 {
		t.Errorf("Difference should be 0, was %v instead", diff3)
	}
}

// Quick benchmark to test how expensive diffing nodes is
func BenchmarkNodeIDDiff(b *testing.B) {
	b.StopTimer()
	n1, err := NodeIDFromBytes(make([]byte, 16))
	if err != nil {
		b.Fatalf(err.Error())
	}
	n2, err := NodeIDFromBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255})
	if err != nil {
		b.Fatalf(err.Error())
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		n1.Diff(n2)
	}
}

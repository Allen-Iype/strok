package domain

import "testing"

func TestFingerShortName(t *testing.T) {
	seen := map[string]bool{}
	for f := Finger(0); int(f) < FingerCount; f++ {
		short := f.ShortName()
		if short == "" || short == "?" {
			t.Errorf("finger %d has empty/unknown ShortName", f)
		}
		if len(short) > len(f.String()) {
			t.Errorf("finger %d ShortName %q longer than full name %q", f, short, f.String())
		}
		if seen[short] {
			t.Errorf("duplicate ShortName %q", short)
		}
		seen[short] = true
	}

	if got := Finger(-1).ShortName(); got != "?" {
		t.Errorf("out-of-range ShortName = %q, want ?", got)
	}
}

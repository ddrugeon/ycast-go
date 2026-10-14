package station

import "testing"

func TestChecksumMatchesPython(t *testing.T) {
	// Expected values generated with the original my_stations.get_checksum().
	cases := [][3]string{
		{"SomaFM Fluid", "http://ice1.somafm.com/fluid-128-mp3", "5A52DE6AB2D1"},
		{"Radio d'Alsace & Co", "http://example.org/a?b=1&c=2", "D4C68FA881E8"},
		{"Café Zoé", "http://exämple.org/stream", "B82502BD21C3"},
	}
	for _, c := range cases {
		if got := MyChecksum(c[0], c[1]); got != c[2] {
			t.Errorf("checksum(%q,%q) = %s, want %s", c[0], c[1], got, c[2])
		}
	}
}

func TestSplitStationID(t *testing.T) {
	for _, id := range []string{"", "RB", "RB_", "../../etc/passwd", "RB/abcdef", "RBxabcdef", "RB_a b c"} {
		if _, _, ok := SplitID(id); ok {
			t.Errorf("SplitID(%q) accepted an invalid id", id)
		}
	}
	p, r, ok := SplitID("RB_1234-abcd")
	if !ok || p != "RB" || r != "1234-abcd" {
		t.Errorf("SplitID valid id: %q %q %v", p, r, ok)
	}
}

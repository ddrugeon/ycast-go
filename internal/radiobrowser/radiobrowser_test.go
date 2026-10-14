package radiobrowser

import "testing"

func TestPyCase(t *testing.T) {
	if got := pyTitle("swiss german-ish"); got != "Swiss German-Ish" {
		t.Errorf("pyTitle = %q", got)
	}
	if got := pyCapitalize("hIP hop"); got != "Hip hop" {
		t.Errorf("pyCapitalize = %q", got)
	}
	if got := pyCapitalize(""); got != "" {
		t.Errorf("pyCapitalize(empty) = %q", got)
	}
}

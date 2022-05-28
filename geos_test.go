package geos

import "testing"

func skipIfVersionLessThan(t *testing.T, major, minor, patch int) {
	t.Helper()
	if !versionEqualOrGreaterThan(major, minor, patch) {
		t.Skip()
	}
}

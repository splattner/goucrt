// Package spectest provides small testing helpers for checking goucrt's hand-maintained protocol
// constants and golden fixtures against the vendored Core-API spec (see internal/spec).
package spectest

import (
	"sort"
	"testing"
)

// AssertSubset fails the test if any value in got is absent from allowed, naming exactly which
// value(s) are extra. It does not fail when allowed contains values absent from got - the spec is
// free to define more than goucrt currently implements; that is tracked separately, not a defect.
//
// Use this for anything the spec defines as a closed, structured enum (currently: entity `features`
// and `device_class`) - a value goucrt declares that the spec doesn't recognize is almost always a
// typo or a stale/renamed value, exactly the class of bug this package exists to catch.
func AssertSubset(t *testing.T, label string, got []string, allowed []string) {
	t.Helper()

	allowedSet := make(map[string]bool, len(allowed))
	for _, v := range allowed {
		allowedSet[v] = true
	}

	var extra []string
	for _, v := range got {
		if !allowedSet[v] {
			extra = append(extra, v)
		}
	}

	if len(extra) > 0 {
		sort.Strings(extra)
		t.Errorf("%s: value(s) %v not present in the spec's enum %v - check for a typo or a stale/renamed value", label, extra, allowed)
	}
}

// LogMissing logs (does not fail) any value present in fromSpec but absent from known, prefixed with
// label. Use this where a gap represents an unimplemented feature rather than broken existing
// behavior - e.g. a cmd_id or attribute a newer spec example uses that goucrt has no constant for
// yet. It surfaces the gap for a human decision without blocking the build on work nobody asked for
// in this pass.
func LogMissing(t *testing.T, label string, fromSpec []string, known map[string]bool) {
	t.Helper()

	var missing []string
	for _, v := range fromSpec {
		if !known[v] {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Logf("%s: spec example uses %v, which goucrt has no matching constant for yet (tracked separately, not a failure)", label, missing)
	}
}

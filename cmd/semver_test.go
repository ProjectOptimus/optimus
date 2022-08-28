package cmd

import "testing"

func TestGetSemver(t *testing.T) {
	var want, got, s string // 's' is the version string from Rhadfile

	s = "1.0.0.0"
	want = "v1.0.0"
	got = getSemver(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}

	s = "1.0-alpha"
	want = "v1.0.0-alpha"
	got = getSemver(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}

	s = "1.0.2+abc"
	want = "v1.0.2+abc"
	got = getSemver(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}

	// No change on ideally-provided semver format
	s = "v1.1.9-prebeta1+abc"
	want = s
	got = getSemver(s)
	if want != got {
		t.Errorf("Expected version string '%s' to become '%s', but got '%s'\n", s, want, got)
	}
}

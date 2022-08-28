package cmd

import "testing"

func TestCleanSectionName(t *testing.T) {
	var want, got string

	want = "[.]"
	got = cleanSectionName("[.]")
	if want != got {
		t.Errorf("Want: %v, Got: %v\n", want, got)
	}

	want = "[.]"
	got = cleanSectionName("[  .  ]")
	if want != got {
		t.Errorf("Want: %v, Got: %v\n", want, got)
	}

	want = "[dir1]"
	got = cleanSectionName("[  ./dir1    ]")
	if want != got {
		t.Errorf("Want: %v, Got: %v\n", want, got)
	}
}

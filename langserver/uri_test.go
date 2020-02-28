package langserver

import (
	"testing"
)

func TestAbsolutePath(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			path:     "file:///C:/Users/acmeuser/git/tf-test/main.tf",
			expected: "C:/Users/acmeuser/git/tf-test/main.tf",
		},
		{
			path:     "file:///etc/fstab",
			expected: "/etc/fstab",
		},
	}

	for _, test := range tests {
		uri, err := absolutePath(test.path)
		if err != nil {
			t.Errorf("Parsing path '%s' resulted in error: %s", test.path, err)
		}
		res := uri.Filename()

		if res != test.expected {
			t.Errorf("Unexpected result. Got '%s', Expected '%s'", res, test.expected)
		}
	}
}

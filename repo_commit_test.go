package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_Repository_getCommitIdOfPackedRef(t *testing.T) {
	tests := []struct {
		packedRefs  string
		refPath     string
		expCommitID string
		expError    error
	}{{
		packedRefs: `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa refs/tags/master
`,
		refPath:     "refs/tags/master",
		expCommitID: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
	}, {
		packedRefs: `aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa refs/tags/master-123
`,
		refPath:  "refs/tags/master",
		expError: RefNotFound("refs/tags/master"),
	}}
	for _, test := range tests {
		func() {
			repoDir := tmpRepoDir(t, test.packedRefs)
			defer os.RemoveAll(repoDir)

			r := &Repository{Path: repoDir}
			commitID, err := r.getCommitIdOfPackedRef(test.refPath)
			if string(commitID) != test.expCommitID {
				t.Errorf("expected commit %q, got %q", test.expCommitID, string(commitID))
			}
			if err != test.expError {
				t.Errorf("expected error %q, got %q", test.expError, err)
			}
		}()
	}
}

func tmpRepoDir(t *testing.T, packedRefs string) string {
	d, err := ioutil.TempDir("", "Test_Repository_getCommitIdOfPackedRef")
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(d, "packed-refs"), []byte(packedRefs), os.FileMode(0777)); err != nil {
		t.Fatal(err)
	}
	return d
}

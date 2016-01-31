package git

import (
	"testing"
	"time"
)

func openTestRepo(t *testing.T) *Repository {
	r, err := OpenRepository("testdata/repo")
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestObject(t *testing.T) {
	r := openTestRepo(t)
	testObject(t, r, "30d74d258442c7c65512eafab474568dd706c430", "test")
	testObject(t, r, "d76bde4f5d1ed609dc82d8cd7d216d893830f1c9", "test unpacked")
}

func testObject(t *testing.T, r *Repository, id string, content string) {
	o, err := r.Object(ObjectIDHex(id))
	if err != nil {
		t.Fatal(err)
	}
	if o.Type != ObjectBlob {
		t.Error("wrong type")
	}
	if o.Size != uint64(len(content)) {
		t.Error("wrong size")
	}
	if string(o.Data) != content {
		t.Error("wrong content")
	}
}

func TestTree(t *testing.T) {
	r := openTestRepo(t)
	tree := NewTree(r, ObjectIDHex("095a057d4a651ec412d06b59e32e9b02871592d5"))
	entries, err := tree.ListEntries()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Error("wrong number of entries")
	}
	e := entries[0]
	if e.Id != ObjectIDHex("30d74d258442c7c65512eafab474568dd706c430") {
		t.Error("wrong id")
	}
	if e.Name() != "test.txt" {
		t.Error("wrong name")
	}
	if e.Size() != 4 {
		t.Error("wrong size")
	}
	if e.EntryMode() != 0100644 {
		t.Error("wrong mode")
	}
}

func TestCommit(t *testing.T) {
	r := openTestRepo(t)
	c, err := r.GetCommit("8b61789a76de9edaa49b2529d3aaa302ba238c0b")
	if err != nil {
		t.Fatal(err)
	}
	if c.Message() != "test commit\n" {
		t.Error("wrong message")
	}
	if c.TreeId() != ObjectIDHex("095a057d4a651ec412d06b59e32e9b02871592d5") {
		t.Error("wrong tree")
	}
	if c.Author.Name != "Test Author" || c.Author.Email != "author@example.com" || !c.Author.When.Equal(parseTimestamp("Thu, 07 Apr 2005 22:13:13 +0200")) {
		t.Error("wrong author")
	}
	if c.Committer.Name != "Test Committer" || c.Committer.Email != "committer@example.com" || !c.Committer.When.Equal(parseTimestamp("Thu, 07 Apr 2005 22:13:14 +0200")) {
		t.Error("wrong committer")
	}
}

func TestBranch(t *testing.T) {
	r := openTestRepo(t)
	c, err := r.GetCommitOfBranch("master")
	if err != nil {
		t.Fatal(err)
	}
	if c.Id != ObjectIDHex("40b7c29973f5ff265a241f29c8154fa05594454f") {
		t.Error("wrong commit")
	}
}

func parseTimestamp(s string) time.Time {
	t, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		panic(err)
	}
	return t
}

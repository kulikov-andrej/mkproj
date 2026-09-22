package installer

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
	"testing"
)

func TestFindAsset(t *testing.T) {
	oldGOOS := runtimeGOOS
	runtimeGOOS = "linux"
	t.Cleanup(func() { runtimeGOOS = oldGOOS })

	got, err := findAsset(release{Assets: []asset{
		{Name: "other", Digest: "sha256:abc"},
		{Name: "mkproj-linux-amd64", URL: "download", Digest: "sha256:abc"},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "mkproj-linux-amd64" {
		t.Fatalf("asset = %#v", got)
	}
}

func TestVerify(t *testing.T) {
	path := t.TempDir() + string(os.PathSeparator) + "download"
	content := []byte("mkproj")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}

	hash := sha256.Sum256(content)
	if err := verify(path, "sha256:"+hex.EncodeToString(hash[:])); err != nil {
		t.Fatal(err)
	}
	if err := verify(path, "sha256:"+strings.Repeat("0", 64)); err == nil {
		t.Fatal("verify() error = nil, want mismatch")
	}
}

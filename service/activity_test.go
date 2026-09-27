package service

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"awesomeProject/model"
)

func TestUnsharedActivityImagePaths(t *testing.T) {
	t.Setenv("DOMAIN_NAME", "localhost:8080")
	cover := "https://mobile.xupt.edu.cn/res/cover.png"
	content := `<img src="/res/body.jpeg?version=1"><img src='http://localhost:8080/res/shared.webp'>`
	shared := "/res/shared.webp"
	activity := &model.ActivityPO{AID: 1, Img: &cover, Content: &content}
	other := &model.ActivityPO{AID: 2, Img: &shared}

	got := unsharedActivityImagePaths(activity, []*model.ActivityPO{activity, other})
	want := map[string]struct{}{
		filepath.Join("res", "cover.png"): {},
		filepath.Join("res", "body.jpeg"): {},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unshared images = %v, want %v", got, want)
	}
}

func TestActivityImagePathsIgnoresUnsafeAndExternalPaths(t *testing.T) {
	content := `<img src="/res/../outside.jpg"><img src="/res/static/default.jpg">` +
		`<img src="https://example.com/res/external.png"><img src="/videos/movie.jpg">` +
		`<img src="/res/document.txt"><img src="/res/valid.jpg">`
	activity := &model.ActivityPO{Content: &content}

	got := activityImagePaths(activity)
	want := map[string]struct{}{filepath.Join("res", "valid.jpg"): {}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("image paths = %v, want %v", got, want)
	}
}

func TestRemoveActivityImages(t *testing.T) {
	directory := t.TempDir()
	removePath := filepath.Join(directory, "remove.jpg")
	keepPath := filepath.Join(directory, "keep.jpg")
	for _, filePath := range []string{removePath, keepPath} {
		if err := os.WriteFile(filePath, []byte("image"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := removeActivityImages(map[string]struct{}{removePath: {}, filepath.Join(directory, "missing.jpg"): {}}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(removePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("removed image still exists: %v", err)
	}
	if _, err := os.Stat(keepPath); err != nil {
		t.Fatalf("unrelated image was removed: %v", err)
	}
}

package blogposts

import (
	"errors"
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

type StubFailingFs struct {

}

func (s StubFailingFs)Open(name string) (fs.File, error) {
	return nil, errors.New("nope, no file for you")
}

func TestNewBlogposts(t *testing.T) {
	const (
		firstBody = `Title: Post 1
Description: Description 1
Tags: tdd, go
---
Hello
World`
		secondBody = `Title: Post 2
Description: Description 2
Tags: borrowchecker, rust
---
A
L
M`)
	fs := fstest.MapFS{
		"hello world.md": {Data: []byte(firstBody)},
		"hello-word2.md": {Data: []byte(secondBody)},
	}

	posts, err := NewPostsFromFs(fs)

	if err != nil {
		t.Fatal(err)
	}

	if len(posts) != len(fs) {
		t.Errorf("got %d, wanted %d posts", len(posts), len(fs))
	}
	got := posts[0]
	want := Post{
		Title: "Post 1", 
		Description: "Description 1", 
		Tags: []string{"tdd", "go"},
		Body: `Hello
World`}

	assertPost(t, got, want)
}

func assertPost(t *testing.T, got Post, want Post) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
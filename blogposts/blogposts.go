package blogposts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
)

const (
	titleSeperator = "Title: "
	descriptionSeperator = "Description: "
	tagSeperator = "Tags: "
)

type Post struct {
	Title string
	Description string
	Tags []string
	Body string
}

func NewPostsFromFs(filesystem fs.FS) ([]Post, error) {
	dir, err := fs.ReadDir(filesystem, ".")
	if err != nil {
		return nil, err
	}
	var posts []Post
	for _, f := range dir {
		post, err := getPost(filesystem, f.Name())
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func getPost(filesystem fs.FS, file string) (Post, error) {
	postFile, err := filesystem.Open(file)
	if err != nil {
		return Post{}, err
	}
	defer postFile.Close()
	return newPost(postFile)
}

func newPost(postFile io.Reader) (Post, error) {
	scanner := bufio.NewScanner(postFile)
	readLine := func(tagName string) string {
		scanner.Scan()
		return strings.TrimPrefix(scanner.Text(), tagName)
	}
	titleLine := readLine(titleSeperator)
	descriptionLine := readLine(descriptionSeperator)
	tags := strings.Split(readLine(tagSeperator), ", ")
	body := readBody(scanner)
	post := Post{Title: titleLine, Description: descriptionLine, Tags: tags, Body: body}
	return post, nil
}

func readBody(scanner *bufio.Scanner) string {
	scanner.Scan() //ignore line ------
	buf := bytes.Buffer{}
	for scanner.Scan() {
		fmt.Fprintln(&buf, scanner.Text())
	}
	body := strings.TrimSuffix(buf.String(), "\n")
	return body
}
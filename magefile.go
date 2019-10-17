//+build mage

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

const (
	postsDir    = "./docs/post"
	previewsDir = "./static"
)

func Publish() error {
	err := Preview()
	if err != nil {
		return err
	}

	err = Content()
	if err != nil {
		return err
	}

	return nil
}

func Content() error {
	verbose("generating content...")

	err := sh.Run("hugo")
	if err != nil {
		return err
	}

	return nil
}

func Preview() error {
	mg.Deps(installHugoPostPreview)

	posts, err := changedPosts()
	if err != nil {
		return err
	}

	for _, p := range posts {
		err = hugoPostPreview(p)
		if err != nil {
			return err
		}
	}

	return nil
}

func installHugoPostPreview() {
	verbose("installing hugo-post-preview...")

	sh.Run("go", "get", "github.com/mtpereira/hugo-post-preview")
	sh.Run("go", "install", "github.com/mtpereira/hugo-post-preview/cmd/hugo-post-preview")
}

func hugoPostPreview(post string) error {
	verbose(fmt.Sprintf("creating preview for post %s", post))

	err := sh.Run("hugo-post-preview", "-filename", previewPath(post), "-post", post, "-timeout", "3s")
	if err != nil {
		return err
	}

	return nil
}

func changedPosts() ([]string, error) {
	var posts []string

	files, err := ioutil.ReadDir(postsDir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}

		file := f.Name()
		changed, err := target.Path(previewPath(file), postPath(file))
		if err != nil {
			return nil, err
		}

		if changed {
			posts = append(posts, file)
		}
	}

	return posts, nil
}

func previewPath(post string) string { return fmt.Sprintf("%s/p_%s.png", previewsDir, post) }

func postPath(post string) string { return fmt.Sprintf("%s/%s/index.html", postsDir, post) }

func verbose(msg string) {
	if mg.Verbose() {
		fmt.Println(msg)
	}
}

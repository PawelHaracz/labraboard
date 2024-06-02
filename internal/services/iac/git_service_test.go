package iac

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"labraboard/internal/models"
	"os"
	"testing"
)

const gitRepoUrl = "https://github.com/Labraboard/testing-repo.git"

func TestGit_Clone_Tag(t *testing.T) {
	t.Run("v0.0.1", func(t *testing.T) {
		var folderPath = t.TempDir()

		var git, err = GitClone(gitRepoUrl, folderPath, "v0.0.1", models.TAG)
		if err != nil {
			t.Error(err)
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "0877d712ce274753f7253e56cb8bcf136b7538ee")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), true)
	})

	t.Run("v0.0.2", func(t *testing.T) {
		var folderPath = t.TempDir()

		var git, err = GitClone(gitRepoUrl, folderPath, "v0.0.2", models.TAG)
		if err != nil {
			t.Error(err)
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "3587b396e4abb7724b63ba03f3dac11e7a2f2aa8")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), true)
	})
}

func TestGit_Clone_Sha(t *testing.T) {
	t.Run("0877d712ce274753f7253e56cb8bcf136b7538ee", func(t *testing.T) {
		var folderPath = t.TempDir()

		var git, err = GitClone(gitRepoUrl, folderPath, "0877d712ce274753f7253e56cb8bcf136b7538ee", models.SHA)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "0877d712ce274753f7253e56cb8bcf136b7538ee")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), true)
	})

	t.Run("2b7ef90f4cc37f6d7c18e67c8673ac4aff71a47d", func(t *testing.T) {
		var folderPath = t.TempDir()

		var git, err = GitClone(gitRepoUrl, folderPath, "2b7ef90f4cc37f6d7c18e67c8673ac4aff71a47d", models.SHA)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "2b7ef90f4cc37f6d7c18e67c8673ac4aff71a47d")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), false)
	})
}

func TestGit_Clone_Branch(t *testing.T) {

	t.Run("lorem", func(t *testing.T) {
		var folderPath = t.TempDir()
		var git, err = GitClone(gitRepoUrl, folderPath, "lorem", models.BRANCH)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "b1ee720bf09c9c6e4415fbaf9fd39e68b45ddc9b")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), true)
	})
	t.Run("test", func(t *testing.T) {
		var folderPath = t.TempDir()
		var git, err = GitClone(gitRepoUrl, folderPath, "test", models.BRANCH)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "69cb49343839a7467c70915791f12438bb5a3c93")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), true)
	})
	t.Run("lorem2", func(t *testing.T) {
		var folderPath = t.TempDir()
		var git, err = GitClone(gitRepoUrl, folderPath, "lorem2", models.BRANCH)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		t.Cleanup(func() {
			err = git.Clear()
			if err != nil {
				t.Fatal()
			}
		})

		assert.Equal(t, git.gitSha, "3587b396e4abb7724b63ba03f3dac11e7a2f2aa8")
		assert.Equal(t, exists(folderPath), true)
		assert.Equal(t, exists(fmt.Sprintf("%s/lorem.md", folderPath)), true)
	})
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

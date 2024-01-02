package internal

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

type repoSlice struct {
	r     *git.Repository
	since *time.Time
	until *time.Time
}

// commits return all commit history
// TODO: cache
func (rs *repoSlice) commits() ([]*object.Commit, error) {
	refs, err := rs.r.Branches()
	if err != nil {
		return nil, err
	}

	var result []*object.Commit
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		commits, err := rs.r.Log(&git.LogOptions{From: ref.Hash(), Since: rs.since, Until: rs.until})
		if err != nil {
			return err
		}
		err = commits.ForEach(func(commit *object.Commit) error {
			result = append(result, commit)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (rs *repoSlice) developer2Commits() (map[string][]*object.Commit, error) {
	commits, err := rs.commits()
	if err != nil {
		return nil, err
	}

	d2c := make(map[string][]*object.Commit)
	for _, commit := range commits {
		d2c[commit.Author.String()] = append(d2c[commit.Author.String()], commit)
	}

	return d2c, nil
}

func (rs *repoSlice) MostActiveDeveloper() (string, int, error) {
	d2c, err := rs.developer2Commits()
	if err != nil {
		return "", 0, err
	}

	var mostActiveDeveloper string
	var maxCommitCount int

	for developer, commits := range d2c {
		if len(commits) > maxCommitCount {
			maxCommitCount = len(commits)
			mostActiveDeveloper = developer
		}
	}

	return mostActiveDeveloper, maxCommitCount, nil
}

func (rs *repoSlice) MostHardworkingDeveloper() (string, int, error) {
	d2c, err := rs.developer2Commits()
	if err != nil {
		return "", 0, err
	}

	var mostHardworkingDeveloper string
	var maxHardworkingCommitCount int

	for developer, commits := range d2c {
		var hardworkingCommits []*object.Commit
		for _, c := range commits {
			hour := c.Committer.When.Hour()
			if hour >= 20 || hour < 6 {
				hardworkingCommits = append(hardworkingCommits, c)
			}
		}
		if len(hardworkingCommits) > maxHardworkingCommitCount {
			maxHardworkingCommitCount = len(hardworkingCommits)
			mostHardworkingDeveloper = developer
		}

	}

	return mostHardworkingDeveloper, maxHardworkingCommitCount, nil
}

func (rs *repoSlice) FirstCommit() (*object.Commit, error) {
	commits, err := rs.commits()
	if err != nil {
		return nil, err
	}

	var firstCommit *object.Commit
	for _, commit := range commits {
		if firstCommit == nil || commit.Committer.When.Before(firstCommit.Committer.When) {
			firstCommit = commit
		}
	}

	return firstCommit, nil
}

func (rs *repoSlice) LastCommit() (*object.Commit, error) {
	commits, err := rs.commits()
	if err != nil {
		return nil, err
	}

	var lastCommit *object.Commit
	for _, commit := range commits {
		if lastCommit == nil || commit.Committer.When.After(lastCommit.Committer.When) {
			lastCommit = commit
		}
	}

	return lastCommit, nil
}

func (rs *repoSlice) RepoName() (string, error) {
	remote, err := rs.r.Remote("origin")
	if err != nil {
		return "", err
	}
	repoName := remote.Config().URLs[0]

	return repoName, nil
}

// newRepoSlice returns a repoSlice
func newRepoSlice(path string, since, until *time.Time) (*repoSlice, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	return &repoSlice{r: r, since: since, until: until}, nil
}

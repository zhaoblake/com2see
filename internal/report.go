package internal

import (
	"github.com/go-git/go-git/v5/plumbing/object"
	"html/template"
	"os"
	"path/filepath"
)

type reportData struct {
	RepoName                       string
	Year                           int
	CommitCount                    int
	DeveloperCount                 int
	FirstCommit                    *object.Commit
	LastCommit                     *object.Commit
	MostActiveDeveloper            string
	MostActiveDeveloperCommitCount int
	MostHardworkingDeveloper       string
	MaxHardworkingCommitCount      int
}

// GenerateHTMLReport generate annual commit report
func GenerateHTMLReport(gitPath string, year int) error {
	tmplPath := filepath.Join("templates", "simple.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	outputFile, err := os.Create("commit-report.html")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	firstDay, lastDay := getYearBounds(year)
	rs, err := newRepoSlice(gitPath, &firstDay, &lastDay)
	if err != nil {
		return err
	}

	data, err := NewReportData(year, rs)
	if err != nil {
		return err
	}

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return err
	}

	return nil
}

func NewReportData(year int, rs *repoSlice) (*reportData, error) {
	var data reportData
	repoName, err := rs.RepoName()
	if err != nil {
		return nil, err
	}
	data.RepoName = repoName

	data.Year = year
	commits, err := rs.commits()
	if err != nil {
		return nil, err
	}
	data.CommitCount = len(commits)

	d2c, err := rs.developer2Commits()
	if err != nil {
		return nil, err
	}
	data.DeveloperCount = len(d2c)

	firstCommit, err := rs.FirstCommit()
	if err != nil {
		return nil, err
	}
	data.FirstCommit = firstCommit

	lastCommit, err := rs.LastCommit()
	if err != nil {
		return nil, err
	}
	data.LastCommit = lastCommit

	mostActiveDeveloper, commitCount, err := rs.MostActiveDeveloper()
	if err != nil {
		return nil, err
	}
	data.MostActiveDeveloper = mostActiveDeveloper
	data.MostActiveDeveloperCommitCount = commitCount

	mostHardworkingDeveloper, maxHardworkingCommitCount, err := rs.MostHardworkingDeveloper()
	if err != nil {
		return nil, err
	}
	data.MostHardworkingDeveloper = mostHardworkingDeveloper
	data.MaxHardworkingCommitCount = maxHardworkingCommitCount

	return &data, nil
}

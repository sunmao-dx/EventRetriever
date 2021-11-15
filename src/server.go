package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	gitee_utils "gitee.com/lizi/test-bot/src/gitee-utils"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

var repo []byte

type RepoInfo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
}

func getToken() []byte {
	return []byte(os.Getenv("gitee_token"))
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Event received.")
	eventType, _, payload, ok, _ := gitee_utils.ValidateWebhook(w, r)
	if !ok {
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "gitee hook is broken",
		}).Info("info log")
		return
	}

	switch eventType {
	case "Issue Hook":
		var ie gitee.IssueEvent
		if err := json.Unmarshal(payload, &ie); err != nil {
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context": "gitee hook is broken",
			}).Info("info log")
			return
		}
		go handleIssueEvent(&ie)
	case "Note Hook":
		var ic gitee.NoteEvent
		if err := json.Unmarshal(payload, &ic); err != nil {
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context": "gitee hook is broken",
			}).Info("info log")
			return
		}
		go handleCommentEvent(&ic)
	default:
		return
	}
}

func handleIssueEvent(i *gitee.IssueEvent) error {
	var issue gitee_utils.Issue
	var repoinfo RepoInfo
	err := json.Unmarshal(repo, &repoinfo)
	if err != nil {
		log.Println("wrong repo", err)
		return err
	}
	issue = _init(issue)
	issue.IssueID = i.Issue.Number
	issue.IssueAction = *(i.Action)
	issue.IssueUser.IssueUserID = i.Issue.User.Login
	issue.IssueUser.IssueUserName = i.Issue.User.Name
	issue.IssueUser.IsOrgUser = 0 //default is 0

	issue.IssueTime = i.Issue.CreatedAt.Format(time.RFC3339)
	issue.IssueUpdateTime = i.Issue.UpdatedAt.Format(time.RFC3339)
	issue.IssueTitle = i.Issue.Title
	issue.IssueContent = i.Issue.Body
	if i.Issue.Assignee == nil {
		issue.IssueAssignee = ""
	} else {
		issue.IssueAssignee = i.Issue.Assignee.Login
	}
	issue.IssueLabel = getLabels(i.Issue.Labels)

	fmt.Println(issue)

	strApi := os.Getenv("api_url")

	c := gitee_utils.NewClient(getToken)

	issue.IssueUser.IsEntUser = isUserInEnt(issue.IssueUser.IssueUserID, repoinfo.Repo, c)

	_, errIssue := c.SendIssue(issue, strApi)
	if err != nil {
		fmt.Println(err.Error())
		return errIssue
	}
	return nil
}

func handleCommentEvent(i *gitee.NoteEvent) {
	switch *(i.NoteableType) {
	case "Issue":
		go handleIssueCommentEvent(i)
	default:
		return
	}
}

func handleIssueCommentEvent(i *gitee.NoteEvent) {
	return
}

func getLabels(initLabels []gitee.LabelHook) []gitee_utils.Label {
	var issueLabel gitee_utils.Label
	var issueLabels []gitee_utils.Label
	for _, label := range initLabels {
		issueLabel.Name = label.Name
		issueLabel.Desciption = label.Name
		issueLabels = append(issueLabels, issueLabel)
	}
	return issueLabels
}

func isUserInEnt(login, entOrigin string, c gitee_utils.Client) int {
	_, err := c.GetUserEnt(entOrigin, login)
	if err != nil {
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "Is not Enterprise member",
		}).Info("info log")
		fmt.Println(err)
		return 0
	} else {
		return 1
	}
}

func loadFile(path, fileType string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		defer jsonFile.Close()
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	switch {
	case fileType == "repo":
		repo = byteValue
	default:
		fmt.Printf("no filetype\n")
	}
	return nil
}

func _init(i gitee_utils.Issue) gitee_utils.Issue {
	i.IssueID = "XXXXXX"
	i.IssueAction = "Open"
	i.IssueUser.IssueUserID = "no_name"
	i.IssueUser.IssueUserName = "NO_NAME"
	i.IssueUser.IsOrgUser = 0 //default is 0
	i.IssueUser.IsEntUser = 1
	i.IssueAssignee = "no_assignee"
	i.IssueLabel = nil

	i.IssueTime = time.Now().Format(time.RFC3339)
	i.IssueUpdateTime = time.Now().Format(time.RFC3339)
	i.IssueTitle = "no_title"
	i.IssueContent = "no_content"
	return i
}

func configFile() {
	loadFile("src/data/repo.json", "repo")
}

func main() {
	configFile()
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":8001", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	gitee_utils "gitee.com/lizi/test-bot/src/gitee-utils"
	"gitee.com/openeuler/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

var repo []byte

type RepoInfo struct {
	Org  string `json:"org"`
	Repo string `json:"repo"`
	Ent  string `json:"ent"`
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
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "gitee hook success",
		}).Info("info log")
		go handleIssueEvent(&ie)
	case "Note Hook":
		var ic gitee.NoteEvent
		if err := json.Unmarshal(payload, &ic); err != nil {
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context": "gitee hook is broken",
			}).Info("info log")
			return
		}
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "gitee hook success",
		}).Info("info log")
		go handleCommentEvent(&ic)
	default:
		return
	}
}

func handleIssueEvent(i *gitee.IssueEvent) error {
	if *(i.Action) != "open" {
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "the hook is not for opening a issue",
		}).Info("info log")
		return nil
	}
	var issue gitee_utils.Issue
	var repoinfo RepoInfo
	var strEnt string
	repoinfo.Org = i.Repository.Namespace
	repoinfo.Repo = i.Repository.Name
	if i.Enterprise == nil {
		strEnt = ""
	} else {
		strEnt = i.Enterprise.Url
		strEnt = strEnt[strings.LastIndex(strEnt, "/")+1:]
	}
	repoinfo.Ent = strEnt
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

	if i.Issue.Number == "I1EL99" {
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "the hook is a test msg",
		}).Info("info log")
		return nil
	}

	if i.Issue.Assignee == nil {
		issue.IssueAssignee = ""
	} else {
		issue.IssueAssignee = i.Issue.Assignee.Login
	}
	issue.IssueLabel = getLabels(i.Issue.Labels)

	fmt.Println(issue)

	strApi := os.Getenv("api_url")

	c := gitee_utils.NewClient(getToken)

	if repoinfo.Ent == "" {
		issue.IssueUser.IsEntUser = 0
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context":   issue.IssueUser.IssueUserID + " is not an Enterprise member _ Name is null",
			"issueID":   issue.IssueID,
			"EntName":   repoinfo.Ent,
			"isEntUser": issue.IssueUser.IsEntUser,
			"issue":     issue,
		}).Info("info log")
	} else {
		issue.IssueUser.IsEntUser = isUserInEnt(issue.IssueUser.IssueUserID, repoinfo.Ent, c)
		if issue.IssueUser.IsEntUser == 0 {
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context":   issue.IssueUser.IssueUserID + " is not an Enterprise member",
				"issueID":   issue.IssueID,
				"EntName":   repoinfo.Ent,
				"isEntUser": issue.IssueUser.IsEntUser,
				"issue":     issue,
			}).Info("info log")
		} else {
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context":   issue.IssueUser.IssueUserID + " is an Enterprise member",
				"issueID":   issue.IssueID,
				"EntName":   repoinfo.Ent,
				"isEntUser": issue.IssueUser.IsEntUser,
				"issue":     issue,
			}).Info("info log")
		}
	}

	_, errIssue := c.SendIssue(issue, strApi)
	if errIssue != nil {
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": "Send issue problem",
		}).Info("info log")
		fmt.Println(errIssue.Error())
		return errIssue
	}
	gitee_utils.LogInstance.WithFields(logrus.Fields{
		"context": "Send issue success",
	}).Info("info log")
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
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		fmt.Println(err.Error() + login + " is not an Ent memeber")
		gitee_utils.LogInstance.WithFields(logrus.Fields{
			"context": err.Error() + " " + login + " is not an Ent memeber",
		}).Info("info log")
		return 0
	} else {
		if err == nil {
			fmt.Println(" is an Ent memeber")
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context": login + " is an Ent memeber",
			}).Info("info log")
			return 1
		} else {
			fmt.Println(err.Error() + "  now, retry...")
			gitee_utils.LogInstance.WithFields(logrus.Fields{
				"context": "  now, retry...",
			}).Info("info log")
			time.Sleep(time.Duration(5) * time.Second)
			return isUserInEnt(login, entOrigin, c)
		}
	}
}

func _init(i gitee_utils.Issue) gitee_utils.Issue {
	i.IssueID = "XXXXXX"
	i.IssueAction = "Open"
	i.IssueUser.IssueUserID = "no_name"
	i.IssueUser.IssueUserName = "NO_NAME"
	i.IssueUser.IsOrgUser = 0
	i.IssueUser.IsEntUser = 1
	i.IssueAssignee = "no_assignee"
	i.IssueLabel = nil

	i.IssueTime = time.Now().Format(time.RFC3339)
	i.IssueUpdateTime = time.Now().Format(time.RFC3339)
	i.IssueTitle = "no_title"
	i.IssueContent = "no_content"
	return i
}

func main() {
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":8001", nil)
}

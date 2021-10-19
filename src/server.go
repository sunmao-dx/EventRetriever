package main

import (
	"encoding/json"
	"fmt"
	gitee_utils "gitee.com/lizi/test-bot/src/gitee-utils"
	"gitee.com/openeuler/go-gitee/gitee"
	"io/ioutil"
	"net/http"
	"os"
)

var token []byte
var apiUrl []byte

func getToken() []byte {
	return token
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Event.")
	eventType, _, payload, ok, _ := gitee_utils.ValidateWebhook(w, r)
	if !ok {
		return
	}

	switch eventType {
	case "Issue Hook":
		var ie gitee.IssueEvent
		if err := json.Unmarshal(payload, &ie); err != nil {
			return
		}
		go handleIssueEvent(&ie)
	case "Note Hook":
		var ic gitee.NoteEvent
		if err := json.Unmarshal(payload, &ic); err != nil {
			return
		}
		go handleCommentEvent(&ic)
	default:
		return
	}
}

func handleIssueEvent(i *gitee.IssueEvent) {
	var issue gitee_utils.Issue
	issue.IssueID = i.Issue.Number
	issue.IssueUser = i.Issue.User.Name
	issue.IssueUserID = i.Issue.User.Login
	issue.IssueTime = i.Issue.CreatedAt.String()
	issue.IssueUpdateTime = i.Issue.UpdatedAt.String()
    issue.IssueAssignee = i.Issue.Assignee.Login
	issue.IssueLabel = getLabels(i.Issue.Labels)

	strApi := string(apiUrl[:])

	c := gitee_utils.NewClient(getToken)
	_, err := c.SendIssue(issue, strApi)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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


func getLabels(initLabels []gitee.LabelHook) []gitee_utils.Label{
	var issueLabel gitee_utils.Label
	var issueLabels []gitee_utils.Label
	for _, label := range initLabels {
		issueLabel.Name = label.Name
		issueLabel.Desciption = label.Name
		issueLabels = append(issueLabels, issueLabel)
	}
	return issueLabels
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
	case fileType == "token" :
		token = byteValue
	case fileType == "APIUrl" :
		apiUrl = byteValue
	default:
		fmt.Printf("no filetype\n" )
	}
	return nil
}

func configFile() {
	loadFile("src/data/token.md", "token")
	loadFile("src/data/ApiUrl.md", "APIUrl")
}

func main() {
	configFile()
	http.HandleFunc("/", ServeHTTP)
	http.ListenAndServe(":8001", nil)
}

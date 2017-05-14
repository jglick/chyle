package decorators

import (
	"fmt"
	"net/http"
)

type jiraIssueConfig struct {
	CREDENTIALS struct {
		URL      string
		USERNAME string
		PASSWORD string
	}
	KEYS map[string]struct {
		DESTKEY string
		FIELD   string
	}
}

// jiraIssue fetch data using jira issue api
type jiraIssue struct {
	client http.Client
	config jiraIssueConfig
}

// Decorate fetch remote jira service if a jiraIssueId is defined to fetch issue datas
func (j jiraIssue) Decorate(commitMap *map[string]interface{}) (*map[string]interface{}, error) {
	var URLpattern string

	switch (*commitMap)["jiraIssueId"].(type) {
	case string:
		URLpattern = "%s/rest/api/2/issue/%s"
	case int64:
		URLpattern = "%s/rest/api/2/issue/%d"
	default:
		return commitMap, nil
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(URLpattern, j.config.CREDENTIALS.URL, (*commitMap)["jiraIssueId"]), nil)

	if err != nil {
		return commitMap, err
	}

	req.SetBasicAuth(j.config.CREDENTIALS.USERNAME, j.config.CREDENTIALS.PASSWORD)
	req.Header.Set("Content-Type", "application/json")

	return jSONResponse{&j.client, req, j.config.KEYS}.Decorate(commitMap)
}

// buildJiraIssue create a new jira ticket decorator
func buildJiraIssue(config jiraIssueConfig) Decorater {
	return jiraIssue{http.Client{}, config}
}

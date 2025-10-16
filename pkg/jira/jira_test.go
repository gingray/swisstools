package jira

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetJiraUrl(t *testing.T) {
	assertion := assert.New(t)
	key := "PRO-123"
	result := getJiraUrl("https://jira.example.com/", key)
	assertion.Equal("https://jira.example.com/browse/PRO-123", result)
}

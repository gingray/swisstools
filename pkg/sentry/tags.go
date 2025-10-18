package sentry

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/charmbracelet/log"
)

type Tag struct {
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	Value     string    `json:"value"`
	Count     int       `json:"count"`
	LastSeen  time.Time `json:"lastSeen"`
	FirstSeen time.Time `json:"firstSeen"`
}

func (c *client) GetTagValues(organization, project, tag string) ([]Tag, error) {
	apiSegment := fmt.Sprintf("projects/%s/%s/tags/%s/values/", organization, project, tag)
	resp, err := c.makeRequest("GET", apiSegment, nil)
	if err != nil {
		log.Error(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}
	allTagValues := make([]Tag, 0)
	tags := make([]Tag, 0)
	err = json.Unmarshal(body, &tags)
	if err != nil {
		log.Error(err)
	}

	allTagValues = append(allTagValues, tags...)
	pagination := parsePagination(resp.Header.Get("Link"))
	for pagination != nil && pagination.NextUrl != "" {
		resp, err := c.do("GET", pagination.NextUrl, nil)
		if err != nil {
			log.Error(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
		}
		err = json.Unmarshal(body, &tags)
		if err != nil {
			log.Error(err)
		}
		allTagValues = append(allTagValues, tags...)
		pagination = parsePagination(resp.Header.Get("Link"))
	}
	return allTagValues, err
}

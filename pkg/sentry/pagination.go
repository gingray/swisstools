package sentry

import (
	"strconv"
	"strings"
)

type pagination struct {
	NextUrl string
	PrevUrl string
}

func parsePagination(header string) *pagination {
	if header == "" {
		return nil
	}
	pageResult := &pagination{}
	links := strings.SplitN(header, ",", 2)
	for _, page := range links {
		data := strings.SplitN(page, ";", 4)

		pagelink := strings.TrimLeft(strings.TrimSpace(data[0]), "<")
		pagelink = strings.TrimRight(pagelink, ">")

		pagetype := strings.Trim(strings.Split(data[1], "=")[1], `"`)
		results, err := strconv.ParseBool(strings.Trim(strings.Split(strings.TrimSpace(data[2]), "=")[1], `"`))
		if err != nil {
			results = false
		}

		if pagetype == "previous" && results {
			pageResult.PrevUrl = pagelink
		}
		if pagetype == "next" && results {
			pageResult.NextUrl = pagelink
		}
	}

	return pageResult
}

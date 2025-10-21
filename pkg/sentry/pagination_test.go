package sentry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPageHeaderParsing(t *testing.T) {
	assertion := assert.New(t)
	header := `<https://sentry.io/api/0/projects/the-interstellar-jurisdiction/pump-station/releases/?&cursor=100:-1:1>; rel="previous"; results="true"; cursor="100:-1:1", <https://sentry.io/api/0/projects/the-interstellar-jurisdiction/pump-station/releases/?&cursor=100:1:0>; rel="next"; results="true"; cursor="100:1:0`

	link := parsePagination(header)
	assertion.Equal("https://sentry.io/api/0/projects/the-interstellar-jurisdiction/pump-station/releases/?&cursor=100:1:0", link.NextUrl)
	assertion.Equal("https://sentry.io/api/0/projects/the-interstellar-jurisdiction/pump-station/releases/?&cursor=100:-1:1", link.PrevUrl)
}

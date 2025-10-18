package sentry

import "fmt"

func generateLinkToDashboard(baseUrl string, organization, project, tagName, tagValue string) string {
	return fmt.Sprintf("%s/organizations/%s/issues/?project=%s&query=%s:%s&statsPeriod=24h", baseUrl, organization, project, tagName, tagValue)
}

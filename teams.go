package asana

import (
	"fmt"
	"strconv"

	sentry "github.com/getsentry/sentry-go"
)

// Team stores Team from Asana
//
type Team struct {
	ID              string        `json:"gid"`
	Name            string        `json:"name"`
	ResourceType    string        `json:"resource_type"`
	Description     string        `json:"description"`
	HTMLDescription string        `json:"html_description"`
	Organization    CompactObject `json:"organization"`
}

// GetTeamsByWorkspaceID returns all teams for a specific team
//
func (i *Asana) GetTeamsByWorkspaceID(workspaceID string) ([]Team, error) {
	return i.GetTeamsInternal(workspaceID)
}

// GetTeamsInternal is the generic function retrieving teams from Asana
//
func (i *Asana) GetTeamsInternal(workspaceID string) ([]Team, error) {
	urlStr := "%sorganizations/%s/teams?limit=%s%s&opt_fields=%s"
	limit := 100
	offset := ""
	//rowCount := limit
	batch := 0

	teams := []Team{}

	for batch == 0 || offset != "" {
		batch++
		//fmt.Printf("Batch %v for WorkspaceID %v\n", batch, workspaceID)

		url := fmt.Sprintf(urlStr, i.ApiURL, workspaceID, strconv.Itoa(limit), offset, GetJSONTaggedFieldNames(Team{}))
		//fmt.Println(url)

		ts := []Team{}

		nextPage, response, err := i.Get(url, &ts)
		if err != nil {
			return nil, err
		}

		if response != nil {
			if response.Errors != nil {
				for _, e := range *response.Errors {
					message := fmt.Sprintf("Error for WorkspaceID %v: %v", workspaceID, e.Message)
					if i.IsLive {
						sentry.CaptureMessage(message)
					}
					fmt.Println(message)
				}
			}
		}

		for _, t := range ts {
			teams = append(teams, t)
		}

		//rowCount = len(ts)
		offset = ""
		if nextPage != nil {
			offset = fmt.Sprintf("&offset=%s", nextPage.Offset)
		}
	}

	if len(teams) == 0 {
		teams = nil
	}

	return teams, nil
}

package bangumi

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jiang-zhexin/animedb/internal/app"
)

const DefaultBaseURL = "https://api.bgm.tv"

var DefaultUserAgent string

func init() {
	DefaultUserAgent = fmt.Sprintf("%s/%s (https://github.com/jiang-zhexin/animedb)", app.Appname, app.Version)
}

func ClientEditor(c *Client) error {
	c.RequestEditors = []RequestEditorFn{UserAgent}
	return nil
}

func UserAgent(ctx context.Context, req *http.Request) error {
	req.Header.Set("User-Agent", DefaultUserAgent)
	return nil
}

func SearchSubject(keyword string) (result []Subject, err error) {
	c, err := NewClientWithResponses(DefaultBaseURL, ClientEditor)
	if err != nil {
		return
	}

	resp, err := c.SearchSubjectsWithResponse(context.Background(), &SearchSubjectsParams{},
		SearchSubjectsJSONRequestBody{
			Keyword: keyword,
			Filter: &struct {
				AirDate     *[]string      "json:\"air_date,omitempty\""
				MetaTags    *[]string      "json:\"meta_tags,omitempty\""
				Nsfw        *bool          "json:\"nsfw,omitempty\""
				Rank        *[]string      "json:\"rank,omitempty\""
				Rating      *[]string      "json:\"rating,omitempty\""
				RatingCount *[]string      "json:\"rating_count,omitempty\""
				Tag         *[]string      "json:\"tag,omitempty\""
				Type        *[]SubjectType "json:\"type,omitempty\""
			}{
				Type: &[]SubjectType{
					Anime,
				},
			},
		},
	)
	if err != nil {
		return
	}
	if resp.JSON200 == nil || resp.JSON200.Data == nil {
		err = fmt.Errorf("the HTTP code is %d", resp.StatusCode())
		return
	}
	result = (*resp.JSON200.Data)
	return
}

func GetSubject(subjectID SubjectID) (result *Subject, err error) {
	c, err := NewClientWithResponses(DefaultBaseURL, ClientEditor)
	if err != nil {
		return
	}

	resp, err := c.GetSubjectByIdWithResponse(context.Background(), PathSubjectId(subjectID))
	if err != nil {
		return
	}
	if resp.JSON200 == nil {
		err = fmt.Errorf("the HTTP code is %d", resp.StatusCode())
		return
	}
	result = resp.JSON200
	return
}

package web

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func InitSuperTokens() error {
	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "https://st-dev-5f1e08c0-252b-11ef-ad47-516b0aeb722e.aws.supertokens.io",
			APIKey:        "2iFbtKJlI9eFn3HHu5LlT1dDzI",
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "Data Monster",
			APIDomain:       "http://dev.local:8080",
			WebsiteDomain:   "http://dev.local:8090",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
		},
	})
	return err
}

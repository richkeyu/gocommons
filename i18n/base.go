package i18n

import (
	"context"
	"strings"

	"github.com/richkeyu/gocommons/client"
)

const (
	baseServiceI18nPath = "/v1/i18n/category/tran"
)

type BaseI18nResponse struct {
	KeyList      map[string]string `json:"key_list"`
	Category     int               `json:"category"`
	LanguageCode string            `json:"language_code"`
}

func getTransFromBase(ctx context.Context, conf Config) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	for _, lang := range strings.Split(conf.Support, ",") {
		resp, err := client.NewClient().WithContext(ctx).Get(conf.BaseUri+baseServiceI18nPath, client.Options{
			Query: map[string]interface{}{
				"category":      conf.Category,
				"language_code": lang,
			},
		})
		if err != nil {
			return nil, err
		}
		var i18nResp BaseI18nResponse
		_, err = resp.MustParseBody(&i18nResp)
		if err != nil {
			return nil, err
		}
		result[lang] = i18nResp.KeyList
	}
	return result, nil
}

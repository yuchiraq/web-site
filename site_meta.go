package main

import (
	"os"
	"strings"
)

type SiteMeta struct {
	GoogleSiteVerification     string
	YandexVerification         string
	FacebookDomainVerification string
	MetaPixelID                string
}

func currentSiteMeta() SiteMeta {
	return SiteMeta{
		GoogleSiteVerification:     strings.TrimSpace(os.Getenv("GOOGLE_SITE_VERIFICATION")),
		YandexVerification:         strings.TrimSpace(os.Getenv("YANDEX_VERIFICATION")),
		FacebookDomainVerification: strings.TrimSpace(os.Getenv("FACEBOOK_DOMAIN_VERIFICATION")),
		MetaPixelID:                strings.TrimSpace(os.Getenv("META_PIXEL_ID")),
	}
}

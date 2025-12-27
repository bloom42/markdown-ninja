package service

import (
	"encoding/json"
	"fmt"
	"html/template"

	"markdown.ninja/pkg/services/contacts"
	"markdown.ninja/pkg/services/content"
	"markdown.ninja/pkg/services/site"
	"markdown.ninja/pkg/services/websites"
)

type pageTemplateData struct {
	Url         template.URL
	Title       string
	Description string
	Language    string
	SocialImage string

	Website           site.Website
	Page              *site.Page
	MarkdownNinjaData template.JS
	Header            template.HTML
	Footer            template.HTML

	Country string
	// ClientMetadata template.HTML
}

// // clientMetadata are rendered on the page for debugging purpose
// type clientMetadata struct {
// 	IP        string    `json:"ip"`
// 	ASN       int64     `json:"asn"`
// 	Country   string    `json:"country"`
// 	UserAgent string    `json:"user_agent"`
// 	Date      time.Time `json:"date"`
// }

func (service *SiteService) convertPageTemplateData(website websites.Website,
	page *site.Page, tags []content.Tag, contact *contacts.Contact, country string) (ret pageTemplateData, err error) {
	// clientMetadata := clientMetadata{
	// 	IP:        httpCtx.Client.IPStr,
	// 	ASN:       httpCtx.Client.ASN,
	// 	Country:   httpCtx.Client.CountryCode,
	// 	UserAgent: httpCtx.Client.UserAgent,
	// 	Date:      now,
	// }
	// clientMetadataJson, err := json.Marshal(clientMetadata)
	// if err != nil {
	// 	return
	// }
	// clientMetadataStr := template.HTML(fmt.Sprintf(`<meta name="markdown_ninja:client" content="%s" />`, base64.StdEncoding.EncodeToString(clientMetadataJson)))

	websiteData := service.convertWebsite(website)
	url := websiteData.Url
	title := websiteData.Name
	description := websiteData.Description
	language := websiteData.Language
	var markdowNinjaData site.MarkdowNinjaData

	markdowNinjaData.Country = country
	markdowNinjaData.Website = websiteData

	if page != nil {
		url = page.Url
		title = page.Title
		// if we are serving the home page, we want to use the website's description
		if page.Path != "/" {
			description = page.Description
		}
		language = page.Language

		markdowNinjaData.Page = page
	}

	if contact != nil {
		markdowNinjaData.Contact = contact
	}

	markdowNinjaDataJson, err := json.Marshal(markdowNinjaData)
	if err != nil {
		err = fmt.Errorf("marshaling markdown ninja data: %w", err)
		return
	}

	ret = pageTemplateData{
		Url:         url,
		Title:       title,
		Description: description,
		Language:    language,
		// SocialImage:       "",
		Website:           websiteData,
		Page:              page,
		MarkdownNinjaData: template.JS(markdowNinjaDataJson),
		Header:            template.HTML(website.Header),
		Footer:            template.HTML(website.Footer),
		Country:           country,
	}
	return
}

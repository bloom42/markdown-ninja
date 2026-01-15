package client

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/skerkour/stdx-go/crypto/blake3"
	"github.com/yuin/goldmark"
	mdparser "github.com/yuin/goldmark/parser"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"markdown.ninja/pkg/markdown"
	"markdown.ninja/pkg/services/content"
)

func (client *Client) readAndParseMarkdownFile(ctx context.Context, filesystem fs.FS, path, realPath string, pagesFromApi map[string]content.PageMetadata) (localPage localPage, err error) {
	fileName := filepath.Base(path)
	extension := filepath.Ext(fileName)
	fileNameWithoutExtension := strings.TrimSuffix(fileName, extension)
	pathWithoutExtension := strings.TrimSuffix(path, extension)

	localPage.LocalPath = realPath

	// fileMeatadata, err := os.Stat(path)
	// if err != nil {
	// 	err = fmt.Errorf("publish: opening markdown file (%s): %w", path, err)
	// 	return
	// }

	fileData, err := fs.ReadFile(filesystem, path)
	if err != nil {
		err = fmt.Errorf("publish: reading markdown file (%s): %w", realPath, err)
		return
	}

	fileContent := strings.TrimSpace(string(fileData))

	markdownConverter := goldmark.New(
		goldmark.WithExtensions(
			markdown.FrontmatterExtension,
		),
	)
	var markdownOuuput bytes.Buffer
	parserCtx := mdparser.NewContext()

	err = markdownConverter.Convert([]byte(fileContent), &markdownOuuput, mdparser.WithContext(parserCtx))
	if err != nil {
		err = fmt.Errorf("publish: parsing markdown file (%s): %w", realPath, err)
		return
	}

	frontmatter, err := markdown.GetFrontmatter(parserCtx)
	if err != nil {
		if err == markdown.ErrFrontmatterIsMissing {
			// client.logger.Warn(fmt.Sprintf("missing frontmatter (%s)", path))
			frontmatter = markdown.NewEmtpyFrontmatter()
			err = nil
		} else {
			err = fmt.Errorf("publish: parsing frontmatter (%s): %w", realPath, err)
			return
		}
	}

	localPage.FrontMatterSource = frontmatter.Source

	localPage.BodyMarkdown = trimFrontmatter(fileContent, localPage.FrontMatterSource)
	localPageBodyMarkdownHash := blake3.Sum256([]byte(localPage.BodyMarkdown))
	localPage.BodyHash = localPageBodyMarkdownHash[:]

	pageUrlInterface := frontmatter.Data["url"]
	if pageUrlInterface != nil {
		pageUrlStr, pageUrlInterfaceIsString := pageUrlInterface.(string)
		if !pageUrlInterfaceIsString {
			err = fmt.Errorf("publish: parsing frontmatter: url is not a string (%s)", realPath)
			return
		}

		localPage.Url = strings.TrimSpace(pageUrlStr)
	} else {
		localPage.Url = pathWithoutExtension
	}
	if !strings.HasPrefix(localPage.Url, "/") {
		err = fmt.Errorf("publish: url should start with / (%s)", realPath)
		return
	}

	var existingPageFromApi *content.PageMetadata
	if existingPageFromApiValue, existingPageFromApiOk := pagesFromApi[localPage.Url]; existingPageFromApiOk {
		existingPageFromApi = &existingPageFromApiValue
	}

	pageTitleInterface := frontmatter.Data["title"]
	if pageUrlInterface != nil {
		pageTitleStr, pageTitleInterfaceIsString := pageTitleInterface.(string)
		if !pageTitleInterfaceIsString {
			err = fmt.Errorf("publish: parsing frontmatter: title is not a string (%s)", realPath)
			return
		}

		localPage.Title = strings.TrimSpace(pageTitleStr)
	} else if existingPageFromApi != nil {
		localPage.Title = existingPageFromApi.Title
	} else {
		localPage.Title = cases.Upper(language.AmericanEnglish).String(fileNameWithoutExtension)
	}

	pageDateInterface := frontmatter.Data["date"]
	if pageDateInterface != nil {
		date, pageDateInterfaceIsTime := pageDateInterface.(time.Time)
		if !pageDateInterfaceIsTime {
			err = fmt.Errorf("publish: parsing frontmatter: date is not time.Time (%s)", realPath)
			return
		}

		localPage.Date = date.UTC()
	} else {
		err = fmt.Errorf("publish: date is missing (%s)", realPath)
		return
		// localPage.Date = fileMeatadata.ModTime().UTC().Add(-6 * time.Hour)
		// localPage.Date = time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC)
	}

	pageTypeInterface := frontmatter.Data["type"]
	if pageTypeInterface != nil {
		pageTypeStr, pageTypeInterfaceIsString := pageTypeInterface.(string)
		if !pageTypeInterfaceIsString {
			err = fmt.Errorf("publish: parsing frontmatter: type is not a string (%s)", realPath)
			return
		}

		pageTypeStr = strings.TrimSpace(pageTypeStr)
		if pageTypeStr != string(content.PageTypePage) && pageTypeStr != string(content.PageTypePost) {
			err = fmt.Errorf("publish: page type (%s) is not valid (%s)", pageTypeStr, realPath)
			return
		}

		localPage.Type = content.PageType(pageTypeStr)
	} else if existingPageFromApi != nil {
		localPage.Type = existingPageFromApi.Type
	} else {
		localPage.Type = content.PageTypePage
	}
	// if localPage.Type == sites.ContentTypeUnknown {
	// 	err = fmt.Errorf("publish (%s): unknown content type", path)
	// 	return
	// }

	draftInterface := frontmatter.Data["draft"]
	if draftInterface != nil {
		draft, draftInterfaceIsBool := draftInterface.(bool)
		if !draftInterfaceIsBool {
			err = fmt.Errorf("publish: parsing frontmatter: draft is not a bool (%s)", realPath)
			return
		}
		localPage.Draft = draft
	} else {
		localPage.Draft = false
	}

	tagsInterface := frontmatter.Data["tags"]
	if tagsInterface != nil {
		tags, tagsInterfaceIsSlice := tagsInterface.([]any)
		if !tagsInterfaceIsSlice {
			err = fmt.Errorf("publish: parsing frontmatter: tags is not a []string (%s)", realPath)
			return
		}

		localPage.Tags = make([]string, len(tags))
		for i, tag := range tags {
			tagStr, tagInterfaceIsString := tag.(string)
			if !tagInterfaceIsString {
				err = fmt.Errorf("publish: parsing frontmatter: tags is not a []string (%s)", realPath)
				return
			}

			localPage.Tags[i] = strings.TrimSpace(tagStr)
		}
	} else {
		localPage.Tags = make([]string, 0)
	}

	langInterface := frontmatter.Data["lang"]
	if langInterface != nil {
		pageLangStr, pageLangInterfaceIsString := langInterface.(string)
		if !pageLangInterfaceIsString {
			err = fmt.Errorf("publish: parsing frontmatter: lang is not a string (%s)", realPath)
			return
		}

		localPage.Language = strings.ToLower(strings.TrimSpace(pageLangStr))
	} else if existingPageFromApi != nil {
		localPage.Language = existingPageFromApi.Language
	} else {
		localPage.Language = content.PageDefaultLanguage
	}

	descriptionInterface := frontmatter.Data["description"]
	if descriptionInterface != nil {
		pageDescriptionStr, pageDescriptionInterfaceIsString := descriptionInterface.(string)
		if !pageDescriptionInterfaceIsString {
			err = fmt.Errorf("publish: parsing frontmatter: description is not a string (%s)", realPath)
			return
		}

		localPage.Description = strings.TrimSpace(pageDescriptionStr)
	} else if existingPageFromApi != nil {
		localPage.Description = existingPageFromApi.Description
	} else {
		localPage.Description = ""
	}

	pageUpdatedAtInterface := frontmatter.Data["updated"]
	if pageUpdatedAtInterface != nil {

		updatedAt, pageUpdatedAtInterfaceIsTime := pageDateInterface.(time.Time)
		if !pageUpdatedAtInterfaceIsTime {
			err = fmt.Errorf("publish: parsing frontmatter: date is not time.Time (%s)", realPath)
			return
		}

		updatedAt = updatedAt.UTC()
		localPage.UpdatedAt = &updatedAt
	}

	newsletterInterface := frontmatter.Data["newsletter"]
	if newsletterInterface != nil {
		sendAsNewsletter, newsletterInterfaceIsBool := newsletterInterface.(bool)
		if !newsletterInterfaceIsBool {
			err = fmt.Errorf("publish: parsing frontmatter: newsletter is not a bool (%s)", realPath)
			return
		}
		localPage.SendAsNewsletter = sendAsNewsletter
	} else {
		localPage.SendAsNewsletter = false
	}

	localPage.MetadataHash = content.HashPageMetadata(localPage.Type, localPage.Url, localPage.Date, localPage.SendAsNewsletter, localPage.Language, localPage.Title, localPage.Description, localPage.Tags)

	return
}

func trimFrontmatter(markdown, frontmatter string) string {
	runeSize := 0
	var rn rune

	for len(markdown) > 0 {
		rn, runeSize = utf8.DecodeRune([]byte(markdown))
		if rn == '+' || rn == '-' || unicode.IsSpace(rn) {
			markdown = markdown[runeSize:]
		} else {
			break
		}
	}

	markdown = strings.TrimPrefix(markdown, frontmatter)

	for len(markdown) > 0 {
		rn, runeSize = utf8.DecodeRune([]byte(markdown))
		if rn == '+' || rn == '-' || unicode.IsSpace(rn) {
			markdown = markdown[runeSize:]
		} else {
			break
		}
	}

	return markdown
}

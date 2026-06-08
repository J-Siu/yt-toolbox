package lib

import (
	"errors"
	"strings"

	"github.com/J-Siu/go-helper/v2/errs"
)

// var 'ytinitialdata' from channel page
type YTInitialData struct {
	Contents struct {
		TwoColumnBrowseResultsRenderer struct {
			Tabs []struct {
				TabRenderer struct {
					Content struct {
						SectionListRenderer struct {
							Contents []struct {
								ItemSectionRenderer struct {
									Contents []struct {
										ShelfRenderer struct {
											Content struct {
												ExpandedShelfContentsRenderer struct {
													Items []struct {
														ChannelRenderer struct {
															ChannelID string `json:"channelId"`
															Title     struct {
																SimpleText string `json:"simpleText"`
															} `json:"title"`
														} `json:"channelRenderer"`
													} `json:"items"`
												} `json:"expandedShelfContentsRenderer"`
											} `json:"content"`
										} `json:"shelfRenderer"`
									} `json:"contents"`
								} `json:"itemSectionRenderer"`
							} `json:"contents"`
						} `json:"sectionListRenderer"`
					} `json:"content"`
				} `json:"tabRenderer"`
			} `json:"tabs"`
		} `json:"twoColumnBrowseResultsRenderer"`
	} `json:"contents"`
}

// map (key,value) => (title,id)
type ChTitleID map[string]string

func (t ChTitleID) New(data *YTInitialData) {
	prefix := "ChTitleID.New"
	var title string
	for _, tab := range data.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		for _, content1 := range tab.TabRenderer.Content.SectionListRenderer.Contents {
			for _, content2 := range content1.ItemSectionRenderer.Contents {
				for _, item := range content2.ShelfRenderer.Content.ExpandedShelfContentsRenderer.Items {
					title = strings.TrimSpace(item.ChannelRenderer.Title.SimpleText)
					if title == "" {
						errs.Queue(prefix, errors.New("Empty channel title: "+item.ChannelRenderer.ChannelID))
					}
					t[title] = strings.TrimSpace(item.ChannelRenderer.ChannelID)
				}
			}
		}
	}
}

// cspell:words ytinitialdata

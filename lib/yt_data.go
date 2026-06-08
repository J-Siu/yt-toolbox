package lib

import "strings"

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

// map key:title, value:id
type ChTitle string
type ChId string
type ChTitleID map[ChTitle]ChId

func (t ChTitleID) New(data *YTInitialData) {
	for _, tab := range data.Contents.TwoColumnBrowseResultsRenderer.Tabs {
		for _, content1 := range tab.TabRenderer.Content.SectionListRenderer.Contents {
			for _, content2 := range content1.ItemSectionRenderer.Contents {
				for _, item := range content2.ShelfRenderer.Content.ExpandedShelfContentsRenderer.Items {
					t[ChTitle(strings.TrimSpace(item.ChannelRenderer.Title.SimpleText))] = ChId(strings.TrimSpace(item.ChannelRenderer.ChannelID))
				}
			}
		}
	}
}

// cspell:words ytinitialdata

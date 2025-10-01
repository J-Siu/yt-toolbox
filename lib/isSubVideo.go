/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package lib

import (
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsSubVideo struct {
	*is.Processor

	Day int
}

func (s *IsSubVideo) New(page *rod.Page, urlStr string, scrollMax int) *IsSubVideo {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	s.Processor = is.New(&property) // Init the base struct
	s.MyType = "IsSubVideo"
	prefix := s.MyType + ".New"

	s.override()

	ezlog.Trace().N(prefix).M("Done")
	return s
}

func (s *IsSubVideo) override() {
	s.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := s.MyType + ".V020_Elements"
		ezlog.Trace().N(prefix).TxtStart().Out()

		var elements rod.Elements
		tagName := "ytd-rich-item-renderer"
		s.Page.MustElement(tagName)
		elements = s.Page.MustElements(tagName)
		ezlog.Debug().N(prefix).N("elements count").M(len(elements)).Out()

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &elements
	}

	// [element] : "ytd-rich-item-renderer" element from V20_elements()
	s.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := s.MyType + ".V030_ElementInfo"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if element != nil {
			var (
				info    YT_Info
				tagName string
			)
			// Tile block("h3"): title and link of the video
			eH3 := element.MustElement("h3")
			info.Title = eH3.MustText()
			info.Url = UrlYT.Base + *eH3.MustElement("a").MustAttribute("href")

			// Meta element: channel info, views and date
			tagName = "yt-content-metadata-view-model"
			eMeta, err := element.Element(tagName)
			if err == nil && eMeta != nil {
				// Meta element -> link(<a>) block
				a := eMeta.MustElement("a")
				info.ChName = a.MustText()
				info.ChUrlShort = UrlDecode(*a.MustAttribute("href"))
				info.ChUrl = UrlYT.Base + info.ChUrlShort
				// Meta element -> elements with [role]='text' attribute
				tagName = "[role='text']"
				eRoles, err := eMeta.Elements(tagName)
				if err == nil {
					excludeText := []string{"views", "watch", "scheduled"}
					for _, eRole := range eRoles {
						text := eRole.MustText()
						if !str.ContainsAnySubStringsBool(&text, &excludeText) {
							info.Text = text
						}
						// search for watching, minutes, hours, day, <date>
					}
				}
			}
			if err != nil {
				// These are shorts with not meta block
				info.Text = "Short"
				// if ezlog.GetLogLevel() == ezlog.TraceLevel {
				// 	ezlog.Trace().N(prefix).Nn("Err element").M(element.MustHTML()).Out()
				// }
			}
			// ---
			ezlog.Debug().Nn(prefix).M(info).Out()
			infoP = &info
		}
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return infoP
	}
}

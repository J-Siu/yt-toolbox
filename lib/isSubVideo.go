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
	"strconv"

	"github.com/J-Siu/go-ezlog"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsSubVideo struct {
	*is.Processor

	Day int
}

func (s *IsSubVideo) New(pageP *rod.Page, urlStr string, scrollMax int) *IsSubVideo {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      pageP,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	s.Processor = is.New(&property) // Init the base struct
	s.MyType = "IsSubscription"
	prefix := s.MyType + ".New"

	s.initFunc()

	ezlog.Trace(prefix + ": Done")
	return s
}

func (s *IsSubVideo) initFunc() {
	s.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := s.MyType + ".V020_Elements"
		ezlog.Trace(prefix + ": Start")

		var elements rod.Elements
		tagName := "ytd-rich-item-renderer"
		s.Page.MustElement(tagName)
		elements = s.Page.MustElements(tagName)
		ezlog.Debug(prefix + ": es count: " + strconv.Itoa(len(elements)))
		ezlog.Trace(prefix + ": End")
		return &elements
	}

	// [element] : "ytd-rich-item-renderer" element from V20_elements()
	s.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := s.MyType + ".V030_ElementInfo"
		ezlog.Trace(prefix + ": Start")

		if element != nil {
			var info YT_Info

			// Tile block("h3"): title and link of the video
			eH3 := element.MustElement("h3")
			info.Title = eH3.MustText()
			info.Url = UrlYT.Base + *eH3.MustElement("a").MustAttribute("href")

			// Meta element: channel info, views and date
			tagName := "yt-content-metadata-view-model"
			eMeta, err := element.Element(tagName)
			if err == nil && eMeta != nil {
				// Meta element -> link(<a>) block
				a := eMeta.MustElement("a")
				info.ChName = a.MustText()
				info.ChUrlShort = UrlDecode(*a.MustAttribute("href"))
				info.ChUrl = UrlYT.Base + info.ChUrlShort
				// Meta element -> elements with [role]='text' attribute
				eRoles, err := eMeta.Elements("[role='text']")
				if err == nil {
					excludeText := []string{"views", "watch", "scheduled"}
					for _, eRole := range eRoles {
						text := eRole.MustText()
						if !MustStrMatchList(text, &excludeText) {
							info.Text = text
						}
						// search for watching, minutes, hours, day, <date>
					}
				}
			}
			if err != nil {
				info.Text = "Short" // These are shorts with not meta block

				if ezlog.LogLevel() == ezlog.TraceLevel {
					ezlog.Trace(prefix + ": err: element")
					ezlog.Trace(element.MustHTML())
				}
			}
			// ---

			// ---
			ezlog.Debug(prefix + info.String())
			infoP = &info
		}
		ezlog.Trace(prefix + ": End")
		return infoP
	}
}

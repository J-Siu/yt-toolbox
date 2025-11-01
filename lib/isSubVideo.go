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
	"regexp"
	"strconv"

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsSubVideo struct {
	*is.Processor

	Scroll bool
	Day    uint
}

func (t *IsSubVideo) New(page *rod.Page, urlStr string, scrollMax int, day uint) *IsSubVideo {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	t.Processor = is.New(&property) // Init the base struct
	t.MyType = "IsSubVideo"
	prefix := t.MyType + ".New"
	t.Day = day
	t.override()

	// ezlog.Trace().N(prefix).M("Done")
	ezlog.Trace().N(prefix).Lm(t).Out()
	return t
}

func (t *IsSubVideo) Run() *IsSubVideo {
	t.Processor.Run()
	return t
}

func (t *IsSubVideo) override() {
	t.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := t.MyType + ".V020_Elements"
		ezlog.Trace().N(prefix).TxtStart().Out()

		var elements rod.Elements
		tagName := "ytd-rich-item-renderer"
		t.Page.MustElement(tagName)
		elements = t.Page.MustElements(tagName)
		ezlog.Debug().N(prefix).N("elements count").M(len(elements)).Out()

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &elements
	}

	// [element] : "ytd-rich-item-renderer" element from V20_elements()
	t.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := t.MyType + ".V030_ElementInfo"
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
						if !str.ContainsAnySubStringsBool(&text, &excludeText, false) {
							info.Text = text
							t.dayScroll(&text)
						}
						// search for watching, minutes, hours, day, <date>
					}
				}
			}
			if err != nil {
				// These are shorts with not meta block
				info.Text = "Short"
				// if ezlog.GetLogLevel() == ezlog.TRACE {
				// 	ezlog.Trace().N(prefix).Nn("Err element").M(element.MustHTML()).Out()
				// }
			}
			// ---
			ezlog.Debug().N(prefix).Lm(info).Out()
			infoP = &info
		}
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return infoP
	}

	t.V100_ScrollLoopEnd = func(state *is.State) {
		if t.Day > 0 && t.Scroll {
			t.ScrollMax = -1
		} else {
			t.ScrollMax = 0
		}
	}
}

func (t *IsSubVideo) dayScroll(text *string) {
	// only calculate if t.Day > 0
	if t.Day > 0 {
		prefix := t.MyType + ".dayScroll"

		var (
			day     uint64
			pattern string
			matches [][]string
			re      *regexp.Regexp
			e       error
		)
		// only update t.scroll if text is time
		ezlog.Trace().N(prefix).N("text").M(*text).Out()

		// Following count as 1 day: hour, now, second
		pattern = `(\d+) (now?|hour?|minute?|second?)`
		re = regexp.MustCompile(pattern)
		matches = re.FindAllStringSubmatch(*text, -1)
		if len(matches) > 0 {
			ezlog.Trace().N(prefix).M("1 day").Out()
			day = 1
		}

		// # day
		if day == 0 {
			pattern = `(\d+) (day.*)`
			re = regexp.MustCompile(pattern)
			matches = re.FindAllStringSubmatch(*text, -1)
			ezlog.Trace().N(prefix).N("day matches").M(matches).Out()
			if len(matches) > 0 {
				ezlog.Trace().N(prefix).M("day").Out()
				day, e = strconv.ParseUint(matches[0][1], 10, 64)
				if e != nil {
					ezlog.Err().M(e).Out()
					day = 0
				}
			}
		}
		// # week
		if day == 0 {
			pattern = `(\d+) (week.*)`
			re = regexp.MustCompile(pattern)
			matches = re.FindAllStringSubmatch(*text, -1)
			ezlog.Trace().N(prefix).N("week matches").M(matches).Out()
			if len(matches) > 0 {
				ezlog.Trace().N(prefix).M("week").Out()
				day, e = strconv.ParseUint(matches[0][1], 10, 64)
				if e == nil {
					day *= 7
				}
			}
		}
		// # month
		if day == 0 {
			pattern = `(\d+) (month.*)`
			re = regexp.MustCompile(pattern)
			matches = re.FindAllStringSubmatch(*text, -1)
			ezlog.Trace().N(prefix).N("month matches").M(matches).Out()
			if len(matches) > 0 {
				ezlog.Trace().N(prefix).M("month").Out()
				day, e = strconv.ParseUint(matches[0][1], 10, 64)
				if e == nil {
					day *= 30
				}
			}
		}
		// # year
		if day == 0 {
			pattern = `(\d+) (year.*)`
			re = regexp.MustCompile(pattern)
			matches = re.FindAllStringSubmatch(*text, -1)
			ezlog.Trace().N(prefix).N("year matches").M(matches).Out()
			if len(matches) > 0 {
				ezlog.Trace().N(prefix).M("year").Out()
				day, e = strconv.ParseUint(matches[0][1], 10, 64)
				if e == nil {
					day *= 365
				}
			}
		}

		t.Scroll = true
		if day > uint64(t.Day) {
			t.Scroll = false
		}
		ezlog.Trace().N(prefix).N("day").M(day).N("scroll").M(t.Scroll).Out()
	}
}

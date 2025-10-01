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
	"net/url"
	"strings"

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsHistoryVideo struct {
	*is.Processor

	ClickSleep  float64 // in second
	Del         bool    // false;
	Deleted     bool    // In Run(), elements loop, current element is deleted or not
	Desc        bool    // false;
	PrintHeader bool    // true;
	Remove      bool    // false;
	Standalone  bool    // false;
	Verbose     bool    // false;
	Filter      []string
}

func (s *IsHistoryVideo) New(property *is.Property, del bool, remove bool, filter *[]string) *IsHistoryVideo {
	s.Processor = is.New(property) // Init the base struct
	s.MyType = "IsHistoryVideo"

	s.Filter = append(s.Filter, *filter...)
	s.Del = del
	s.Remove = remove
	s.override()

	return s
}

func (s *IsHistoryVideo) override() {
	s.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := s.MyType + ".V020_Elements"
		ezlog.Trace().N(prefix).TxtStart().Out()

		s.V021_ElementsRemoveShorts(element)
		tagName := "ytd-video-renderer"
		var elements rod.Elements
		elements, s.Err = element.Elements(tagName)
		if s.Err != nil {
			ezlog.Err().N(prefix).M(s.Err).Out()
		}
		ezlog.Debug().N(prefix).N("elements count").M(len(elements)).Out()
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &elements
	}

	s.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := s.MyType + ".V030_ElementInfo"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if element != nil {
			var info YT_Info
			byId := "#video-title"
			elementTitle := element.MustElement(byId) // by id
			if ezlog.GetLogLevel() == ezlog.TraceLevel {
				ezlog.Trace().N(prefix).Mn(byId).M(elementTitle.MustHTML()).Out()
			}
			info.Title = strings.TrimSpace(elementTitle.MustText())
			if len(info.Title) != 0 {
				info.Url = *elementTitle.MustAttribute("href")
				info.Text = strings.TrimSpace(element.MustElement("#description-text").MustText()) // by id

				meta := element.MustElement("#metadata") // by id
				a := meta.MustElement("a")               // by name
				info.ChName = strings.TrimSpace(a.MustText())

				if a != nil {
					chUrlP := a.MustAttribute("href")
					if chUrlP != nil {
						info.ChUrl = *chUrlP
						parsedUrl, err := url.Parse(info.ChUrl)
						if err == nil {
							info.ChUrlShort = parsedUrl.Path
						}
					} else {
						ezlog.Debug().N(prefix).M("chUrlP == nil").Out()
					}
				}
			}
			infoP = &info
			ezlog.Debug().N(prefix).N("info").M(infoP.String()).Out()
		}

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return infoP
	}

	s.V040_ElementMatch = func(element *rod.Element, index int, infoP is.IInfo) (matched bool, matchedStr string) {
		prefix := s.MyType + ".V040_ElementMatch"
		ezlog.Trace().N(prefix).TxtStart().Out()

		info := infoP.(*YT_Info)
		chkStr := info.Title + " " + info.Text + " " + info.ChName + " " + info.ChUrlShort

		matched, matchedStr = str.ContainsAnySubStrings(&chkStr, &s.Filter)

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return matched, matchedStr
	}

	s.V050_ElementProcessMatched = func(element *rod.Element, index int, infoP is.IInfo) {
		prefix := s.MyType + ".V050_ElementProcessMatched"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if s.Del {
			s.Deleted = true
			element.MustElement("#top-level-buttons-computed").MustClick()
			// time.Sleep(time.Duration(s.ClickSleep))
			ezlog.Debug().N(prefix).M("Deleted").Out()
		}

		ezlog.Trace().N(prefix).TxtEnd().Out()
	}

	s.V060_ElementProcessUnmatch = func(element *rod.Element, index int, infoP is.IInfo) {
		prefix := s.MyType + ".V060_ElementProcessUnmatch"
		ezlog.Trace().N(prefix).M("Done").Out()
	}

	s.V080_ElementScrollable = func(element *rod.Element, index int, infoP is.IInfo) bool {
		prefix := s.MyType + ".V080_ElementScrollable"
		ezlog.Trace().N(prefix).TxtStart().Out()

		info := infoP.(*YT_Info)
		notScrollable := s.Deleted
		if (!notScrollable && len(info.Title) == 0) || !element.MustVisible() {
			notScrollable = true
		}
		s.Deleted = false

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return !notScrollable
	}

	s.V100_ScrollLoopEnd = func(state *is.State) {
		prefix := s.MyType + ".V100_ScrollLoopEnd"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if s.Remove {
			if state.Elements != nil {
				removed := 0
				for _, e := range *state.Elements {
					e.Remove()
					removed++
				}
				state.ElementCountLast = 0
				state.ElementLast = nil
			}
			es := s.Page.MustElements("ytd-item-section-renderer")
			if es != nil {
				state.ElementLast = es.Last()
				if state.ElementLastScroll != nil && state.ElementLast != nil && (state.ElementLastScroll.Object.ObjectID != state.ElementLast.Object.ObjectID) {
					state.Scroll = true
				} else if state.ElementLastScroll == nil && state.ElementLast != nil {
					state.Scroll = true
				} else {
					state.Scroll = false
				}
			} else {
				state.Scroll = false
			}
			if state.ElementLastScroll != nil {
				ezlog.Trace().N(prefix).N("state.ElementLastScroll").M(string(state.ElementLastScroll.Object.ObjectID)).Out()
			}
			if state.ElementLast != nil {
				ezlog.Trace().N(prefix).N("state.ElementLast").M(string(state.ElementLast.Object.ObjectID)).Out()
			}
		}
		ezlog.Trace().N(prefix).TxtEnd().Out()
	}

}

func (s *IsHistoryVideo) V021_ElementsRemoveShorts(element *rod.Element) *IsHistoryVideo {
	prefix := s.MyType + ".V020_ElementsRemoveShorts"
	ezlog.Trace().N(prefix).TxtStart().Out()

	var tagName = "ytd-reel-shelf-renderer"
	var es rod.Elements
	es, s.Err = element.Elements(tagName)
	if s.Err != nil {
		ezlog.Err().N(prefix).M(s.Err).Out()
	}
	es_count := len(es)
	for i := range es {
		es[es_count-i-1].Remove()
	}

	ezlog.Trace().N(prefix).TxtEnd().Out()
	return s
}

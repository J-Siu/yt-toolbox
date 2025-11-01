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

func (t *IsHistoryVideo) New(property *is.Property, del bool, remove bool, filter *[]string) *IsHistoryVideo {
	t.Processor = is.New(property) // Init the base struct
	t.MyType = "IsHistoryVideo"

	t.Filter = append(t.Filter, *filter...)
	t.Del = del
	t.Remove = remove
	t.override()

	return t
}

func (t *IsHistoryVideo) Run() *IsHistoryVideo {
	t.Processor.Run()
	return t
}

func (t *IsHistoryVideo) override() {
	t.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := t.MyType + ".V020_Elements"
		ezlog.Trace().N(prefix).TxtStart().Out()

		t.V021_ElementsRemoveShorts(element)
		tagName := "ytd-video-renderer"
		var elements rod.Elements
		elements, t.Err = element.Elements(tagName)
		if t.Err != nil {
			ezlog.Err().N(prefix).M(t.Err).Out()
		}
		ezlog.Debug().N(prefix).N("elements count").M(len(elements)).Out()
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &elements
	}

	t.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := t.MyType + ".V030_ElementInfo"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if element != nil {
			var info YT_Info
			byId := "#video-title"
			elementTitle := element.MustElement(byId) // by id
			if ezlog.GetLogLevel() == ezlog.TRACE {
				ezlog.Trace().N(prefix).M(byId).Lm(elementTitle.MustHTML()).Out()
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

	t.V040_ElementMatch = func(element *rod.Element, index int, infoP is.IInfo) (matched bool, matchedStr string) {
		prefix := t.MyType + ".V040_ElementMatch"
		ezlog.Trace().N(prefix).TxtStart().Out()

		info := infoP.(*YT_Info)
		chkStr := info.Title + " " + info.Text + " " + info.ChName + " " + info.ChUrlShort

		matched, matchedStr = str.ContainsAnySubStrings(&chkStr, &t.Filter, false)

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return matched, matchedStr
	}

	t.V050_ElementProcessMatched = func(element *rod.Element, index int, infoP is.IInfo) {
		prefix := t.MyType + ".V050_ElementProcessMatched"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if t.Del {
			t.Deleted = true
			element.MustElement("#top-level-buttons-computed").MustClick()
			// time.Sleep(time.Duration(s.ClickSleep))
			ezlog.Debug().N(prefix).M("Deleted").Out()
		}

		ezlog.Trace().N(prefix).TxtEnd().Out()
	}

	t.V060_ElementProcessUnmatch = func(element *rod.Element, index int, infoP is.IInfo) {
		prefix := t.MyType + ".V060_ElementProcessUnmatch"
		ezlog.Trace().N(prefix).M("Done").Out()
	}

	t.V080_ElementScrollable = func(element *rod.Element, index int, infoP is.IInfo) bool {
		prefix := t.MyType + ".V080_ElementScrollable"
		ezlog.Trace().N(prefix).TxtStart().Out()

		info := infoP.(*YT_Info)
		notScrollable := t.Deleted
		if (!notScrollable && len(info.Title) == 0) || !element.MustVisible() {
			notScrollable = true
		}
		t.Deleted = false

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return !notScrollable
	}

	t.V100_ScrollLoopEnd = func(state *is.State) {
		prefix := t.MyType + ".V100_ScrollLoopEnd"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if t.Remove {
			if state.Elements != nil {
				removed := 0
				for _, e := range *state.Elements {
					e.Remove()
					removed++
				}
				state.ElementCountLast = 0
				state.ElementLast = nil
			}
			es := t.Page.MustElements("ytd-item-section-renderer")
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

func (t *IsHistoryVideo) V021_ElementsRemoveShorts(element *rod.Element) *IsHistoryVideo {
	prefix := t.MyType + ".V020_ElementsRemoveShorts"
	ezlog.Trace().N(prefix).TxtStart().Out()

	var tagName = "ytd-reel-shelf-renderer"
	var es rod.Elements
	es, t.Err = element.Elements(tagName)
	if t.Err != nil {
		ezlog.Err().N(prefix).M(t.Err).Out()
	}
	es_count := len(es)
	for i := range es {
		es[es_count-i-1].Remove()
	}

	ezlog.Trace().N(prefix).TxtEnd().Out()
	return t
}

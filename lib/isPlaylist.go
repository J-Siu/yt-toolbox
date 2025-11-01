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

type IsPlaylist struct {
	*is.Processor

	Exclude *[]string `json:"Exclude"`
	Include *[]string `json:"Include"`
}

func (t *IsPlaylist) New(page *rod.Page, urlStr string, scrollMax int, exclude *[]string, include *[]string) *IsPlaylist {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	t.Processor = is.New(&property) // Init the base struct
	t.MyType = "IsPlaylist"
	prefix := t.MyType + ".New"
	if t.Err != nil {
		ezlog.Err().N(prefix).M(t.Err).Out()
	}
	t.Exclude = exclude
	t.Include = include

	ezlog.Debug().N(prefix).Lm(t).Out()
	t.override()

	return t
}

func (t *IsPlaylist) Run() *IsPlaylist {
	t.Processor.Run()
	return t
}

func (t *IsPlaylist) override() {
	t.V010_Container = func() *rod.Element {
		prefix := t.MyType + ".V010_ElementsContainer"
		ezlog.Trace().N(prefix).TxtStart().Out()
		byId := "#contents"
		e := t.Page.MustElement(byId) // by id
		if ezlog.GetLogLevel() == ezlog.TRACE {
			ezlog.Trace().N(prefix).M(byId).Lm(e.MustHTML()).Out()
		}
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return e
	}

	t.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := t.MyType + ".V020_Elements"
		ezlog.Trace().N(prefix).TxtStart().Out()
		tagName := "ytd-rich-item-renderer"
		elements := element.MustElements(tagName)
		ezlog.Debug().N(prefix).N(tagName).N("element count").M(len(elements)).Out()

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &elements
	}

	t.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := t.MyType + ".V030_ElementInfo"
		ezlog.Trace().N(prefix).TxtStart().Out()

		if element != nil {
			var info YT_Info
			if ezlog.GetLogLevel() == ezlog.TRACE {
				ezlog.Trace().N(prefix).N("element").Lm(element.MustHTML()).Out()
			}
			h3 := element.MustElement("h3")
			if h3 != nil {
				info.Title = *(h3.MustAttribute("title"))
				ezlog.Debug().N(prefix).N("Title").M(info.Title).Out()
			}
			tagName := "a"
			es := element.MustElements(tagName)
			for _, s := range es {
				if s.MustText() == "View full playlist" {
					if ezlog.GetLogLevel() == ezlog.TRACE {
						ezlog.Trace().N(prefix).N(tagName).Lm(s.MustHTML()).Out()
					}
					info.Url = UrlYT.Base + *(s.MustAttribute("href"))
				}
			}
			infoP = &info
		}

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return infoP
	}

	t.V040_ElementMatch = func(element *rod.Element, index int, info is.IInfo) (matched bool, matchedStr string) {
		prefix := t.MyType + ".V040_ElementMatch"
		yt_Info := info.(*YT_Info)
		matched = true // default to matched
		if len(*t.Include) != 0 {
			matched, matchedStr = str.ContainsAnySubStrings(&yt_Info.Title, t.Include, false)
		}
		// Exclude override Include
		if len(*t.Exclude) != 0 {
			matched, matchedStr = str.ContainsAnySubStrings(&yt_Info.Title, t.Exclude, false)
			if matched {
				matched = false
			}
		}
		ezlog.Trace().N(prefix).N("matched").M(matched).N("matchedStr").M(matchedStr).Out()
		return matched, matchedStr
	}
}

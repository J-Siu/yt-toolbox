/*
Copyright © 2026 John, Sing Dao, Siu <john.sd.siu@gmail.com>

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
	"github.com/J-Siu/go-is/v3/is"
	"github.com/go-rod/rod"
	"github.com/yosssi/gohtml"
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
	t.V010_Container = t.override_V010_Container
	t.V020_Elements = t.override_V020_Elements
	t.V030_ElementInfo = t.override_V030_ElementInfo
	t.V040_ElementMatch = t.override_V040_ElementMatch
}

func (t *IsPlaylist) override_V010_Container() {
	prefix := t.MyType + ".V010_ElementsContainer"
	t.StateCurr.Name = prefix
	byId := "#contents"
	t.Container = t.Page.MustElement(byId) // by id
	if ezlog.GetLogLevel() == ezlog.TRACE {
		ezlog.Trace().N(prefix).M(byId).Lm(gohtml.Format(t.Container.MustHTML())).Out()
	}
}

func (t *IsPlaylist) override_V020_Elements() {
	prefix := t.MyType + ".V020_Elements"
	t.StateCurr.Name = prefix
	tagName := "ytd-rich-item-renderer"
	t.StateCurr.Elements = t.Container.MustElements(tagName)
	ezlog.Debug().N(prefix).N(tagName).N("element count").M(len(t.StateCurr.Elements)).Out()
}

func (t *IsPlaylist) override_V030_ElementInfo() {
	prefix := t.MyType + ".V030_ElementInfo"
	t.StateCurr.Name = prefix
	if t.StateCurr.Element != nil {
		var info YT_Info
		if ezlog.GetLogLevel() == ezlog.TRACE {
			ezlog.Trace().N(prefix).N("element").Lm(gohtml.Format(t.StateCurr.Element.MustHTML())).Out()
		}
		h3 := t.StateCurr.Element.MustElement("h3")
		if h3 != nil {
			info.Title = *(h3.MustAttribute("title"))
			ezlog.Debug().N(prefix).N("Title").M(info.Title).Out()
		}
		tagName := "a"
		es := t.StateCurr.Element.MustElements(tagName)
		for _, s := range es {
			if s.MustText() == "View full playlist" {
				if ezlog.GetLogLevel() == ezlog.TRACE {
					ezlog.Trace().N(prefix).N(tagName).Lm(gohtml.Format(s.MustHTML())).Out()
				}
				info.Url = YT_FullUrl(*s.MustAttribute("href"))
			}
		}
		t.StateCurr.ElementInfo = &info
	}
}

func (t *IsPlaylist) override_V040_ElementMatch() {
	prefix := t.MyType + ".V040_ElementMatch"
	t.StateCurr.Name = prefix
	var (
		yt_Info         = t.StateCurr.ElementInfo.(*YT_Info)
		matched    bool = true // start with matched
		matchedStr string
	)
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
	t.StateCurr.ElementInfo.SetMatched(matched)
	t.StateCurr.ElementInfo.SetMatchedStr(matchedStr)
	ezlog.Trace().N(prefix).N("matched").M(matched).N("matchedStr").M(matchedStr).Out()
}

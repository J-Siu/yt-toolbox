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
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsPlaylistVideo struct {
	*is.Processor
}

func (t *IsPlaylistVideo) New(page *rod.Page, urlStr string, scrollMax int) *IsPlaylistVideo {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	t.Processor = is.New(&property) // Init the base struct
	t.MyType = "IsPlaylistVideo"

	t.override()

	return t
}

func (t *IsPlaylistVideo) Run() *IsPlaylistVideo {
	t.Processor.Run()
	return t
}

func (t *IsPlaylistVideo) override() {
	t.V010_Container = func() *rod.Element {
		prefix := t.MyType + ".V010_ElementsContainer"
		ezlog.Trace().N(prefix).TxtStart().Out()
		tagName := "ytd-playlist-video-list-renderer"
		element := t.Page.MustElement(tagName)
		ezlog.Debug().N(prefix).N(tagName).Lm(element).Out()

		ezlog.Trace().N(prefix).TxtEnd().Out()
		return element
	}

	t.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := t.MyType + ".V020_Elements"
		ezlog.Trace().N(prefix).TxtStart().Out()
		tagName := "ytd-playlist-video-renderer"
		es := element.MustElements(tagName)
		ezlog.Debug().N(prefix).N(tagName).N("count").M(len(es)).Out()
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return &es
	}

	t.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := t.MyType + ".V030_ElementInfo"
		ezlog.Trace().N(prefix).TxtStart().Out()
		if element != nil {
			var info YT_Info
			e := element.MustElement("#video-title")
			info.Title = e.MustText()
			info.Url = *e.MustAttribute("href")
			infoP = &info
		}
		ezlog.Trace().N(prefix).TxtEnd().Out()
		return infoP
	}
}

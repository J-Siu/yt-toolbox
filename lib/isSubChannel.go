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
	"github.com/J-Siu/go-ezlog/v2"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsSubChannel struct {
	*is.Processor
}

func (s *IsSubChannel) New(page *rod.Page, urlStr string, scrollMax int) *IsSubChannel {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	s.Processor = is.New(&property) // Init the base struct
	s.MyType = "IsSubChannel"

	s.override()
	return s
}

func (s *IsSubChannel) override() {
	s.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := s.MyType + ".V020_Elements"
		ezlog.Trace().Name(prefix).Msg("Start").Out()

		var elements rod.Elements
		s.Page.MustElement("#content-section").MustWaitVisible()
		elements = s.Page.MustElements("#content-section")

		ezlog.Trace().Name(prefix).Msg("End").Out()
		return &elements
	}

	s.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := s.MyType + ".V030_ElementInfo"
		ezlog.Trace().Name(prefix).Msg("Start").Out()

		if element != nil {
			var info YT_Info
			var title string
			var urlPath *string
			title = element.MustElement("#text").MustText()
			urlPath = element.MustElement("#main-link").MustAttribute("href")
			info.Title = title
			info.Url = UrlYT.Base + *urlPath
			ezlog.Debug().Name(prefix).Msg(info.String()).Out()
			infoP = &info
		}
		ezlog.Trace().Name(prefix).Msg("End").Out()
		return infoP
	}
}

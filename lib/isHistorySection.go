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

	"github.com/J-Siu/go-ezlog/v2"
	"github.com/J-Siu/go-is"
	"github.com/go-rod/rod"
)

type IsHistorySection struct {
	*is.Processor

	Del         bool // false;
	Deleted     bool // false;
	Desc        bool // false;
	PrintHeader bool // true;
	Remove      bool // false;
	Standalone  bool // false;
	Verbose     bool // false;
	Filter      []string
}

func (s *IsHistorySection) New(page *rod.Page, urlStr string, remove bool, scrollMax int, verbose bool) *IsHistorySection {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	s.Processor = is.New(&property) // Init the base struct
	s.MyType = "IsHistorySection"

	s.Remove = remove
	s.Verbose = verbose
	s.override()

	return s
}

func (s *IsHistorySection) override() {
	s.V020_Elements = func(element *rod.Element) *rod.Elements {
		prefix := s.MyType + ".V020_Elements"
		ezlog.Trace().Name(prefix).Msg("Start").Out()

		var tagName = "ytd-item-section-renderer"
		e := s.Page.MustElement(tagName)
		ezlog.Trace().Name(prefix).Name("MustWaitDOMStable").Msg("Start").Out()
		e.MustWaitVisible()
		ezlog.Trace().Name(prefix).Name("MustWaitDOMStable").Msg("End").Out()
		elements := s.Page.MustElements(tagName)
		ezlog.Trace().Name(prefix).Name(tagName).Name("element count").Msg(len(elements)).Out()

		ezlog.Trace().Name(prefix).Msg("End").Out()
		return &elements
	}

	s.V030_ElementInfo = func(element *rod.Element, index int) (infoP is.IInfo) {
		prefix := s.MyType + ".V030_ElementInfo"
		ezlog.Trace().Name(prefix).Msg("Start").Out()

		if element != nil {
			var info YT_Info
			var titles rod.Elements
			titles, s.Err = element.Elements("#title") //by id
			if s.Err != nil {
				ezlog.Err().Name(prefix).Name("Err").Msg(s.Err).Out()
			}

			for j := range len(titles) {
				title := titles[j].MustText()
				ezlog.Log().Ln().Name("## Section[" + strconv.Itoa(index) + "] Title[" + strconv.Itoa(j) + "]").Msg(title).Out()
				info.Titles = append(info.Titles, title)
			}
			if len(info.Titles) == 0 {
				info.Titles = []string{""}
			}
			infoP = &info
		}

		ezlog.Trace().Name(prefix).Msg("End").Out()
		return infoP
	}

	s.V070_ElementProcess = func(element *rod.Element, index int, infoP is.IInfo) {
		prefix := s.MyType + ".V070_ElementProcess"
		ezlog.Trace().Name(prefix).Msg("Start").Out()

		titles := &infoP.(*YT_Info).Titles
		if len(*titles) != 0 && len((*titles)[0]) != 0 {
			var isHistoryVideo IsHistoryVideo
			property := is.Property{
				Container: element,
				IInfoList: new(is.IInfoList),
				Page:      s.Page,
			}

			isHistoryVideo.New(&property, s.Del, s.Remove, &s.Filter).Run()
			var mode is.IInfoListPrintMode
			if s.Verbose {
				mode = is.PrintAll
			} else {
				mode = is.PrintMatched
			}
			isHistoryVideo.IInfoList.Print(mode)
		}

		ezlog.Trace().Name(prefix).Msg("End").Out()
	}

	s.V100_ScrollLoopEnd = func(state *is.State) {
		prefix := s.MyType + "V100_ScrollLoopEnd"
		ezlog.Trace().Name(prefix).Msg("Start").Out()

		if s.Remove {
			s.removeSpinningWheel()
			if state.Elements != nil {
				for _, e := range *state.Elements {
					e.Remove()
				}
			}
			state.ElementCountLast = 0
			state.ElementLast = nil
		}
		state.Scroll = true
		ezlog.Trace().Name(prefix).Name("MustWaitLoad").Msg("Start").Out()
		s.Page.MustWaitLoad()
		ezlog.Trace().Name(prefix).Name("MustWaitLoad").Msg("End").Out()

		ezlog.Trace().Name(prefix).Msg("End").Out()
	}
}

func (s *IsHistorySection) removeSpinningWheel() {
	es := s.Page.MustElements("ytd-continuation-item-renderer") // tag name
	if es != nil {
		count := len(es)
		var removed int
		for _, e := range es {
			removed++
			if removed == count {
				break
			}
			e.Remove()
		}
	}
}

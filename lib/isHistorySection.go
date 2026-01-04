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
	"github.com/J-Siu/go-helper/v2/strany"
	"github.com/J-Siu/go-is/v2/is"
	"github.com/go-rod/rod"
)

// Process YT History by sections
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

func (t *IsHistorySection) New(page *rod.Page, urlStr string, remove bool, scrollMax int, verbose bool) *IsHistorySection {
	property := is.Property{
		IInfoList: new(is.IInfoList),
		Page:      page,
		ScrollMax: scrollMax,
		UrlLoad:   true,
		UrlStr:    urlStr,
	}
	t.Processor = is.New(&property) // Init the base struct
	t.MyType = "IsHistorySection"

	t.Remove = remove
	t.Verbose = verbose
	t.override()

	return t
}

func (t *IsHistorySection) Run() *IsHistorySection {
	t.Processor.Run()
	return t
}

func (t *IsHistorySection) override() {
	t.V020_Elements = t.override_V020_Elements
	t.V030_ElementInfo = t.override_V030_ElementInfo
	t.V070_ElementProcess = t.override_V070_ElementProcess
	t.V100_ScrollLoopEnd = t.override_V100_ScrollLoopEnd
}

func (t *IsHistorySection) override_V020_Elements(element *rod.Element) *rod.Elements {
	prefix := t.MyType + ".V020_Elements"
	ezlog.Debug().N(prefix).TxtStart().Out()

	var tagName = "ytd-item-section-renderer"
	e := t.Page.MustElement(tagName)
	ezlog.Trace().N(prefix).N("MustWaitDOMStable").TxtStart().Out()
	e.MustWaitVisible()
	ezlog.Trace().N(prefix).N("MustWaitDOMStable").TxtEnd().Out()
	elements := t.Page.MustElements(tagName)
	ezlog.Trace().N(prefix).N(tagName).N("element count").M(len(elements)).Out()

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return &elements
}

func (t *IsHistorySection) override_V030_ElementInfo() (infoP is.IInfo) {
	prefix := t.MyType + ".V030_ElementInfo"
	ezlog.Debug().N(prefix).TxtStart().Out()

	if t.StateCurr.Element != nil {
		var (
			info     YT_Info
			elements rod.Elements
			byId     = "#title"
		)
		elements, t.Err = t.StateCurr.Element.Elements(byId)
		if t.Err != nil {
			ezlog.Err().N(prefix).M(t.Err).Out()
		}

		for j, item := range elements {
			title := item.MustText()
			tmp := "## Section[" + *strany.Any(t.StateCurr.ElementIndex) + "] Title[" + *strany.Any(j) + "]"
			ezlog.Log().L().N(tmp).M(title).Out()
			info.Titles = append(info.Titles, title)
		}
		if len(info.Titles) == 0 {
			info.Titles = []string{""}
		}
		infoP = &info
	}

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return infoP
}

func (t *IsHistorySection) override_V070_ElementProcess() {
	prefix := t.MyType + ".V070_ElementProcess"
	ezlog.Debug().N(prefix).TxtStart().Out()

	titles := t.StateCurr.ElementInfo.(*YT_Info).Titles
	if len(titles) != 0 && len((titles)[0]) != 0 {
		var (
			isHistoryEntry IsHistoryEntry
			mode           is.IInfoListPrintMode
			property       = is.Property{
				Container: t.StateCurr.Element,
				IInfoList: new(is.IInfoList),
				Page:      t.Page,
			}
		)
		isHistoryEntry.New(&property, t.Del, t.Remove, &t.Filter).Run()
		if t.Verbose {
			mode = is.PrintAll
		} else {
			mode = is.PrintMatched
		}
		ezlog.Log().M("#|match|video|ch|desc").Out()
		ezlog.Log().M("--|--|--|--|--").Out()
		isHistoryEntry.IInfoList.Print(mode)
	}

	ezlog.Debug().N(prefix).TxtEnd().Out()
}

func (t *IsHistorySection) override_V100_ScrollLoopEnd() {
	prefix := t.MyType + "V100_ScrollLoopEnd"
	ezlog.Debug().N(prefix).TxtStart().Out()

	if t.Remove {
		t.removeSpinningWheel()
		if t.StateCurr.Elements != nil {
			for _, e := range *t.StateCurr.Elements {
				e.Remove()
			}
		}
		t.StateCurr.ElementsCount = 0
		t.StateCurr.Element = nil
		t.StateCurr.ScrollableElement = nil
	}
	t.StateCurr.Scroll = true
	ezlog.Trace().N(prefix).N("MustWaitLoad").TxtStart().Out()
	t.Page.MustWaitLoad()
	ezlog.Trace().N(prefix).N("MustWaitLoad").TxtEnd().Out()

	ezlog.Debug().N(prefix).TxtEnd().Out()
}

func (t *IsHistorySection) removeSpinningWheel() {
	es := t.Page.MustElements("ytd-continuation-item-renderer") // tag name
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

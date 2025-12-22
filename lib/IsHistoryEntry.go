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
	"math/rand/v2"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-is/v2/is"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/yosssi/gohtml"
)

// process YT history by entry
type IsHistoryEntry struct {
	*is.Processor

	ClickSleep  float64 // in second
	Del         bool    // delete entry from history
	Deleted     bool    // In Run(), elements loop, current element is deleted or not
	Desc        bool    // false;
	PrintHeader bool    // true;
	Remove      bool    // remove entry from screen
	Standalone  bool    // false;
	Verbose     bool    // false;
	Filter      []string
}

func (t *IsHistoryEntry) New(property *is.Property, del bool, remove bool, filter *[]string) *IsHistoryEntry {
	t.Processor = is.New(property) // Init the base struct
	t.MyType = "IsHistoryEntry"

	t.Filter = append(t.Filter, *filter...)
	t.Del = del
	t.Remove = remove
	t.override()

	return t
}

func (t *IsHistoryEntry) Run() *IsHistoryEntry {
	t.Processor.Run()
	return t
}

func (t *IsHistoryEntry) override() {
	t.V020_Elements = t.override_V020_Elements

	t.V030_ElementInfo = t.override_V030_ElementInfo

	t.V040_ElementMatch = t.override_V040_ElementMatch

	t.V050_ElementProcessMatched = t.override_V050_ElementProcessMatched

	t.V060_ElementProcessUnmatch = t.override_V060_ElementProcessUnmatch

	t.V080_ElementScrollable = t.override_V080_ElementScrollable

	t.V100_ScrollLoopEnd = t.override_V100_ScrollLoopEnd

}

func (t *IsHistoryEntry) V021_ElementsRemoveShorts(element *rod.Element) *IsHistoryEntry {
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

func (t *IsHistoryEntry) override_V020_Elements(element *rod.Element) *rod.Elements {
	prefix := t.MyType + ".V020_Elements"
	ezlog.Trace().N(prefix).TxtStart().Out()

	t.V021_ElementsRemoveShorts(element)
	var (
		tagNames = []string{"ytd-video-renderer", "yt-lockup-view-model"}
		elements rod.Elements
	)
	elements, t.Err = element.Elements(strings.Join(tagNames, ",")) // multiple tag names separate by comma
	if t.Err != nil {
		ezlog.Err().N(prefix).M(t.Err).Out()
	}
	ezlog.Debug().N(prefix).N("elements count").M(len(elements)).Out()
	ezlog.Trace().N(prefix).TxtEnd().Out()
	return &elements
}

func (t *IsHistoryEntry) override_V030_ElementInfo() (infoP is.IInfo) {
	prefix := t.MyType + ".V030_ElementInfo"
	ezlog.Trace().N(prefix).TxtStart().Out()

	if t.StateCurr.Element != nil {
		var (
			by                string
			info              YT_Info
			err               error
			elementMeta       *rod.Element
			elementsText      rod.Elements
			elementsTextCount int
		)
		by = "#video-title"                                // by ID
		elementMeta, err = t.StateCurr.Element.Element(by) // by id
		if err == nil {
			// -- trace
			if ezlog.GetLogLevel() == ezlog.TRACE {
				ezlog.Info().N(prefix).
					// Ln("element").Lm(t.StateCurr.Element).
					// Ln("html").Lm(gohtml.Format(t.StateCurr.Element.MustHTML())).
					Ln(by).Lm(gohtml.Format(elementMeta.MustHTML())).
					Out()
			}
			info.Title = strings.TrimSpace(elementMeta.MustText())
			if len(info.Title) != 0 {
				info.Url = *elementMeta.MustAttribute("href")
				info.Text = strings.TrimSpace(t.StateCurr.Element.MustElement("#description-text").MustText()) // by id

				meta := t.StateCurr.Element.MustElement("#metadata") // by id
				a := meta.MustElement("a")                           // by name
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
		} else {
			by = ".yt-lockup-metadata-view-model__text-container" //by class
			elementMeta, err = t.StateCurr.Element.Element(by)
			if err == nil {
				// -- trace
				// if ezlog.GetLogLevel() == ezlog.TRACE {
				// 	traceElement(prefix, by, elementMeta)
				// }
				// -- title
				info.Title = elementMeta.MustElement("a").MustText()
				info.Url = *elementMeta.MustElement("a").MustAttribute("href")
				by = "[role='text']"
				elementsText = elementMeta.MustElements(by)
				elementsTextCount = len(elementsText)
				ezlog.Info().N(prefix).N(by).N("elementsText len").M(elementsTextCount).Out()
				if elementsTextCount != 4 {
					ezlog.Warning().N(prefix).N(by).M("YT changed format").Out()
				}
				for i, e := range elementsText {
					ezlog.Info().N(prefix).N(by).N(i).M(e.MustText()).Out()
				}
				info.Title = elementsText[0].MustText()
				info.ChName = elementsText[1].MustText()
				switch elementsTextCount {
				case 3:
					// member video don't have views
					info.Text = elementsText[2].MustText()
				case 4:
					info.Text = elementsText[3].MustText()
				default:
					ezlog.Err().N(prefix).N(by).
						Ln("element").Lm(t.StateCurr.Element).
						Ln("html").Lm(gohtml.Format(t.StateCurr.Element.MustHTML())).
						Out()
				}
			} else {
				ezlog.Err().N(prefix).N("YT format not recognize").
					Ln("element").Lm(t.StateCurr.Element).
					Ln("html").Lm(gohtml.Format(t.StateCurr.Element.MustHTML())).
					Out()
			}
		}

		infoP = &info
		ezlog.Debug().N(prefix).N("info").M(infoP.String()).Out()
	}

	ezlog.Trace().N(prefix).TxtEnd().Out()
	return infoP
}

func (t *IsHistoryEntry) override_V040_ElementMatch() (matched bool, matchedStr string) {
	prefix := t.MyType + ".V040_ElementMatch"
	ezlog.Trace().N(prefix).TxtStart().Out()

	t.Deleted = false
	info := t.StateCurr.ElementInfo.(*YT_Info)
	chkStr := info.Title + " " + info.Text + " " + info.ChName + " " + info.ChUrlShort

	matched, matchedStr = str.ContainsAnySubStrings(&chkStr, &t.Filter, false)

	ezlog.Trace().N(prefix).TxtEnd().Out()
	return matched, matchedStr
}

func (t *IsHistoryEntry) override_V050_ElementProcessMatched() {
	prefix := t.MyType + ".V050_ElementProcessMatched"
	ezlog.Trace().N(prefix).TxtStart().Out()
	var (
		button         *rod.Element
		err            error
		menu           *rod.Element
		menuItems      rod.Elements
		menuItemsCount int
		tag            string
		txt            string
		loop           bool
		loopCount      int
	)
	if t.Del {
		// -- Old X-button -- START
		// tag = "#top-level-buttons-computed" // X button
		// button, err = t.StateCurr.Element.Element(tag)
		// -- Old X-button -- END

		ezlog.Trace().N(prefix).M("page stable - wait").Out()
		t.Page.MustWaitStable()
		ezlog.Trace().N(prefix).M("page stable").Out()
		tag = "button"                                 // X button gone, 3-dot button
		button, err = t.StateCurr.Element.Element(tag) // class, the 3-dot button
		if err == nil {
			button.MustVisible()
			ezlog.Trace().N(prefix).M("3-dot clicked - before").Out()
			button.MustClick() // click to bring up contextual menu
			ezlog.Trace().N(prefix).M("3-dot clicked - done").Out()
			tag = "tp-yt-paper-listbox,tp-yt-iron-dropdown"
			menu, err = t.Page.Element(tag) // tag
			if err == nil {
				tag = "ytd-menu-service-item-renderer,yt-list-item-view-model"
				loop = true
				for loop {
					loopCount++
					menuItems, err = menu.Elements(tag) // tag
					if err != nil {
						loop = false
					} else {
						menuItemsCount = len(menuItems)
						loop = menuItemsCount == 0
						if loop && loopCount == 100 {
							ezlog.Trace().N(prefix).N("loopCount").M(loopCount).N("menuItems").M(menuItemsCount).Out()
							os.Exit(1)
						}
						if !loop {
							ezlog.Trace().N(prefix).N("menuItems").M(menuItemsCount).Out()
						}
						for _, b := range menuItems {
							txt = b.MustText()
							// traceElement(prefix, txt, b)
							ezlog.Trace().N(prefix).N("menuItems").M(txt).Out()
							if strings.EqualFold(txt, "Remove from watch history") {
								elementRandomClick(t.Page, b, 250, 500)
								ezlog.Trace().N(prefix).N(txt).M("clicked").Out()
								t.Deleted = true
								ezlog.Debug().N(prefix).N("Deleted").M(t.Deleted).Out()
								break
							}
						}
					}
				}
			}
		}
	}
	if err != nil {
		ezlog.Err().N(prefix).M(err).Out()
	}
	ezlog.Trace().N(prefix).TxtEnd().Out()
}

func (t *IsHistoryEntry) override_V060_ElementProcessUnmatch() {
	prefix := t.MyType + ".V060_ElementProcessUnmatch"
	ezlog.Trace().N(prefix).M("Done").Out()
}

func (t *IsHistoryEntry) override_V080_ElementScrollable() (scrollable bool) {
	prefix := t.MyType + ".V080_ElementScrollable"
	ezlog.Trace().N(prefix).TxtStart().Out()

	info := t.StateCurr.ElementInfo.(*YT_Info)
	scrollable = !t.Deleted && (len(info.Title) == 0 || !t.StateCurr.Element.MustVisible())
	ezlog.Trace().N(prefix).TxtEnd().Out()
	return
}

func (t *IsHistoryEntry) override_V100_ScrollLoopEnd() {
	prefix := t.MyType + ".V100_ScrollLoopEnd"
	ezlog.Trace().N(prefix).TxtStart().Out()

	if t.Remove {
		if t.StateCurr.Elements != nil {
			for _, e := range *t.StateCurr.Elements {
				e.Remove()
			}
			t.StateCurr.ElementsCount = 0
		}
		es := t.Page.MustElements("ytd-item-section-renderer")
		if es != nil {
			t.StateCurr.Element = es.Last()
			if t.StateCurr.ScrollableElement != nil && t.StateCurr.Element != nil && (t.StateCurr.ScrollableElement.Object.ObjectID != t.StateCurr.Element.Object.ObjectID) {
				t.StateCurr.Scroll = true
			} else if t.StateCurr.ScrollableElement == nil && t.StateCurr.Element != nil {
				t.StateCurr.Scroll = true
			} else {
				t.StateCurr.Scroll = false
			}
		} else {
			t.StateCurr.Scroll = false
		}
		if t.StateCurr.ScrollableElement != nil {
			ezlog.Trace().N(prefix).N("t.StateCurr.ScrollableElement").M(string(t.StateCurr.ScrollableElement.Object.ObjectID)).Out()
		}
		if t.StateCurr.Element != nil {
			ezlog.Trace().N(prefix).N("t.StateCurr.Element").M(string(t.StateCurr.Element.Object.ObjectID)).Out()
		}
	}
	ezlog.Trace().N(prefix).TxtEnd().Out()
}

func traceElement(prefix, tag string, e *rod.Element) {
	ezlog.Trace().N(prefix).N(tag).Lm(gohtml.Format(e.MustHTML())).Out()
}

func elementRandomClick(page *rod.Page, e *rod.Element, min, max int64) {
	prefix := "elementRandomClick"
	ezlog.Trace().N(prefix).TxtStart().Out()
	var (
		interval int64
	)
	traceElement(prefix, "", e)
	ezlog.Trace().N(prefix).M("page stable - wait").Out()
	page.MustWaitStable()
	ezlog.Trace().N(prefix).M("page stable").Out()
	interval = min + rand.Int64N(max-min)
	time.Sleep(time.Duration(interval) * time.Millisecond / 2)
	ezlog.Trace().N(prefix).N("sleep(ms) - before").M(interval).Out()
	{
		var (
			x, y float64
		)
		box := e.MustShape().Box()
		ezlog.Trace().N(prefix).N("box").M(box).Out()
		x = box.X + 1 + rand.Float64()*(box.Width-2)
		y = box.Y + 1 + rand.Float64()*(box.Height-2)
		page.Mouse.MustMoveTo(x, y).MustClick(proto.InputMouseButtonLeft)
	}
	{
		// e.MustClick()
	}
	interval = min + rand.Int64N(max-min)
	time.Sleep(time.Duration(interval) * time.Millisecond / 2)
	ezlog.Trace().N(prefix).N("sleep(ms) - after").M(interval).Out()
	ezlog.Trace().N(prefix).TxtEnd().Out()
}

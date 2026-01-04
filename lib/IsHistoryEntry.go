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
	"errors"
	"math/rand/v2"
	"net/url"
	"strings"
	"time"

	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/state"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-is/v2/is"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
	"github.com/yosssi/gohtml"
)

// process YT history by entry
type IsHistoryEntry struct {
	is.Processor

	ClickSleep  float64 // in second
	Del         bool    // delete entry from history
	Deleted     bool    // In Run(), elements loop, current element is deleted or not
	Desc        bool    // false;
	PrintHeader bool    // true;
	Remove      bool    // remove entry from screen
	Standalone  bool    // false;
	Verbose     bool    // false;
	Filter      []string

	state state.State[V050_StateData]
}

func (t *IsHistoryEntry) New(property *is.Property, del bool, remove bool, filter *[]string) *IsHistoryEntry {
	t.Processor.New(property) // Init the base struct
	t.MyType = "IsHistoryEntry"

	t.Filter = append(t.Filter, *filter...)
	t.Del = del
	t.Remove = remove
	t.override()

	t.state = state.State[V050_StateData]{
		Data: V050_StateData{
			SleepMax: 300,
			SleepMin: 100,
		},
		Pre:   t.V051_FuncPre,
		OnErr: V051_OnErrFunc,
	}

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

func (t *IsHistoryEntry) override_V020_Elements(element *rod.Element) *rod.Elements {
	prefix := t.MyType + ".V020_Elements"
	ezlog.Debug().N(prefix).TxtStart().Out()

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
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return &elements
}

func (t *IsHistoryEntry) override_V030_ElementInfo() (infoP is.IInfo) {
	prefix := t.MyType + ".V030_ElementInfo"
	ezlog.Debug().N(prefix).TxtStart().Out()

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
				info.Url = YT_FullUrl(*elementMeta.MustAttribute("href"))
				info.Text = strings.TrimSpace(t.StateCurr.Element.MustElement("#description-text").MustText()) // by id

				meta := t.StateCurr.Element.MustElement("#metadata") // by id
				a := meta.MustElement("a")                           // by name
				info.ChName = strings.TrimSpace(a.MustText())

				if a != nil {
					chUrlP := a.MustAttribute("href")
					if chUrlP != nil {
						info.ChUrl = YT_FullUrl(*chUrlP)
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
				info.Url = YT_FullUrl(*elementMeta.MustElement("a").MustAttribute("href"))
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
				// info.Title = elementsText[0].MustText()
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

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return infoP
}

func (t *IsHistoryEntry) override_V040_ElementMatch() (matched bool, matchedStr string) {
	prefix := t.MyType + ".V040_ElementMatch"
	ezlog.Debug().N(prefix).TxtStart().Out()

	t.Deleted = false
	info := t.StateCurr.ElementInfo.(*YT_Info)
	chkStr := info.Title + " " + info.Text + " " + info.ChName + " " + info.ChUrlShort

	matched, matchedStr = str.ContainsAnySubStrings(&chkStr, &t.Filter, false)

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return matched, matchedStr
}

func (t *IsHistoryEntry) override_V050_ElementProcessMatched() {
	prefix := t.MyType + ".V050_ElementProcessMatched"
	ezlog.Debug().N(prefix).TxtStart().Out()
	if t.Del && t.StateCurr.Element.MustVisible() {
		// -- Old X-button -- START
		// tag = "#top-level-buttons-computed" // X button
		// button, err = t.StateCurr.Element.Element(tag)
		// -- Old X-button -- END
		t.state.Run(t.V0511_WaitStable, true)
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
}

func (t *IsHistoryEntry) override_V060_ElementProcessUnmatch() {
	prefix := t.MyType + ".V060_ElementProcessUnmatch"
	ezlog.Trace().N(prefix).M("Done").Out()
}

func (t *IsHistoryEntry) override_V080_ElementScrollable() (scrollable bool) {
	prefix := t.MyType + ".V080_ElementScrollable"
	ezlog.Debug().N(prefix).TxtStart().Out()

	info := t.StateCurr.ElementInfo.(*YT_Info)
	scrollable = !t.Deleted && (len(info.Title) == 0 || !t.StateCurr.Element.MustVisible())
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return
}

func (t *IsHistoryEntry) override_V100_ScrollLoopEnd() {
	prefix := t.MyType + ".V100_ScrollLoopEnd"
	ezlog.Debug().N(prefix).TxtStart().Out()

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
	ezlog.Debug().N(prefix).TxtEnd().Out()
}

func (t *IsHistoryEntry) V021_ElementsRemoveShorts(element *rod.Element) *IsHistoryEntry {
	prefix := t.MyType + ".V020_ElementsRemoveShorts"
	ezlog.Debug().N(prefix).TxtStart().Out()

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

	ezlog.Debug().N(prefix).TxtEnd().Out()
	return t
}

type V050_StateData struct {
	Element  *rod.Element
	SleepMax float64
	SleepMin float64
}

func (t *IsHistoryEntry) V051_FuncPre(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	prefix := t.MyType + ".V051_FuncPre"
	var (
		duration float64
	)
	state.Name = prefix
	ezlog.Debug().N(prefix).TxtStart().Out()
	duration = state.Data.SleepMin + rand.Float64()*(state.Data.SleepMax-state.Data.SleepMin)
	ezlog.Info().N(prefix).N("sleep(ms)").M(int64(duration)).Out()
	time.Sleep(time.Duration(duration) * time.Millisecond)
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return state
}

func V051_OnErrFunc(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	errs.Queue(state.Name, state.Err)
	return state
}

func (t *IsHistoryEntry) V0511_WaitStable(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	prefix := t.MyType + ".V0511"
	state.Name = prefix
	ezlog.Debug().N(prefix).TxtStart().Out()
	state.Err = t.Page.Keyboard.Press(input.Escape)
	if state.Err == nil {
		t.Page.Mouse.Scroll(0, 12, 3)
		// WaitPageStable(prefix, t.Page)
		state.Next = t.V0512_3DotClick
	} else {
		state.Next = t.V0511_WaitStable
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return state
}

// Click the 3-dot button to open popup menu
func (t *IsHistoryEntry) V0512_3DotClick(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	prefix := t.MyType + ".V0512"
	state.Name = prefix
	ezlog.Debug().N(prefix).TxtStart().Out()
	var (
		tags = []string{
			"button",
			".yt-lockup-metadata-view-model__menu-button"}
	)
	// select 3-dot
	for _, tag := range tags {
		state.Data.Element, state.Err = t.StateCurr.Element.Element(tag)
		if state.Err == nil {
			break
		}
	}
	if state.Err != nil {
		state.Next = nil
	}
	ezlog.Debug().N(prefix).N("state.Err").M(state.Err).Out()
	if state.Err == nil {
		// make 3-dot within screen
		// state.Data.Element.MustVisible()
		// click 3-dot
		state.Err = state.Data.Element.Click(proto.InputMouseButtonLeft, 1)
		ezlog.Debug().N(prefix).M("clicked").Out()
	}
	if state.Err == nil {
		state.Next = t.V0513_MenuSelect
	} else {
		state.Next = t.V0511_WaitStable
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return state
}

// Select the popup menu
func (t *IsHistoryEntry) V0513_MenuSelect(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	prefix := t.MyType + ".V0513"
	state.Name = prefix
	ezlog.Debug().N(prefix).TxtStart().Out()
	var (
		elements rod.Elements
		// tag      = "tp-yt-paper-listbox,tp-yt-iron-dropdown"
		tag = "#contentWrapper"
	)
	// select menu
	state.Data.Element = nil
	elements, state.Err = t.Page.Elements(tag) // tag
	if state.Err == nil {
		for _, element := range elements {
			if element.MustVisible() {
				state.Data.Element = element
				state.Next = t.V0514_MenuRead
				break
			}
		}
		if state.Data.Element == nil {
			state.Err = errors.New("cannot select menu")
		}
	}
	if state.Err != nil {
		state.Next = t.V0511_WaitStable
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return state
}

// Read menu items
func (t *IsHistoryEntry) V0514_MenuRead(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	prefix := t.MyType + ".V0514"
	state.Name = prefix
	ezlog.Debug().N(prefix).TxtStart().Out()
	// TraceElement(prefix, "", state.Element)
	var (
		matched         bool
		menuItems       rod.Elements
		menuItemText    string
		menuItemTextReq = "Remove from watch history"
		tag             = "ytd-menu-service-item-renderer,yt-list-item-view-model"
	)
	menuItems, state.Err = state.Data.Element.Elements(tag) // tag
	if state.Err == nil {
		if len(menuItems) > 0 {
			for _, item := range menuItems {
				menuItemText = strings.TrimSpace(item.MustText())
				ezlog.Trace().N(prefix).N("menuItems").M("'" + menuItemText + "'").Out()
				if strings.EqualFold(menuItemText, menuItemTextReq) {
					state.Data.Element = item
					state.Next = t.V0515_MenuClick
					matched = true
					break
				}
			}
			if !matched {
				TraceElement(prefix, "", state.Data.Element)
				state.Err = errors.New("unmatch: " + menuItemTextReq)
			}
		} else {
			state.Err = errors.New("0 menu item")
		}
	}
	if state.Err != nil {
		state.Next = t.V0511_WaitStable
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return state
}

// Click the item, use state.Element from V0513
func (t *IsHistoryEntry) V0515_MenuClick(state *state.State[V050_StateData]) *state.State[V050_StateData] {
	prefix := t.MyType + ".V0515"
	state.Name = prefix
	ezlog.Debug().N(prefix).TxtStart().Out()
	TraceElement(prefix, "", state.Data.Element)
	var (
		x, y  float64
		box   *proto.DOMRect
		shape *proto.DOMGetContentQuadsResult
	)
	{
		// -- just click
		// state.Err = state.Data.Element.Click(proto.InputMouseButtonLeft, 1)
	}
	{
		// -- random position click
		shape, state.Err = state.Data.Element.Shape()
		if state.Err == nil {
			if shape == nil {
				state.Err = errors.New("nil shape")
			} else {
				box = shape.Box()
				ezlog.Trace().N(prefix).N("box").M(box).Out()
				if box == nil {
					state.Err = errors.New("nil box")
				} else {
					x = box.X + 1 + rand.Float64()*(box.Width-2)
					y = box.Y + 1 + rand.Float64()*(box.Height-2)
					t.Page.Mouse.MustMoveTo(x, y).MustClick(proto.InputMouseButtonLeft)
					t.Deleted = true
				}
			}
		}
	}
	if state.Err == nil {
		state.Next = nil
	} else {
		// TraceElement(prefix, state.Err.Error(), state.Element)
		state.Next = t.V0511_WaitStable
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return state
}

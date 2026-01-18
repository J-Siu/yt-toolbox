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
	"errors"
	"math/rand/v2"
	"net/url"
	"strings"
	"time"

	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/state"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/J-Siu/go-is/v3/is"
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
		OnErr:         t.V051_OnErrFunc,
		OnErrContinue: true,
		Pre:           t.V051_FuncPre,
	}
	t.state.MyType = t.MyType + ".state"

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

func (t *IsHistoryEntry) override_V020_Elements() {
	prefix := t.MyType + ".V020_Elements"
	t.StateCurr.Name = prefix
	var (
		tagNames = []string{"ytd-video-renderer", "yt-lockup-view-model"}
	)
	ezlog.Debug().N(prefix).N("Container").M(t.Container).Out()
	t.V021_ElementsRemoveShorts(t.Container)
	t.StateCurr.Elements, t.Err = t.Container.Elements(strings.Join(tagNames, ",")) // multiple tag names separate by comma
	if t.Err != nil {
		ezlog.Err().N(prefix).M(t.Err).Out()
	}
	ezlog.Debug().N(prefix).N("elements count").M(len(t.StateCurr.Elements)).Out()
}

func (t *IsHistoryEntry) override_V030_ElementInfo() {
	prefix := t.MyType + ".V030_ElementInfo"
	t.StateCurr.Name = prefix
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
		t.StateCurr.ElementInfo = &info
		ezlog.Debug().N(prefix).N("info").M(info.String()).Out()
	}
}

func (t *IsHistoryEntry) override_V040_ElementMatch() {
	prefix := t.MyType + ".V040_ElementMatch"
	var (
		matched    bool
		matchedStr string
	)
	t.StateCurr.Name = prefix
	t.Deleted = false
	info := t.StateCurr.ElementInfo.(*YT_Info)
	chkStr := info.Title + " " + info.Text + " " + info.ChName + " " + info.ChUrlShort
	matched, matchedStr = str.ContainsAnySubStrings(&chkStr, &t.Filter, false)
	t.StateCurr.ElementInfo.SetMatched(matched)
	t.StateCurr.ElementInfo.SetMatchedStr(matchedStr)
}

func (t *IsHistoryEntry) override_V050_ElementProcessMatched() {
	prefix := t.MyType + ".V050_ElementProcessMatched"
	t.StateCurr.Name = prefix
	if t.Del && t.StateCurr.Element.MustVisible() {
		t.state.Run(t.V0511_WaitStable)
	}
}

func (t *IsHistoryEntry) override_V060_ElementProcessUnmatch() {
	prefix := t.MyType + ".V060_ElementProcessUnmatch"
	t.StateCurr.Name = prefix
}

func (t *IsHistoryEntry) override_V080_ElementScrollable() {
	prefix := t.MyType + ".V080_ElementScrollable"
	t.StateCurr.Name = prefix
	info := t.StateCurr.ElementInfo.(*YT_Info)
	t.StateCurr.ElementScrollable = !t.Deleted && (len(info.Title) == 0 || !t.StateCurr.Element.MustVisible())
}

func (t *IsHistoryEntry) override_V100_ScrollLoopEnd() {
	prefix := t.MyType + ".V100_ScrollLoopEnd"
	t.StateCurr.Name = prefix
	if t.Remove {
		if t.StateCurr.Elements != nil {
			for _, e := range t.StateCurr.Elements {
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
}

func (t *IsHistoryEntry) V021_ElementsRemoveShorts(element *rod.Element) *IsHistoryEntry {
	prefix := t.MyType + ".V020_ElementsRemoveShorts"
	t.StateCurr.Name = prefix
	if element != nil {
		var (
			tagName = "ytd-reel-shelf-renderer"
			es      rod.Elements
		)
		es, t.Err = element.Elements(tagName)
		if t.Err != nil {
			ezlog.Err().N(prefix).M(t.Err).Out()
		}
		es_count := len(es)
		for i := range es {
			es[es_count-i-1].Remove()
		}
	}
	return t
}

type V050_StateData struct {
	Element  *rod.Element
	SleepMax int64
	SleepMin int64
}

// sleep between min and max ms before each state
func (t *IsHistoryEntry) V051_FuncPre() *state.State[V050_StateData] {
	prefix := t.MyType + ".V051_FuncPre"
	t.state.Name = prefix
	var (
		duration int64
	)
	duration = t.state.Data.SleepMin + rand.Int64N(t.state.Data.SleepMax-t.state.Data.SleepMin)
	ezlog.Info().N(prefix).N("sleep(ms)").M(int64(duration)).Out()
	time.Sleep(time.Millisecond * time.Duration(duration))
	return &t.state
}

func (t *IsHistoryEntry) V051_OnErrFunc() *state.State[V050_StateData] {
	errs.Queue(t.state.Name, t.state.Err)
	return &t.state
}

func (t *IsHistoryEntry) V0511_WaitStable() *state.State[V050_StateData] {
	prefix := t.MyType + ".V0511"
	t.state.Name = prefix
	t.state.Err = t.Page.Keyboard.Press(input.Escape)
	if t.state.Err == nil {
		t.Page.Mouse.Scroll(0, 12, 3)
		// WaitPageStable(prefix, t.Page)
		t.state.Next = t.V0512_3DotClick
	} else {
		t.state.Next = t.V0511_WaitStable
	}
	return &t.state
}

// Click the 3-dot button to open popup menu
func (t *IsHistoryEntry) V0512_3DotClick() *state.State[V050_StateData] {
	prefix := t.MyType + ".V0512"
	t.state.Name = prefix
	var (
		tags = []string{
			"button",
			".yt-lockup-metadata-view-model__menu-button"}
	)
	// select 3-dot
	for _, tag := range tags {
		t.state.Data.Element, t.state.Err = t.StateCurr.Element.Element(tag)
		if t.state.Err == nil {
			break
		}
	}
	if t.state.Err != nil {
		t.state.Next = nil
	}
	ezlog.Debug().N(prefix).N("t.state.Err").M(t.state.Err).Out()
	if t.state.Err == nil {
		// make 3-dot within screen
		// t.state.Data.Element.MustVisible()
		// click 3-dot
		t.state.Err = t.state.Data.Element.Click(proto.InputMouseButtonLeft, 1)
		ezlog.Debug().N(prefix).M("clicked").Out()
	}
	if t.state.Err == nil {
		t.state.Next = t.V0513_MenuSelect
	} else {
		t.state.Next = t.V0511_WaitStable
	}
	return &t.state
}

// Select the popup menu
func (t *IsHistoryEntry) V0513_MenuSelect() *state.State[V050_StateData] {
	prefix := t.MyType + ".V0513"
	t.state.Name = prefix
	var (
		elements rod.Elements
		// tag      = "tp-yt-paper-listbox,tp-yt-iron-dropdown"
		tag = "#contentWrapper"
	)
	// select menu
	t.state.Data.Element = nil
	elements, t.state.Err = t.Page.Elements(tag) // tag
	if t.state.Err == nil {
		for _, element := range elements {
			if element.MustVisible() {
				t.state.Data.Element = element
				t.state.Next = t.V0514_MenuRead
				break
			}
		}
		if t.state.Data.Element == nil {
			t.state.Err = errors.New("cannot select menu")
		}
	}
	if t.state.Err != nil {
		t.state.Next = t.V0511_WaitStable
	}
	return &t.state
}

// Read menu items
func (t *IsHistoryEntry) V0514_MenuRead() *state.State[V050_StateData] {
	prefix := t.MyType + ".V0514"
	t.state.Name = prefix
	// TraceElement(prefix, "", t.state.Element)
	var (
		matched         bool
		menuItems       rod.Elements
		menuItemText    string
		menuItemTextReq = "Remove from watch history"
		tag             = "ytd-menu-service-item-renderer,yt-list-item-view-model"
	)
	menuItems, t.state.Err = t.state.Data.Element.Elements(tag) // tag
	if t.state.Err == nil {
		if len(menuItems) > 0 {
			for _, item := range menuItems {
				menuItemText = strings.TrimSpace(item.MustText())
				ezlog.Trace().N(prefix).N("menuItems").M("'" + menuItemText + "'").Out()
				if strings.EqualFold(menuItemText, menuItemTextReq) {
					t.state.Data.Element = item
					t.state.Next = t.V0515_MenuClick
					matched = true
					break
				}
			}
			if !matched {
				TraceElement(prefix, "", t.state.Data.Element)
				t.state.Err = errors.New("unmatch: " + menuItemTextReq)
			}
		} else {
			t.state.Err = errors.New("0 menu item")
		}
	}
	if t.state.Err != nil {
		t.state.Next = t.V0511_WaitStable
	}
	return &t.state
}

// Click the item, use t.state.Element from V0513
func (t *IsHistoryEntry) V0515_MenuClick() *state.State[V050_StateData] {
	prefix := t.MyType + ".V0515"
	t.state.Name = prefix
	TraceElement(prefix, "", t.state.Data.Element)
	var (
		x, y  float64
		box   *proto.DOMRect
		shape *proto.DOMGetContentQuadsResult
	)
	{
		// -- just click
		// t.state.Err = t.state.Data.Element.Click(proto.InputMouseButtonLeft, 1)
	}
	{
		// -- random position click
		shape, t.state.Err = t.state.Data.Element.Shape()
		if t.state.Err == nil {
			if shape == nil {
				t.state.Err = errors.New("nil shape")
			} else {
				box = shape.Box()
				ezlog.Trace().N(prefix).N("box").M(box).Out()
				if box == nil {
					t.state.Err = errors.New("nil box")
				} else {
					x = box.X + 1 + rand.Float64()*(box.Width-2)
					y = box.Y + 1 + rand.Float64()*(box.Height-2)
					t.Page.Mouse.MustMoveTo(x, y).MustClick(proto.InputMouseButtonLeft)
					t.Deleted = true
				}
			}
		}
	}
	if t.state.Err == nil {
		t.state.Next = nil
	} else {
		// TraceElement(prefix, t.state.Err.Error(), t.state.Element)
		t.state.Next = t.V0511_WaitStable
	}
	return &t.state
}

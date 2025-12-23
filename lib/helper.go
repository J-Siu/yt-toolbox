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

	"github.com/J-Siu/go-dtquery/dq"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/go-rod/rod"
	"github.com/yosssi/gohtml"
)

func GetTab(host string, port int) (page *rod.Page) {
	prefix := "GetTab"
	ezlog.Debug().N(prefix).TxtStart().Out()
	var (
		browser *rod.Browser
		err     error
		pages   rod.Pages
	)
	devtools := dq.Get(host, port)
	err = devtools.Err
	if err == nil {
		browser = rod.New().ControlURL(devtools.DT_Url)
		err = browser.Connect()
	}
	if err == nil {
		pages, err = browser.NoDefaultDevice().Pages()
	}
	if err == nil {
		page = pages.First()
		if page != nil {
			page.Activate()
		}
	}

	if err != nil {
		ezlog.Err().N(prefix).M(err).Out()
	}
	ezlog.Debug().N(prefix).TxtEnd().Out()
	return
}

// log at trace level, format element html
func TraceElement(prefix, tag string, e *rod.Element) {
	ezlog.Trace()
	if len(prefix) > 0 {
		ezlog.N(prefix)
	}
	if len(tag) > 0 {
		ezlog.N(tag)
	}
	ezlog.Lm(e)
	if e != nil {
		ezlog.Lm(gohtml.Format(e.MustHTML()))
	}
	ezlog.Out()
}

func UrlCleanup(urlIn string) (urlOut string, err error) {
	// - domain to lowercase
	// - uUnquote url
	err = nil
	urlOut = urlIn
	if len(urlOut) != 0 {
		var u *url.URL
		u, err = url.Parse(urlIn)
		if err == nil {
			urlOut = u.String()
		}
	}
	return
}

func UrlDecode(urlIn string) (urlOut string) {
	urlOut, err := url.QueryUnescape(urlIn)
	if err != nil {
		urlOut = urlIn
	}
	return
}

func WaitPageStable(prefix string, page *rod.Page) {
	prefix = prefix + ".WaitPageStable"
	ezlog.Debug().N(prefix).TxtStart().Out()
	page.MustWaitStable()
	ezlog.Debug().N(prefix).TxtEnd().Out()
}

// Add YT base if missing, unescape query path
func YT_FullUrl(urlIn string) (urlOut string) {
	var (
		err error
	)
	if !strings.HasPrefix(urlIn, YT_Base) {
		urlOut, err = url.JoinPath(YT_Base, urlIn)
		if err != nil {
			urlOut = urlIn
		}
	}
	return UrlDecode(urlOut)
}

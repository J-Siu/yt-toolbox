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

	"github.com/J-Siu/go-dtquery/dq"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/go-rod/rod"
)

func UrlDecode(urlIn string) (urlOut string) {
	urlOut, err := url.QueryUnescape(urlIn)
	if err != nil {
		urlOut = urlIn
	}
	return urlOut
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
	return urlOut, err
}

func GetTab(host string, port int) *rod.Page {
	prefix := "GetTab"
	ezlog.Trace().N(prefix).TxtStart().Out()
	var (
		browser *rod.Browser
		err     error
		page    *rod.Page
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
	ezlog.Trace().N(prefix).TxtEnd().Out()
	return page
}

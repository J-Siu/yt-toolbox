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
	is "github.com/J-Siu/go-is/v3/is"
	"github.com/J-Siu/yt-toolbox/v2/global"
)

// Embed [is.InfoBase] for [is.IInfo] interface
type YT_Info struct {
	is.InfoBase

	// --- Channel info
	ChName     string `json:"ChName,omitempty"`
	ChUrl      string `json:"ChUrl,omitempty"`
	ChUrlShort string `json:"ChUrlShort,omitempty"`

	// --- Video info
	Text   string   `json:"Text,omitempty"`
	Title  string   `json:"Title,omitempty"`
	Titles []string `json:"Titles,omitempty"`
	Url    string   `json:"Url,omitempty"`
}

func (t *YT_Info) String() string {
	if global.Flag.Desc {
		return "[" + t.Title + "](" + UrlDecode(t.Url) + ") | [" + t.ChName + "](" + t.ChUrl + ") | " + t.Text
	} else {
		return "[" + t.Title + "](" + UrlDecode(t.Url) + ") | [" + t.ChName + "](" + t.ChUrl + ")"
	}
}

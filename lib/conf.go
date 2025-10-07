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
	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/spf13/viper"
)

var ConfDefault = TypeConf{
	FileConf: "$HOME/.config/yt-toolbox.json",

	DevtoolsHost: "localhost",
	DevtoolsPort: 9222,
}

type TypeConf struct {
	basestruct.Base

	FileConf string `json:"FileConf"`

	HistoryFilter []string `json:"HistoryFilter"`

	DevtoolsHost string `json:"DevtoolsHost"`
	DevtoolsPort int    `json:"DevtoolsPort"`
}

func (t *TypeConf) New() {
	t.Initialized = true
	t.MyType = "TypeConf"
	prefix := t.MyType + ".New"

	t.setDefault()
	ezlog.Debug().N(prefix).Nn("Default").M(t).Out()

	t.readFileConf()
	ezlog.Debug().N(prefix).Nn("Raw").M(t).Out()

	t.expand()
	ezlog.Debug().N(prefix).Nn("Expand").M(t).Out()
}

func (t *TypeConf) readFileConf() {
	prefix := t.MyType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(file.TildeEnvExpand(t.FileConf))
	viper.AutomaticEnv()
	t.Err = viper.ReadInConfig()

	if t.Err == nil {
		t.Err = viper.Unmarshal(&t)
	} else {
		ezlog.Debug().N(prefix).M(t.Err).Out()
	}
}

// Should be called before reading config file
func (t *TypeConf) setDefault() {
	if t.FileConf == "" {
		t.FileConf = ConfDefault.FileConf
	}
	t.DevtoolsHost = ConfDefault.DevtoolsHost
	t.DevtoolsPort = ConfDefault.DevtoolsPort
}

func (t *TypeConf) expand() {
	t.FileConf = file.TildeEnvExpand(t.FileConf)
}

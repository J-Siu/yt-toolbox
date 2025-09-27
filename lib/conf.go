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
	"os"
	"regexp"

	"github.com/J-Siu/go-ezlog/v2"
	"github.com/spf13/viper"
)

var ConfDefault = TypeConf{
	FileConf: "$HOME/.config/yt-toolbox.json",

	DevtoolsHost: "localhost",
	DevtoolsPort: 9222,
}

type TypeConf struct {
	Err    error
	init   bool
	myType string

	FileConf string `json:"FileConf"`

	HistoryFilter []string `json:"HistoryFilter"`

	DevtoolsHost string `json:"DevtoolsHost"`
	DevtoolsPort int    `json:"DevtoolsPort"`
}

func (c *TypeConf) New() {
	c.init = true
	c.myType = "TypeConf"
	prefix := c.myType + ".New"

	c.setDefault()
	ezlog.Debug().Name(prefix).NameLn("Default").Msg(c).Out()

	c.readFileConf()
	ezlog.Debug().Name(prefix).NameLn("Raw").Msg(c).Out()

	// TODO: add flag

	c.expand()
	ezlog.Debug().Name(prefix).NameLn("Expand").Msg(c).Out()
}

func (c *TypeConf) readFileConf() {
	prefix := c.myType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(TildeEnvExpand(c.FileConf))
	viper.AutomaticEnv()
	c.Err = viper.ReadInConfig()

	if c.Err == nil {
		c.Err = viper.Unmarshal(&c)
	} else {
		ezlog.Debug().Name(prefix).Name("Config file").Msg(c.Err).Out()
	}
}

// Should be called before reading config file
func (c *TypeConf) setDefault() {
	if c.FileConf == "" {
		c.FileConf = ConfDefault.FileConf
	}
	c.DevtoolsHost = ConfDefault.DevtoolsHost
	c.DevtoolsPort = ConfDefault.DevtoolsPort
}

func (c *TypeConf) expand() {
	c.FileConf = TildeEnvExpand(c.FileConf)
}

// Expand Linux `~` and environment variable in string
func TildeEnvExpand(strIn string) (strOut string) {
	if strIn == "~" {
		strOut = "$HOME"
	} else {
		re := regexp.MustCompile(`^~/`)
		strOut = re.ReplaceAllString(strIn, "$$HOME/")
	}
	return os.ExpandEnv(strOut)
}

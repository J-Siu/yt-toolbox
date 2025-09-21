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
package cmd

import (
	"os"

	"github.com/J-Siu/go-ezlog"
	"github.com/J-Siu/yt-toolbox/global"
	"github.com/J-Siu/yt-toolbox/lib"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-yt-toolbox",
	Aliases: []string{"gyt"},
	Short:   "YouTube toolbox",
	Version: global.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DebugLevel)
		}
		if global.Flag.Trace {
			ezlog.SetLogLevel(ezlog.TraceLevel)
		}
		ezlog.Debug("Version: " + global.Version)
		ezlog.Debug("Flag:")
		ezlog.DebugP(lib.MustToJsonStrP(&global.Flag))
		global.Conf.New()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cmd := rootCmd
	cmd.PersistentFlags().BoolVarP(&global.Flag.Debug, "debug", "", false, "Enable debug")
	cmd.PersistentFlags().BoolVarP(&global.Flag.Trace, "trace", "", false, "Enable trace (include debug)")
	cmd.PersistentFlags().BoolVarP(&global.Flag.Verbose, "verbose", "v", false, "Verbose")
	cmd.PersistentFlags().StringVarP(&global.Conf.FileConf, "config", "c", lib.ConfDefault.FileConf, "Config file")

	cmd.PersistentFlags().IntVarP(&global.Flag.ScrollMax, "scroll-max", "s", 0, "Unlimited -1")

	cmd.PersistentFlags().String("DevToolsHost", "", "Devtools Host")
	cmd.PersistentFlags().String("DevToolsUrl", "", "Devtools Url (override Devtools host, ver and port)")
	cmd.PersistentFlags().String("DevToolsVer", "", "Devtools Version")
	cmd.PersistentFlags().Int("DevToolsPort", 9222, "Devtools Port")
}

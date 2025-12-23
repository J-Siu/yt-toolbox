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
	"github.com/J-Siu/yt-toolbox/v2/global"
	"github.com/J-Siu/yt-toolbox/v2/lib"
	"github.com/spf13/cobra"
)

// historyCmd represents the history command
var historyCmd = &cobra.Command{
	Use:     "history",
	Aliases: []string{"h", "hist"},
	Short:   "Get Youtube History",
	Run: func(cmd *cobra.Command, args []string) {
		page := lib.GetTab(global.Conf.DevtoolsHost, global.Conf.DevtoolsPort)

		isHistorySection := new(lib.IsHistorySection).
			New(
				page,
				lib.YT_History,
				!global.FlagHistory.NoRemove,
				global.Flag.ScrollMax,
				global.Flag.Verbose)
		isHistorySection.Del = global.FlagHistory.Del
		isHistorySection.Filter = append(isHistorySection.Filter, global.Conf.HistoryFilter...)
		isHistorySection.
			Run()
	},
}

func init() {
	cmd := historyCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().BoolVarP(&global.FlagHistory.Del, "del", "", false, "Perform actual deletion. [default: Dry run]")
	cmd.PersistentFlags().BoolVarP(&global.FlagHistory.Desc, "desc", "", false, "Show description")
	cmd.PersistentFlags().BoolVarP(&global.FlagHistory.NoRemove, "no-remove", "n", false, "No removal of screen element. (Not history deletion!) [default: Remove screen element.]")
}

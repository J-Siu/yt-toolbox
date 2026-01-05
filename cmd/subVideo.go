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

package cmd

import (
	"github.com/J-Siu/go-is/v3/is"
	"github.com/J-Siu/yt-toolbox/v2/global"
	"github.com/J-Siu/yt-toolbox/v2/lib"
	"github.com/spf13/cobra"
)

// videosCmd represents the videos command
var subVideoCmd = &cobra.Command{
	Use:     "video",
	Aliases: []string{"v", "videos"},
	Short:   "Get YT Subscription Videos",
	Run: func(cmd *cobra.Command, args []string) {
		page := lib.GetTab(global.Conf.DevtoolsHost, global.Conf.DevtoolsPort)

		isSubVideo := new(lib.IsSubVideo).
			New(
				page,
				lib.YT_SubVideos,
				global.Flag.ScrollMax,
				global.FlagSub.Day,
			).
			Run()
		if isSubVideo.Err == nil {
			isSubVideo.IInfoList.Print(is.PrintAll)
		}
	},
}

func init() {
	cmd := subVideoCmd
	subscriptionsCmd.AddCommand(cmd)

	cmd.Flags().UintVarP(&global.FlagSub.Day, "day", "", 0, "number of days (override scroll)")
}

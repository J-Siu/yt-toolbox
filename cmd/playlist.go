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
	"github.com/J-Siu/go-ezlog/v2"
	"github.com/J-Siu/go-is"
	"github.com/J-Siu/yt-toolbox/global"
	"github.com/J-Siu/yt-toolbox/lib"
	"github.com/go-rod/rod"
	"github.com/spf13/cobra"
)

// playlistCmd represents the playlist command
var playlistCmd = &cobra.Command{
	Use:     "playlist",
	Aliases: []string{"p", "pl"},
	Short:   "Get Youtube Playlist",
	Run: func(cmd *cobra.Command, args []string) {
		page := lib.GetTab(global.Conf.DevtoolsHost, global.Conf.DevtoolsPort)
		isPlaylist := new(lib.IsPlaylist)
		isPlaylist.
			New(
				page,
				lib.UrlYT.Playlists,
				global.Flag.ScrollMax,
				&global.FlagPlaylist.Exclude,
				&global.FlagPlaylist.Include).
			Run()

		if isPlaylist.Err == nil {
			ezlog.Log().Name("Playlist").Out()
			isPlaylist.IInfoList.Print(is.PrintMatched)
			if global.FlagPlaylist.GetList {
				for _, info := range *isPlaylist.IInfoList {
					if info.Matched() {
						processVideoList(info, page)
					}
				}
			}
		}
	},
}

func init() {
	cmd := playlistCmd
	rootCmd.AddCommand(cmd)

	cmd.PersistentFlags().BoolVarP(&global.FlagPlaylist.GetList, "get-list", "g", false, "Get individual playlist")
	cmd.PersistentFlags().StringArrayVarP(&global.FlagPlaylist.Exclude, "exclude", "e", []string{}, "Exclude play list containing string (Override Include)")
	cmd.PersistentFlags().StringArrayVarP(&global.FlagPlaylist.Include, "include", "i", []string{}, "Include play list containing string")
}

func processVideoList(iinfo is.IInfo, page *rod.Page) {
	info := iinfo.(*lib.YT_Info)
	var isVideoList lib.IsPlaylistVideo
	isVideoList.
		New(
			page,
			info.Url,
			global.Flag.ScrollMax).
		Run()
	ezlog.Log().Name(info.Title).Out()
	isVideoList.IInfoList.Print(is.PrintAll)
}

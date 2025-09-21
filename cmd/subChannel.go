/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/J-Siu/go-is"
	"github.com/J-Siu/yt-toolbox/global"
	"github.com/J-Siu/yt-toolbox/lib"
	"github.com/spf13/cobra"
)

// subChannelCmd represents the channel command
var subChannelCmd = &cobra.Command{
	Use:     "channel",
	Aliases: []string{"c", "ch"},
	Short:   "Get YT Subscription Channels",
	Run: func(cmd *cobra.Command, args []string) {
		page := lib.GetTab(global.Conf.DevtoolsHost, global.Conf.DevtoolsPort)

		isSubCh := new(lib.IsSubChannel)
		isSubCh.
			New(
				page,
				lib.UrlYT.SubChannels,
				global.Flag.ScrollMax).
			Run()
		if isSubCh.Err == nil {
			isSubCh.IInfoList.Print(is.PrintAll)
		}
	},
}

func init() {
	cmd := subChannelCmd
	subscriptionsCmd.AddCommand(cmd)
}

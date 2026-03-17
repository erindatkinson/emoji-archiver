// /*
// Copyright © 2025 Erin Atkinson
// */
package cmd

import (
	"log/slog"
	"os"
	"slices"

	"github.com/erindatkinson/slack-emojinator/internal/cache"
	"github.com/erindatkinson/slack-emojinator/internal/slack"
	"github.com/erindatkinson/slack-emojinator/internal/utilities"
	"github.com/gammazero/workerpool"
	"github.com/spf13/cobra"
)

var outputDir string
var concurrency int

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Pull all emoji from a given slack team",
	Run: func(cmd *cobra.Command, args []string) {
		if browser == "" || profile == "" || subdomain == "" {
			slog.Error("error reading configs from env, config, or flags")
			return
		}

		logger := utilities.NewLogger(
			cmd.Flag("log-level").Value.String(),
			"team", subdomain, "dir", outputDir)

		logger.Info("creating export directory")
		os.MkdirAll(outputDir, 0755)
		client, err := slack.NewSlackClient(cmd.Context(), browser, profile, subdomain)
		if err != nil {
			logger.Error("unable to create slack client", "error", err)
			return
		}
		logger.Debug("client setup complete")

		logger.Info("retrieving list of current emoji")
		currentEmoji, err := client.ListEmoji()
		if err != nil {
			logger.Error("error retrieving current emoji list", "error", err)
			return
		}
		cached, err := cache.ListDownloadedEmojis(outputDir)
		if err != nil {
			logger.Error("unable to get cached emojis", "error", err)
			return
		}

		wp := workerpool.New(concurrency)

		for _, emoji := range currentEmoji {
			request := emoji
			wp.Submit(func() {
				loopLog := logger.With("name", request.Name)
				if slices.ContainsFunc(cached, func(e cache.EmojiItem) bool {
					return e.Name == request.Name
				}) {
					loopLog.Debug("already downloaded, skipping")
					return
				}

				loopLog.Debug("exporting emoji")
				if err := client.ExportEmoji(request, outputDir); err != nil {
					loopLog.Error("error exporting", "error", err)
				}

			})
		}

		wp.StopWait()
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringVarP(&outputDir, "directory", "d", "./export/", "the directory to use to export")
	exportCmd.Flags().IntVar(&concurrency, "concurrency", 1, "concurrency to use to download")
	exportCmd.Flags().String("log-level", "info", "enable debug logging")
}

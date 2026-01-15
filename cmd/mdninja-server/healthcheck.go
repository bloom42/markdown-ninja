package main

import (
	"fmt"
	"io"
	"os"

	"github.com/skerkour/stdx-go/cobra"
	"github.com/skerkour/stdx-go/httpx"
	"github.com/skerkour/stdx-go/log/slogx"
	"markdown.ninja/cmd/mdninja-server/config"
)

var healthcheckCmd = &cobra.Command{
	Use:   "healthcheck",
	Short: "Performs healthchecks for Docker. exit(1) if mdninja-server is not answering HTTP requests.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		ctx := cmd.Context()

		logger, logLevel, lokiWriter := newLogger(ctx)
		// inject the logger into ctx
		ctx = slogx.ToCtx(ctx, logger)

		conf, err := config.Load(ctx, flagServerConfigPath)
		if err != nil {
			return
		}
		logLevel.Set(conf.Logs.Level)
		if conf.Logs.LokiEndpoint != nil {
			lokiWriter.SetEndpoint(*conf.Logs.LokiEndpoint)
		}

		httpClient := httpx.DefaultClient()

		res, err := httpClient.Get(fmt.Sprintf("http://localhost:%d/__markdown_ninja/healthcheck", conf.HTTP.Port))
		if err != nil {
			logger.Error(fmt.Sprintf("healthcheck: %s", err.Error()))
			os.Exit(1)
		}
		io.Copy(io.Discard, res.Body)
		res.Body.Close()

		if res.StatusCode != 200 {
			logger.Error(fmt.Sprintf("healthcheck: received HTTP status code != 200: %s", res.Status))
			os.Exit(1)
		}

		return
	},
}

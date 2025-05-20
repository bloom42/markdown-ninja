package main

import (
	"github.com/bloom42/stdx-go/cobra"
	"markdown.ninja/cmd/mdninja-ebook/ebook"
)

var coverCmd = &cobra.Command{
	Use:   "cover source.png dest.pdf",
	Short: "Convert a PNG or JPG cover to PDF",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if len(args) != 2 {
			return cmd.Help()
		}

		source := args[0]
		dest := args[1]

		err = ebook.ConvertCoverToPdf(cmd.Context(), dest, source)
		return err
	},
}

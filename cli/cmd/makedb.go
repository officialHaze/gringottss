/*
Copyright © 2026 Gringottss CLI <moinak.dey8@gmail.com>
*/
package cmd

import (
	"github.com/officialHaze/gringottss/cli/helpers"
	"github.com/spf13/cobra"
)

// buildDbCmd represents the buildDb command
var makedbCmd = &cobra.Command{
	Use:   "makedb",
	Short: "Call Gringottss Engine to build the DB once.",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.BuildDB(rootCmd.Context())
	},
}

func init() {
	rootCmd.AddCommand(makedbCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildDbCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildDbCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

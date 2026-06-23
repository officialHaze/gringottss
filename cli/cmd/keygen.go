/*
Copyleft © 2026 Gringottss CLI <moinak.dey8@gmail.com>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/officialHaze/gringottss/cli/helpers"
	"github.com/officialHaze/gringottss/cli/logger"
	"github.com/spf13/cobra"
)

var (
	keytype string
)

// encryptionkeygenCmd represents the encryptionkeygen command
var keygenCmd = &cobra.Command{
	Use:   "keygen",
	Short: "Generate keys used by gringottss engine. Engine restart required.",
	// 	Long: `A longer description that spans multiple lines and likely contains examples
	// and usage of using your command. For example:

	// Cobra is a CLI library for Go that empowers applications.
	// This application is a tool to generate the needed files
	// to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		switch strings.ToLower(strings.TrimSpace(keytype)) {
		case "encryption":
			helpers.GenerateEncryptionKeys(rootCmd.Context())
		default:
			logger.ERROR().Println("Unsupported keygen requested!")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(keygenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptionkeygenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	keygenCmd.Flags().StringVarP(&keytype, "type", "t", "", "Mention type of keys to generate. Eg:- encryption")

	// Mark required flags here
	keygenCmd.MarkFlagRequired("type")
}

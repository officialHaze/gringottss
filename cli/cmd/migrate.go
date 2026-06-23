/*
Copyleft © 2026 Gringottss CLI <moinak.dey8@gmail.com>
*/
package cmd

import (
	"github.com/officialHaze/gringottss/cli/helpers"
	"github.com/spf13/cobra"
)

var (
	dbname string
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate old Gringottss DB into current.",
	Long: `
	Migrate all necessary tables from an old Gringottss DB (might be from another release)
	to the DB in current release. Make sure the old DB is present inside <PROJECT_ROOT>/engine/migrate
	Eg:- <PROJECT_ROOT>/engine/migrate/old.db
	`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.MigrateOldDB(rootCmd.Context(), dbname)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	migrateCmd.Flags().StringVar(&dbname, "dbname", "", "Old DB name (Required)")

	// Mark required flags here
	migrateCmd.MarkFlagRequired("dbname")
}

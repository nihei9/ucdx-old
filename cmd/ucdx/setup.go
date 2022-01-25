package main

import (
	"os"
	"path/filepath"

	"github.com/nihei9/ucdx/db"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Set up the database ucdx refereces",
		Long:  `setup downloads the UCD's data files and parses them. The parsed data files are saved to the ${HOME}/.ucdx directory in JSON format.`,
		Args:  cobra.NoArgs,
		RunE:  runSetup,
	}
	rootCmd.AddCommand(cmd)
}

func runSetup(cmd *cobra.Command, args []string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	appDirPath := filepath.Join(homeDir, ".ucdx")
	err = os.Mkdir(appDirPath, 0744)
	if err != nil && !os.IsExist(err) {
		return err
	}
	return db.MakeDB(&db.DBConfig{
		AppDirPath: appDirPath,
	})
}

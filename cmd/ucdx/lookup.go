package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nihei9/ucdx/db"
	"github.com/nihei9/ucdx/ucd"
	"github.com/spf13/cobra"
)

var lookupOutputSet = []string{
	"table",
	"json",
}

type lookupFlagSet struct {
	output *string
}

func (f *lookupFlagSet) validate() error {
	passed := false
	for _, o := range lookupOutputSet {
		if *f.output == o {
			passed = true
			break
		}
	}
	if !passed {
		var b strings.Builder
		fmt.Fprint(&b, lookupOutputSet[0])
		for _, o := range lookupOutputSet[1:] {
			fmt.Fprint(&b, ", ", o)
		}
		return fmt.Errorf("--output doesn't support %v, allowed values are: %v", *f.output, b.String())
	}

	return nil
}

var lookupFlags = &lookupFlagSet{}

func init() {
	cmd := &cobra.Command{
		Use:   "lookup",
		Short: "Look up the properties of a code point",
		Long:  `lookup looks up the properties of a code point specified in hexadecimal notation.`,
		Example: `  ucdx lookup 1F63A`,
		Args:  cobra.ExactArgs(1),
		RunE:  runLookup,
	}
	lookupFlags.output = cmd.Flags().StringP("output", "o", "table", "Output format. One of: json|table")
	rootCmd.AddCommand(cmd)
}

func runLookup(cmd *cobra.Command, args []string) error {
	err := lookupFlags.validate()
	if err != nil {
		return err
	}

	var u *ucd.UCD
	{
		homeDirPath, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		appDirPath := filepath.Join(homeDirPath, ".ucdx")

		u, err = db.OpenDB(appDirPath)
		if err != nil {
			return err
		}
	}

	n, err := strconv.ParseInt(args[0], 16, 32)
	if err != nil {
		return fmt.Errorf("invalid code point: %v", err)
	}
	c := rune(n)
	if c < '\u0000' || c > '\U0010FFFF' {
		return fmt.Errorf("%X is an invalid code point. A code point must be in the range of U+0000 to U+10FFFF.", c)
	}

	result := u.AnalizeCodePoint(c)

	switch *lookupFlags.output {
	case "table":
		printPropertySetAsTable([]*ucd.PropertySet{result})
	case "json":
		b, err := json.Marshal(result)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}

	return nil
}

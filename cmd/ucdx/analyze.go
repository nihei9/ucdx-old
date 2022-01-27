package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/nihei9/ucdx/db"
	"github.com/nihei9/ucdx/ucd"
	"github.com/spf13/cobra"
)

func init() {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze characters and print their properties",
		Long:  `analyze analyzes characters and print their properties.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAnalyze,
	}
	rootCmd.AddCommand(cmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
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

	var src io.Reader
	if len(args) > 0 {
		src = strings.NewReader(args[0])
	} else {
		src = os.Stdin
	}
	r := bufio.NewReader(src)
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if c == unicode.ReplacementChar {
			continue
		}

		props := u.AnalizeCodePoint(c)

		fmt.Println(string(c), fmt.Sprintf("U+%X", c))
		var opts []string
		if len(props.GeneralCategoryGroups) > 0 {
			var gs strings.Builder
			fmt.Fprint(&gs, props.GeneralCategoryGroups[0])
			for _, g := range props.GeneralCategoryGroups[1:] {
				fmt.Fprintf(&gs, ", %v", g)
			}
			opts = []string{
				fmt.Sprintf("(%v)", gs.String()),
			}
		}
		printProperty(props.Lookup(ucd.PropNameGeneralCategory), opts...)
		printProperty(props.Lookup(ucd.PropNameAlphabetic))
		printProperty(props.Lookup(ucd.PropNameLowercase))
		printProperty(props.Lookup(ucd.PropNameUppercase))
		printProperty(props.Lookup(ucd.PropNameOtherAlphabetic))
		printProperty(props.Lookup(ucd.PropNameOtherLowercase))
		printProperty(props.Lookup(ucd.PropNameOtherUppercase))
		printProperty(props.Lookup(ucd.PropNameWhiteSpace))
	}
}

func printProperty(prop *ucd.Property, opts ...string) {
	fmt.Printf("%-20v: %v", prop.Name, prop.Value)
	for _, opt := range opts {
		fmt.Printf(" %v", opt)
	}
	fmt.Print("\n")
}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/nihei9/ucdx/db"
	"github.com/nihei9/ucdx/ucd"
	"github.com/nihei9/ucdx/ucd/property"
	"github.com/spf13/cobra"
)

var analyzeOutputSet = []string{
	"table",
	"json",
}

type analyzeFlagSet struct {
	output *string
}

func (f *analyzeFlagSet) validate() error {
	passed := false
	for _, o := range analyzeOutputSet {
		if *f.output == o {
			passed = true
			break
		}
	}
	if !passed {
		var b strings.Builder
		fmt.Fprint(&b, analyzeOutputSet[0])
		for _, o := range analyzeOutputSet[1:] {
			fmt.Fprint(&b, ", ", o)
		}
		return fmt.Errorf("--output doesn't support %v, allowed values are: %v", *f.output, b.String())
	}

	return nil
}

var analyzeFlags = &analyzeFlagSet{}

func init() {
	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "Analyze characters and print their properties",
		Long:  `analyze analyzes characters and print their properties.`,
		Args:  cobra.MaximumNArgs(1),
		RunE:  runAnalyze,
	}
	analyzeFlags.output = cmd.Flags().StringP("output", "o", "table", "Output format. One of: json|table")
	rootCmd.AddCommand(cmd)
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	err := analyzeFlags.validate()
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

	var src io.Reader
	if len(args) > 0 {
		src = strings.NewReader(args[0])
	} else {
		src = os.Stdin
	}
	r := bufio.NewReader(src)
	results := []*ucd.PropertySet{}
	for {
		c, _, err := r.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if c == unicode.ReplacementChar {
			continue
		}

		props := u.AnalizeCodePoint(c)
		results = append(results, props)
	}

	switch *analyzeFlags.output {
	case "table":
		printPropertySetAsTable(results)
	case "json":
		b, err := json.Marshal(results)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
	}

	return nil
}

func printPropertySetAsTable(ps []*ucd.PropertySet) {
	for _, p := range ps {
		fmt.Println(string(p.CP), fmt.Sprintf("U+%X", p.CP))
		var opts []string
		if len(p.GeneralCategoryGroups) > 0 {
			var gs strings.Builder
			fmt.Fprint(&gs, p.GeneralCategoryGroups[0])
			for _, g := range p.GeneralCategoryGroups[1:] {
				fmt.Fprintf(&gs, ", %v", g)
			}
			opts = []string{
				fmt.Sprintf("(%v)", gs.String()),
			}
		}
		printProperty(p.Lookup(property.PropNameName))
		printProperty(p.Lookup(property.PropNameNameAlias))
		printProperty(p.Lookup(property.PropNameGeneralCategory), opts...)
		printProperty(p.Lookup(property.PropNameAlphabetic))
		printProperty(p.Lookup(property.PropNameLowercase))
		printProperty(p.Lookup(property.PropNameUppercase))
		printProperty(p.Lookup(property.PropNameIDStart))
		printProperty(p.Lookup(property.PropNameIDContinue))
		printProperty(p.Lookup(property.PropNameXIDStart))
		printProperty(p.Lookup(property.PropNameXIDContinue))
		printProperty(p.Lookup(property.PropNameWhiteSpace))
	}
}

func printProperty(prop *property.Property, opts ...string) {
	fmt.Printf("%-20v: %v", prop.Name, prop.Value)
	for _, opt := range opts {
		fmt.Printf(" %v", opt)
	}
	fmt.Print("\n")
}

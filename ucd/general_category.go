package ucd

import "sort"

// See section 5.7.1 General Category Values in [UAX44].
var generalCategoryGroups = map[string][]string{
	// Cased_Letter
	"lc": {"lu", "ll", "lt"},
	// Letter
	"l": {"lu", "ll", "lt", "lm", "lo"},
	// Mark
	"m": {"mm", "mc", "me"},
	// Number
	"n": {"nd", "nl", "no"},
	// Punctuation
	"p": {"pc", "pd", "ps", "pi", "pe", "pf", "po"},
	// Symbol
	"s": {"sm", "sc", "sk", "so"},
	// Separator
	"z": {"zs", "zl", "zp"},
	// Other
	"c": {"cc", "cf", "cs", "co", "cn"},
}

func lookupGCGroups(gc string) []string {
	// A General_Category may belong to one or more groups.
	var groups []string
	for group, gcs := range generalCategoryGroups {
		for _, g := range gcs {
			if g == gc {
				groups = append(groups, group)
				break
			}
		}
	}
	if len(groups) > 0 {
		sort.Slice(groups, func(i, j int) bool {
			return groups[i] < groups[j]
		})
	}

	return groups
}

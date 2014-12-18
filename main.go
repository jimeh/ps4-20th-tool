package main

import (
	"fmt"
	"os"

	"github.com/jimeh/ps4-20th-tool/auto"
	"github.com/jimeh/ps4-20th-tool/brute"
	"github.com/jimeh/ps4-20th-tool/find"
)

func help() {
	fmt.Println(`usage: ps4-20th-tool <command> [<args>]

Commands:
   find   Lookup the SP (redirect code) and the secret URL.
   brute  Attempt to a brute force attack against the redirect page, trying
          every possible combination of 2 and 3 characters.
   auto   Checks for new secret URL and clues being tweeted every 0.5 seconds,
          when new secret URL is up and clues have been posted, the web form
          is automatically submitted. Requires path to config JSON file as
          argument, example: ps4-20th-tool auto auto-config.json
`)
}

func main() {
	cmd := ""
	subCmd := ""

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if len(os.Args) > 2 {
		subCmd = os.Args[2]
	}

	switch cmd {
	case "find":
		switch subCmd {
		case "source":
			find.Source()
		case "sp":
			find.Sp("")
		case "redirect":
			find.RedirectURL("")
		case "secret":
			find.Secret("")
		default:
			find.Details()
		}
	case "brute":
		brute.Do()
	case "auto":
		if subCmd == "" {
			help()
			os.Exit(1)
		}
		auto.Do(subCmd)
	default:
		help()
		os.Exit(1)
	}
}

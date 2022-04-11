package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"flag"

	"github.com/cli/go-gh"
)

var sudo = flag.Bool("sudo", false, "Execute possibly disasterous command.")
var dryrun = flag.Bool("dry-run", false, "Print the sub-command that will be run to stdout, but do not execute it.")

func selectFromMap(matchmap map[string]int) int {
	i := 0
	titles := make([]string, len(matchmap))
	for title, number := range matchmap {
		titles[i] = title
		fmt.Printf("\t%v. [%v] %v\n", i, number, title)
	}
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	index, err := strconv.Atoi(string(char))
	if err != nil {
		fmt.Println("Error reading from stdin.")
	}
	return matchmap[titles[index]]
}

func main() {
	repo, err := gh.CurrentRepository()
	if err != nil {
		fmt.Println("Could not fetch repository details. Does this have a GitHub upstream?")
		return
	}
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println("Usage: core_command \"search term\"")
		return
	}

	if repo.Host() != "github.com" {
		fmt.Printf("Error: Only github.com supported. Found host %v", repo.Host())
	}

	querier := CreateGQLPRQuerier(Repository{repo.Owner(), repo.Name()}, 5*time.Second, 50)

	search := flag.Args()[1]
	if len(search) < 3 {
		fmt.Printf("Error: Search '%v' is too short. Search requires at least 3 characters.\n", search)
		return
	}

	matchmap := FindMatchingPRs(querier, search)

	if len(matchmap) > 9 {
		fmt.Println("Error: Search resulted in too many matches. Try a more specific search string.", search)
		return
	} else if len(matchmap) == 0 {
		fmt.Println("Error: Search resulted in no results. Try again!")
		return
	}

	var number int
	if len(matchmap) > 1 {
		fmt.Println("Multiple matches found. Did you mean...?")
		number = selectFromMap(matchmap)
	} else {
		for _, n := range matchmap {
			number = n
			break
		}
	}

	basecmd := flag.Args()[0]
	if basecmd == "merge" || basecmd == "close" && !*sudo {
		fmt.Printf("Error: basecmd %v considered dangerous. Re-run with --sudo to execute.", basecmd)
		return
	}

	args := append([]string{"pr", basecmd, strconv.Itoa(number)}, flag.Args()[2:]...)

	if *dryrun {
		fmt.Println("Would execute `gh " + strings.Join(args, " ") + "`")
		return
	} else {
		fmt.Println("Executing `gh " + strings.Join(args, " ") + "`")
	}

	gh.Exec(args...)
}

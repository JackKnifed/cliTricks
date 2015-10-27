package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/JackKnifed/cliTricks"
)

// Define a type named "stringSlice" as a slice of strings
type stringSlice []string

// Now, for our new type, implement the two methods of
// the flag.Value interface...
// The first method is String() string
func (s *stringSlice) String() string {
	return fmt.Sprintf("%s", *s)
}

// The second method is Set(value string) error
func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func jsonDecoder(in io.Reader, out io.Writer, t [][]interface{}) (err error) {
	var requestData, item interface{}
	var line []string

	decoder := json.NewDecoder(in)

	for decoder.More() {
		err = decoder.Decode(&requestData)
		if err != nil {
			return err
		}
		for _, oneTarget := range t {
			item, err = cliTricks.GetItem(requestData, oneTarget)
			if err != nil {
				return err
			}
			line = append(line, fmt.Sprintf("%s", item))

		}
		out.Write([]byte(strings.Join(line, " ")))
	}

	out.Write([]byte("\n"))

	return err
}

func main() {
	var targetStrings stringSlice
	flag.Var(&targetStrings, "target", "locations to pluck from input")

	flag.Parse()

	var targets [][]interface{}
	for _, oneTarget := range targetStrings {
		targets = append(targets, cliTricks.BreakupArray(oneTarget))
	}

	if err := jsonDecoder(os.Stdin, os.Stdout, targets); err != nil {
		log.Print(err)
	}
}
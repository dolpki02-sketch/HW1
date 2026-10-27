package main

import (
	"flag"
	"fmt"
	"os"
)

type Config struct {
	Timeout int
	URLs    []string
}

func parseConfig() *Config {
	var timeoutLong int
	var timeoutShort int

	flag.IntVar(&timeoutLong, "timeout", 15, "Timeout in seconds")
	flag.IntVar(&timeoutShort, "t", 15, "Timeout in seconds (shorthand)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <url1> <url2> ...\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	finalTimeout := 15
	if timeoutLong != 15 {
		finalTimeout = timeoutLong
	} else if timeoutShort != 15 {
		finalTimeout = timeoutShort
	}

	urls := flag.Args()
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "Ошибка: не указаны URL.")
		flag.Usage()
		os.Exit(1)
	}

	return &Config{
		Timeout: finalTimeout,
		URLs:    urls,
	}
}

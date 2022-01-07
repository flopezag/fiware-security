package cmd

import (
	commandLine "github.com/urfave/cli/v2"

	"fmt"
	"log"
	"os"
	"time"
)

func Init() {
	//var file string
	//var project string
	var pull bool

	app := &commandLine.App{
		Name:      "scan",
		Usage:     "Security Scan of FIWARE Generic Enablers",
		Version:   "v0.1.0",
		Compiled:  time.Now(),
		Copyright: "(c) 2022 FIWARE Foundation, e.V.",

		Description: `This program searchs Docker Images vulnerabilities in the FIWARE Generic Enablers based on Anchore 
and Clair tools and provide a set of best practices of a running instance of them based on a docker 
compose file.

If there is no arguments, the analysis will be over all the content defined in the enablers.json.
In case that a specific FIWARE Generic Enabler is specified, the analysis will be developed only 
on this component.`,

		Action: func(c *commandLine.Context) error {
			fmt.Printf("Hello %q", c.Args().Get(0))

			commandLine.ShowAppHelp(c)

			return nil
		},

		UsageText: "scan [global options] [args and such]",

		Flags: []commandLine.Flag{
			&commandLine.BoolFlag{
				Name:        "pull",
				Usage:       "Pull images bedore running scans",
				Aliases:     []string{"p"},
				Destination: &pull,
			},
		},

		ArgsUsage: "[(Optional) FIWARE GE]",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	/*
		initialize()

		Security_analysis("fiware/orion-ld:latest")
		Docker_bench_security("fiware/orion-ld:latest")
		Anchore("fiware/orion-ld:latest")

		clean()
	*/
}

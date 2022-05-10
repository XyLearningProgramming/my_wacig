package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"monkey/my_engine"
	repl "monkey/my_repl"
	"os"
	"os/user"
)

var (
	engineFlag = flag.String(
		"engine",
		"vm",
		"engine to execute code; possible options: vm, eval; default to vm",
	)
	vmEngine   = my_engine.NewVMEngine()
	evalEngine = my_engine.NewEvalEngine()
)

func main() {
	flag.Parse()

	args := flag.Args()
	switch len(args) {
	case 0:
		// interactive repl
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the Monkey programming language!\n",
			user.Username)
		fmt.Printf("Feel free to type in commands\n")
		switch *engineFlag {
		case "eval":
			repl.Start(os.Stdout, evalEngine)
		case "vm":
			fallthrough
		default:
			repl.Start(os.Stdout, vmEngine)
		}
	case 1:
		// output file result
		code, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Printf(
				"Sorry, Monkey doesn't know how to compile %s: %s",
				args[0],
				err.Error(),
			)
			os.Exit(1)
		}
		switch *engineFlag {
		case "eval":
			res, err := evalEngine.Evaluate(string(code))
			if err != nil {
				repl.PrintErrors(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Print(res.String())
		case "vm":
			fallthrough
		default:
			res, err := vmEngine.Evaluate(string(code))
			if err != nil {
				repl.PrintErrors(os.Stderr, err)
				os.Exit(1)
			}
			fmt.Print(res.String())
		}
	default:
		// error
		fmt.Printf("Sorry, Monkey only knows how to compile one file at a time for now\n")
		os.Exit(1)
	}
}

package my_repl

import (
	"fmt"
	"io"
	"monkey/my_engine"
	"os"

	// "monkey/evaluator"
	// "monkey/parser"

	console "github.com/xingshuo/console/src"
)

const PROMPT = ">> "

func Start(out io.Writer, engine my_engine.Engine) {
	con := console.NewConsole()
	con.SetKeyDownHook('\x03', func(c *console.Console, s string) {
		fmt.Fprintln(out, "keyboard interupt.")
		os.Exit(0)
	})
	con.Init(func(c *console.Console, line string) {
		switch line {
		case "exit":
			fallthrough
		case "exit()":
			fmt.Fprintln(out)
			os.Exit(0)
		default:
			evaluated, err := engine.Evaluate(line)
			if err != nil {
				PrintErrors(out, err)
				return
			}
			if evaluated != nil {
				io.WriteString(out, "\n")
				io.WriteString(out, evaluated.String())
			}
		}
	})
	con.LoopCmd()
}

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func PrintErrors(out io.Writer, err error) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	io.WriteString(out, fmt.Sprintf("\t%s\n", err.Error()))
}

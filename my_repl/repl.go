package my_repl

import (
	"fmt"
	"io"
	"monkey/my_ast"
	"monkey/my_compiler"
	lexer "monkey/my_lexer"
	"monkey/my_vm"
	"os"

	evaluator "monkey/my_evaluator"
	parser "monkey/my_parser"

	// "monkey/evaluator"
	// "monkey/parser"
	object "monkey/my_object"

	console "github.com/xingshuo/console/src"
)

const PROMPT = ">> "

func Start(out io.Writer) {
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
			compiler_call(line, out)
		}
	})
	con.LoopCmd()
}

func compiler_call(line string, out io.Writer) {
	program := parse(line, out)
	if program == nil {
		return
	}
	comp := my_compiler.New()
	err := comp.Compile(program)
	if err != nil {
		printParserErrors(out, []string{err.Error()})
		return
	}
	virtualMachine := my_vm.New(comp.ByteCode())
	err = virtualMachine.Run()
	if err != nil {
		printParserErrors(out, []string{err.Error()})
		return
	}
	stackTop := virtualMachine.LastPoppedStackItem()
	io.WriteString(out, "\n")
	io.WriteString(out, stackTop.String())
}

var env = object.NewEnvironment()

func eval_call(line string, out io.Writer) {
	program := parse(line, out)
	if program == nil {
		return
	}
	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, "\n")
		io.WriteString(out, evaluated.String())
	}
}

func parse(line string, out io.Writer) *my_ast.Program {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.Parse()
	if p.Error() != nil {
		printParserErrors(out, []string{p.Error().Error()})
		return nil
	}
	return program
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

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

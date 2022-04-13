# Writing a compiler in go - My improved version

## What is repo for?

This repo is my version of implementing the compiler of the "monkey" language described in the book "writing a compiler in go" by Thorsten Ball.

It is based on my version of the [first part](https://github.com/XyLearningProgramming/my_writing_an_interpreter_in_go_playground) of the sequel, with tweaks to the monkey language such as the `float` structure, python like array indexing, etc.

## How to run?

- Use command line calling `repl`:

    ```bash
    go run ./
    ```

- Run with .monkey source code file (TODO):

## Features

### Improvements based on the part II (TODO)

### Improvements based on the part I

- `float` type
- implicit conversions among `boolean`, `integer`, `float` constants
- `repl` enabling backspace, history tracing, exit() command
- `string` literals now can start with either ' or ", support backslash escapes

    ```bash
    "boo\"foo"
    'boo""foo'
    'boo\\"foo'
    '\"'
    "\n\t\r"
    ```

- `array` allow python-like indexing including striding

    ```bash
    [0,1][1:] yields [1]
    [0,1,2][::-1] yields [2,1,0]
    ```

- loops: `for` && `while` && `do while` and keywords `continue` and `break`

    ```bash
    let a = 1;while(a==1){let a = 2; break;};a; 
    # output: 1
    ```

- reassigning values using `=`

    ```bash
    let a = 1; do{let a= 2}while(false);a=3;a;
    # output: 3
    ```

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var SuccesCode int = 0
var FailedCode int = 1

var ConstantsString = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var Identifier int = 0

func UNUSED(x ...interface{}) {}

// opcodes
const (
	OpCodes_Increment int = iota
	OpCodes_Decrement
	OpCodes_PointerIncrement
	OpCodes_PointerDecrement
	OpCodes_PrintChar
)

// token struct
type Token_T struct {
	value      string // value of token
	type_token int    // token type
}

// lexer struct
type Lexer_T struct {
	tokens   []Token_T // array of tokens
	filename string    // filename
}

// opcode struct
type OpCodes_T struct {
	opcode int // opcode
}

// parser struct
type Parser_T struct {
	opcodes []OpCodes_T
}

// function that's create a token and return it
func token_create(value string, type_token int) Token_T {
	var token Token_T

	token.type_token = type_token
	token.value = value

	return token
}

func opcode_create(opcode int) OpCodes_T {
	var opcoderet OpCodes_T

	opcoderet.opcode = opcode

	return opcoderet
}

// initialize parser
func parser_init(lexer Lexer_T) Parser_T {
	var parser Parser_T
	var _p_error bool = false

	for i := 0; i < len(lexer.tokens); i++ {
		if lexer.tokens[i].value == "Wzium" {
			parser.opcodes = append(parser.opcodes, opcode_create(OpCodes_Increment))
		} else if lexer.tokens[i].value == "wziuM" {
			parser.opcodes = append(parser.opcodes, opcode_create(OpCodes_Decrement))
		} else if lexer.tokens[i].value == "Wziwzium" {
			parser.opcodes = append(parser.opcodes, opcode_create(OpCodes_PointerIncrement))
		} else if lexer.tokens[i].value == "wziwziuM" {
			parser.opcodes = append(parser.opcodes, opcode_create(OpCodes_PointerDecrement))
		} else if lexer.tokens[i].value == "Wziumnij" {
			parser.opcodes = append(parser.opcodes, opcode_create(OpCodes_PrintChar))
		} else {
			fmt.Printf("Unknown keyword: \"%s\"", lexer.tokens[i].value)
			_p_error = true
		}
	}

	UNUSED(_p_error)

	if _p_error == true {
		os.Exit(FailedCode)
	}

	return parser
}

// initialize lexer
func lexer_init(value string, filename string) Lexer_T {
	var lexer Lexer_T

	lexer.filename = filename

	var lexer_row_ int = 0
	var lexer_col_ int = 0
	var lexer_space_ bool = false
	var _l_error bool = false

	for i := 0; i < len(value); i++ {
		a_char := string(value[i])

		if a_char == "\n" {
			lexer_row_++
			lexer_col_ = 0
			lexer_space_ = false
		} else if a_char == " " {
			lexer_space_ = false
		} else if strings.ContainsAny(a_char, ConstantsString) {

			if lexer_space_ == false {
				lexer.tokens = append(lexer.tokens, token_create(a_char, Identifier))
				lexer_space_ = true
			} else {
				lexer.tokens[len(lexer.tokens)-1].value += a_char
				lexer_space_ = true
			}
		} else if value[i] == 13 {
			continue
		} else {
			fmt.Printf("%s:%d:%d: Nie znany znak \"%s\"\n %d", filename, lexer_row_, lexer_col_, a_char, value[i])
			_l_error = true
		}
	}

	UNUSED(lexer_col_, lexer_row_, lexer_space_)

	if _l_error == true {
		os.Exit(FailedCode)
	}

	return lexer
}

// initialize eval
func eval_init(opcodes []OpCodes_T) {
	var stack [30000]int8
	var stack_pointer int = 0
	var error_eval bool = false
	var str string = ""

	for i := 0; i < len(opcodes); i++ {
		switch opcodes[i].opcode {
		case OpCodes_Increment:
			stack[stack_pointer]++
			break
		case OpCodes_Decrement:
			stack[stack_pointer]--
			break
		case OpCodes_PointerIncrement:
			stack_pointer++
			stack[stack_pointer] = 0
			break
		case OpCodes_PointerDecrement:
			if stack_pointer != 0 {
				stack[stack_pointer]--
			}
			break
		case OpCodes_PrintChar:
			str += string(int(stack[stack_pointer]))
		}
	}

	println("Output: " + str)

	UNUSED(stack, stack_pointer, error_eval)
}

// check if file exists
func file_exists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

// get file content
func get_file_content(filename string) string {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("File not exists")
	}

	return string(content)
}

// print out help
func help() {
	fmt.Println("WziumFuck - Prosty język podobny do brainfuck. Jest on oparty na stacku. Interpreter został napisany w GoLangu\n\t- --help - Pokazuje tą wiadomość\n\t- <nazwa_pliku.wf> - Interpretuje plik")
}

// main function
func main() {
	argv := os.Args

	if len(argv) < 2 {
		// TODO: Shell
	} else {
		if argv[1] == "--help" || argv[1] == "--h" || argv[1] == "-h" || argv[1] == "-help" {
			help()
			os.Exit(SuccesCode)
		}

		var fileexst bool = file_exists(argv[1])

		if fileexst == false {
			fmt.Println(argv[1] + " <-- Ten plik nie istnieje")
			os.Exit(FailedCode)
		}

		var file_content string = get_file_content(argv[1])

		var lexer Lexer_T = lexer_init(file_content, argv[1])
		var parser Parser_T = parser_init(lexer)
		eval_init(parser.opcodes)
	}
}

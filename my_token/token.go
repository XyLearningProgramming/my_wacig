package my_token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT" // add, foobar, x, y, ...
	INT    = "INT"   // 1343456
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	REASSIGN = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	COLON    = ":"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	DO       = "DO"
	WHILE    = "WHILE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	NULL     = "NULL"
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":       FUNCTION,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"do":       DO,
	"while":    WHILE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
	"null":     NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

var kwreversed = map[TokenType]string{
	FUNCTION: "fn",
	LET:      "let",
	TRUE:     "true",
	FALSE:    "false",
	IF:       "if",
	ELSE:     "else",
	RETURN:   "return",
	DO:       "do",
	WHILE:    "while",
	FOR:      "for",
	BREAK:    "break",
	CONTINUE: "continue",
	NULL:     "null",
}

func LookupKeywords(t TokenType) string {
	if s, sok := kwreversed[t]; sok {
		return s
	}
	return "unknown"
}

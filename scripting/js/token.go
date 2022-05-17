package js

import "github.com/evocert/kwe/scripting/lang"

const (
	_ lang.Token = iota

	ILLEGAL
	EOF
	COMMENT

	STRING
	NUMBER

	PLUS      // +
	MINUS     // -
	MULTIPLY  // *
	EXPONENT  // **
	SLASH     // /
	REMAINDER // %

	AND                  // &
	OR                   // |
	EXCLUSIVE_OR         // ^
	SHIFT_LEFT           // <<
	SHIFT_RIGHT          // >>
	UNSIGNED_SHIFT_RIGHT // >>>

	ADD_ASSIGN       // +=
	SUBTRACT_ASSIGN  // -=
	MULTIPLY_ASSIGN  // *=
	EXPONENT_ASSIGN  // **=
	QUOTIENT_ASSIGN  // /=
	REMAINDER_ASSIGN // %=

	AND_ASSIGN                  // &=
	OR_ASSIGN                   // |=
	EXCLUSIVE_OR_ASSIGN         // ^=
	SHIFT_LEFT_ASSIGN           // <<=
	SHIFT_RIGHT_ASSIGN          // >>=
	UNSIGNED_SHIFT_RIGHT_ASSIGN // >>>=

	LOGICAL_AND // &&
	LOGICAL_OR  // ||
	COALESCE    // ??
	INCREMENT   // ++
	DECREMENT   // --

	EQUAL        // ==
	Strict_EQUAL // ===
	LESS         // <
	GREATER      // >
	ASSIGN       // =
	NOT          // !

	BITWISE_NOT // ~

	NOT_EQUAL        // !=
	Strict_NOT_EQUAL // !==
	LESS_OR_EQUAL    // <=
	GREATER_OR_EQUAL // >=

	LEFT_PARENTHESIS // (
	LEFT_BRACKET     // [
	LEFT_BRACE       // {
	COMMA            // ,
	PERIOD           // .

	RIGHT_PARENTHESIS // )
	RIGHT_BRACKET     // ]
	RIGHT_BRACE       // }
	SEMICOLON         // ;
	COLON             // :
	QUESTION_MARK     // ?
	QUESTION_DOT      // ?.
	ARROW             // =>
	ELLIPSIS          // ...
	BACKTICK          // `

	// Tokens below (and only them) are syntactically valid identifiers

	IDENTIFIER
	KEYWORD
	BOOLEAN
	NULL

	IF
	IN
	OF
	DO

	VAR
	LET
	FOR
	NEW
	TRY

	THIS
	ELSE
	CASE
	VOID
	WITH

	CONST
	WHILE
	BREAK
	CATCH
	THROW

	RETURN
	TYPEOF
	DELETE
	SWITCH

	DEFAULT
	FINALLY

	FUNCTION
	CONTINUE
	DEBUGGER

	INSTANCEOF
)

var Token2string = [...]string{
	ILLEGAL:                     "ILLEGAL",
	EOF:                         "EOF",
	COMMENT:                     "COMMENT",
	KEYWORD:                     "KEYWORD",
	STRING:                      "STRING",
	BOOLEAN:                     "BOOLEAN",
	NULL:                        "NULL",
	NUMBER:                      "NUMBER",
	IDENTIFIER:                  "IDENTIFIER",
	PLUS:                        "+",
	MINUS:                       "-",
	EXPONENT:                    "**",
	MULTIPLY:                    "*",
	SLASH:                       "/",
	REMAINDER:                   "%",
	AND:                         "&",
	OR:                          "|",
	EXCLUSIVE_OR:                "^",
	SHIFT_LEFT:                  "<<",
	SHIFT_RIGHT:                 ">>",
	UNSIGNED_SHIFT_RIGHT:        ">>>",
	ADD_ASSIGN:                  "+=",
	SUBTRACT_ASSIGN:             "-=",
	MULTIPLY_ASSIGN:             "*=",
	EXPONENT_ASSIGN:             "**=",
	QUOTIENT_ASSIGN:             "/=",
	REMAINDER_ASSIGN:            "%=",
	AND_ASSIGN:                  "&=",
	OR_ASSIGN:                   "|=",
	EXCLUSIVE_OR_ASSIGN:         "^=",
	SHIFT_LEFT_ASSIGN:           "<<=",
	SHIFT_RIGHT_ASSIGN:          ">>=",
	UNSIGNED_SHIFT_RIGHT_ASSIGN: ">>>=",
	LOGICAL_AND:                 "&&",
	LOGICAL_OR:                  "||",
	COALESCE:                    "??",
	INCREMENT:                   "++",
	DECREMENT:                   "--",
	EQUAL:                       "==",
	Strict_EQUAL:                "===",
	LESS:                        "<",
	GREATER:                     ">",
	ASSIGN:                      "=",
	NOT:                         "!",
	BITWISE_NOT:                 "~",
	NOT_EQUAL:                   "!=",
	Strict_NOT_EQUAL:            "!==",
	LESS_OR_EQUAL:               "<=",
	GREATER_OR_EQUAL:            ">=",
	LEFT_PARENTHESIS:            "(",
	LEFT_BRACKET:                "[",
	LEFT_BRACE:                  "{",
	COMMA:                       ",",
	PERIOD:                      ".",
	RIGHT_PARENTHESIS:           ")",
	RIGHT_BRACKET:               "]",
	RIGHT_BRACE:                 "}",
	SEMICOLON:                   ";",
	COLON:                       ":",
	QUESTION_MARK:               "?",
	QUESTION_DOT:                "?.",
	ARROW:                       "=>",
	ELLIPSIS:                    "...",
	BACKTICK:                    "`",
	IF:                          "if",
	IN:                          "in",
	OF:                          "of",
	DO:                          "do",
	VAR:                         "var",
	LET:                         "let",
	FOR:                         "for",
	NEW:                         "new",
	TRY:                         "try",
	THIS:                        "this",
	ELSE:                        "else",
	CASE:                        "case",
	VOID:                        "void",
	WITH:                        "with",
	CONST:                       "const",
	WHILE:                       "while",
	BREAK:                       "break",
	CATCH:                       "catch",
	THROW:                       "throw",
	RETURN:                      "return",
	TYPEOF:                      "typeof",
	DELETE:                      "delete",
	SWITCH:                      "switch",
	DEFAULT:                     "default",
	FINALLY:                     "finally",
	FUNCTION:                    "function",
	CONTINUE:                    "continue",
	DEBUGGER:                    "debugger",
	INSTANCEOF:                  "instanceof",
}

var keywordMap = lang.KeywordMap{
	"if": {
		Token: IF,
	},
	"in": {
		Token: IN,
	},
	"do": {
		Token: DO,
	},
	"var": {
		Token: VAR,
	},
	"for": {
		Token: FOR,
	},
	"new": {
		Token: NEW,
	},
	"try": {
		Token: TRY,
	},
	"this": {
		Token: THIS,
	},
	"else": {
		Token: ELSE,
	},
	"case": {
		Token: CASE,
	},
	"void": {
		Token: VOID,
	},
	"with": {
		Token: WITH,
	},
	"while": {
		Token: WHILE,
	},
	"break": {
		Token: BREAK,
	},
	"catch": {
		Token: CATCH,
	},
	"throw": {
		Token: THROW,
	},
	"return": {
		Token: RETURN,
	},
	"typeof": {
		Token: TYPEOF,
	},
	"delete": {
		Token: DELETE,
	},
	"switch": {
		Token: SWITCH,
	},
	"default": {
		Token: DEFAULT,
	},
	"finally": {
		Token: FINALLY,
	},
	"function": {
		Token: FUNCTION,
	},
	"continue": {
		Token: CONTINUE,
	},
	"debugger": {
		Token: DEBUGGER,
	},
	"instanceof": {
		Token: INSTANCEOF,
	},
	"const": {
		Token: CONST,
	},
	"class": {
		Token:         KEYWORD,
		FutureKeyword: true,
	},
	"enum": {
		Token:         KEYWORD,
		FutureKeyword: true,
	},
	"export": {
		Token:         KEYWORD,
		FutureKeyword: true,
	},
	"extends": {
		Token:         KEYWORD,
		FutureKeyword: true,
	},
	"import": {
		Token:         KEYWORD,
		FutureKeyword: true,
	},
	"super": {
		Token:         KEYWORD,
		FutureKeyword: true,
	},
	"implements": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
	"interface": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
	"let": {
		Token:  LET,
		Strict: true,
	},
	"package": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
	"private": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
	"protected": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
	"public": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
	"static": {
		Token:         KEYWORD,
		FutureKeyword: true,
		Strict:        true,
	},
}

func init() {
	lang.RegisterLangKeyWordMap("js", keywordMap)
}

package formula_engine

import (
	"fmt"
	"strings"
)

// lexer 词法分析器
type lexer struct {
	FStr        string
	Idx         int
	CurrentChar uint8 // 当前的字符
}

func newLexer(fStr string) *lexer {
	l := &lexer{
		FStr:        fStr,
		Idx:         -1,
		CurrentChar: 0,
	}
	l.advance()
	return l
}

// MakeTokens 获取tokens
func (l *lexer) MakeTokens() ([]*token, error) {
	tokens := make([]*token, 0)
	for l.CurrentChar != 0 {
		switch {
		case InSlice([]uint8{' ', '\t'}, l.CurrentChar):
			l.advance()
		case IsDigit(l.CurrentChar):
			token, err := l.makeNumber()
			if err != nil {
				return nil, err // todo: 错误处理
			}
			tokens = append(tokens, token)
		case l.CurrentChar == '+':
			tokens = append(tokens, l.makeCharacter(TTPlus))
		case l.CurrentChar == '-':
			tokens = append(tokens, l.makeCharacter(TTMinus))
		case l.CurrentChar == '*':
			tokens = append(tokens, l.makeCharacter(TTMul))
		case l.CurrentChar == '/':
			tokens = append(tokens, l.makeCharacter(TTDiv))
		case l.CurrentChar == '^':
			tokens = append(tokens, l.makeCharacter(TTPow))
		case l.CurrentChar == '(':
			tokens = append(tokens, l.makeCharacter(TTLparen))
		case l.CurrentChar == ')':
			tokens = append(tokens, l.makeCharacter(TTRparen))
		case l.CurrentChar == '&':
			tokens = append(tokens, l.makeCharacter(TTAnd))
		case l.CurrentChar == '|':
			tokens = append(tokens, l.makeCharacter(TTOr))
		case l.CurrentChar == '!':
			tokens = append(tokens, l.makeNot())
		case l.CurrentChar == '=':
			tokens = append(tokens, l.makeCharacter(TTEq))
		case l.CurrentChar == ',':
			tokens = append(tokens, l.makeCharacter(TTComma))
		case l.CurrentChar == '>':
			tokens = append(tokens, l.makeCompare(TTGt))
		case l.CurrentChar == '<':
			tokens = append(tokens, l.makeCompare(TTLt))
		case l.CurrentChar == '{':
			token, err := l.makeIdentifier()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		case IsAlpha(l.CurrentChar):
			token, err := l.makeFunction()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, token)
		default:
			// 没有匹配到，非法字符错误
			return nil, l.makeErr(illegalCharErrMsg, fmt.Sprintf("UnExpected character '%c'", l.CurrentChar))
		}
	}
	tokens = append(tokens, newToken(TTEof, TTEof, l.Idx, l.Idx))
	return tokens, nil
}

// makeCharacter 处理字符，通用组装token方法
func (l *lexer) makeCharacter(type_ TT) *token {
	token := newToken(type_, string(l.CurrentChar), l.Idx, l.Idx)
	l.advance()
	return token
}

// makeNumber 处理数字：整数、小数
func (l *lexer) makeNumber() (*token, error) {
	var numBuilder strings.Builder
	dotCount := 0
	start := l.Idx
	//for InSlice(append(Digits, '.'), l.CurrentChar) {
	for IsDigit(l.CurrentChar) || l.CurrentChar == '.' {
		if l.CurrentChar == '.' {
			if dotCount == 1 {
				return nil, l.makeErr(illegalCharErrMsg, "UnExpected character '.', Only one '.' is allowed")
			}
			dotCount += 1
			numBuilder.WriteByte('.')
		} else {
			numBuilder.WriteByte(l.CurrentChar)
		}
		l.advance()
	}

	str := numBuilder.String()

	// 判断是否有错误
	if !IsDigit(str[len(str)-1]) {
		// 如果数字最后一位不是 0~9 ,即有可能是 .
		return nil, l.makeErr(illegalCharErrMsg, fmt.Sprintf("UnExpected character '%c', expected '0'~'9'", str[len(str)-1]))
	}

	return newToken(TTNum, numBuilder.String(), start, l.Idx-1), nil
}

// makeNot 处理非 ! 或者不等于 !=
func (l *lexer) makeNot() *token {
	var str strings.Builder
	begin := l.Idx
	str.WriteByte(l.CurrentChar)
	l.advance()
	// 判断是否是 !=
	if l.CurrentChar == '=' {
		str.WriteByte(l.CurrentChar)
		l.advance()
		return newToken(TTNeq, str.String(), begin, l.Idx-1)
	}
	return newToken(TTNot, str.String(), begin, l.Idx-1)
}

// makeCompare 处理大于号或者小于号 > >= < <=
func (l *lexer) makeCompare(type_ TT) *token {
	var str strings.Builder

	begin := l.Idx
	str.WriteByte(l.CurrentChar)
	l.advance()
	// 判断是否是 >= 或者 <=
	if l.CurrentChar == '=' {
		str.WriteByte(l.CurrentChar)
		if type_ == TTGt {
			type_ = TTGte
		} else {
			type_ = TTLte
		}
		l.advance()
	}
	return newToken(type_, str.String(), begin, l.Idx-1)
}

// makeIdentifier 处理变量
func (l *lexer) makeIdentifier() (*token, error) {
	var (
		str strings.Builder
	)
	l.advance()
	start := l.Idx
	// 变量开头不是'字母'或者'_',报错
	if !(IsAlpha(l.CurrentChar) || l.CurrentChar == '_') {
		return nil, l.makeErr(illegalCharErrMsg, fmt.Sprintf("UnExpected Initial '%c', expected letter or '_'", l.CurrentChar))
	}
	// 字符是字母、数字或'_'
	for IsAlpha(l.CurrentChar) || IsDigit(l.CurrentChar) || l.CurrentChar == '_' {
		str.WriteByte(l.CurrentChar)
		l.advance()
	}

	if l.CurrentChar != '}' {
		return nil, l.makeErr(illegalCharErrMsg, fmt.Sprintf("UnExpected character '%c', expected '}' after an Identifier", l.CurrentChar))
	}
	l.advance()
	return newToken(TTIdentifier, str.String(), start, l.Idx-2), nil
}

// makeFunction 处理函数
func (l *lexer) makeFunction() (*token, error) {
	var strBuilder strings.Builder
	start := l.Idx
	for IsAlpha(l.CurrentChar) || l.CurrentChar == '.' {
		strBuilder.WriteByte(l.CurrentChar)
		l.advance()
	}

	str := strings.ToUpper(strBuilder.String())
	_, ok := FuncMap[str]
	if !ok {
		return nil, l.makeErr(illegalCharErrMsg, fmt.Sprintf("UnKnow function name %s", str))
	}
	return newToken(TTFunction, str, start, l.Idx-1), nil
}

// advance 预读
func (l *lexer) advance() {
	l.Idx += 1
	if l.Idx < len(l.FStr) {
		l.CurrentChar = l.FStr[l.Idx]
	} else {
		l.CurrentChar = 0
	}
}

// makeErrWithIdx 组装错误
func (l *lexer) makeErr(errName, details string) error {
	return makeStrErr(l.Idx, l.FStr, errName, details)
}

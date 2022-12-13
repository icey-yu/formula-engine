package formula_engine

import "github.com/shopspring/decimal"

//var (
//	Digits  = []uint8{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
//	LoAlpha = []uint8{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
//	UpAlpha = []uint8{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
//)

type TT string

// TT => token Type
const (
	TTNum    TT = "NUM"    // 数字类型
	TTPlus      = "PLUS"   // + 加号
	TTMinus     = "MINUS"  // - 减号
	TTMul       = "MUL"    // * 乘号
	TTDiv       = "DIV"    // / 除号
	TTPow       = "POW"    // ^ 乘方
	TTLparen    = "LPAREN" // ( 左括号
	TTRparen    = "RPAREN" // ) 右括号
	TTAnd       = "AND"    // & 与
	TTOr        = "OR"     // | 或
	TTNot       = "NOT"    // ! 非
	TTEq        = "EQ"     // = 等于
	TTNeq       = "NEQ"    // != 不等于
	TTGt        = "GT"     // > 大于
	TTLt        = "LT"     // < 小于
	TTGte       = "GTE"    // >= 大于等于
	TTLte       = "LTE"    // <= 小于等于
	TTComma     = "COMMA"  // , 逗号

	TTIdentifier = "IDENTIFIER" // 变量名
	TTFunction   = "FUNCTION"   // 函数
	TTEof        = "EOF"        // 结束符
)

// 报错信息
const (
	illegalCharErrMsg   = "Illegal Character"
	illegalSyntaxErrMsg = "Illegal Syntax"
	illegalCalErrMsg    = "Illegal Calculation"
	systemErrMsg        = "System Err"
)

// 表示函数参数个数的特殊值
const (
	Abt    = -1 // arbitrarily 任意函数参数个数
	GtZero = -2 // greater than 0，大于0个
)

var (
	// FuncMap 函数map，规定函数调用哪个方法
	FuncMap = map[string]func(...*decimal.Decimal) (*decimal.Decimal, error){
		"MAX": max,
		"MIN": min,
		"IF":  if_,
	}

	// FuncParNumMap 函数参数个数map，用于校验。
	FuncParNumMap = map[string]int{
		"MAX": GtZero,
		"MIN": GtZero,
		"IF":  3,
	}
)

const (
	zeroStr = "0"
)

package formula_engine

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

// convertBool 将bool值转化为Decimal类型。0为假，1为真
func convertBool(p bool) decimal.Decimal {
	if p {
		return decimal.NewFromInt(1)
	} else {
		return decimal.Zero
	}
}

// convertToBool 将Decimal值转化为bool类型。0为假，1为真
func convertToBool(p *decimal.Decimal) bool {
	if p.Equal(decimal.Zero) {
		return false
	} else {
		return true
	}
}

// checkParNum 函数计算前进行参数个数校验
func checkParNum(parNum int, length int) error {
	switch parNum {
	case Abt:
		return nil
	case GtZero:
		if length == 0 {
			return makeErr(illegalSyntaxErrMsg, fmt.Sprintf("Required Greater than 0 params,but got %d.\n", length))
		}
	default:
		if parNum != length {
			return makeErr(illegalSyntaxErrMsg, fmt.Sprintf("Required %d params,but got %d.\n", parNum, length))
		}
	}
	return nil
}

// plus 加
func plus(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	res := p1.Add(*p2)
	return &res, nil
}

// minus 减
func minus(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	res := p1.Sub(*p2)
	return &res, nil
}

// unPlusOrMinus 处理 +1 ++1 类似情况
func unPlus(p *decimal.Decimal) (*decimal.Decimal, error) {
	return p, nil
}

// unPlusOrMinus 处理 -1 --1 类似情况
func unMinus(p *decimal.Decimal) (*decimal.Decimal, error) {
	res := p.Mul(decimal.NewFromInt(-1))
	return &res, nil
}

// mul 乘
func mul(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	res := p1.Mul(*p2)
	return &res, nil
}

// div 除
func div(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	if p2.Equal(decimal.Zero) {
		return nil, makeErr(illegalCalErrMsg, "Cannot divide by 0")
	}
	res := p1.Div(*p2)
	return &res, nil
}

// pow 乘方
func pow(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	res := p1.Pow(*p2)
	return &res, nil
}

// and 与
func and(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b1 := convertToBool(p1)
	b2 := convertToBool(p2)
	b := b1 && b2
	res := convertBool(b)
	return &res, nil
}

// or 或
func or(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b1 := convertToBool(p1)
	b2 := convertToBool(p2)
	b := b1 || b2
	res := convertBool(b)
	return &res, nil
}

// not 非
func not(p *decimal.Decimal) (*decimal.Decimal, error) {
	b := convertToBool(p)
	b = !b
	res := convertBool(b)
	return &res, nil
}

// eq 等于
func eq(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b := p1.Equal(*p2)
	res := convertBool(b)
	return &res, nil
}

// neq 不等于
func neq(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b := p1.Equal(*p2)
	res := convertBool(!b)
	return &res, nil
}

// gt 大于
func gt(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b := p1.GreaterThan(*p2)
	res := convertBool(b)
	return &res, nil
}

// lt 小于
func lt(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b := p1.LessThan(*p2)
	res := convertBool(b)
	return &res, nil
}

// gte 大于等于
func gte(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b := p1.GreaterThanOrEqual(*p2)
	res := convertBool(b)
	return &res, nil
}

// lte 小于等于
func lte(p1 *decimal.Decimal, p2 *decimal.Decimal) (*decimal.Decimal, error) {
	b := p1.LessThanOrEqual(*p2)
	res := convertBool(b)
	return &res, nil
}

// function 函数
func function(funcName string, ps ...*decimal.Decimal) (*decimal.Decimal, error) {
	parNum, ok := FuncParNumMap[funcName]
	if !ok {
		return nil, makeErr(illegalCharErrMsg, fmt.Sprintf("UnKnow function name %s in FuncParNumMap", funcName))
	}
	err := checkParNum(parNum, len(ps))
	if err != nil {
		return nil, errors.Wrapf(err, fmt.Sprintf("Function name: %s", funcName))
	}
	fun, ok := FuncMap[funcName]
	if !ok {
		return nil, makeErr(illegalCharErrMsg, fmt.Sprintf("UnKnow function name %s in FuncMap", funcName))
	}
	res, err := fun(ps...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

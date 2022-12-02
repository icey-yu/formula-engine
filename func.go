// ----------------------------------------------------------------------------------------------------------------
// func ，函数处理
// ----------------------------------------------------------------------------------------------------------------

package formula_engine

import "github.com/shopspring/decimal"

// max 返回最大值。
func max(ps ...*decimal.Decimal) (*decimal.Decimal, error) {
	m := ps[0]
	for _, p := range ps {
		if p.GreaterThan(*m) {
			m = p
		}
	}
	return m, nil
}

// min 返回最小值。
func min(ps ...*decimal.Decimal) (*decimal.Decimal, error) {
	m := ps[0]
	for _, p := range ps {
		if p.LessThan(*m) {
			m = p
		}
	}
	return m, nil
}

// if_ IF函数,必须为三个参数,IF(term, r1, r2). 若term为真，返回r1，否则返回r2。
// em:
//
//	IF(2>1, 3, 4) --> return 3
//	IF(2<1, 3, 4) --> return 4
func if_(ps ...*decimal.Decimal) (*decimal.Decimal, error) {
	// 0 表示为假
	if ps[0].Equal(decimal.Zero) {
		return ps[2], nil
	}
	// 否则为真
	return ps[1], nil
}

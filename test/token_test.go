package test

import (
	formulaengine "e.coding.net/oiine/backend/formula-engine"
	"testing"
)

func TestDemo(t *testing.T) {
	str := "2^(-2^2)+{ds}"
	node, err := formulaengine.GetAstTreeByString(str)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := formulaengine.CalByAstTree(node, nil)
	if err != nil {
		t.Error(err)
		return
	}
	println(res.String())
	println(res.Float64())
	println(res.IntPart())
}

package formula_engine

import (
	"fmt"
	"github.com/pkg/errors"
)

// getTokPos 通过token获得位置信息
func getTokPos(tok *token) string {
	return fmt.Sprintf("Report at token:Type:%s,value:%s,start:%d,end:%d", tok.Type, tok.Value, tok.Start, tok.End)
}

// makeErr 组装错误
func makeErr(errName, details string) error {
	message := fmt.Sprintf("err:%s:%s", errName, details)
	return errors.New(message)
}

// makeStrErr 主要用于词法分析器报错
func makeStrErr(idx int, str, errName, details string) error {
	message := fmt.Sprintf("err:%s:%s,Report at :%s, index:%d\n", errName, details, str, idx)
	return errors.New(message)
}

// makeErrWithIdx 使用idx定位报错
func makeErrWithIdx(idx int, errName, details string) error {
	message := fmt.Sprintf("err:%s:%s,Report at index:%d\n", errName, details, idx)
	return errors.New(message)
}

// makeErrWithToken 使用token定位报错
func makeErrWithToken(tok *token, errName, details string) error {
	message := fmt.Sprintf("err:%s:%s,Report at token:Type:%s,value:%s,start:%d,end:%d\n", errName, details, tok.Type, tok.Value, tok.Start, tok.End)
	return errors.New(message)
}

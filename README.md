## **规则**
建议使用其他md查看器查看。

### **基础运算**

基础运算符表：

| 符号 |   意义   |
| :--: | :------: |
|  +   |    加    |
|  -   |    减    |
|  *   |    乘    |
|  /   |    除    |
| ^ | 乘方 |
|  (   |  左括号  |
|  )  |  右括号  |
|  &   |    与    |
|  \|  |    或    |
|  ！  |    非    |
|  =   |   等于   |
|  >   |   大于   |
|  >=  | 大于等于 |
|  <   |   小于   |
|  <=  | 小于等于 |



### **变量**

变量使用`{}`包裹，如`{count}`表示`count`变量。变量由字母、数字、`_`组成。首写字符必须是`_`或者字母。

### **函数**

只有函数库中存在的函数才可以使用。函数不区分大小写。如：`MAX` `max` `Max` `mAx`都表示函数`MAX`。

函数在构建ast树、运算时都会进行参数个数校验，保证运行安全。

#### **增加函数步骤**

1. 在`func.go`中添加函数方法，参数、返回体必须是固定的格式
2. 在`vars.go/FuncMap`中添加函数名和对应方法
3. 在`vars.go/FuncParNumMap`中添加函数的参数个数。

#### **函数参数特定值**

- `Abt` 任意个
- `GtZero` 大于0个

## **BNF**

此处记录BNF式。

> BNF 规则：
>
> 在双引号中的字("word")代表着这些字符本身。而double_quote用来代表双引号。
>
> 在双引号外的字（有可能有下划线）代表着语法部分。
>
> 尖括号( < > )内包含的为必选项。
>
> 方括号( [ ] )内包含的为可选项。
>
> 大括号( { } )内包含的为可重复0至无数次的项。
>
> 竖线( | )表示在其左右两边任选一项，相当于"OR"的意思。
>
> ::= 是“被定义为”的意思。

优先级(低到高)：`|,&,!,{>,<,=,>=,<=},{+,-},{*,/},^,func,()`

越底层优先级越高

```BNF
<expr> ::= <and_term> { OR <expr> }
<and_term> ::= <not_expr> { AND <and_expr>}
<not_term> ::= { NOT } <com_expr>
<com_term> ::= <pri_pre> { GT|LT|EQ|NEQ|GTE|LTE <pri_pre> }                 // compare term
<pri_ope> ::= { <pri_ope> +|- } <sec_ope>                                   // Primary operation
<sec_ope> ::= { <sec_ope> *|/ } <ter_ope>                                   // Secondary operation
<ter_ope> ::= <factor> { ^ <ter_ope> }                                      // Tertiary operation
<factor> ::= NUM|
			FUNCTION LPAREN [ expr { COMMA expr }] RPAREN|
			IDENTIFIER|
			{ PLUS | MINUS } factor|
			LPAREN expr RPAREN
```

## **使用**

在`enter.go`中提供了`GetAstTreeByString()`和`CalByAstTree`方法。

- `GetAstTreeByString`:参数为字符串，返回ast根节点和err。
- `CalByAstTree`：参数为ast树节点，返回`*decimal.Decimal`类型数据和err

样例：

```golang
func TestDemo(t *testing.T) {
	str := "IF(1+5*9=8/8|2-1<5*1&1+2!=8*5,5+(7*7+2),9998848.24425)"
	node, err := calculator.GetAstTreeByString(str)
	if err != nil {
		t.Error(err)
		return
	}
	res, err := calculator.CalByAstTree(node, nil)
	println(res.String())
	println(res.Float64())
	println(res.IntPart())
}
```

print:

```bash
=== RUN   TestDemo
56
+5.600000e+001 true
56
--- PASS: TestDemo (0.00s)
PASS
```


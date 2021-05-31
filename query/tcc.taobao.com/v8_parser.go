package taobao

import (
	"encoding/json"
	"rogchap.com/v8go"
)

type V8Parser struct {
	ctx *v8go.Context
}

func NewV8Parser() *V8Parser {
	ctx, _ := v8go.NewContext()
	return &V8Parser{ctx: ctx}
}

func (v *V8Parser) Parse(body []byte) PhoneLoc {
	text := string(body)
	val, _ := v.ctx.RunScript(text, "")
	encoded, _ := val.MarshalJSON()
	var p PhoneLoc
	json.Unmarshal(encoded, &p)
	return p
}

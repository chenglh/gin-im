package util

import (
	"github.com/go-playground/validator/v10"
	"reflect"
)

func GetValidMsg(err error, obj interface{}) string {
	getObj := reflect.TypeOf(obj)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			if f, exist := getObj.Elem().FieldByName(e.Field()); exist {
				//错误信息不需要全部返回，当找到第一个错误的信息时，就可以结束
				return f.Tag.Get("msg")
			}
		}
	}
	return err.Error()
}

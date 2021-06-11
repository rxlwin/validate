package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Rule 参数规则
type Rule struct {
	Key          string      //map对应的Key
	T            string      //参数最终要求的type类型
	Fun          Fun         //参数需要验证的函数信息
	DefaultValue interface{} //参数默认值，如果为nil时，要求必填
}

// Fun 验证函数规则
type Fun map[string][]interface{}

// Param 参数格式
type Param map[string]interface{}

// GetRuleParam 获取格式化后的参数
func GetRuleParam(c Param, rules []Rule) (Param, error) {
	ruleParam := Param{}
	var originValue interface{}
	var ok bool
	for _, rule := range rules {
		//检测默认值
		if !checkDefault(rule) {
			return nil, errors.New(rule.Key + "默认值与规则类型不符")
		}
		//设置参数值
		originValue, ok = c[rule.Key]
		if !ok {
			if rule.DefaultValue == false {
				ruleParam[rule.Key] = nil
				continue
			}
			originValue = rule.DefaultValue //设置默认值
		}
		if originValue == nil {
			return nil, errors.New(rule.Key + "必传")
		}
		//格式化参数
		value, err := formatValue(originValue, rule)
		if err != nil {
			return nil, err
		}
		//参数校验
		err = checkFun(value, rule)
		if err != nil {
			return nil, err
		}
		//收集合法参数
		ruleParam[rule.Key] = value
	}

	return ruleParam, nil
}

func checkDefault(r Rule) bool {
	if r.DefaultValue == nil || r.DefaultValue == false {
		return true
	}
	v := reflect.TypeOf(r.DefaultValue).String()
	if v == r.T {
		return true
	}
	return false
}

func checkFun(value interface{}, rule Rule) error {
	f := F{}
	for name, params := range rule.Fun {
		fun := reflect.ValueOf(f).MethodByName(name)
		funParamNum := fun.Type().NumIn()
		ft1 := fun.Type().In(0).Name()
		if ft1 != "" && ft1 != (rule.T) {
			return errors.New(rule.Key + " 规则参数错误 方法" + name + "()中 " + rule.Key + "需" + (fun.Type().In(0).Name()) + "类型")
		}
		in := []reflect.Value{reflect.ValueOf(value)}
		if len(params)+1 != funParamNum {
			return errors.New(rule.Key + " 规则参数错误 fun:" + name + " 需" + val2string(funParamNum-1) + "个参数")
		}
		for i, param := range params {
			rt := fun.Type().In(i + 1)
			if rt != reflect.TypeOf(param) {
				return errors.New(rule.Key + " 规则参数错误 fun:" + name + " 参" + val2string(i+1) + " 需" + rt.Name() + "类型")
			}
			in = append(in, reflect.ValueOf(param))
		}
		res := fun.Call(in)[0]
		if !res.Bool() {
			return errors.New(rule.Key + " 校验失败 fun:" + name + "()")
		}
	}
	return nil
}

func formatValue(val interface{}, rule Rule) (interface{}, error) {
	t := rule.T
	switch t {
	case "string":
		return val2string(val), nil
	case "int":
		return val2int(val), nil
	default:
		return nil, errors.New(rule.Key + " 规则参数错误 未知类型: " + t)
	}
}

func val2string(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func val2int(val interface{}) int {
	switch val.(type) {
	case int:
		return val.(int)
	case string:
		v1 := val.(string)
		v2, err := strconv.Atoi(v1)
		if err != nil {
			return 0
		}
		return v2
	case float32:
		v1 := val.(float32)
		return int(v1)
	case float64:
		v1 := val.(float64)
		return int(v1)
	default:
		return 0
	}
}

package validate

import (
	"regexp"
	"unicode/utf8"
)

type F struct{}

func (f F) Str(v string, min int, max int) bool {
	runeCountInString := utf8.RuneCountInString(v)
	if runeCountInString < min {
		return false
	}
	if max > 0 && runeCountInString > max {
		return false
	}
	return true
}

func (f F) Int(v int, min int, max int) bool {
	if v < min {
		return false
	}
	if max > 0 && v > max {
		return false
	}
	return true
}

func (f F) Phone(phone string) bool {
	reg := `^1\d{10}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(phone)
}

func (f F) Email(email string) bool {
	//reg:=`^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
	reg := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z].){1,4}[a-z]{2,4}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(email)
}

func (f F) Date(date string) bool {
	reg := `^[0-9]{4}-[0-9]{2}-[0-9]{2}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(date)
}

func (f F) Enum(val interface{}, list []interface{}) bool {
	for _, l := range list {
		if val == l {
			return true
		}
	}
	return false
}

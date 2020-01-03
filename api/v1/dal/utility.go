package dal

import (
	"strings"
)

func ParseQueryStringParams(params map[string][]string) ([]string, []interface{}) {
	var wc []string
	var vals []interface{}
	c := 0

	for key, element := range params {
		v := element[0]
		k := key
		v_op := "="
		k_op := ""

		if c > 0 {
			k_op = "AND"
		}
		if len(strings.Split(key, ":")) > 1 {
			k = strings.Split(key, ":")[1]
			k_op = strings.Split(key, ":")[0]
			k_op = translateOperator(k_op)
		}
		if len(strings.Split(element[0], ":")) > 1 {
			v = strings.Split(element[0], ":")[1]
			v_op = strings.Split(element[0], ":")[0]
			v_op = translateOperator(v_op)
		}

		k = translateFieldName(k)
		wc = append(wc, k_op+k+v_op+"?")
		vals = append(vals, v)
		c++
	}

	return wc, vals
}

func translateOperator(op string) string {
	switch o := op; o {
	case "lt":
		return "<"
	case "gt":
		return ">"
	case "||":
		return "OR"
	case "ne":
		return "<>"
	case "lte":
		return "<="
	case "gte":
		return ">="
	case "lk":
		return " LIKE "
	}
	return ""
}

func translateFieldName(fn string) string {
	switch n := fn; strings.ToLower(n) {
	case "fn":
		return "first_name"
	case "ln":
		return "last_name"
	case "bn":
		return "business_name"
	case "a1":
		return "address_1"
	case "a2":
		return "address_2"
	case "pc":
		return "postal_code"
	case "em":
		return "email"
	case "mp":
		return "mobile"
	case "da":
		return "date_added"
	case "dm":
		return "date_modified"
	case "ct":
		return "city"
	case "st":
		return "state"
	case "cn":
		return "country"
	case "sn":
		return "service_name"
	case "sd":
		return "service.description"
	}
	return ""
}

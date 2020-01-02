package dal

import (
	"strings"
)

func ParseQueryStringParams(params map[string][]string) ([]string, []string) {
	var wc []string
	var vals []string
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
	case "eq":
		return "="
	case "lte":
		return "<="
	case "gte":
		return ">="
	case "btw":
		return "BETWEEN"
	case "lk":
		return "LIKE"
	case "n":
		return "IN"
	}
	return ""
}

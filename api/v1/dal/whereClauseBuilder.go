package dal

import (
	"strings"
)

func ParseQueryStringParams(params map[string][]string) ([]string, []interface{}) {
	var wc []string
	var vals []interface{}
	c := 0

	for dbField, value := range params {
		dbf := dbField
		v_op := "="
		dbf_op := ""

		if c > 0 {
			dbf_op = " AND "
		}

		if len(strings.Split(dbf, ":")) > 1 {
			dbf = strings.Split(dbf, ":")[1]
			dbf_op = strings.Split(dbf, ":")[0]
			dbf_op = translateOperator(dbf_op)
		}

		v := value[0]
		//There is an operator in the value 
		if len(strings.Split(v, ":")) > 1 { 
			v_op = strings.Split(v, ":")[0]
			v = strings.Split(v, ":")[1]
		
			// There are multiple values for a single key
			// such as for IN and BETWEEN operators
			if len(strings.Split(v, ",")) > 1 {
				vv := strings.Split(v, ",")
				for i := 0; i < len(vv); i++ {
					vals = append(vals, vv[i])
				}
			} else {
				vals = append(vals, v)
			}
		} else {
			vals = append(vals, v)
		}

		dbf = translateDbFieldName(dbf)
		wc = append(wc, buildWhereClause(dbf_op, dbf, v_op, vals))
		
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
		return "LIKE"
	case "btw":
		return "BETWEEN"
	case "in":
		return "IN"
	}
	return ""
}

func buildWhereClause(dbf_op string, dbf string, v_op string, vals []interface{}) string {
	v_op = translateOperator(v_op)
	switch o := v_op; o {
	case "BETWEEN":
		return buildBetweenWhereClause(dbf_op, dbf, o)
	case "IN":
		return buildInWhereClause(dbf_op, dbf, o, vals)
	}
	return dbf_op + dbf + " " + v_op + " ?"
}

func buildBetweenWhereClause(dbf_op string, dbf string, v_op string) string {
	return  dbf_op + dbf + " " + v_op + " ? AND ?"
}

func buildInWhereClause(dbf_op string, dbf string, v_op string, vals []interface{}) string {
	bindvars := ""
	for i := 0; i < len(vals); i++ {
		bindvars += "?,"
	}
	bindvars = strings.TrimSuffix(bindvars, ",")
	return dbf_op + dbf + " " + v_op + " (" + bindvars + ")"
}

func translateDbFieldName(fn string) string {
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

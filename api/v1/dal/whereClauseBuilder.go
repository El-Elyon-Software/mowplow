package dal

import (
	"net/url"
	"strings"
)

/*
Used to build a where claused, leveraging the API's query syntax

Feilds:
fn = first_name
ln = last_name
bn = business_name
a1 = address_1
a2 = address_2
pc = postal_code
em = email
mp = mobile
da = date_added
dm = date_mpdified
ct = city
st = state
cn = country
sn = service_name
sd = service.description
spsd = service_provider_service.description
ecid = end_customer_service.end_customer_id

Operators:
= (default)
lt = <
gt = >
|| = OR
ne = <>
lte = <=
gte = >=
lk = LIKE
btw = BETWEEN
in = IN

Example:
/api/v1/endCustomer?fn=benny&ln=pooper
*/
func ParseQueryStringParams(qs string) ([]string, []interface{}) {
	var wc []string
	var vals []interface{}

	c := 0
	for _, val := range sliceAndConditionQueryString(qs) {
		kv := strings.Split(val, "=")
		dbField := kv[0]
		v := kv[1]
		v_op := ""
		dbf_op := "AND"

		if strings.Contains(dbField, "||") && c > 0 {
			dbf_op = "OR"
			dbField = strings.Replace(dbField, "||", "", 1)
		}

		dbField = translateDbFieldName(dbField)

		//Multiple query string parameters.
		//If the dbField is not in the WHERE clause add the operator
		//and complete the loop, else continue to next dbField.
		if c > 0 && !strings.Contains(strings.Join(wc, ""), dbField) {
			wc = append(wc, dbf_op)
		} else if c > 0 {
			continue
		}

		//There is an operator in the value
		if len(strings.Split(v, ":")) > 1 {
			v_op = strings.Split(v, ":")[0]
			v = strings.Split(v, ":")[1]

			// There are multiple values for a single key
			// think IN or BETWEEN operators
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

		wc = append(wc, buildWhereClause(dbField, v_op, vals))
		c++
	}

	return wc, vals
}

//Creates a slice from the query string and
//unescapes any URL encoding via the net/url package
func sliceAndConditionQueryString(qs string) []string {
	slc := strings.Split(qs, "&")

	for i, val := range slc {
		slc[i], _ = url.QueryUnescape(val)
	}

	return slc
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
	return "="
}

func buildWhereClause(dbf string, v_op string, vals []interface{}) string {
	v_op = translateOperator(v_op)
	switch o := v_op; o {
	case "BETWEEN":
		return buildBetweenWhereClause(dbf, o)
	case "IN":
		return buildInWhereClause(dbf, o, vals)
	}
	return dbf + " " + v_op + " ?"
}

func buildBetweenWhereClause(dbf string, v_op string) string {
	return dbf + " " + v_op + " ? AND ?"
}

func buildInWhereClause(dbf string, v_op string, vals []interface{}) string {
	bindvars := ""
	for i := 0; i < len(vals); i++ {
		bindvars += "?,"
	}
	bindvars = strings.TrimSuffix(bindvars, ",")
	return dbf + " " + v_op + " (" + bindvars + ")"
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
		return "DATE_FORMAT(date_added, '%Y-%m-%d')"
	case "dm":
		return "DATE_FORMAT(date_modified, '%Y-%m-%d')"
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
	case "spsd":
		return "service_provider_service.description"
	case "ecid":
		return "end_customer_service.end_customer_id"
	}
	return ""
}

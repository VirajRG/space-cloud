package sql

import (
	"errors"

	goqu "gopkg.in/doug-martin/goqu.v4"
)

func generator(find map[string]interface{}) goqu.Expression {
	orTemp, isOr := find["$or"]
	if isOr {
		orArray := orTemp.([]interface{})
		orFinalArray := []goqu.Expression{}

		for _, item := range orArray {
			exp := generator(item.(map[string]interface{}))
			orFinalArray = append(orFinalArray, exp)
		}

		return goqu.Or(orFinalArray...)
	}

	array := []goqu.Expression{}
	for k, v := range find {
		val, isObj := v.(map[string]interface{})
		if isObj {
			for k2, v2 := range val {
				switch k2 {
				case "$eq":
					array = append(array, goqu.I(k).Eq(v2))
				case "$ne":
					array = append(array, goqu.I(k).Neq(v2))

				case "$gt":
					array = append(array, goqu.I(k).Gt(v2))

				case "$gte":
					array = append(array, goqu.I(k).Gte(v2))

				case "$lt":
					array = append(array, goqu.I(k).Lt(v2))

				case "$lte":
					array = append(array, goqu.I(k).Lte(v2))

				case "$in":
					array = append(array, goqu.I(k).In(v2))

				case "$nin":
					array = append(array, goqu.I(k).NotIn(v2))
				}
			}
		} else {
			array = append(array, goqu.I(k).Eq(v))
		}
	}
	return goqu.And(array...)
}
func generateWhereClause(q *goqu.Dataset, find map[string]interface{}) (query *goqu.Dataset, err error) {
	query = q
	err = nil
	if len(find) == 0 {
		return
	}
	exp := generator(find)
	query = query.Where(exp)
	return
}

func generateRecord(temp interface{}) (goqu.Record, error) {
	insertObj, ok := temp.(map[string]interface{})
	if !ok {
		return nil, errors.New("Incorrect insert object provided")
	}

	record := make(goqu.Record, len(insertObj))
	for k, v := range insertObj {
		record[k] = v
	}
	return record, nil
}

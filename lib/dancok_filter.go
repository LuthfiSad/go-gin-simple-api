package lib

import (
	"strings"

	"github.com/bytesaddict/dancok"
)

// FilterParam represents a single filter condition
type FilterParam struct {
	Field    string
	Value    any
	Operator string
}

// FilterParams is a collection of filter parameters
type FilterParams []FilterParam

// Define filter operators
const (
	IsEqual        = "equals"
	IsNotEqual     = "notequals"
	IsGreaterThan  = "greaterthan"
	IsGreaterEqual = "greaterthanorequal"
	IsLessThan     = "lessthan"
	IsLessEqual    = "lessthanorequal"
	IsContain      = "contains"
	IsBeginWith    = "startswith"
	IsEndWith      = "endswith"
	IsIn           = "in"
)

// Convert FilterParams to Dancok FilterDescriptors
func ConvertToDancokFilters(params FilterParams) []dancok.FilterDescriptor {
	filters := make([]dancok.FilterDescriptor, len(params))

	for i, param := range params {
		filters[i] = dancok.FilterDescriptor{
			FieldName: param.Field,
			Operator:  getDancokOperator(param.Operator),
			Value:     param.Value,
		}
	}

	return filters
}

// Parse filter string into FilterParams
func ParseFilterString(filterStr string) FilterParams {
	if filterStr == "" {
		return FilterParams{}
	}

	params := FilterParams{}
	parts := strings.Split(filterStr, "&")

	for _, part := range parts {
		fields := strings.Split(part, ":")
		if len(fields) == 3 {
			params = append(params, FilterParam{
				Field:    fields[0],
				Value:    fields[1],
				Operator: fields[2],
			})
		}
	}

	return params
}

// Convert string operator to dancok.Operator
func getDancokOperator(op string) dancok.Operator {
	switch op {
	case IsEqual:
		return dancok.IsEqual
	case IsNotEqual:
		return dancok.IsNotEqual
	case IsGreaterThan:
		return dancok.IsMoreThan
	case IsGreaterEqual:
		return dancok.IsMoreThanOrEqual
	case IsLessThan:
		return dancok.IsLessThan
	case IsLessEqual:
		return dancok.IsLessThanOrEqual
	case IsContain:
		return dancok.IsContain
	case IsBeginWith:
		return dancok.IsBeginWith
	case IsEndWith:
		return dancok.IsEndWith
	case IsIn:
		return dancok.IsIn
	default:
		return dancok.IsEqual
	}
}

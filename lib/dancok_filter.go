package lib

import (
	"go-gin-simple-api/utils"
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
	parts := strings.Split(filterStr, "|")

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

func ValidateFilterGeneric[T any](filter string, data T) map[string]string {
	// Konversi data generik ke map
	mapData, err := utils.ConvertToMap(data)
	if err != nil {
		return map[string]string{"convert_error": "Gagal mengkonversi data ke format yang dapat divalidasi"}
	}

	return ValidateFilter(filter, mapData)
}

func ValidateFilter(filter string, data map[string]interface{}) map[string]string {
	errors := make(map[string]string)

	// Parse filter string
	filterParams := ParseFilterString(filter)

	// Jika tidak ada filter, anggap valid
	if len(filterParams) == 0 {
		return errors // Kosong berarti tidak ada error
	}

	// Coba jalankan setiap filter
	for _, param := range filterParams {
		// Cek apakah field ada dalam data
		if _, exists := data[param.Field]; !exists {
			errors[param.Field] = "Field tidak ditemukan dalam data"
			continue
		}

		// descriptor := dancok.FilterDescriptor{
		// 	FieldName: param.Field,
		// 	Operator:  getDancokOperator(param.Operator),
		// 	Value:     param.Value,
		// }

		// err := validateFilterDescriptorWithError(descriptor, data)
		// if err != "" {
		// 	errors[param.Field] = err
		// }
	}

	return errors
}

// validateFilterDescriptorWithError mengembalikan string error daripada boolean
// func validateFilterDescriptorWithError(filter dancok.FilterDescriptor, data map[string]interface{}) string {
// 	fieldValue, exists := data[filter.FieldName]
// 	if !exists {
// 		return "Field tidak ditemukan"
// 	}

// 	switch filter.Operator {
// 	case dancok.IsEqual:
// 		if fieldValue != filter.Value {
// 			return fmt.Sprintf("Nilai %v tidak sama dengan %v", fieldValue, filter.Value)
// 		}
// 	case dancok.IsNotEqual:
// 		if fieldValue == filter.Value {
// 			return fmt.Sprintf("Nilai %v sama dengan %v", fieldValue, filter.Value)
// 		}
// 	case dancok.IsContain:
// 		str, ok := fieldValue.(string)
// 		if !ok {
// 			return "Field bukan bertipe string"
// 		}
// 		valueStr, ok := filter.Value.(string)
// 		if !ok {
// 			return "Nilai filter bukan bertipe string"
// 		}
// 		if !strings.Contains(str, valueStr) {
// 			return fmt.Sprintf("'%s' tidak mengandung '%s'", str, valueStr)
// 		}
// 	case dancok.IsBeginWith:
// 		str, ok := fieldValue.(string)
// 		if !ok {
// 			return "Field bukan bertipe string"
// 		}
// 		valueStr, ok := filter.Value.(string)
// 		if !ok {
// 			return "Nilai filter bukan bertipe string"
// 		}
// 		if !strings.HasPrefix(str, valueStr) {
// 			return fmt.Sprintf("'%s' tidak diawali dengan '%s'", str, valueStr)
// 		}
// 	case dancok.IsEndWith:
// 		str, ok := fieldValue.(string)
// 		if !ok {
// 			return "Field bukan bertipe string"
// 		}
// 		valueStr, ok := filter.Value.(string)
// 		if !ok {
// 			return "Nilai filter bukan bertipe string"
// 		}
// 		if !strings.HasSuffix(str, valueStr) {
// 			return fmt.Sprintf("'%s' tidak diakhiri dengan '%s'", str, valueStr)
// 		}
// 	// Tambahkan case lain sesuai kebutuhan
// 	default:
// 		return fmt.Sprintf("Operator %v tidak didukung", filter.Operator)
// 	}

// 	return "" // String kosong berarti tidak ada error
// }

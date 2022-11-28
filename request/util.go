package request

import "strconv"

func stringToAny(value string, out any) bool {
	if out == nil {
		return false
	}

	switch a := any(out).(type) {
	case *bool:
		result, err := strconv.ParseBool(value)
		if err != nil {
			return false
		}
		*a = result
	case *string:
		*a = value
	case *int, *int8, *int16, *int32, *int64:

		switch b := a.(type) {
		case *int:
			result, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return false
			}
			*b = int(result)
		case *int8:
			result, err := strconv.ParseInt(value, 10, 8)
			if err != nil {
				return false
			}
			*b = int8(result)
		case *int16:
			result, err := strconv.ParseInt(value, 10, 16)
			if err != nil {
				return false
			}
			*b = int16(result)
		case *int32:
			result, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return false
			}
			*b = int32(result)
		case *int64:
			result, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return false
			}
			*b = result
		}

	case *uint, *uint8, *uint16, *uint32, *uint64:

		switch b := a.(type) {
		case *uint:
			result, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return false
			}
			*b = uint(result)
		case *uint8:
			result, err := strconv.ParseUint(value, 10, 8)
			if err != nil {
				return false
			}
			*b = uint8(result)
		case *uint16:
			result, err := strconv.ParseUint(value, 10, 16)
			if err != nil {
				return false
			}
			*b = uint16(result)
		case *uint32:
			result, err := strconv.ParseUint(value, 10, 32)
			if err != nil {
				return false
			}
			*b = uint32(result)
		case *uint64:
			result, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return false
			}
			*b = result
		}
	case *float32:
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		*a = float32(result)
	case *float64:
		result, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return false
		}
		*a = result
	default:
		return false
	}

	return true
}

package common

func getValueWithType(m map[string]interface{}, key string) (interface{}, bool) {
	val, ok := m[key]
	if !ok {
		return nil, false
	}
	return val, true
}

func GetIntValue(m map[string]interface{}, key string) (int, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return 0, false
	}
	if iVal, ok := val.(int); ok {
		return iVal, true
	}
	return 0, false
}

func GetFloatValue(m map[string]interface{}, key string) (float64, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return 0, false
	}
	if fVal, ok := val.(float64); ok {
		return fVal, true
	}
	return 0, false
}

func GetStringValue(m map[string]interface{}, key string) (string, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return "", false
	}
	if sVal, ok := val.(string); ok {
		return sVal, true
	}
	return "", false
}

func GetBoolValue(m map[string]interface{}, key string) (bool, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return false, false
	}
	if bVal, ok := val.(bool); ok {
		return bVal, true
	}
	return false, false
}

func GetStringSlice(m map[string]interface{}, key string) ([]string, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return nil, false
	}
	if ssVal, ok := val.([]string); ok {
		return ssVal, true
	}
	return nil, false
}

func GetInterfaceSlice(m map[string]interface{}, key string) ([]interface{}, bool) {
	val, ok := getValueWithType(m, key)
	if !ok {
		return nil, false
	}
	if ssVal, ok := val.([]interface{}); ok {
		return ssVal, true
	}
	return nil, false
}

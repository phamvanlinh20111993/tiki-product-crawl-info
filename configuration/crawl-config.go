package configuration

import (
	"reflect"
	"slices"
)

// TODO we can use yaml to extract the value from yaml file but i challenge my self a bit
// i just want to deeply learn reflection in golang

func setFieldValue(field reflect.Value, value interface{}) {
	switch field.Kind() {

	case reflect.String:
		if v, ok := value.(string); ok {
			field.SetString(v)
		}

	case reflect.Bool:
		if v, ok := value.(bool); ok {
			field.SetBool(v)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := value.(type) {
		case int:
			field.SetInt(int64(v))
		case int64:
			field.SetInt(v)
		case float64:
			field.SetInt(int64(v))
		}

	case reflect.Float32, reflect.Float64:
		if v, ok := value.(float64); ok {
			field.SetFloat(v)
		}

	default:
		getLogger().Info("func setFieldValue() Unsupported type:", field.Kind())
	}
}

// TODO need to handle more type if needed
func handleSliceStructElement(fieldType reflect.StructField, fieldValue reflect.Value, value interface{}) {
	if fieldType.Type.Kind() != reflect.Slice {
		return
	}

	v, ok := value.([]interface{})
	if !ok {
		getLogger().Debug("values get from map not mapping for type ", fieldType.Type.Kind())
		return
	}

	if !slices.Contains([]reflect.Kind{reflect.String, reflect.Int, reflect.Float32, reflect.Float64, reflect.Bool}, fieldType.Type.Elem().Kind()) {
		getLogger().Debug("Not support for this type, please change the config data field ", fieldType.Type.Elem().Kind())
		return
	}

	if fieldType.Type.Elem().Kind() == reflect.String {
		var strArrValue []string
		for i := 0; i < len(v); i++ {
			str, ok := v[i].(string)
			if ok {
				strArrValue = append(strArrValue, str)
			} else {
				getLogger().Debug("Failed to convert to string")
				return
			}
		}
		fieldValue.Set(reflect.ValueOf(strArrValue))
	}

	if fieldType.Type.Elem().Kind() == reflect.Int {
		var intArrValue []int
		for i := 0; i < len(v); i++ {
			iV, ok := v[i].(int)
			if ok {
				intArrValue = append(intArrValue, iV)
			} else {
				getLogger().Debug("Failed to convert to int")
				return
			}
		}
		fieldValue.Set(reflect.ValueOf(intArrValue))
	}

	if fieldType.Type.Elem().Kind() == reflect.Float64 {
		var floatArrValue []float64
		for i := 0; i < len(v); i++ {
			f64, ok := v[i].(float64)
			if ok {
				floatArrValue = append(floatArrValue, f64)
			} else {
				getLogger().Debug("Failed to convert to float64")
				return
			}
		}
		fieldValue.Set(reflect.ValueOf(floatArrValue))
	}

	if fieldType.Type.Elem().Kind() == reflect.Float32 {
		var floatArrValue []float32
		for i := 0; i < len(v); i++ {
			f, ok := v[i].(float32)
			if ok {
				floatArrValue = append(floatArrValue, f)
			} else {
				getLogger().Debug("Failed to convert to float32")
				return
			}
		}
		fieldValue.Set(reflect.ValueOf(floatArrValue))
	}

	if fieldType.Type.Elem().Kind() == reflect.Bool {
		var boolArrValue []bool
		for i := 0; i < len(v); i++ {
			str, ok := v[i].(bool)
			if ok {
				boolArrValue = append(boolArrValue, str)
			} else {
				getLogger().Debug("Failed to convert to bool")
				return
			}
		}
		fieldValue.Set(reflect.ValueOf(boolArrValue))
	}

}

// TODO will need to improve when have new type to config
func getStructTypeConfig(target interface{}) {
	t := reflect.TypeOf(target)
	v := reflect.ValueOf(target)
	if t.Kind() != reflect.Ptr {
		getLogger().Error("Must be a pointer to a struct type:", t)
	}

	t = t.Elem()
	v = v.Elem()
	if t.Kind() != reflect.Struct {
		getLogger().Error("Provided type is not a struct")
	}
	configuration := loadConfiguration()
	for i := 0; i < t.NumField(); i++ {
		fieldType := t.Field(i)
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			getLogger().Error("Can't set field ", t.Field(i).Name)
			continue
		}
		confKey := fieldType.Tag.Get(DescriptionConfigStruct)

		// Handle nested struct recursively
		if fieldType.Type.Kind() == reflect.Struct {
			nested := fieldValue.Addr().Interface()
			getStructTypeConfig(nested)
			continue
		}

		if confKey == "" {
			getLogger().Debug("confKey is not filled yet, can not apply to type: ", confKey)
			continue
		}
		rawValue, ok := configuration[confKey]
		if !ok {
			getLogger().Debug("value for confKey in map is empty: ", confKey)
			continue
		}

		if fieldType.Type.Kind() == reflect.Slice {
			handleSliceStructElement(fieldType, fieldValue, rawValue)
			continue
		}

		setFieldValue(fieldValue, rawValue)
	}
}

func getTikiPageConfig() TikiPageConfig {
	pageConfig := &TikiPageConfig{}
	getStructTypeConfig(pageConfig)
	return *pageConfig
}

func getOpenSearchConfig() OpenSearchConfig {
	openSearchConfig := &OpenSearchConfig{}
	getStructTypeConfig(openSearchConfig)
	return *openSearchConfig
}

func getLoggerConfig() LoggerConfig {
	loggerConfig := &LoggerConfig{}
	getStructTypeConfig(loggerConfig)
	return *loggerConfig
}

func getPostgresConfig() PostgresConfig {
	postgresConfig := &PostgresConfig{}
	getStructTypeConfig(postgresConfig)
	return *postgresConfig
}

func getFileConfig() FileConfig {
	fileConfig := &FileConfig{}
	getStructTypeConfig(fileConfig)
	return *fileConfig
}

var GetTikiPageConfig = getTikiPageConfig
var GetOpenSearchConfig = getOpenSearchConfig
var GetLoggerConfig = getLoggerConfig
var GetPostgresConfig = getPostgresConfig
var GetFileConfig = getFileConfig

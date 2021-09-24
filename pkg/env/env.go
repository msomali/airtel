package env

import (
	"os"
	"strconv"
)

//const (
//	tagName = "env"
//)
//
//type (
//	Config struct {
//		Command     string `env:"name"`
//		Port     int    `env:"port"`
//		Password string `env:"password"`
//	}
//)
//
//// Unmarshal read environmental variables and set them into a struct
////	type Config struct {
////		Command     string `env:"name"`
////		Port     int    `env:"port"`
////		Password string `env:"password"`
////	}
////  conf := new(Config)
//// For a struct like this Unmarshal("POSTGRES",conf) will look for
//// POSTGRES_NAME, POSTGRES_PORT and POSTGRES_PASSWORD
//func Unmarshal(prefix string, v interface{}) error {
//	return nil
//}
//
//func Set(prefix string, v interface{}) error {
//
//	// TypeOf returns the reflection Type that represents the dynamic type of variable.
//	// If variable is a nil interface value, TypeOf returns nil.
//	t := reflect.TypeOf(v)
//
//	// Get the type and kind of our user variable
//	fmt.Println("Type:", t.Command())
//	fmt.Println("Kind:", t.Kind())
//
//	// Iterate over all available fields and read the tag value
//	for i := 0; i < t.NumField(); i++ {
//		// Get the field, returns https://golang.org/pkg/reflect/#StructField
//		field := t.Field(i)
//
//		// Get the field tag value
//		tag := field.Tag.Get(tagName)
//
//		fmt.Printf("%d. %v (%v), tag: '%v'\n", i+1, field.Command, field.Type.Command(), tag)
//	}
//
//	return nil
//}

func Get(key string, defaultValue interface{}) interface{} {
	var strValue string
	if strValue = os.Getenv(key); strValue == "" {
		return defaultValue
	}

	switch defaultValue.(type) {
	case string:

		return strValue

	case int, int64, int8, int16, int32:
		retValue, err := strconv.Atoi(strValue)
		if err != nil {
			return defaultValue
		}

		return retValue

	case bool:
		retValue, err := strconv.ParseBool(strValue)
		if err != nil {

			return defaultValue
		}

		return retValue

	case float32:
		retValue, err := strconv.ParseFloat(strValue, 32)
		if err != nil {

			return defaultValue
		}

		return retValue

	case float64:
		retValue, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			return defaultValue
		}

		return retValue

	default:
		return strValue
	}
}

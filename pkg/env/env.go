/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package env

import (
	"fmt"
	"os"
	"strconv"
)

func Get(key string, defaultValue interface{}) interface{} {
	var strValue string
	if strValue = os.Getenv(key); strValue == "" {
		return defaultValue
	}

	switch defaultValue.(type) {
	case string:

		return strValue

	case int:
		retValue, err := strconv.Atoi(strValue)
		if err != nil {
			return defaultValue
		}

		return retValue

	case int64:
		retValue, err := strconv.ParseInt(strValue, 10, 64)
		if err != nil {
			return defaultValue
		}

		return retValue

	case int32:
		retValue, err := strconv.ParseInt(strValue, 10, 32)
		if err != nil {
			return defaultValue
		}

		return retValue

	case int8:
		retValue, err := strconv.ParseInt(strValue, 10, 8)
		if err != nil {
			return defaultValue
		}

		return retValue

	case int16:
		retValue, err := strconv.ParseInt(strValue, 10, 16)
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

func String(key string, defaultValue interface{}) string {
	i := Get(key, defaultValue)
	return fmt.Sprintf("%v\n", i)
}

func Bool(key string, defaultValue interface{}) bool {
	i := String(key, defaultValue)
	parseBool, _ := strconv.ParseBool(i)
	return parseBool
}

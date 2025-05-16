package config_handler

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_CONFIG_PATH = "../resources/config.conf"
)

type ExpectedType int

const (
	TYPE_STRING ExpectedType = iota
	TYPE_INT
	TYPE_BOOL
	TYPE_FLOAT
	TYPE_TIME
)

var expected_type_map = map[ExpectedType]reflect.Type{
	TYPE_STRING: reflect.TypeOf(""),
	TYPE_INT:    reflect.TypeOf(0),
	TYPE_BOOL:   reflect.TypeOf(true),
	TYPE_FLOAT:  reflect.TypeOf(0.0),
	TYPE_TIME:   reflect.TypeOf(time.Duration(0)),
}

func LoadConsts(path string) (map[string]any, error) {
	result := make(map[string]any)
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.Trim(strings.TrimSpace(parts[0]), `"`)
		value := strings.Trim(strings.TrimSpace(parts[1]), `"`)

		if int64_val, err := strconv.ParseInt(value, 10, 64); err == nil && strings.Contains(key, "INTERVAL") {
			result[key] = time.Duration(int64_val)
		} else if bool_val, err := strconv.ParseBool(value); err == nil {
			result[key] = bool_val
		} else if int_val, err := strconv.Atoi(value); err == nil {
			result[key] = int_val
		} else {
			result[key] = value
		}
	}
	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func LoadReceivers(path string) ([]string, error) {
	result := make([]string, 0)
	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if strings.ToLower(strings.Trim(strings.TrimSpace(parts[0]), `"`)) == "send_to" {
			result = append(result, strings.Trim(strings.TrimSpace(parts[1]), `"`))
		}
	}
	return result, nil
}

func GetData(m map[string]any, key string, expected_type ExpectedType) (any, error) {
	val, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("[ERROR] Key \"%s\" not found", key)
	}
	ex_type, ok := expected_type_map[expected_type]
	if !ok {
		return nil, fmt.Errorf("[ERROR] Unknown expected type for key \"%s\"", key)
	}
	ac_type := reflect.TypeOf(val)
	if ac_type != ex_type {
		return nil, fmt.Errorf("[ERROR] Key \"%s\", actual type = \"%s\", expected type = \"%s\"", key, ac_type, ex_type)
	}
	return val, nil
}

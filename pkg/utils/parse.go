package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseIntSafe convierte un string a int eliminando espacios en blanco.
// Retorna un error descriptivo si la conversión falla.
func ParseIntSafe(valStr string) (int, error) {
	cleanVal := strings.TrimSpace(valStr)
	if cleanVal == "" {
		return 0, fmt.Errorf("value is empty, expected a number")
	}

	val, err := strconv.Atoi(cleanVal)
	if err != nil {
		return 0, err
	}
	return val, nil
}

// ParseInt64Safe convierte un string a int64 (ideal para el NIT).
// Retorna un error descriptivo si la conversión falla.
func ParseInt64Safe(valStr string) (int64, error) {
	cleanVal := strings.TrimSpace(valStr)
	if cleanVal == "" {
		return 0, fmt.Errorf("value is empty, expected a number")
	}

	val, err := strconv.ParseInt(cleanVal, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

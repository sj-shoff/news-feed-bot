package botkit

import (
	"encoding/json"
)

// ParseJSON распарсит json в структуру с помощью стандартного json пакета
func ParseJSON[T any](src string) (T, error) {
	var args T

	if err := json.Unmarshal([]byte(src), &args); err != nil {
		return *(new(T)), err
	}

	return args, nil
}

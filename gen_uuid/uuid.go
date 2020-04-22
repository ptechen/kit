package gen_uuid

import (
	"encoding/json"
	"github.com/google/uuid"
	)

func GenUuid(params interface{}) (uuidStr string, err error) {
	data, err := json.Marshal(params)
	if err != nil {
		return uuidStr, err
	}
	uuidStr = uuid.NewMD5(uuid.UUID{}, data).String()
	return uuidStr, err
}
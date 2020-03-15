package utils

import (
	"encoding/json"
	"os"
)

func SaveToJson(fpath string, pcm []float64) error {
	out, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer out.Close()

	data := floatPcmToInt(pcm)
	jsonData, _ := json.Marshal(&data)

	out.Write([]byte(jsonData))

	return nil
}

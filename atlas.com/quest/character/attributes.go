package character

type attributes struct {
	JobId uint16 `json:"jobId"`
	MapId uint32 `json:"mapId"`
	Level byte   `json:"level"`
	Fame  int16  `json:"fame"`
	Meso  uint32 `json:"meso"`
}

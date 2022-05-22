package character

type Model struct {
	id    uint32
	jobId uint16
	mapId uint32
	level byte
	fame  int16
	meso  uint32
}

func (a Model) Id() uint32 {
	return a.id
}

func (a Model) JobId() uint16 {
	return a.jobId
}

func (a Model) MapId() uint32 {
	return a.mapId
}

func (a Model) Level() byte {
	return a.level
}

func (a Model) Fame() int16 {
	return a.fame
}

func (a Model) Meso() uint32 {
	return a.meso
}

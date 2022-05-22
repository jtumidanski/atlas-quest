package character

import (
	"atlas-quest/model"
	"atlas-quest/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ByIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(id uint32) model.Provider[Model] {
	return func(id uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(id), makeModel)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		return ByIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(ca requests.DataBody[attributes]) (Model, error) {
	cid, err := strconv.ParseUint(ca.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := ca.Attributes
	r := Model{
		id:    uint32(cid),
		jobId: att.JobId,
		mapId: att.MapId,
		level: att.Level,
		fame:  att.Fame,
		meso:  att.Meso,
	}
	return r, nil
}

type Criteria func(Model) bool

func MeetsCriteria(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, criteria ...Criteria) bool {
	return func(characterId uint32, criteria ...Criteria) bool {
		c, err := GetById(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d for criteria check.", characterId)
			return false
		}
		for _, check := range criteria {
			if ok := check(c); !ok {
				return false
			}
		}
		return true
	}
}

func IsJob(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, options []uint16) bool {
	return func(characterId uint32, options []uint16) bool {
		return MeetsCriteria(l, span)(characterId, IsJobCriteria(options))
	}
}

func IsJobCriteria(options []uint16) Criteria {
	return func(c Model) bool {
		for _, id := range options {
			if id == c.JobId() {
				return true
			}
		}
		return false
	}
}

func InMap(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, mapId uint32) bool {
	return func(characterId uint32, mapId uint32) bool {
		return MeetsCriteria(l, span)(characterId, InMapCriteria(mapId))
	}
}

func InMapCriteria(mapId uint32) Criteria {
	return func(c Model) bool {
		return c.MapId() == mapId
	}
}

func HasItems(_ logrus.FieldLogger, _ opentracing.Span) func(characterId uint32, items map[uint32]uint32) bool {
	return func(characterId uint32, items map[uint32]uint32) bool {
		//TODO
		return false
	}
}

func IsMinimalLevel(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, level byte) bool {
	return func(characterId uint32, level byte) bool {
		return MeetsCriteria(l, span)(characterId, MinimalLevelCriteria(level))
	}
}

func MinimalLevelCriteria(level byte) Criteria {
	return func(c Model) bool {
		return c.Level() >= level
	}
}

func IsLevel(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, level byte) bool {
	return func(characterId uint32, level byte) bool {
		return MeetsCriteria(l, span)(characterId, IsLevelCriteria(level))
	}
}

func IsLevelCriteria(level byte) Criteria {
	return func(c Model) bool {
		return c.Level() == level
	}
}

func IsPopularityLevel(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, popularity int16) bool {
	return func(characterId uint32, popularity int16) bool {
		return MeetsCriteria(l, span)(characterId, MinimalPopularityCriteria(popularity))
	}
}

func MinimalPopularityCriteria(popularity int16) Criteria {
	return func(c Model) bool {
		return c.Fame() >= popularity
	}
}

func IsMaximalLevel(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, level byte) bool {
	return func(characterId uint32, level byte) bool {
		return MeetsCriteria(l, span)(characterId, MaximalLevelCriteria(level))
	}
}

func MaximalLevelCriteria(level byte) Criteria {
	return func(c Model) bool {
		return level >= c.Level()
	}
}

func IsMorphed(_ logrus.FieldLogger, _ opentracing.Span) func(characterId uint32, morph uint32) bool {
	return func(characterId uint32, morph uint32) bool {
		//TODO
		return false
	}
}

func HasSkill(_ logrus.FieldLogger, _ opentracing.Span) func(characterId uint32, skills map[uint32]uint32) bool {
	return func(characterId uint32, skills map[uint32]uint32) bool {
		//TODO
		return false
	}
}

func HasBuff(_ logrus.FieldLogger, _ opentracing.Span) func(characterId uint32, buffId int) bool {
	return func(characterId uint32, buffId int) bool {
		//TODO
		return false
	}
}

func GetMonsterBook(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (map[uint32]uint32, error) {
	return func(characterId uint32) (map[uint32]uint32, error) {
		//TODO
		return make(map[uint32]uint32), nil
	}
}

func MonsterBookCount(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) uint32 {
	return func(characterId uint32) uint32 {
		mb, err := GetMonsterBook(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve monster book for character %d.", characterId)
			return 0
		}
		return uint32(len(mb))
	}
}

func HasMinimalMeso(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, meso uint32) bool {
	return func(characterId uint32, meso uint32) bool {
		return MeetsCriteria(l, span)(characterId, MinimalMesoCriteria(meso))
	}
}

func MinimalMesoCriteria(meso uint32) Criteria {
	return func(c Model) bool {
		return c.Meso() >= meso
	}
}

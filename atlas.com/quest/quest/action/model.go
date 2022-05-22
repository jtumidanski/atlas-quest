package action

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	TypeExperience      = "EXP"
	TypeMoney           = "MESO"
	TypeItem            = "ITEM"
	TypeSkill           = "SKILL"
	TypeNextQuest       = "NEXT_QUEST"
	TypePopularity      = "POP"
	TypeBuffItemId      = "BUFF"
	TypePetSkill        = "PET_SKILL"
	TypeNo              = "NO"
	TypeYes             = "YES"
	TypeNPC             = "NPC"
	TypeMinimumLevel    = "MIN_LEVEL"
	TypeNormalAutoStart = "NORMAL_AUTO_START"
	TypePetTameness     = "PET_TAMENESS"
	TypePetSpeed        = "PET_SPEED"
	TypeInfo            = "INFO"
	Type0               = "ZERO"
)

type Model struct {
	theType Type
	check   CheckFunc
	run     RunFunc
}

type Type string

type CheckFunc func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, extSelection int) bool

type RunFunc func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32, extSelection int)

func (m Model) Type() Type {
	return m.theType
}

func (m Model) Check() CheckFunc {
	return m.check
}

func (m Model) Run() RunFunc {
	return m.run
}

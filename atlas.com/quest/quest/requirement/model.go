package requirement

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	TypeJob                  = "JOB"
	TypeQuest                = "OTHER_QUEST"
	TypeItem                 = "ITEM"
	TypeMinimumLevel         = "MIN_LEVEL"
	TypeMaximumLevel         = "MAX_LEVEL"
	TypeEndDate              = "END_DATE"
	TypeMob                  = "MOB"
	TypeNPC                  = "NPC"
	TypeFieldEnter           = "FIELD_ENTER"
	TypeInterval             = "INTERVAL"
	TypeStartScript          = "SCRIPT"
	TypeEndScript            = "SCRIPT"
	TypePet                  = "PET"
	TypePetTamenessMinimum   = "MIN_PET_TAMENESS"
	TypeMonsterBook          = "MONSTER_BOOK"
	TypeNormalAutoStart      = "NORMAL_AUTO_START"
	TypeInfoNumber           = "INFO_NUMBER"
	TypeInfoEx               = "INFO_EX"
	TypeQuestComplete        = "COMPLETED_QUEST"
	TypeStart                = "START"
	TypeDayByDay             = "DAY_BY_DAY"
	TypeMoney                = "MESO"
	TypeBuff                 = "BUFF"
	TypeExceptBuff           = "EXCEPT_BUFF"
	TypeEquipAllNeed         = "EQUIP_ALL_NEED"
	TypeEquipSelectNeed      = "EQUIP_SELECT_NEED"
	TypeSkill                = "SKILL"
	TypeInfo                 = "INFO"
	TypeMonsterBookCard      = "MONSTER_BOOK_CARD"
	TypeWorldMin             = "WORLD_MIN"
	TypeWorldMax             = "WORLD_MAX"
	TypeMorph                = "MORPH"
	TypePopularity           = "POPULARITY"
	TypeEndMeso              = "END_MESO"
	TypeLevel                = "LEVEL"
	TypePartyQuestS          = "PARTY_QUEST_S"
	TypeUserInteract         = "USER_INTERACT"
	TypePetRecallLimit       = "PET_RECALL_LIMIT"
	TypePetAutoSpeakingLimit = "PET_AUTO_SPEAKING_LIMIT"
	TypeTamingMobLevelMin    = "TAMING_MOB_LEVEL_MIN"
)

type Type string

type CheckFunc func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool

type Model struct {
	typeString   Type
	relevantMobs []uint32
	check        CheckFunc
}

func (m Model) Type() Type {
	return m.typeString
}

func (m Model) RelevantMobs() []uint32 {
	return m.relevantMobs
}

func (m Model) Check() CheckFunc {
	return m.check
}

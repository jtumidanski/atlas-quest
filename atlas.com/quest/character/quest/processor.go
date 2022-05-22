package quest

import (
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func ForCharacter(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		//TODO
		return nil, nil
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, questId uint16) (Model, error) {
	return func(characterId uint32, questId uint16) (Model, error) {
		//TODO
		return Model{}, errors.New("not found")
	}
}

func HasMetMonsterRequirement(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, questId uint16, counts map[uint32]uint32) bool {
	return func(characterId uint32, questId uint16, counts map[uint32]uint32) bool {
		//TODD
		return false
	}
}

func QuestsByStatus(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, status string) ([]Model, error) {
	return func(characterId uint32, status string) ([]Model, error) {
		//TODO
		return nil, nil
	}
}

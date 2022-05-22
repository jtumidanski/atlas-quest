package quest

import (
	"atlas-quest/quest/action"
	"atlas-quest/quest/requirement"
)

type Model struct {
	id                   uint16
	name                 string
	parent               string
	timeLimit            uint32
	timeLimit2           uint32
	autoStart            bool
	autoPreComplete      bool
	autoComplete         bool
	repeatable           bool
	medalId              uint32
	startRequirements    map[requirement.Type]requirement.CheckFunc
	completeRequirements map[requirement.Type]requirement.CheckFunc
	startActions         map[action.Type]Action
	completeActions      map[action.Type]Action
	relevantMobs         []uint32
}

func (m *Model) Id() uint16 {
	return m.id
}

type ModelBuilder struct {
	id                   uint16
	name                 string
	parent               string
	timeLimit            uint32
	timeLimit2           uint32
	autoStart            bool
	autoPreComplete      bool
	autoComplete         bool
	repeatable           bool
	medalId              uint32
	startRequirements    map[requirement.Type]requirement.CheckFunc
	completeRequirements map[requirement.Type]requirement.CheckFunc
	startActions         map[action.Type]Action
	completeActions      map[action.Type]Action
	relevantMobs         []uint32
}

type Action struct {
	check action.CheckFunc
	run   action.RunFunc
}

func NewBuilder(id uint16) *ModelBuilder {
	return &ModelBuilder{
		id:                   id,
		startRequirements:    make(map[requirement.Type]requirement.CheckFunc),
		completeRequirements: make(map[requirement.Type]requirement.CheckFunc),
		startActions:         make(map[action.Type]Action),
		completeActions:      make(map[action.Type]Action),
		relevantMobs:         make([]uint32, 0),
	}
}

func (m *ModelBuilder) SetRepeatable(value bool) *ModelBuilder {
	m.repeatable = value
	return m
}

func (m *ModelBuilder) AddStartingRequirement(t requirement.Type, rcf requirement.CheckFunc) *ModelBuilder {
	m.startRequirements[t] = rcf
	return m
}

func (m *ModelBuilder) AddCompletionRequirement(t requirement.Type, rcf requirement.CheckFunc) *ModelBuilder {
	m.completeRequirements[t] = rcf
	return m
}

func (m *ModelBuilder) AddRelevantMob(id uint32) *ModelBuilder {
	m.relevantMobs = append(m.relevantMobs, id)
	return m
}

func (m *ModelBuilder) Build() Model {
	return Model{
		id:                   m.id,
		name:                 m.name,
		parent:               m.parent,
		timeLimit:            m.timeLimit,
		timeLimit2:           m.timeLimit2,
		autoStart:            m.autoStart,
		autoPreComplete:      m.autoPreComplete,
		autoComplete:         m.autoComplete,
		repeatable:           m.repeatable,
		medalId:              m.medalId,
		startRequirements:    m.startRequirements,
		completeRequirements: m.completeRequirements,
		startActions:         m.startActions,
		completeActions:      m.completeActions,
		relevantMobs:         m.relevantMobs,
	}
}

func (m *ModelBuilder) SetName(name string) {
	m.name = name
}

func (m *ModelBuilder) SetParent(parent string) {
	m.parent = parent
}

func (m *ModelBuilder) SetTimeLimit(limit uint32) {
	m.timeLimit = limit
}

func (m *ModelBuilder) SetTimeLimit2(limit2 uint32) {
	m.timeLimit2 = limit2
}

func (m *ModelBuilder) SetAutoStart(start bool) {
	m.autoStart = start
}

func (m *ModelBuilder) SetAutoPreComplete(complete bool) {
	m.autoPreComplete = complete
}

func (m *ModelBuilder) SetAutoComplete(complete bool) {
	m.autoComplete = complete
}

func (m *ModelBuilder) SetMedalItem(value uint32) {
	m.medalId = value
}

func (m *ModelBuilder) AddStartingAction(t action.Type, check action.CheckFunc, run action.RunFunc) {
	m.startActions[t] = Action{check: check, run: run}
}

func (m *ModelBuilder) AddCompletionAction(t action.Type, check action.CheckFunc, run action.RunFunc) {
	m.completeActions[t] = Action{check: check, run: run}
}

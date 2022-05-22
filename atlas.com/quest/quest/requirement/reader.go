package requirement

import (
	"atlas-quest/character"
	"atlas-quest/character/quest"
	"atlas-quest/xml"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func get(questId uint16, root xml.Noder, nodeName string) ([]Model, error) {
	questData, ok := root.(xml.Parent)
	if !ok {
		return nil, errors.New("invalid xml structure")
	}

	requirementsRoot, err := questData.ChildByName(nodeName)
	if err != nil {
		return nil, errors.New("invalid xml structure")
	}

	rootAsParent, ok := requirementsRoot.(xml.Parent)
	if !ok {
		return nil, errors.New("invalid xml structure")
	}

	results := make([]Model, 0)
	for _, req := range rootAsParent.Children() {
		reqType, err := getByWZName(req.Name())
		if err != nil {
			return nil, err
		}

		m := Model{typeString: reqType}
		if reqType == TypeMob {
			rms, err := getRelevantMobs(req)
			if err != nil {
				return nil, err
			}
			m.relevantMobs = rms
		}
		check, err := getCheckProducer(questId, reqType, req)()
		if err != nil {
			return nil, err
		}
		m.check = check
		results = append(results, m)
	}
	return results, nil
}

func getRelevantMobs(req xml.Noder) ([]uint32, error) {
	reqAsParent, ok := req.(xml.Parent)
	if !ok {
		return nil, errors.New("invalid xml structure")
	}

	results := make([]uint32, 0)
	for _, md := range reqAsParent.Children() {
		mob, ok := md.(xml.Parent)
		if !ok {
			return nil, errors.New("invalid xml structure")
		}
		mid, err := xml.GetInteger(mob, "id")
		if err != nil {
			return nil, err
		}
		results = append(results, uint32(mid))
	}
	return results, nil
}

type checkProducer func() (CheckFunc, error)

func getCheckProducer(questId uint16, rt Type, sr xml.Noder) checkProducer {
	switch rt {
	case TypeEndDate:
		return endDateRequirement(sr)
	case TypeJob:
		return jobRequirement(sr)
	case TypeQuest:
		return otherQuestRequirement(sr)
	case TypeFieldEnter:
		return fieldEnterRequirement(sr)
	case TypeInfoNumber:
		return infoNumberRequirement(sr)
	case TypeInfoEx:
		return infoExRequirement(sr)
	case TypeInterval:
		return intervalRequirement(questId, sr)
	case TypeQuestComplete:
		return questCompleteRequirement(sr)
	case TypeItem:
		return itemRequirement(sr)
	case TypeMaximumLevel:
		return maxLevelRequirement(sr)
	case TypeMoney:
		return mesoRequirement(sr)
	case TypeMinimumLevel:
		return minLevelRequirement(sr)
	case TypePetTamenessMinimum:
		return petTamenessMinimumRequirement(sr)
	case TypeMob:
		return monsterRequirement(questId, sr)
	case TypeMonsterBook:
		return monsterBookRequirement(sr)
	case TypeNPC:
		return npcRequirement(sr)
	case TypePet:
		return petRequirement(sr)
	case TypeBuff:
		return buffRequirement(sr)
	case TypeExceptBuff:
		return exceptBuffRequirement(sr)
	case TypeStartScript:
		return scriptRequirement(sr)
	//case TypeEndScript:
	//	return scriptRequirement(sr)
	case TypeEquipAllNeed:
		return validRequirementProducer(invalidCheck)
	case TypeEquipSelectNeed:
		return validRequirementProducer(invalidCheck)
	case TypeSkill:
		return skillRequirement(sr)
	case TypeInfo:
		return validRequirementProducer(invalidCheck)
	case TypeMonsterBookCard:
		return monsterBookCardRequirement(sr)
	case TypeNormalAutoStart:
		return normalAutoStartRequirement(sr)
	case TypeStart:
		return startRequirement(sr)
	case TypeDayByDay:
		return dayByDayRequirement(sr)
	case TypeWorldMin:
		return worldMinRequirement(sr)
	case TypeWorldMax:
		return worldMaxRequirement(sr)
	case TypeMorph:
		return morphRequirement(sr)
	case TypePopularity:
		return popularityRequirement(sr)
	case TypeEndMeso:
		return endMesoRequirement(sr)
	case TypeLevel:
		return levelRequirement(sr)
	case TypePartyQuestS:
		return partyQuestSRequirement(sr)
	case TypeUserInteract:
		return userInteractRequirement(sr)
	case TypePetRecallLimit:
		return petRecallLimitRequirement(sr)
	case TypePetAutoSpeakingLimit:
		return petAutoSpeakingLimitRequirement(sr)
	case TypeTamingMobLevelMin:
		return tamingMobLevelMinRequirement(sr)
	}
	return errorCheckProducer(errors.New("requirement type not found"))
}

func validCheck(_ logrus.FieldLogger, _ opentracing.Span, _ *gorm.DB) func(_ uint32, _ uint32) bool {
	return func(_ uint32, _ uint32) bool {
		return true
	}
}

func invalidCheck(_ logrus.FieldLogger, _ opentracing.Span, _ *gorm.DB) func(_ uint32, _ uint32) bool {
	return func(_ uint32, _ uint32) bool {
		return false
	}
}

func errorCheckProducer(err error) checkProducer {
	return func() (CheckFunc, error) {
		return invalidCheck, err
	}
}

func validRequirementProducer(cf CheckFunc) checkProducer {
	return func() (CheckFunc, error) {
		return cf, nil
	}
}

func tamingMobLevelMinRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func petAutoSpeakingLimitRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func petRecallLimitRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func userInteractRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func partyQuestSRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func levelRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkLevel(byte(val)))
}

func checkLevel(level byte) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.IsLevel(l, span)(characterId, level)
		}
	}
}

func endMesoRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func popularityRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkPopularity(int16(val)))
}

func checkPopularity(pop int16) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.IsPopularityLevel(l, span)(characterId, pop)
		}
	}
}

func morphRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkMorph(uint32(val)))
}

func checkMorph(morph uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.IsMorphed(l, span)(characterId, morph)
		}
	}
}

func worldMaxRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func worldMinRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func dayByDayRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func startRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func normalAutoStartRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func monsterBookCardRequirement(r xml.Noder) checkProducer {
	mins := make(map[uint32]uint32)
	mbrs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, mbr := range mbrs.Children() {
		mbd, ok := mbr.(xml.Parent)
		if !ok {
			continue
		}

		id, err := xml.GetInteger(mbd, "id")
		if err != nil {
			continue
		}
		min, err := xml.GetInteger(mbd, "min")
		if err != nil {
			continue
		}
		mins[uint32(id)] = uint32(min)
	}
	return validRequirementProducer(checkMinMonsterBookCard(mins))
}

func checkMinMonsterBookCard(mins map[uint32]uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			mb, err := character.GetMonsterBook(l, span)(characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve monster book for character %d. Assuming check fails.", characterId)
				return false
			}
			for id, qty := range mins {
				if val, ok := mb[id]; ok {
					if val >= qty {
						return true
					}
				}
			}
			return false
		}
	}
}

func skillRequirement(r xml.Noder) checkProducer {
	skills := make(map[uint32]uint32)
	srs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, sr := range srs.Children() {
		sd, ok := sr.(xml.Parent)
		if !ok {
			continue
		}

		id, err := xml.GetInteger(sd, "id")
		if err != nil {
			continue
		}
		acquire := xml.GetIntegerWithDefault(sd, "acquire", 0)
		skills[uint32(id)] = uint32(acquire)
	}
	return validRequirementProducer(checkSkills(skills))
}

func checkSkills(skills map[uint32]uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.HasSkill(l, span)(characterId, skills)
		}
	}
}

func scriptRequirement(_ xml.Noder) checkProducer {
	//TODO identify what this is supposed to be doing
	return validRequirementProducer(validCheck)
}

func exceptBuffRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromStringNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkBuffExcept(val * -1))
}

func checkBuffExcept(buffId int) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return !character.HasBuff(l, span)(characterId, buffId)
		}
	}
}

func buffRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromStringNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkBuff(val * -1))
}

func checkBuff(buffId int) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.HasBuff(l, span)(characterId, buffId)
		}
	}
}

func petRequirement(r xml.Noder) checkProducer {
	petIds := make([]uint32, 0)
	prs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, pr := range prs.Children() {
		pd, ok := pr.(xml.Parent)
		if !ok {
			continue
		}

		id, err := xml.GetInteger(pd, "id")
		if err != nil {
			continue
		}
		petIds = append(petIds, uint32(id))
	}
	return validRequirementProducer(checkPets(petIds))
}

func checkPets(ids []uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			//TODO
			return false
		}
	}
}

func npcRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkNpc(uint32(val)))
}

func checkNpc(reqNpc uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			return npcId == reqNpc
		}
	}
}

func monsterBookRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkMonsterBookCount(uint32(val)))
}

func checkMonsterBookCount(count uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.MonsterBookCount(l, span)(characterId) >= count
		}
	}
}

func monsterRequirement(questId uint16, r xml.Noder) checkProducer {
	monsters := make(map[uint32]uint32)
	mrs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, mr := range mrs.Children() {
		md, ok := mr.(xml.Parent)
		if !ok {
			continue
		}

		id, err := xml.GetInteger(md, "id")
		if err != nil {
			continue
		}
		count, err := xml.GetInteger(md, "count")
		if err != nil {
			continue
		}
		monsters[uint32(id)] = uint32(count)
	}
	return validRequirementProducer(checkMonster(questId, monsters))
}

func checkMonster(questId uint16, monsters map[uint32]uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			return quest.HasMetMonsterRequirement(l, span, db)(characterId, questId, monsters)
		}
	}
}

func petTamenessMinimumRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkMinTameness(val))
}

func checkMinTameness(tameness int) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			//TODO
			return false
		}
	}
}

func minLevelRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkMinLevel(byte(val)))
}

func checkMinLevel(level byte) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.IsMinimalLevel(l, span)(characterId, level)
		}
	}
}

func mesoRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkMinMeso(uint32(val)))
}

func checkMinMeso(meso uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.HasMinimalMeso(l, span)(characterId, meso)
		}
	}
}

func maxLevelRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkMaxLevel(byte(val)))
}

func checkMaxLevel(level byte) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.IsMaximalLevel(l, span)(characterId, level)
		}
	}
}

func itemRequirement(r xml.Noder) checkProducer {
	items := make(map[uint32]uint32)
	irs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, ir := range irs.Children() {
		ic, ok := ir.(xml.Parent)
		if !ok {
			continue
		}
		id, err := xml.GetInteger(ic, "id")
		if err != nil {
			continue
		}
		count, err := xml.GetInteger(ic, "count")
		if err != nil {
			continue
		}
		items[uint32(id)] = uint32(count)
	}
	return validRequirementProducer(checkItems(items))
}

func checkItems(items map[uint32]uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.HasItems(l, span)(characterId, items)
		}
	}
}

func questCompleteRequirement(sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkCompletedQuest(val))
}

func checkCompletedQuest(requiredQuest int) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			qs, err := quest.QuestsByStatus(l, span, db)(characterId, quest.StatusCompleted)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve completed quests for character, assuming none.")
				return false
			}
			return len(qs) >= requiredQuest
		}
	}
}

func intervalRequirement(questId uint16, sr xml.Noder) checkProducer {
	val, err := xml.IntFromIntegerNode(sr)
	if err != nil {
		return errorCheckProducer(err)
	}
	return validRequirementProducer(checkInterval(questId, int64(val)*60*1000))
}

func checkInterval(questId uint16, interval int64) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			cq, err := quest.GetById(l, span, db)(characterId, questId)
			if err != nil {
				l.WithError(err).Errorf("Unable to locate quest %d information for character %d. Assuming check fails.", questId, characterId)
				return false
			}
			if cq.Status() != quest.StatusCompleted {
				return true
			}

			if cq.Completion().UnixMilli() <= time.Now().UnixMilli()-interval {
				return true
			}

			//TODO emit PINK_TEXT saying This quest will become available again in approximately xxx.
			return false
		}
	}
}

func infoExRequirement(_ xml.Noder) checkProducer {
	return validRequirementProducer(validCheck)
}

func infoNumberRequirement(_ xml.Noder) checkProducer {
	return validRequirementProducer(validCheck)
}

func fieldEnterRequirement(r xml.Noder) checkProducer {
	mapId := uint32(0)
	fr, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	zf, err := xml.GetInteger(fr, "0")
	if err == nil {
		mapId = uint32(zf)
	}
	return validRequirementProducer(checkMap(mapId))
}

func checkMap(mapId uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.InMap(l, span)(characterId, mapId)
		}
	}
}

func otherQuestRequirement(r xml.Noder) checkProducer {
	quests := make(map[uint16]uint32)
	qrs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, qr := range qrs.Children() {
		qd, ok := qr.(xml.Parent)
		if !ok {
			continue
		}
		id, err := xml.GetInteger(qd, "id")
		if err != nil {
			continue
		}
		state, err := xml.GetInteger(qd, "state")
		if err != nil {
			continue
		}
		quests[uint16(id)] = uint32(state)
	}
	return validRequirementProducer(checkOtherQuests(quests))
}

func checkOtherQuests(quests map[uint16]uint32) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			cqs, err := quest.ForCharacter(l, span, db)(characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve quests for character %d. Assuming criteria is not met.", characterId)
				return false
			}
			qm := make(map[uint16]quest.Model)
			for _, cq := range cqs {
				qm[cq.Id()] = cq
			}

			for questId, statusId := range quests {
				expectedStatus := getStatusById(statusId)
				if expectedStatus == quest.StatusNotStarted {
					continue
				}

				if q, ok := qm[questId]; ok {
					if expectedStatus != q.Status() {
						return false
					}
				} else {
					return false
				}
			}
			return true
		}
	}
}

func getStatusById(id uint32) string {
	switch id {
	case 0:
		return quest.StatusNotStarted
	case 1:
		return quest.StatusStarted
	case 2:
		return quest.StatusCompleted
	default:
		return quest.StatusUndefined
	}
}

func jobRequirement(r xml.Noder) checkProducer {
	var ids []uint16
	jrs, ok := r.(xml.Parent)
	if !ok {
		return errorCheckProducer(errors.New("invalid xml structure"))
	}

	for _, jr := range jrs.Children() {
		jd, ok := jr.(*xml.IntegerNode)
		if !ok {
			continue
		}
		id, err := strconv.Atoi(jd.Value())
		if err != nil {
			return errorCheckProducer(err)
		}
		ids = append(ids, uint16(id))
	}
	return validRequirementProducer(checkJobs(ids))
}

func checkJobs(ids []uint16) CheckFunc {
	return func(l logrus.FieldLogger, span opentracing.Span, _ *gorm.DB) func(characterId uint32, _ uint32) bool {
		return func(characterId uint32, _ uint32) bool {
			return character.IsJob(l, span)(characterId, ids)
		}
	}
}

func endDateRequirement(r xml.Noder) checkProducer {
	val, err := xml.StringFromStringNode(r)
	if err != nil {
		return errorCheckProducer(err)
	}

	year, err := strconv.Atoi(val[0:4])
	if err != nil {
		return errorCheckProducer(err)
	}
	month, err := strconv.Atoi(val[4:6])
	if err != nil {
		return errorCheckProducer(err)
	}
	day, err := strconv.Atoi(val[6:8])
	if err != nil {
		return errorCheckProducer(err)
	}
	hod, err := strconv.Atoi(val[8:10])
	if err != nil {
		return errorCheckProducer(err)
	}
	endDate := time.Date(year, time.Month(month), day, hod, 0, 0, 0, time.Now().Location())
	return validRequirementProducer(func(l logrus.FieldLogger, span opentracing.Span, db *gorm.DB) func(characterId uint32, npcId uint32) bool {
		return func(characterId uint32, npcId uint32) bool {
			return endDate.After(time.Now())
		}
	})
}

func getByWZName(name string) (Type, error) {
	switch name {
	case "job":
		return TypeJob, nil
	case "quest":
		return TypeQuest, nil
	case "item":
		return TypeItem, nil
	case "lvmin":
		return TypeMinimumLevel, nil
	case "lvmax":
		return TypeMaximumLevel, nil
	case "end":
		return TypeEndDate, nil
	case "mob":
		return TypeMob, nil
	case "npc":
		return TypeNPC, nil
	case "fieldEnter":
		return TypeFieldEnter, nil
	case "interval":
		return TypeInterval, nil
	case "startscript":
		return TypeStartScript, nil
	case "endscript":
		return TypeEndScript, nil
	case "pet":
		return TypePet, nil
	case "pettamenessmin":
		return TypePetTamenessMinimum, nil
	case "mbmin":
		return TypeMonsterBook, nil
	case "normalAutoStart":
		return TypeNormalAutoStart, nil
	case "infoNumber":
		return TypeInfoNumber, nil
	case "infoex":
		return TypeInfoEx, nil
	case "questComplete":
		return TypeQuestComplete, nil
	case "start":
		return TypeStart, nil
	case "dayByDay":
		return TypeDayByDay, nil
	case "money":
		return TypeMoney, nil
	case "buff":
		return TypeBuff, nil
	case "exceptbuff":
		return TypeExceptBuff, nil
	case "equipAllNeed":
		return TypeEquipAllNeed, nil
	case "equipSelectNeed":
		return TypeEquipSelectNeed, nil
	case "skill":
		return TypeSkill, nil
	case "info":
		return TypeInfo, nil
	case "mbcard":
		return TypeMonsterBookCard, nil
	case "worldmin":
		return TypeWorldMin, nil
	case "worldmax":
		return TypeWorldMax, nil
	case "morph":
		return TypeMorph, nil
	case "pop":
		return TypePopularity, nil
	case "endmeso":
		return TypeEndMeso, nil
	case "level":
		return TypeLevel, nil
	case "partyQuest_S":
		return TypePartyQuestS, nil
	case "userInteract":
		return TypeUserInteract, nil
	case "petRecallLimit":
		return TypePetRecallLimit, nil
	case "petAutoSpeakingLimit":
		return TypePetAutoSpeakingLimit, nil
	case "tamingmoblevelmin":
		return TypeTamingMobLevelMin, nil
	}
	return "", errors.New("unknown type")
}

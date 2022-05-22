package quest

import (
	"atlas-quest/quest/action"
	"atlas-quest/quest/requirement"
	"atlas-quest/wz"
	"atlas-quest/xml"
	"errors"
	"strconv"
)

func readQuests() ([]Model, error) {
	qi, err := getQuestInfo()
	if err != nil {
		return nil, err
	}

	ci, err := getCheckInfo(err)
	if err != nil {
		return nil, err
	}

	ai, err := getActInfo(err)
	if err != nil {
		return nil, err
	}

	results := make([]Model, 0)
	for _, cn := range qi.Children() {
		questId, err := strconv.Atoi(cn.Name())
		if err != nil {
			return nil, err
		}
		q, err := createQuest(uint16(questId), cn, ci, ai)
		if err != nil {
			return nil, err
		}
		results = append(results, q)
	}
	return results, nil
}

func getCheckInfo(err error) (xml.Parent, error) {
	fe, err := wz.GetFileCache().GetFile("Check.img.xml")
	if err != nil {
		return nil, err
	}
	ci, err := xml.Read(fe.Path())
	if err != nil {
		return nil, err
	}
	return ci, nil
}

func getActInfo(err error) (xml.Parent, error) {
	fe, err := wz.GetFileCache().GetFile("Act.img.xml")
	if err != nil {
		return nil, err
	}
	ci, err := xml.Read(fe.Path())
	if err != nil {
		return nil, err
	}
	return ci, nil
}

func getQuestInfo() (xml.Parent, error) {
	fe, err := wz.GetFileCache().GetFile("QuestInfo.img.xml")
	if err != nil {
		return nil, err
	}
	fn, err := xml.Read(fe.Path())
	if err != nil {
		return nil, err
	}
	return fn, nil
}

func createQuest(questId uint16, cn xml.Noder, ci xml.Parent, ai xml.Parent) (Model, error) {
	modelBuilder := NewBuilder(questId)

	rd, err := ci.ChildByName(strconv.Itoa(int(questId)))
	if err != nil {
		// most likely infoEx
		return modelBuilder.Build(), nil
	}

	qi, ok := cn.(xml.Parent)
	if !ok {
		return modelBuilder.Build(), errors.New("invalid xml structure")
	}
	name, err := xml.GetString(qi, "name")
	if err != nil {
		return modelBuilder.Build(), err
	}
	modelBuilder.SetName(name)
	parent, err := xml.GetString(qi, "parent")
	if err == nil {
		modelBuilder.SetParent(parent)
	}
	timeLimit, err := xml.GetInteger(qi, "timeLimit")
	if err == nil {
		modelBuilder.SetTimeLimit(uint32(timeLimit))
	}
	timeLimit2, err := xml.GetInteger(qi, "timeLimit2")
	if err == nil {
		modelBuilder.SetTimeLimit2(uint32(timeLimit2))
	}
	autoStart, err := xml.GetBoolean(qi, "autoStart")
	if err == nil {
		modelBuilder.SetAutoStart(autoStart)
	}
	autoPreComplete, err := xml.GetBoolean(qi, "autoPreComplete")
	if err == nil {
		modelBuilder.SetAutoPreComplete(autoPreComplete)
	}
	autoComplete, err := xml.GetBoolean(qi, "autoComplete")
	if err == nil {
		modelBuilder.SetAutoComplete(autoComplete)
	}
	viewMedalItem, err := xml.GetInteger(qi, "viewMedalItem")
	if err == nil {
		modelBuilder.SetMedalItem(uint32(viewMedalItem))
	}

	// load starting requirements
	srs, err := requirement.GetStarting(questId, rd)
	if err != nil {
		return modelBuilder.Build(), err
	}
	for _, sr := range srs {
		if sr.Type() == requirement.TypeInterval {
			modelBuilder.SetRepeatable(true)
		} else if sr.Type() == requirement.TypeMob {
			for _, rm := range sr.RelevantMobs() {
				modelBuilder.AddRelevantMob(rm)
			}
		}
		modelBuilder.AddStartingRequirement(sr.Type(), sr.Check())
	}

	// load completion requirements
	ers, err := requirement.GetEnding(questId, rd)
	if err != nil {
		return modelBuilder.Build(), err
	}
	for _, er := range ers {
		if er.Type() == requirement.TypeInterval {
			modelBuilder.SetRepeatable(true)
		} else if er.Type() == requirement.TypeMob {
			for _, rm := range er.RelevantMobs() {
				modelBuilder.AddRelevantMob(rm)
			}
		}
		modelBuilder.AddCompletionRequirement(er.Type(), er.Check())
	}

	ad, err := ai.ChildByName(strconv.Itoa(int(questId)))
	if ad == nil || err != nil {
		return modelBuilder.Build(), nil
	}

	sas, err := action.GetStarting(questId, ad)
	if err != nil {
		return modelBuilder.Build(), err
	}
	for _, sa := range sas {
		modelBuilder.AddStartingAction(sa.Type(), sa.Check(), sa.Run())
	}

	cas, err := action.GetEnding(questId, ad)
	if err != nil {
		return modelBuilder.Build(), err
	}
	for _, sa := range cas {
		modelBuilder.AddCompletionAction(sa.Type(), sa.Check(), sa.Run())
	}

	return modelBuilder.Build(), nil
}

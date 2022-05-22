package action

import (
	"atlas-quest/xml"
	"errors"
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
		actType, err := getByWZName(req.Name())
		if err != nil {
			return nil, err
		}

		m := Model{theType: actType}
		check, run, err := getActionProducer(questId, actType, req)()
		if err != nil {
			return nil, err
		}
		m.check = check
		m.run = run
		results = append(results, m)
	}
	return results, nil
}

type actionProducer func() (CheckFunc, RunFunc, error)

func getActionProducer(questId uint16, actType Type, req xml.Noder) actionProducer {
	switch actType {

	}
}

func getByWZName(name string) (Type, error) {
	switch name {
	case "exp":
		return TypeExperience, nil
	case "money":
		return TypeMoney, nil
	case "item":
		return TypeItem, nil
	case "skill":
		return TypeSkill, nil
	case "nextQuest":
		return TypeNextQuest, nil
	case "pop":
		return TypePopularity, nil
	case "buffItemID":
		return TypeBuffItemId, nil
	case "petskill":
		return TypePetSkill, nil
	case "no":
		return TypeNo, nil
	case "yes":
		return TypeYes, nil
	case "npc":
		return TypeNPC, nil
	case "lvmin":
		return TypeMinimumLevel, nil
	case "normalAutoStart":
		return TypeNormalAutoStart, nil
	case "pettameness":
		return TypePetTameness, nil
	case "petspeed":
		return TypePetSpeed, nil
	case "info":
		return TypeInfo, nil
	case "0":
		return Type0, nil
	}
	return "", errors.New("unknown type")
}

package requirement

import (
	"atlas-quest/xml"
)

func GetStarting(questId uint16, root xml.Noder) ([]Model, error) {
	return get(questId, root, "0")
}

func GetEnding(questId uint16, root xml.Noder) ([]Model, error) {
	return get(questId, root, "1")
}

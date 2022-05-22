package quest

import "github.com/sirupsen/logrus"

func GetById(l logrus.FieldLogger) func(id uint32) (*Model, error) {
	return func(id uint32) (*Model, error) {
		return nil, nil
	}
}

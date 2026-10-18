package service

import (
	"fmt"
	"strings"

	"hw1/domain"
)

// ServiceCenter - техосмотр: решает, можно ли принять транспорт в парк
type ServiceCenter interface {
	Inspect(v domain.Vehicle) (bool, string)
}

// BasicServiceCenter - обычный техосмотр с простыми проверками
type BasicServiceCenter struct {
	maxDailyEnergy float64
}

func NewBasicServiceCenter(maxDailyEnergy float64) *BasicServiceCenter {
	return &BasicServiceCenter{maxDailyEnergy: maxDailyEnergy}
}

func (s *BasicServiceCenter) Inspect(v domain.Vehicle) (bool, string) {
	if strings.TrimSpace(v.Name()) == "" {
		return false, "не указано наименование"
	}
	if v.Simplicity() < 1 || v.Simplicity() > 10 {
		return false, fmt.Sprintf("простота должна быть от 1 до 10, а указано %d", v.Simplicity())
	}

	consumer, ok := v.(domain.EnergyConsumer)
	if ok {
		if consumer.DailyEnergy() < 0 {
			return false, "расход энергии не может быть отрицательным"
		}
		if consumer.DailyEnergy() > s.maxDailyEnergy {
			return false, fmt.Sprintf("расход %.2f кВт·ч больше допустимого %.2f кВт·ч",
				consumer.DailyEnergy(), s.maxDailyEnergy)
		}
	}
	return true, ""
}

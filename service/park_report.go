package service

import "hw1/domain"

const MinBeginnerSimplicity = 6

type KindCount struct {
	Kind  string
	Count int
}

func (p *Park) VehicleCountByKind() []KindCount {
	var result []KindCount
	index := make(map[string]int)

	for _, v := range p.vehicles {
		i, ok := index[v.Kind()]

		if !ok {
			i = len(result)
			index[v.Kind()] = i
			result = append(result, KindCount{Kind: v.Kind()})
		}
		result[i].Count++
	}

	return result
}

// TotalEnergy - суточный расход энергии всего на балансе
func (p *Park) TotalEnergy() float64 {
	total := 0.0

	for _, obj := range p.AllInventory() {
		consumer, ok := obj.(domain.EnergyConsumer)
		if ok {
			total += consumer.DailyEnergy()
		}
	}

	return total
}

// BeginnerVehicles - транспорт, который можно выдавать новичкам
func (p *Park) BeginnerVehicles() []domain.Vehicle {
	var result []domain.Vehicle

	for _, v := range p.vehicles {
		if v.Simplicity() >= MinBeginnerSimplicity {
			result = append(result, v)
		}
	}

	return result
}

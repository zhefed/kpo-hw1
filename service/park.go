package service

import (
	"errors"
	"strings"

	"hw1/domain"
)

type Park struct {
	serviceCenter ServiceCenter
	vehicles      []domain.Vehicle
	things        []domain.Thing
	lastNumber    int
}

func NewPark(serviceCenter ServiceCenter) *Park {
	return &Park{serviceCenter: serviceCenter}
}

// AddVehicle отправляет транспорт на техосмотр и если он пройден,
// ставит его на баланс и выдаёт инвентарный номер
func (p *Park) AddVehicle(v domain.Vehicle) error {
	ok, reason := p.serviceCenter.Inspect(v)

	if !ok {
		return errors.New("техосмотр не пройден: " + reason)
	}

	v.SetInventoryNumber(p.newNumber())
	p.vehicles = append(p.vehicles, v)
	return nil
}

// AddThing проверяет вещь и если с ней всё в порядке,
// ставит её на баланс и выдаёт инвентарный номер
func (p *Park) AddThing(t domain.Thing) error {
	if strings.TrimSpace(t.Name()) == "" {
		return errors.New("не указано наименование")
	}

	consumer, ok := t.(domain.EnergyConsumer)
	if ok && consumer.DailyEnergy() < 0 {
		return errors.New("расход энергии не может быть отрицательным")
	}

	t.SetInventoryNumber(p.newNumber())
	p.things = append(p.things, t)
	return nil
}

func (p *Park) newNumber() int {
	p.lastNumber++
	return p.lastNumber
}

func (p *Park) Vehicles() []domain.Vehicle {
	return p.vehicles
}

func (p *Park) Things() []domain.Thing {
	return p.things
}

func (p *Park) AllInventory() []domain.Inventoried {
	var result []domain.Inventoried
	for _, v := range p.vehicles {
		result = append(result, v)
	}
	for _, t := range p.things {
		result = append(result, t)
	}
	return result
}

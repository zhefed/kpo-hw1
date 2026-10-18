package domain

// Inventoried - всё, что ставится на баланс и имеет инвентарный номер
type Inventoried interface {
	Name() string
	InventoryNumber() int
	SetInventoryNumber(number int)
}

// EnergyConsumer - всё, что потребляет электроэнергию
type EnergyConsumer interface {
	DailyEnergy() float64 // суточный расход в кВт·ч
}

// item - общая часть для всего, что стоит на балансе
type item struct {
	name   string
	number int // 0 - номер ещё не присвоен
}

func (i *item) Name() string {
	return i.name
}

func (i *item) InventoryNumber() int {
	return i.number
}

func (i *item) SetInventoryNumber(number int) {
	i.number = number
}

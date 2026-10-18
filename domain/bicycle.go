package domain

// Bicycle - обычный велосипед. Энергию не потребляет, поэтому не реализует EnergyConsumer
type Bicycle struct {
	baseVehicle
}

func NewBicycle(name string, simplicity int) *Bicycle {
	b := &Bicycle{}
	b.name = name
	b.simplicity = simplicity
	return b
}

func (b *Bicycle) Kind() string {
	return "Велосипед"
}

func init() {
	RegisterVehicleKind("bicycle", func(name string, simplicity int, _ float64) Vehicle {
		return NewBicycle(name, simplicity)
	})
}

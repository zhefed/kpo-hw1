package domain

// Vehicle - прокатный транспорт
type Vehicle interface {
	Inventoried
	Kind() string
	Simplicity() int
}

// baseVehicle - общие поля всех видов транспорта
type baseVehicle struct {
	item
	simplicity int
}

func (v *baseVehicle) Simplicity() int {
	return v.simplicity
}

// VehicleFactory - создаёт транспорт конкретного вида
type VehicleFactory func(name string, simplicity int, energy float64) Vehicle

// VehicleKind - вид транспорта в реестре. HumanName и UsesEnergy, берутся у самого типа при регистрации
type VehicleKind struct {
	Slug       string
	HumanName  string
	UsesEnergy bool
	Create     VehicleFactory
}

var vehicleKinds []VehicleKind
var vehicleBySlug = map[string]VehicleKind{}

// RegisterVehicleKind - добавляет новый вид транспорта. Вызывается из init() в файле конкретного вида.
// Создаёт пустой образец, чтобы узнать у типа его название и потребляет ли он энергию.
func RegisterVehicleKind(slug string, create VehicleFactory) {
	if _, exists := vehicleBySlug[slug]; exists {
		panic("вид транспорта уже зарегистрирован: " + slug)
	}

	sample := create("", 0, 0)
	_, usesEnergy := sample.(EnergyConsumer)

	kind := VehicleKind{
		Slug:       slug,
		HumanName:  sample.Kind(),
		UsesEnergy: usesEnergy,
		Create:     create,
	}
	vehicleBySlug[slug] = kind
	vehicleKinds = append(vehicleKinds, kind)
}

func VehicleKinds() []VehicleKind {
	return vehicleKinds
}

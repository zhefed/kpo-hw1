package service_test

import (
	"strings"
	"testing"

	"hw1/domain"
	"hw1/service"
)

// stubServiceCenter - stub техосмотра: всегда даёт заранее заданный ответ
type stubServiceCenter struct {
	ok     bool
	reason string
}

func (s stubServiceCenter) Inspect(domain.Vehicle) (bool, string) {
	return s.ok, s.reason
}

func acceptingPark() *service.Park {
	return service.NewPark(stubServiceCenter{ok: true})
}

func TestAddVehicle_AcceptedVehicleJoinsParkWithNextNumber(t *testing.T) {
	park := acceptingPark()
	if err := park.AddThing(domain.NewHelmet("Шлем", "M")); err != nil {
		t.Fatal(err)
	}
	scooter := domain.NewScooter("Ninebot", 7, 0.35)

	err := park.AddVehicle(scooter)

	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	vehicles := park.Vehicles()
	if len(vehicles) != 1 || vehicles[0] != scooter {
		t.Fatalf("ожидали в парке только добавленный самокат, получили %v", vehicles)
	}
	if scooter.InventoryNumber() != 2 {
		t.Fatalf("ожидали номер 2 после вещи с номером 1, получили %d", scooter.InventoryNumber())
	}
}

func TestAddVehicle_RejectedVehicleReturnsReasonAndIsNotAdded(t *testing.T) {
	park := service.NewPark(stubServiceCenter{ok: false, reason: "сломан тормоз"})

	err := park.AddVehicle(domain.NewScooter("Ninebot", 7, 0.35))

	if err == nil {
		t.Fatal("ожидали ошибку")
	}
	if !strings.Contains(err.Error(), "сломан тормоз") {
		t.Fatalf("в ошибке нет причины отказа: %v", err)
	}
	if len(park.Vehicles()) != 0 {
		t.Fatalf("отклонённый транспорт попал в парк: %v", park.Vehicles())
	}
}

func TestAddThing_RejectsInvalidThing(t *testing.T) {
	tests := []struct {
		name  string
		thing domain.Thing
	}{
		{"пустое имя", domain.NewHelmet("", "M")},
		{"отрицательный расход док-станции", domain.NewDockStation("Док", -1)},
		{"отрицательный расход зарядного шкафа", domain.NewChargingCabinet("Шкаф", -1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			park := acceptingPark()

			err := park.AddThing(tt.thing)

			if err == nil {
				t.Fatal("ожидали ошибку")
			}
			if len(park.Things()) != 0 {
				t.Fatalf("некорректная вещь попала на баланс: %v", park.Things())
			}
		})
	}
}

func TestAddThing_AcceptsValidThing(t *testing.T) {
	tests := []struct {
		name  string
		thing domain.Thing
	}{
		{"шлем", domain.NewHelmet("Шлем", "L")},
		{"док-станция", domain.NewDockStation("Док", 1.2)},
		{"зарядный шкаф с нулевым расходом", domain.NewChargingCabinet("Шкаф", 0)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			park := acceptingPark()

			err := park.AddThing(tt.thing)

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			things := park.Things()
			if len(things) != 1 || things[0] != tt.thing {
				t.Fatalf("ожидали в Things() только добавленную вещь, получили %v", things)
			}
		})
	}
}

func TestInventoryNumbers_SequentialWithoutGaps(t *testing.T) {
	park := service.NewPark(service.NewBasicServiceCenter(testMaxDailyEnergy))
	scooter := domain.NewScooter("Ninebot", 7, 0.35)
	greedyScooter := domain.NewScooter("Самосбор", 5, 10) // не пройдёт техосмотр
	helmet := domain.NewHelmet("Шлем", "M")
	nameless := domain.NewHelmet("", "S") // не пройдёт проверку вещи
	bicycle := domain.NewBicycle("Stels", 9)

	_ = park.AddVehicle(scooter)
	_ = park.AddVehicle(greedyScooter)
	_ = park.AddThing(helmet)
	_ = park.AddThing(nameless)
	_ = park.AddVehicle(bicycle)

	want := []struct {
		name   string
		item   domain.Inventoried
		number int
	}{
		{"самокат", scooter, 1},
		{"отклонённый самокат", greedyScooter, 0},
		{"шлем", helmet, 2},
		{"отклонённый шлем", nameless, 0},
		{"велосипед", bicycle, 3},
	}
	for _, w := range want {
		if got := w.item.InventoryNumber(); got != w.number {
			t.Errorf("%s: ожидали номер %d, получили %d", w.name, w.number, got)
		}
	}
}

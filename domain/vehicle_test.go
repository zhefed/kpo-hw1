package domain_test

import (
	"testing"

	"hw1/domain"
)

func TestVehicleKinds_CreateMatchesRegistration(t *testing.T) {
	kinds := domain.VehicleKinds()
	if len(kinds) == 0 {
		t.Fatal("реестр видов транспорта пуст")
	}

	for _, kind := range kinds {
		t.Run(kind.Slug, func(t *testing.T) {
			v := kind.Create("Модель", 7, 1.0)

			if v.Name() != "Модель" || v.Simplicity() != 7 {
				t.Fatalf("ожидали имя «Модель» и простоту 7, получили «%s» и %d", v.Name(), v.Simplicity())
			}
			if v.Kind() != kind.HumanName {
				t.Fatalf("вид в реестре «%s», а у объекта «%s»", kind.HumanName, v.Kind())
			}
			consumer, ok := v.(domain.EnergyConsumer)
			if ok && consumer.DailyEnergy() != 1.0 {
				t.Fatalf("ожидали расход 1.0, получили %v", consumer.DailyEnergy())
			}
		})
	}
}

func TestRegisterVehicleKind_DuplicateSlugPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("ожидали панику при повторной регистрации")
		}
	}()

	domain.RegisterVehicleKind("bicycle", func(name string, simplicity int, _ float64) domain.Vehicle {
		return domain.NewBicycle(name, simplicity)
	})
}

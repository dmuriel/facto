package facto_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/paganotoni/facto"
)

type User struct {
	Name string
}

func Test_Build(t *testing.T) {
	facto.Register("User", func(f facto.Helper) facto.Product {
		u := User{
			Name: "Wawandco",
		}
		return facto.Product(u)
	})

	userProduct := facto.Build("User").(User)
	if userProduct.Name != "Wawandco" {
		t.Errorf("expected '%s' but got '%s'", "Wawandco", userProduct.Name)
	}
}

func Test_BuildN(t *testing.T) {
	facto.Register("Users", func(f facto.Helper) facto.Product {
		u := User{
			Name: fmt.Sprintf("Wawandco %d", f.Index),
		}
		return facto.Product(u)
	})

	usersProduct := facto.BuildN("Users", 5).([]User)

	for i := 0; i < 5; i++ {
		if fmt.Sprintf("Wawandco %d", i+1) != usersProduct[i].Name {
			t.Errorf("expected '%s' but got '%s'", fmt.Sprintf("Wawandco %d", i+1), usersProduct[i].Name)
		}
	}
}

func Test_Build_Concurrently(t *testing.T) {
	tcases := []struct {
		factoryName string
		factory     facto.Factory
		expected    string
	}{
		{
			factoryName: "UserNumberOne",
			factory: func(f facto.Helper) facto.Product {
				u := User{
					Name: "Wawandco",
				}
				return facto.Product(u)
			},
			expected: "Wawandco",
		},

		{
			factoryName: "UserNumberTwo",
			factory: func(f facto.Helper) facto.Product {
				u := User{
					Name: "Wawandco 2",
				}
				return facto.Product(u)
			},
			expected: "Wawandco 2",
		},

		{
			factoryName: "UserNumberThree",
			factory: func(f facto.Helper) facto.Product {
				u := User{
					Name: fmt.Sprintf("Wawandco %d", f.Index),
				}
				return facto.Product(u)
			},
			expected: "Wawandco 1",
		},
	}

	var wgreg sync.WaitGroup
	for i := range tcases {
		wgreg.Add(1)
		gr := func(fname string, factory facto.Factory) {
			defer wgreg.Done()

			facto.Register(fname, factory)
		}

		go gr(tcases[i].factoryName, tcases[i].factory)
	}
	wgreg.Wait()

	var wgbuild sync.WaitGroup
	for i := range tcases {
		wgbuild.Add(1)

		gr := func(name, expected string, index int) {
			defer wgbuild.Done()

			userProduct, ok := facto.Build(name).(User)
			if !ok {
				t.Fatalf("Should have got user but got %v", userProduct)
			}

			if expected != userProduct.Name {
				t.Errorf("expected '%s' but got '%s' in '%s'", expected, userProduct.Name, fmt.Sprintf("case %d", i+1))
			}
		}

		go gr(tcases[i].factoryName, tcases[i].expected, i)
	}
	wgbuild.Wait()
}

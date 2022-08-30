package vee_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/cs-clarence/vee"
	"github.com/cs-clarence/vee/is"
)

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type Customer struct {
	Name    string
	Gender  string
	Email   string
	Address Address
}

func (a Address) Validate() error {
	return vee.ValidateStruct(
		&a,
		// Street cannot be empty, and the length must between 5 and 50
		vee.Field(&a.Street, vee.Required, vee.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		vee.Field(&a.City, vee.Required, vee.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		vee.Field(
			&a.State,
			vee.Required,
			vee.Match(regexp.MustCompile("^[A-Z]{2}$")),
		),
		// State cannot be empty, and must be a string consisting of five digits
		vee.Field(
			&a.Zip,
			vee.Required,
			vee.Match(regexp.MustCompile("^[0-9]{5}$")),
		),
	)
}

func (c Customer) Validate() error {
	return vee.ValidateStruct(&c,
		// Name cannot be empty, and the length must be between 5 and 20.
		vee.Field(&c.Name, vee.Required, vee.Length(5, 20)),
		// Gender is optional, and should be either "Female" or "Male".
		vee.Field(&c.Gender, vee.In("Female", "Male")),
		// Email cannot be empty and should be in a valid email format.
		vee.Field(&c.Email, vee.Required, is.Email),
		// Validate Address using its own validation rules
		vee.Field(&c.Address),
	)
}

func Example() {
	c := Customer{
		Name:  "Qiang Xue",
		Email: "q",
		Address: Address{
			Street: "123 Main Street",
			City:   "Unknown",
			State:  "Virginia",
			Zip:    "12345",
		},
	}

	err := c.Validate()
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format.); Email: must be a valid email address.
}

func Example_second() {
	data := "example"
	err := vee.Validate(data,
		vee.Required,       // not empty
		vee.Length(5, 100), // length between 5 and 100
		is.URL,             // is a valid URL
	)
	fmt.Println(err)
	// Output:
	// must be a valid URL
}

func Example_third() {
	addresses := []Address{
		{State: "MD", Zip: "12345"},
		{Street: "123 Main St", City: "Vienna", State: "VA", Zip: "12345"},
		{City: "Unknown", State: "NC", Zip: "123"},
	}
	err := vee.Validate(addresses)
	fmt.Println(err)
	// Output:
	// 0: (City: cannot be blank; Street: cannot be blank.); 2: (Street: cannot be blank; Zip: must be in a valid format.).
}

func Example_four() {
	c := Customer{
		Name:  "Qiang Xue",
		Email: "q",
		Address: Address{
			State: "Virginia",
		},
	}

	err := vee.Errors{
		"name":  vee.Validate(c.Name, vee.Required, vee.Length(5, 20)),
		"email": vee.Validate(c.Name, vee.Required, is.Email),
		"zip": vee.Validate(
			c.Address.Zip,
			vee.Required,
			vee.Match(regexp.MustCompile("^[0-9]{5}$")),
		),
	}.Filter()
	fmt.Println(err)
	// Output:
	// email: must be a valid email address; zip: cannot be blank.
}

func Example_five() {
	type Employee struct {
		Name string
	}

	type Manager struct {
		Employee
		Level int
	}

	m := Manager{}
	err := vee.ValidateStruct(&m,
		vee.Field(&m.Name, vee.Required),
		vee.Field(&m.Level, vee.Required),
	)
	fmt.Println(err)
	// Output:
	// Level: cannot be blank; Name: cannot be blank.
}

type contextKey int

func Example_six() {
	key := contextKey(1)
	rule := vee.WithContext(func(ctx context.Context, value interface{}) error {
		s, _ := value.(string)
		if ctx.Value(key) == s {
			return nil
		}
		return errors.New("unexpected value")
	})
	ctx := context.WithValue(context.Background(), key, "good sample")

	err1 := vee.ValidateWithContext(ctx, "bad sample", rule)
	fmt.Println(err1)

	err2 := vee.ValidateWithContext(ctx, "good sample", rule)
	fmt.Println(err2)

	// Output:
	// unexpected value
	// <nil>
}

func Example_seven() {
	c := map[string]interface{}{
		"Name":  "Qiang Xue",
		"Email": "q",
		"Address": map[string]interface{}{
			"Street": "123",
			"City":   "Unknown",
			"State":  "Virginia",
			"Zip":    "12345",
		},
	}

	err := vee.Validate(c,
		vee.Map(
			// Name cannot be empty, and the length must be between 5 and 20.
			vee.Key("Name", vee.Required, vee.Length(5, 20)),
			// Email cannot be empty and should be in a valid email format.
			vee.Key("Email", vee.Required, is.Email),
			// Validate Address using its own validation rules
			vee.Key("Address", vee.Map(
				// Street cannot be empty, and the length must between 5 and 50
				vee.Key("Street", vee.Required, vee.Length(5, 50)),
				// City cannot be empty, and the length must between 5 and 50
				vee.Key("City", vee.Required, vee.Length(5, 50)),
				// State cannot be empty, and must be a string consisting of two letters in upper case
				vee.Key(
					"State",
					vee.Required,
					vee.Match(regexp.MustCompile("^[A-Z]{2}$")),
				),
				// State cannot be empty, and must be a string consisting of five digits
				vee.Key(
					"Zip",
					vee.Required,
					vee.Match(regexp.MustCompile("^[0-9]{5}$")),
				),
			)),
		),
	)
	fmt.Println(err)
	// Output:
	// Address: (State: must be in a valid format; Street: the length must be between 5 and 50.); Email: must be a valid email address.
}

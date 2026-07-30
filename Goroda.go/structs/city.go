package structs

import (
	"fmt"
)

type City struct {
	name        string
	temperature float64
}

func NewCity(name string) *City {
	return &City{name: name}
}

func (c *City) SetTemperature(temperature float64) {
	c.temperature = temperature
}

func (c *City) GetTemperature() float64 {
	return c.temperature
}

func (c *City) Print() string {
	return fmt.Sprintf("the weather in %s is %f", c.name, c.temperature)
}

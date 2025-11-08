package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

const (
	MASK_H uint16 = 0b111111_0000000000
	MASK_L uint16 = 0b000000_1111111111
	MASK_F uint16 = 0b111111_1111111111
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		person.name = *(*[42]byte)([]byte(name)[:42])
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x, person.y, person.z = int32(x), int32(y), int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		if gold < 0 {
			gold = 0
		}
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		if mana > 1_000 {
			mana = 1_000
		}
		if mana < 0 {
			mana = 0
		}

		person.respect_mana = (person.respect_mana & MASK_H) | uint16(mana)

	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		if health > 1_000 {
			health = 1_000
		}
		if health < 0 {
			health = 0
		}
		person.strength_health = (person.strength_health & MASK_H) | uint16(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		if respect > 10 {
			respect = 10
		}
		if respect < 0 {
			respect = 0
		}

		r := uint16(respect) << 10
		person.respect_mana = (person.respect_mana & MASK_L) | r
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		if strength > 10 {
			strength = 10
		}
		if strength < 0 {
			strength = 0
		}

		r := uint16(strength) << 10
		person.strength_health = (person.strength_health & MASK_L) | r
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		if experience > 10 {
			experience = 10
		}
		if experience < 0 {
			experience = 0
		}

		r := uint8(experience) << 4
		person.experience_level = (person.experience_level & 0b0000_1111) | r
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		if level > 10 {
			level = 10
		}
		if level < 0 {
			level = 0
		}

		person.experience_level = (person.experience_level & 0b1111_0000) | uint8(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.options |= 1 << 0
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.options |= 1 << 1
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.options |= 1 << 2
	}
}

func WithType(personType PersonType) func(*GamePerson) {
	return func(person *GamePerson) {
		// need to implement
		switch personType {
		case BuilderGamePersonType:
			person.options = (person.options & 0b0000_1111) | 0b0001_0000
		case BlacksmithGamePersonType:
			person.options = (person.options & 0b0000_1111) | 0b0010_0000
		case WarriorGamePersonType:
			person.options = (person.options & 0b0000_1111) | 0b0100_0000
		}
	}
}

type PersonType int

const (
	BuilderGamePersonType PersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

var _ uintptr = 64 - unsafe.Sizeof(GamePerson{})

type GamePerson struct {
	x, y, z int32
	gold    uint32
	// 6     10		6       10
	respect_mana, strength_health uint16
	// 	4 		4
	experience_level uint8

	options uint8

	name [42]byte
}

func NewGamePerson(options ...Option) GamePerson {
	ret := GamePerson{}

	for _, opt := range options {
		opt(&ret)
	}

	return ret
}

func (p *GamePerson) Name() string {
	return string(p.name[:])
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

// Магическая сила (мана) [0…1000] значений
func (p *GamePerson) Mana() int {
	r := p.respect_mana & MASK_L
	return int(r)
}

// Здоровье [0…1000] значений
func (p *GamePerson) Health() int {
	r := p.strength_health & MASK_L
	return int(r)
}

// Уважение [0…10] значений
func (p *GamePerson) Respect() int {
	r := p.respect_mana >> 10
	return int(r)
}

// Сила [0…10] значений
func (p *GamePerson) Strength() int {
	r := p.strength_health >> 10
	return int(r)

}

// Опыт [0…10] значений
func (p *GamePerson) Experience() int {
	r := p.experience_level >> 4
	return int(r)
}

// Уровень [0…10] значений
func (p *GamePerson) Level() int {
	r := p.experience_level & 0b0000_1111
	return int(r)
}

func (p *GamePerson) HasHouse() bool {
	var b uint8 = 1 << 0
	r := p.options & b
	if r > 0 {
		return true
	}
	return false
}

func (p *GamePerson) HasGun() bool {
	var b uint8 = 1 << 1
	r := p.options & b
	if r > 0 {
		return true
	}
	return false
}

func (p *GamePerson) HasFamilty() bool {
	var b uint8 = 1 << 2
	r := p.options & b
	if r > 0 {
		return true
	}
	return false
}

func (p *GamePerson) Type() PersonType {
	// need to implement
	r := p.options >> 4
	switch r {
	case 0b0001:
		return BuilderGamePersonType
	case 0b0010:
		return BlacksmithGamePersonType
	case 0b0100:
		return WarriorGamePersonType
	}

	return 0
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

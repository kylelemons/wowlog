package combatlog

import (
	"strings"
)

type UnitFlags uint64
const (
	UnitSelf UnitFlags = 0x1
	UnitParty UnitFlags = 0x2
	UnitRaid UnitFlags = 0x4
	UnitEnemy UnitFlags = 0x40
)
func (f UnitFlags) String() string {
	flags := []string{}
	if f | UnitSelf != 0 {
		flags = append(flags, "Self")
	}
	if f | UnitParty != 0 {
		flags = append(flags, "Party")
	}
	if f | UnitRaid != 0 {
		flags = append(flags, "Raid")
	}
	if f | UnitEnemy != 0 {
		flags = append(flags, "Enemy")
	}
	return strings.Join(flags, "|")
}

type SpellSchool uint32
const (
	SchoolPhysical SpellSchool = 1 << iota
	SchoolHoly
	SchoolFire
	SchoolNature
	SchoolFrost
	SchoolShadow
	SchoolArcane
)
var schoolNames = [...]string{
	"Physical", "Holy", "Fire", "Nature", "Frost", "Shadow", "Arcane",
}
func (s SpellSchool) String() string {
	schools := []string{}
	for i, name := range schoolNames {
		if s | (1 << uint(i)) != 0 {
			schools = append(schools, name)
		}
	}
	return strings.Join(schools, "/")
}

type PowerType int32
const (
	PowerHealth = -2
	PowerMana = iota
	PowerRage
	PowerFocus
	PowerEnergy
	PowerHappiness
	PowerRunes
	PowerRunic
	PowerSoulShard
	PowerEclipse // -solar +lunar
	PowerHoly
	PowerSound // Astramedes
)
var powerNames = [...]string{
	"Mana", "Rage", "Focus", "Energy", "Happiness",
	"Runes", "Runic", "SoulShard", "Eclipse", "Holy",
	"Sound",
}
func (p PowerType) String() string {
	if p == PowerHealth {
		return "Health"
	}
	if p >= 0 && int(p) < len(powerNames) {
		return powerNames[p]
	}
	return "Unknown"
}

const (
	MissMiss = "MISS"
	MissAbsorb = "ABSORB"
	MissBlock = "BLOCK"
	MissDeflect = "DEFLECT"
	MissDodge = "DODGE"
	MissEvade = "EVADE"
	MissImmune = "IMMUNE"
	MissParry = "PARRY"
	MissReflect = "REFLECT"
	MissResist = "RESIST"
)

const (
	AuraBuff = "BUFF"
	AuraDebuff = "DEBUFF"
)

const (
	EnvDrowning = "DROWNING"
	EnvFalling = "FALLING"
	EnvFatigue = "FATIGUE"
	EnvFire = "FIRE"
	EnvLava = "LAVA"
	EnvSlime = "SLIME"
)

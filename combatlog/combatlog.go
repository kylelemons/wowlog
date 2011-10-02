package combatlog

type Unit struct {
	ID    uint64
	Name  string
	Flags uint64
	Null  int32
}

type Common struct {
	Source Unit
	Dest   Unit
}
func (c Common) GetSource() Unit {
	return c.Source
}
func (c Common) GetDest() Unit {
	return c.Dest
}

type Spell struct {
	ID     uint64
	Name   string
	School int32
}
func (s Spell) GetSpell() Spell {
	return s
}

type Damage struct {
	Amount   int64
	Unknown  int32
	School   int32
	Resisted int64
	Blocked  int64
	Absorbed int64
	Critical bool
	Glancing bool
	Crushing bool
}
func (d Damage) GetDamage() Damage {
	return d
}

type Heal struct {
	Amount   int64
	Overheal int64
	Unknown  int64
	Critical bool
}
func (h Heal) GetHeal() Heal {
	return h
}

type Miss struct {
	Type    string
	Unknown int64  `combatlog:"optional"`
}

type Shield struct {
	Amount   int64
	Unknown  int64
	Unknown2 int64
}

type Aura struct {
	Type   string
	Shield       `combatlog:"optional"`
}

type Power struct {
	Amount int64
	Type   int32
}

type Item struct {
	ID   int64
	Item string
}

// ENVIRONMENTAL_DAMAGE
type EnvironmentalDamage struct {
	Common
	Type   string
	Damage
}

// SWING_DAMAGE
type SwingDamage struct {
	Common
	Damage
}

// SWING_MISSED
type SwingMissed struct {
	Common
	Miss
}

// RANGE_DAMAGE
type RangeDamage struct {
	Common
	Spell
	Damage
}

// RANGE_MISSED
type RangeMissed struct {
	Common
	Spell
	Miss
}

// SPELL_CAST_START
type SpellCastStart struct {
	Common
	Spell
}

// SPELL_CAST_SUCCESS
type SpellCastSuccess struct {
	Common
	Spell
}

// SPELL_CAST_FAILED
type SpellCastFailed struct {
	Common
	Spell
	Miss
}

// SPELL_MISSED
type SpellMissed struct {
	Common
	Spell
	Miss
}

// SPELL_DAMAGE
type SpellDamage struct {
	Common
	Spell
	Damage
}

// SPELL_BUILDING_DAMAGE
type SpellBuildingDamage struct {
	Common
	Spell
	Damage
}

// SPELL_HEAL
type SpellHeal struct {
	Common
	Spell
	Heal
}

// SPELL_ENERGIZE
type SpellEnergize struct {
	Common
	Spell
	Power
}

// SPELL_DRAIN
type SpellDrain struct {
	Common
	Spell
	Power
	Drained int64
}

// SPELL_LEECH
type SpellLeech struct {
	Common
	Spell
	Power
	Leeched int64
}

// SPELL_INSTAKILL		
type SpellInstakill struct {
	Common
	Spell
}

// SPELL_INTERRUPT
type SpellInterrupt struct {
	Common
	Spell
	Interrupted Spell
}

// SPELL_DISPEL
type SpellDispel struct {
	Common
	Spell
	Dispelled Spell
	Aura
}

// SPELL_EXTRA_ATTACKS
type SpellExtraAttacks struct {
	Common
	Spell
	Amount int64
}

// SPELL_DURABILITY_DAMAGE
type SpellDurabilityDamage struct {
	Common
	Spell
}

// SPELL_DURABILITY_DAMAGE_ALL
type SpellDurabilityDamageAll struct {
	Common
	Spell
}

// SPELL_AURA_APPLIED
type SpellAuraApplied struct {
	Common
	Spell
	Aura
}

// SPELL_AURA_REFRESH
type SpellAuraRefresh struct {
	Common
	Spell
	Aura
}

// SPELL_AURA_BROKEN_SPELL
type SpellAuraBrokenSpell struct {
	Source  Unit
	Dest    Unit
	Broken  Spell
	Breaker Spell
	Aura
}

// SPELL_AURA_APPLIED_DOSE
type SpellAuraAppliedDose struct {
	Common
	Spell
	Aura
	Amount int64
}

// SPELL_AURA_REMOVED
type SpellAuraRemoved struct {
	Common
	Spell
	Aura
}

// SPELL_AURA_REMOVED_DOSE
type SpellAuraRemovedDose struct {
	Common
	Spell
	Aura
	Amount int64
}

// SPELL_AURA_DISPELLED
type SpellAuraDispelled struct {
	Common
	Spell
	Dispelled Spell
}

// SPELL_AURA_STOLEN
type SpellAuraStolen struct {
	Common
	Spell
	Stolen Spell
}

// SPELL_STOLEN
type SpellStolen struct {
	Common
	Spell
	Stolen Spell
	Aura
}

// ENCHANT_APPLIED	
type EnchantApplied struct {
	Source  Unit
	Dest    Unit
	Enchant string
	Item
}

// ENCHANT_REMOVED	
type EnchantRemoved struct {
	Source  Unit
	Dest    Unit
	Enchant string
	Item
}

// SPELL_PERIODIC_MISSED
type SpellPeriodicMissed struct {
	Common
	Spell
	Miss
}

// SPELL_PERIODIC_DAMAGE
type SpellPeriodicDamage struct {
	Common
	Spell
	Damage
}

// SPELL_PERIODIC_HEAL
type SpellPeriodicHeal struct {
	Common
	Spell
	Heal
}

// SPELL_PERIODIC_ENERGIZE
type SpellPeriodicEnergize struct {
	Common
	Spell
	Power
}

// SPELL_PERIODIC_DRAIN
type SpellPeriodicDrain struct {
	Common
	Spell
	Power
	Drained int64
}

// SPELL_PERIODIC_LEECH
type SpellPeriodicLeech struct {
	Common
	Spell
	Power
	Leeched int64
}

// SPELL_DISPEL_FAILED
type SpellDispelFailed struct {
	Common
	Spell
	Power
}

// DAMAGE_SHIELD
type DamageShield struct {
	Common
	Spell
	Damage
}

// DAMAGE_SHIELD_MISSED
type DamageShieldMissed struct {
	Common
	Spell
	Miss
}

// DAMAGE_SPLIT
type DamageSplit struct {
	Common
	Spell
	Damage
}

// SPELL_SUMMON
type SpellSummon struct {
	Common
	Spell
}

// SPELL_RESURRECT
type SpellResurrect struct {
	Common
	Spell
}

// SPELL_CREATE
type SpellCreate struct {
	Common
	Spell
}

// PARTY_KILL		
type PartyKill struct {
	Killer Unit
	Victim Unit
}

// UNIT_DIED		
type UnitDied struct {
	Victim Unit
	Empty  Unit
}

// UNIT_DESTROYED		
type UnitDestroyed struct {
	Victim Unit
	Empty  Unit
}

var eventTypes = map[string]eventFactory{
	"ENVIRONMENTAL_DAMAGE":        compile(&EnvironmentalDamage{}),
	"SWING_DAMAGE":                compile(&SwingDamage{}),
	"SWING_MISSED":                compile(&SwingMissed{}),
	"RANGE_DAMAGE":                compile(&RangeDamage{}),
	"RANGE_MISSED":                compile(&RangeMissed{}),
	"SPELL_CAST_START":            compile(&SpellCastStart{}),
	"SPELL_CAST_SUCCESS":          compile(&SpellCastSuccess{}),
	"SPELL_CAST_FAILED":           compile(&SpellCastFailed{}),
	"SPELL_MISSED":                compile(&SpellMissed{}),
	"SPELL_DAMAGE":                compile(&SpellDamage{}),
	"SPELL_BUILDING_DAMAGE":       compile(&SpellBuildingDamage{}),
	"SPELL_HEAL":                  compile(&SpellHeal{}),
	"SPELL_ENERGIZE":              compile(&SpellEnergize{}),
	"SPELL_DRAIN":                 compile(&SpellDrain{}),
	"SPELL_LEECH":                 compile(&SpellLeech{}),
	"SPELL_INSTAKILL":             compile(&SpellInstakill{}),
	"SPELL_INTERRUPT":             compile(&SpellInterrupt{}),
	"SPELL_DISPEL":                compile(&SpellDispel{}),
	"SPELL_EXTRA_ATTACKS":         compile(&SpellExtraAttacks{}),
	"SPELL_DURABILITY_DAMAGE":     compile(&SpellDurabilityDamage{}),
	"SPELL_DURABILITY_DAMAGE_ALL": compile(&SpellDurabilityDamageAll{}),
	"SPELL_AURA_APPLIED":          compile(&SpellAuraApplied{}),
	"SPELL_AURA_APPLIED_DOSE":     compile(&SpellAuraAppliedDose{}),
	"SPELL_AURA_REFRESH":          compile(&SpellAuraRefresh{}),
	"SPELL_AURA_BROKEN_SPELL":     compile(&SpellAuraBrokenSpell{}),
	"SPELL_AURA_REMOVED":          compile(&SpellAuraRemoved{}),
	"SPELL_AURA_REMOVED_DOSE":     compile(&SpellAuraRemovedDose{}),
	"SPELL_AURA_DISPELLED":        compile(&SpellAuraDispelled{}),
	"SPELL_AURA_STOLEN":           compile(&SpellAuraStolen{}),
	"SPELL_STOLEN":                compile(&SpellStolen{}),
	"ENCHANT_APPLIED":             compile(&EnchantApplied{}),
	"ENCHANT_REMOVED":             compile(&EnchantRemoved{}),
	"SPELL_PERIODIC_MISSED":       compile(&SpellPeriodicMissed{}),
	"SPELL_PERIODIC_DAMAGE":       compile(&SpellPeriodicDamage{}),
	"SPELL_PERIODIC_HEAL":         compile(&SpellPeriodicHeal{}),
	"SPELL_PERIODIC_ENERGIZE":     compile(&SpellPeriodicEnergize{}),
	"SPELL_PERIODIC_DRAIN":        compile(&SpellPeriodicDrain{}),
	"SPELL_PERIODIC_LEECH":        compile(&SpellPeriodicLeech{}),
	"SPELL_DISPEL_FAILED":         compile(&SpellDispelFailed{}),
	"DAMAGE_SHIELD":               compile(&DamageShield{}),
	"DAMAGE_SHIELD_MISSED":        compile(&DamageShieldMissed{}),
	"DAMAGE_SPLIT":                compile(&DamageSplit{}),
	"SPELL_RESURRECT":             compile(&SpellResurrect{}),
	"SPELL_SUMMON":                compile(&SpellSummon{}),
	"SPELL_CREATE":                compile(&SpellCreate{}),
	"PARTY_KILL":                  compile(&PartyKill{}),
	"UNIT_DIED":                   compile(&UnitDied{}),
	"UNIT_DESTROYED":              compile(&UnitDestroyed{}),
}

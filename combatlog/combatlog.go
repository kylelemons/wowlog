package combatlog

type Unit struct {
	ID    uint64
	Name  string
	Flags uint64
	Null  int32
}

type Spell struct {
	ID     uint64
	Name   string
	School int32
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

type Heal struct {
	Amount   int64
	Overheal int64
	Unknown  int64
	Critical bool
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
	Source Unit
	Dest   Unit
	Type   string
	Damage
}

// SWING_DAMAGE
type SwingDamage struct {
	Source Unit
	Dest   Unit
	Damage
}

// SWING_MISSED
type SwingMissed struct {
	Source Unit
	Dest   Unit
	Miss
}

// RANGE_DAMAGE
type RangeDamage struct {
	Source Unit
	Dest   Unit
	Spell
	Damage
}

// RANGE_MISSED
type RangeMissed struct {
	Source Unit
	Dest   Unit
	Spell
	Miss
}

// SPELL_CAST_START
type SpellCastStart struct {
	Source Unit
	Dest   Unit
	Spell
}

// SPELL_CAST_SUCCESS
type SpellCastSuccess struct {
	Source Unit
	Dest   Unit
	Spell
}

// SPELL_CAST_FAILED
type SpellCastFailed struct {
	Source Unit
	Dest   Unit
	Spell
	Miss
}

// SPELL_MISSED
type SpellMissed struct {
	Source Unit
	Dest   Unit
	Spell
	Miss
}

// SPELL_DAMAGE
type SpellDamage struct {
	Source Unit
	Dest   Unit
	Spell
	Damage
}

// SPELL_BUILDING_DAMAGE
type SpellBuildingDamage struct {
	Source Unit
	Dest   Unit
	Spell
	Damage
}

// SPELL_HEAL
type SpellHeal struct {
	Source Unit
	Dest   Unit
	Spell
	Heal
}

// SPELL_ENERGIZE
type SpellEnergize struct {
	Source Unit
	Dest   Unit
	Spell
	Power
}

// SPELL_DRAIN
type SpellDrain struct {
	Source Unit
	Dest   Unit
	Spell
	Power
	Drained int64
}

// SPELL_LEECH
type SpellLeech struct {
	Source Unit
	Dest   Unit
	Spell
	Power
	Leeched int64
}

// SPELL_INSTAKILL		
type SpellInstakill struct {
	Killer Unit
	Victim Unit
	Spell
}

// SPELL_INTERRUPT
type SpellInterrupt struct {
	Source Unit
	Dest   Unit
	Spell
	Interrupted Spell
}

// SPELL_DISPEL
type SpellDispel struct {
	Source Unit
	Dest   Unit
	Spell
	Dispelled Spell
	Aura
}

// SPELL_EXTRA_ATTACKS
type SpellExtraAttacks struct {
	Source Unit
	Dest   Unit
	Spell
	Amount int64
}

// SPELL_DURABILITY_DAMAGE
type SpellDurabilityDamage struct {
	Source Unit
	Dest   Unit
	Spell
}

// SPELL_DURABILITY_DAMAGE_ALL
type SpellDurabilityDamageAll struct {
	Source Unit
	Dest   Unit
	Spell
}

// SPELL_AURA_APPLIED
type SpellAuraApplied struct {
	Source Unit
	Dest   Unit
	Spell
	Aura
}

// SPELL_AURA_REFRESH
type SpellAuraRefresh struct {
	Source Unit
	Dest   Unit
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
	Source Unit
	Dest   Unit
	Spell
	Aura
	Amount int64
}

// SPELL_AURA_REMOVED
type SpellAuraRemoved struct {
	Source Unit
	Dest   Unit
	Spell
	Aura
}

// SPELL_AURA_REMOVED_DOSE
type SpellAuraRemovedDose struct {
	Source Unit
	Dest   Unit
	Spell
	Aura
	Amount int64
}

// SPELL_AURA_DISPELLED
type SpellAuraDispelled struct {
	Source Unit
	Dest   Unit
	Spell
	Dispelled Spell
}

// SPELL_AURA_STOLEN
type SpellAuraStolen struct {
	Source Unit
	Dest   Unit
	Spell
	Stolen Spell
}

// SPELL_STOLEN
type SpellStolen struct {
	Source Unit
	Dest   Unit
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
	Source Unit
	Dest   Unit
	Spell
	Miss
}

// SPELL_PERIODIC_DAMAGE
type SpellPeriodicDamage struct {
	Source Unit
	Dest   Unit
	Spell
	Damage
}

// SPELL_PERIODIC_HEAL
type SpellPeriodicHeal struct {
	Source Unit
	Dest   Unit
	Spell
	Heal
}

// SPELL_PERIODIC_ENERGIZE
type SpellPeriodicEnergize struct {
	Source Unit
	Dest   Unit
	Spell
	Power
}

// SPELL_PERIODIC_DRAIN
type SpellPeriodicDrain struct {
	Source Unit
	Dest   Unit
	Spell
	Power
	Drained int64
}

// SPELL_PERIODIC_LEECH
type SpellPeriodicLeech struct {
	Source Unit
	Dest   Unit
	Spell
	Power
	Leeched int64
}

// SPELL_DISPEL_FAILED
type SpellDispelFailed struct {
	Source Unit
	Dest   Unit
	Spell
	Power
}

// DAMAGE_SHIELD
type DamageShield struct {
	Source Unit
	Dest   Unit
	Spell
	Damage
}

// DAMAGE_SHIELD_MISSED
type DamageShieldMissed struct {
	Source Unit
	Dest   Unit
	Spell
	Miss
}

// DAMAGE_SPLIT
type DamageSplit struct {
	Source Unit
	Dest   Unit
	Spell
	Damage
}

// SPELL_SUMMON
type SpellSummon struct {
	Source Unit
	Dest   Unit
	Spell
}

// SPELL_RESURRECT
type SpellResurrect struct {
	Source Unit
	Dest   Unit
	Spell
}

// SPELL_CREATE
type SpellCreate struct {
	Source Unit
	Dest Unit
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

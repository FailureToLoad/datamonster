// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/failuretoload/datamonster/ent/settlement"
	"github.com/failuretoload/datamonster/ent/survivor"
)

// SurvivorCreate is the builder for creating a Survivor entity.
type SurvivorCreate struct {
	config
	mutation *SurvivorMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (sc *SurvivorCreate) SetName(s string) *SurvivorCreate {
	sc.mutation.SetName(s)
	return sc
}

// SetBorn sets the "born" field.
func (sc *SurvivorCreate) SetBorn(i int) *SurvivorCreate {
	sc.mutation.SetBorn(i)
	return sc
}

// SetNillableBorn sets the "born" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableBorn(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetBorn(*i)
	}
	return sc
}

// SetGender sets the "gender" field.
func (sc *SurvivorCreate) SetGender(s survivor.Gender) *SurvivorCreate {
	sc.mutation.SetGender(s)
	return sc
}

// SetNillableGender sets the "gender" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableGender(s *survivor.Gender) *SurvivorCreate {
	if s != nil {
		sc.SetGender(*s)
	}
	return sc
}

// SetHuntxp sets the "huntxp" field.
func (sc *SurvivorCreate) SetHuntxp(i int) *SurvivorCreate {
	sc.mutation.SetHuntxp(i)
	return sc
}

// SetNillableHuntxp sets the "huntxp" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableHuntxp(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetHuntxp(*i)
	}
	return sc
}

// SetSurvival sets the "survival" field.
func (sc *SurvivorCreate) SetSurvival(i int) *SurvivorCreate {
	sc.mutation.SetSurvival(i)
	return sc
}

// SetNillableSurvival sets the "survival" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableSurvival(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetSurvival(*i)
	}
	return sc
}

// SetMovement sets the "movement" field.
func (sc *SurvivorCreate) SetMovement(i int) *SurvivorCreate {
	sc.mutation.SetMovement(i)
	return sc
}

// SetNillableMovement sets the "movement" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableMovement(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetMovement(*i)
	}
	return sc
}

// SetAccuracy sets the "accuracy" field.
func (sc *SurvivorCreate) SetAccuracy(i int) *SurvivorCreate {
	sc.mutation.SetAccuracy(i)
	return sc
}

// SetNillableAccuracy sets the "accuracy" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableAccuracy(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetAccuracy(*i)
	}
	return sc
}

// SetStrength sets the "strength" field.
func (sc *SurvivorCreate) SetStrength(i int) *SurvivorCreate {
	sc.mutation.SetStrength(i)
	return sc
}

// SetNillableStrength sets the "strength" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableStrength(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetStrength(*i)
	}
	return sc
}

// SetEvasion sets the "evasion" field.
func (sc *SurvivorCreate) SetEvasion(i int) *SurvivorCreate {
	sc.mutation.SetEvasion(i)
	return sc
}

// SetNillableEvasion sets the "evasion" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableEvasion(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetEvasion(*i)
	}
	return sc
}

// SetLuck sets the "luck" field.
func (sc *SurvivorCreate) SetLuck(i int) *SurvivorCreate {
	sc.mutation.SetLuck(i)
	return sc
}

// SetNillableLuck sets the "luck" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableLuck(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetLuck(*i)
	}
	return sc
}

// SetSpeed sets the "speed" field.
func (sc *SurvivorCreate) SetSpeed(i int) *SurvivorCreate {
	sc.mutation.SetSpeed(i)
	return sc
}

// SetNillableSpeed sets the "speed" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableSpeed(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetSpeed(*i)
	}
	return sc
}

// SetSystemicpressure sets the "systemicpressure" field.
func (sc *SurvivorCreate) SetSystemicpressure(i int) *SurvivorCreate {
	sc.mutation.SetSystemicpressure(i)
	return sc
}

// SetNillableSystemicpressure sets the "systemicpressure" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableSystemicpressure(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetSystemicpressure(*i)
	}
	return sc
}

// SetTorment sets the "torment" field.
func (sc *SurvivorCreate) SetTorment(i int) *SurvivorCreate {
	sc.mutation.SetTorment(i)
	return sc
}

// SetNillableTorment sets the "torment" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableTorment(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetTorment(*i)
	}
	return sc
}

// SetInsanity sets the "insanity" field.
func (sc *SurvivorCreate) SetInsanity(i int) *SurvivorCreate {
	sc.mutation.SetInsanity(i)
	return sc
}

// SetNillableInsanity sets the "insanity" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableInsanity(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetInsanity(*i)
	}
	return sc
}

// SetLumi sets the "lumi" field.
func (sc *SurvivorCreate) SetLumi(i int) *SurvivorCreate {
	sc.mutation.SetLumi(i)
	return sc
}

// SetNillableLumi sets the "lumi" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableLumi(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetLumi(*i)
	}
	return sc
}

// SetCourage sets the "courage" field.
func (sc *SurvivorCreate) SetCourage(i int) *SurvivorCreate {
	sc.mutation.SetCourage(i)
	return sc
}

// SetNillableCourage sets the "courage" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableCourage(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetCourage(*i)
	}
	return sc
}

// SetUnderstanding sets the "understanding" field.
func (sc *SurvivorCreate) SetUnderstanding(i int) *SurvivorCreate {
	sc.mutation.SetUnderstanding(i)
	return sc
}

// SetNillableUnderstanding sets the "understanding" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableUnderstanding(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetUnderstanding(*i)
	}
	return sc
}

// SetSettlementID sets the "settlement_id" field.
func (sc *SurvivorCreate) SetSettlementID(i int) *SurvivorCreate {
	sc.mutation.SetSettlementID(i)
	return sc
}

// SetNillableSettlementID sets the "settlement_id" field if the given value is not nil.
func (sc *SurvivorCreate) SetNillableSettlementID(i *int) *SurvivorCreate {
	if i != nil {
		sc.SetSettlementID(*i)
	}
	return sc
}

// SetSettlement sets the "settlement" edge to the Settlement entity.
func (sc *SurvivorCreate) SetSettlement(s *Settlement) *SurvivorCreate {
	return sc.SetSettlementID(s.ID)
}

// Mutation returns the SurvivorMutation object of the builder.
func (sc *SurvivorCreate) Mutation() *SurvivorMutation {
	return sc.mutation
}

// Save creates the Survivor in the database.
func (sc *SurvivorCreate) Save(ctx context.Context) (*Survivor, error) {
	sc.defaults()
	return withHooks(ctx, sc.sqlSave, sc.mutation, sc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sc *SurvivorCreate) SaveX(ctx context.Context) *Survivor {
	v, err := sc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sc *SurvivorCreate) Exec(ctx context.Context) error {
	_, err := sc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sc *SurvivorCreate) ExecX(ctx context.Context) {
	if err := sc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (sc *SurvivorCreate) defaults() {
	if _, ok := sc.mutation.Born(); !ok {
		v := survivor.DefaultBorn
		sc.mutation.SetBorn(v)
	}
	if _, ok := sc.mutation.Gender(); !ok {
		v := survivor.DefaultGender
		sc.mutation.SetGender(v)
	}
	if _, ok := sc.mutation.Huntxp(); !ok {
		v := survivor.DefaultHuntxp
		sc.mutation.SetHuntxp(v)
	}
	if _, ok := sc.mutation.Survival(); !ok {
		v := survivor.DefaultSurvival
		sc.mutation.SetSurvival(v)
	}
	if _, ok := sc.mutation.Movement(); !ok {
		v := survivor.DefaultMovement
		sc.mutation.SetMovement(v)
	}
	if _, ok := sc.mutation.Accuracy(); !ok {
		v := survivor.DefaultAccuracy
		sc.mutation.SetAccuracy(v)
	}
	if _, ok := sc.mutation.Strength(); !ok {
		v := survivor.DefaultStrength
		sc.mutation.SetStrength(v)
	}
	if _, ok := sc.mutation.Evasion(); !ok {
		v := survivor.DefaultEvasion
		sc.mutation.SetEvasion(v)
	}
	if _, ok := sc.mutation.Luck(); !ok {
		v := survivor.DefaultLuck
		sc.mutation.SetLuck(v)
	}
	if _, ok := sc.mutation.Speed(); !ok {
		v := survivor.DefaultSpeed
		sc.mutation.SetSpeed(v)
	}
	if _, ok := sc.mutation.Systemicpressure(); !ok {
		v := survivor.DefaultSystemicpressure
		sc.mutation.SetSystemicpressure(v)
	}
	if _, ok := sc.mutation.Torment(); !ok {
		v := survivor.DefaultTorment
		sc.mutation.SetTorment(v)
	}
	if _, ok := sc.mutation.Insanity(); !ok {
		v := survivor.DefaultInsanity
		sc.mutation.SetInsanity(v)
	}
	if _, ok := sc.mutation.Lumi(); !ok {
		v := survivor.DefaultLumi
		sc.mutation.SetLumi(v)
	}
	if _, ok := sc.mutation.Courage(); !ok {
		v := survivor.DefaultCourage
		sc.mutation.SetCourage(v)
	}
	if _, ok := sc.mutation.Understanding(); !ok {
		v := survivor.DefaultUnderstanding
		sc.mutation.SetUnderstanding(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sc *SurvivorCreate) check() error {
	if _, ok := sc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "Survivor.name"`)}
	}
	if v, ok := sc.mutation.Name(); ok {
		if err := survivor.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf(`ent: validator failed for field "Survivor.name": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Born(); !ok {
		return &ValidationError{Name: "born", err: errors.New(`ent: missing required field "Survivor.born"`)}
	}
	if v, ok := sc.mutation.Born(); ok {
		if err := survivor.BornValidator(v); err != nil {
			return &ValidationError{Name: "born", err: fmt.Errorf(`ent: validator failed for field "Survivor.born": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Gender(); !ok {
		return &ValidationError{Name: "gender", err: errors.New(`ent: missing required field "Survivor.gender"`)}
	}
	if v, ok := sc.mutation.Gender(); ok {
		if err := survivor.GenderValidator(v); err != nil {
			return &ValidationError{Name: "gender", err: fmt.Errorf(`ent: validator failed for field "Survivor.gender": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Huntxp(); !ok {
		return &ValidationError{Name: "huntxp", err: errors.New(`ent: missing required field "Survivor.huntxp"`)}
	}
	if v, ok := sc.mutation.Huntxp(); ok {
		if err := survivor.HuntxpValidator(v); err != nil {
			return &ValidationError{Name: "huntxp", err: fmt.Errorf(`ent: validator failed for field "Survivor.huntxp": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Survival(); !ok {
		return &ValidationError{Name: "survival", err: errors.New(`ent: missing required field "Survivor.survival"`)}
	}
	if v, ok := sc.mutation.Survival(); ok {
		if err := survivor.SurvivalValidator(v); err != nil {
			return &ValidationError{Name: "survival", err: fmt.Errorf(`ent: validator failed for field "Survivor.survival": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Movement(); !ok {
		return &ValidationError{Name: "movement", err: errors.New(`ent: missing required field "Survivor.movement"`)}
	}
	if v, ok := sc.mutation.Movement(); ok {
		if err := survivor.MovementValidator(v); err != nil {
			return &ValidationError{Name: "movement", err: fmt.Errorf(`ent: validator failed for field "Survivor.movement": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Accuracy(); !ok {
		return &ValidationError{Name: "accuracy", err: errors.New(`ent: missing required field "Survivor.accuracy"`)}
	}
	if v, ok := sc.mutation.Accuracy(); ok {
		if err := survivor.AccuracyValidator(v); err != nil {
			return &ValidationError{Name: "accuracy", err: fmt.Errorf(`ent: validator failed for field "Survivor.accuracy": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Strength(); !ok {
		return &ValidationError{Name: "strength", err: errors.New(`ent: missing required field "Survivor.strength"`)}
	}
	if v, ok := sc.mutation.Strength(); ok {
		if err := survivor.StrengthValidator(v); err != nil {
			return &ValidationError{Name: "strength", err: fmt.Errorf(`ent: validator failed for field "Survivor.strength": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Evasion(); !ok {
		return &ValidationError{Name: "evasion", err: errors.New(`ent: missing required field "Survivor.evasion"`)}
	}
	if v, ok := sc.mutation.Evasion(); ok {
		if err := survivor.EvasionValidator(v); err != nil {
			return &ValidationError{Name: "evasion", err: fmt.Errorf(`ent: validator failed for field "Survivor.evasion": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Luck(); !ok {
		return &ValidationError{Name: "luck", err: errors.New(`ent: missing required field "Survivor.luck"`)}
	}
	if v, ok := sc.mutation.Luck(); ok {
		if err := survivor.LuckValidator(v); err != nil {
			return &ValidationError{Name: "luck", err: fmt.Errorf(`ent: validator failed for field "Survivor.luck": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Speed(); !ok {
		return &ValidationError{Name: "speed", err: errors.New(`ent: missing required field "Survivor.speed"`)}
	}
	if v, ok := sc.mutation.Speed(); ok {
		if err := survivor.SpeedValidator(v); err != nil {
			return &ValidationError{Name: "speed", err: fmt.Errorf(`ent: validator failed for field "Survivor.speed": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Systemicpressure(); !ok {
		return &ValidationError{Name: "systemicpressure", err: errors.New(`ent: missing required field "Survivor.systemicpressure"`)}
	}
	if v, ok := sc.mutation.Systemicpressure(); ok {
		if err := survivor.SystemicpressureValidator(v); err != nil {
			return &ValidationError{Name: "systemicpressure", err: fmt.Errorf(`ent: validator failed for field "Survivor.systemicpressure": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Torment(); !ok {
		return &ValidationError{Name: "torment", err: errors.New(`ent: missing required field "Survivor.torment"`)}
	}
	if v, ok := sc.mutation.Torment(); ok {
		if err := survivor.TormentValidator(v); err != nil {
			return &ValidationError{Name: "torment", err: fmt.Errorf(`ent: validator failed for field "Survivor.torment": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Insanity(); !ok {
		return &ValidationError{Name: "insanity", err: errors.New(`ent: missing required field "Survivor.insanity"`)}
	}
	if v, ok := sc.mutation.Insanity(); ok {
		if err := survivor.InsanityValidator(v); err != nil {
			return &ValidationError{Name: "insanity", err: fmt.Errorf(`ent: validator failed for field "Survivor.insanity": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Lumi(); !ok {
		return &ValidationError{Name: "lumi", err: errors.New(`ent: missing required field "Survivor.lumi"`)}
	}
	if v, ok := sc.mutation.Lumi(); ok {
		if err := survivor.LumiValidator(v); err != nil {
			return &ValidationError{Name: "lumi", err: fmt.Errorf(`ent: validator failed for field "Survivor.lumi": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Courage(); !ok {
		return &ValidationError{Name: "courage", err: errors.New(`ent: missing required field "Survivor.courage"`)}
	}
	if v, ok := sc.mutation.Courage(); ok {
		if err := survivor.CourageValidator(v); err != nil {
			return &ValidationError{Name: "courage", err: fmt.Errorf(`ent: validator failed for field "Survivor.courage": %w`, err)}
		}
	}
	if _, ok := sc.mutation.Understanding(); !ok {
		return &ValidationError{Name: "understanding", err: errors.New(`ent: missing required field "Survivor.understanding"`)}
	}
	if v, ok := sc.mutation.Understanding(); ok {
		if err := survivor.UnderstandingValidator(v); err != nil {
			return &ValidationError{Name: "understanding", err: fmt.Errorf(`ent: validator failed for field "Survivor.understanding": %w`, err)}
		}
	}
	return nil
}

func (sc *SurvivorCreate) sqlSave(ctx context.Context) (*Survivor, error) {
	if err := sc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sc.mutation.id = &_node.ID
	sc.mutation.done = true
	return _node, nil
}

func (sc *SurvivorCreate) createSpec() (*Survivor, *sqlgraph.CreateSpec) {
	var (
		_node = &Survivor{config: sc.config}
		_spec = sqlgraph.NewCreateSpec(survivor.Table, sqlgraph.NewFieldSpec(survivor.FieldID, field.TypeInt))
	)
	if value, ok := sc.mutation.Name(); ok {
		_spec.SetField(survivor.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := sc.mutation.Born(); ok {
		_spec.SetField(survivor.FieldBorn, field.TypeInt, value)
		_node.Born = value
	}
	if value, ok := sc.mutation.Gender(); ok {
		_spec.SetField(survivor.FieldGender, field.TypeEnum, value)
		_node.Gender = value
	}
	if value, ok := sc.mutation.Huntxp(); ok {
		_spec.SetField(survivor.FieldHuntxp, field.TypeInt, value)
		_node.Huntxp = value
	}
	if value, ok := sc.mutation.Survival(); ok {
		_spec.SetField(survivor.FieldSurvival, field.TypeInt, value)
		_node.Survival = value
	}
	if value, ok := sc.mutation.Movement(); ok {
		_spec.SetField(survivor.FieldMovement, field.TypeInt, value)
		_node.Movement = value
	}
	if value, ok := sc.mutation.Accuracy(); ok {
		_spec.SetField(survivor.FieldAccuracy, field.TypeInt, value)
		_node.Accuracy = value
	}
	if value, ok := sc.mutation.Strength(); ok {
		_spec.SetField(survivor.FieldStrength, field.TypeInt, value)
		_node.Strength = value
	}
	if value, ok := sc.mutation.Evasion(); ok {
		_spec.SetField(survivor.FieldEvasion, field.TypeInt, value)
		_node.Evasion = value
	}
	if value, ok := sc.mutation.Luck(); ok {
		_spec.SetField(survivor.FieldLuck, field.TypeInt, value)
		_node.Luck = value
	}
	if value, ok := sc.mutation.Speed(); ok {
		_spec.SetField(survivor.FieldSpeed, field.TypeInt, value)
		_node.Speed = value
	}
	if value, ok := sc.mutation.Systemicpressure(); ok {
		_spec.SetField(survivor.FieldSystemicpressure, field.TypeInt, value)
		_node.Systemicpressure = value
	}
	if value, ok := sc.mutation.Torment(); ok {
		_spec.SetField(survivor.FieldTorment, field.TypeInt, value)
		_node.Torment = value
	}
	if value, ok := sc.mutation.Insanity(); ok {
		_spec.SetField(survivor.FieldInsanity, field.TypeInt, value)
		_node.Insanity = value
	}
	if value, ok := sc.mutation.Lumi(); ok {
		_spec.SetField(survivor.FieldLumi, field.TypeInt, value)
		_node.Lumi = value
	}
	if value, ok := sc.mutation.Courage(); ok {
		_spec.SetField(survivor.FieldCourage, field.TypeInt, value)
		_node.Courage = value
	}
	if value, ok := sc.mutation.Understanding(); ok {
		_spec.SetField(survivor.FieldUnderstanding, field.TypeInt, value)
		_node.Understanding = value
	}
	if nodes := sc.mutation.SettlementIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   survivor.SettlementTable,
			Columns: []string{survivor.SettlementColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(settlement.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.SettlementID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// SurvivorCreateBulk is the builder for creating many Survivor entities in bulk.
type SurvivorCreateBulk struct {
	config
	err      error
	builders []*SurvivorCreate
}

// Save creates the Survivor entities in the database.
func (scb *SurvivorCreateBulk) Save(ctx context.Context) ([]*Survivor, error) {
	if scb.err != nil {
		return nil, scb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(scb.builders))
	nodes := make([]*Survivor, len(scb.builders))
	mutators := make([]Mutator, len(scb.builders))
	for i := range scb.builders {
		func(i int, root context.Context) {
			builder := scb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SurvivorMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, scb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, scb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, scb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (scb *SurvivorCreateBulk) SaveX(ctx context.Context) []*Survivor {
	v, err := scb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (scb *SurvivorCreateBulk) Exec(ctx context.Context) error {
	_, err := scb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (scb *SurvivorCreateBulk) ExecX(ctx context.Context) {
	if err := scb.Exec(ctx); err != nil {
		panic(err)
	}
}

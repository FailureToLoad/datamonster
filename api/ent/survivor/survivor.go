// Code generated by ent, DO NOT EDIT.

package survivor

import (
	"fmt"
	"io"
	"strconv"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the survivor type in the database.
	Label = "survivor"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldBorn holds the string denoting the born field in the database.
	FieldBorn = "born"
	// FieldGender holds the string denoting the gender field in the database.
	FieldGender = "gender"
	// FieldHuntxp holds the string denoting the huntxp field in the database.
	FieldHuntxp = "huntxp"
	// FieldSurvival holds the string denoting the survival field in the database.
	FieldSurvival = "survival"
	// FieldMovement holds the string denoting the movement field in the database.
	FieldMovement = "movement"
	// FieldAccuracy holds the string denoting the accuracy field in the database.
	FieldAccuracy = "accuracy"
	// FieldStrength holds the string denoting the strength field in the database.
	FieldStrength = "strength"
	// FieldEvasion holds the string denoting the evasion field in the database.
	FieldEvasion = "evasion"
	// FieldLuck holds the string denoting the luck field in the database.
	FieldLuck = "luck"
	// FieldSpeed holds the string denoting the speed field in the database.
	FieldSpeed = "speed"
	// FieldSystemicpressure holds the string denoting the systemicpressure field in the database.
	FieldSystemicpressure = "systemicpressure"
	// FieldTorment holds the string denoting the torment field in the database.
	FieldTorment = "torment"
	// FieldInsanity holds the string denoting the insanity field in the database.
	FieldInsanity = "insanity"
	// FieldLumi holds the string denoting the lumi field in the database.
	FieldLumi = "lumi"
	// FieldCourage holds the string denoting the courage field in the database.
	FieldCourage = "courage"
	// FieldUnderstanding holds the string denoting the understanding field in the database.
	FieldUnderstanding = "understanding"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldStatusChangeYear holds the string denoting the status_change_year field in the database.
	FieldStatusChangeYear = "status_change_year"
	// FieldSettlementID holds the string denoting the settlement_id field in the database.
	FieldSettlementID = "settlement_id"
	// EdgeSettlement holds the string denoting the settlement edge name in mutations.
	EdgeSettlement = "settlement"
	// Table holds the table name of the survivor in the database.
	Table = "survivors"
	// SettlementTable is the table that holds the settlement relation/edge.
	SettlementTable = "survivors"
	// SettlementInverseTable is the table name for the Settlement entity.
	// It exists in this package in order to avoid circular dependency with the "settlement" package.
	SettlementInverseTable = "settlements"
	// SettlementColumn is the table column denoting the settlement relation/edge.
	SettlementColumn = "settlement_id"
)

// Columns holds all SQL columns for survivor fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldBorn,
	FieldGender,
	FieldHuntxp,
	FieldSurvival,
	FieldMovement,
	FieldAccuracy,
	FieldStrength,
	FieldEvasion,
	FieldLuck,
	FieldSpeed,
	FieldSystemicpressure,
	FieldTorment,
	FieldInsanity,
	FieldLumi,
	FieldCourage,
	FieldUnderstanding,
	FieldStatus,
	FieldStatusChangeYear,
	FieldSettlementID,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// DefaultBorn holds the default value on creation for the "born" field.
	DefaultBorn int
	// BornValidator is a validator for the "born" field. It is called by the builders before save.
	BornValidator func(int) error
	// DefaultHuntxp holds the default value on creation for the "huntxp" field.
	DefaultHuntxp int
	// HuntxpValidator is a validator for the "huntxp" field. It is called by the builders before save.
	HuntxpValidator func(int) error
	// DefaultSurvival holds the default value on creation for the "survival" field.
	DefaultSurvival int
	// SurvivalValidator is a validator for the "survival" field. It is called by the builders before save.
	SurvivalValidator func(int) error
	// DefaultMovement holds the default value on creation for the "movement" field.
	DefaultMovement int
	// MovementValidator is a validator for the "movement" field. It is called by the builders before save.
	MovementValidator func(int) error
	// DefaultAccuracy holds the default value on creation for the "accuracy" field.
	DefaultAccuracy int
	// AccuracyValidator is a validator for the "accuracy" field. It is called by the builders before save.
	AccuracyValidator func(int) error
	// DefaultStrength holds the default value on creation for the "strength" field.
	DefaultStrength int
	// StrengthValidator is a validator for the "strength" field. It is called by the builders before save.
	StrengthValidator func(int) error
	// DefaultEvasion holds the default value on creation for the "evasion" field.
	DefaultEvasion int
	// EvasionValidator is a validator for the "evasion" field. It is called by the builders before save.
	EvasionValidator func(int) error
	// DefaultLuck holds the default value on creation for the "luck" field.
	DefaultLuck int
	// LuckValidator is a validator for the "luck" field. It is called by the builders before save.
	LuckValidator func(int) error
	// DefaultSpeed holds the default value on creation for the "speed" field.
	DefaultSpeed int
	// SpeedValidator is a validator for the "speed" field. It is called by the builders before save.
	SpeedValidator func(int) error
	// DefaultSystemicpressure holds the default value on creation for the "systemicpressure" field.
	DefaultSystemicpressure int
	// SystemicpressureValidator is a validator for the "systemicpressure" field. It is called by the builders before save.
	SystemicpressureValidator func(int) error
	// DefaultTorment holds the default value on creation for the "torment" field.
	DefaultTorment int
	// TormentValidator is a validator for the "torment" field. It is called by the builders before save.
	TormentValidator func(int) error
	// DefaultInsanity holds the default value on creation for the "insanity" field.
	DefaultInsanity int
	// InsanityValidator is a validator for the "insanity" field. It is called by the builders before save.
	InsanityValidator func(int) error
	// DefaultLumi holds the default value on creation for the "lumi" field.
	DefaultLumi int
	// LumiValidator is a validator for the "lumi" field. It is called by the builders before save.
	LumiValidator func(int) error
	// DefaultCourage holds the default value on creation for the "courage" field.
	DefaultCourage int
	// CourageValidator is a validator for the "courage" field. It is called by the builders before save.
	CourageValidator func(int) error
	// DefaultUnderstanding holds the default value on creation for the "understanding" field.
	DefaultUnderstanding int
	// UnderstandingValidator is a validator for the "understanding" field. It is called by the builders before save.
	UnderstandingValidator func(int) error
	// DefaultStatusChangeYear holds the default value on creation for the "status_change_year" field.
	DefaultStatusChangeYear int
)

// Gender defines the type for the "gender" enum field.
type Gender string

// GenderFemale is the default value of the Gender enum.
const DefaultGender = GenderFemale

// Gender values.
const (
	GenderMale   Gender = "M"
	GenderFemale Gender = "F"
)

func (ge Gender) String() string {
	return string(ge)
}

// GenderValidator is a validator for the "gender" field enum values. It is called by the builders before save.
func GenderValidator(ge Gender) error {
	switch ge {
	case GenderMale, GenderFemale:
		return nil
	default:
		return fmt.Errorf("survivor: invalid enum value for gender field: %q", ge)
	}
}

// Status defines the type for the "status" enum field.
type Status string

// StatusAlive is the default value of the Status enum.
const DefaultStatus = StatusAlive

// Status values.
const (
	StatusAlive         Status = "alive"
	StatusDead          Status = "dead"
	StatusCeasedToExist Status = "ceased_to_exist"
	StatusRetired       Status = "retired"
	StatusSkipHunt      Status = "skip_hunt"
)

func (s Status) String() string {
	return string(s)
}

// StatusValidator is a validator for the "status" field enum values. It is called by the builders before save.
func StatusValidator(s Status) error {
	switch s {
	case StatusAlive, StatusDead, StatusCeasedToExist, StatusRetired, StatusSkipHunt:
		return nil
	default:
		return fmt.Errorf("survivor: invalid enum value for status field: %q", s)
	}
}

// OrderOption defines the ordering options for the Survivor queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByBorn orders the results by the born field.
func ByBorn(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBorn, opts...).ToFunc()
}

// ByGender orders the results by the gender field.
func ByGender(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldGender, opts...).ToFunc()
}

// ByHuntxp orders the results by the huntxp field.
func ByHuntxp(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHuntxp, opts...).ToFunc()
}

// BySurvival orders the results by the survival field.
func BySurvival(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSurvival, opts...).ToFunc()
}

// ByMovement orders the results by the movement field.
func ByMovement(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMovement, opts...).ToFunc()
}

// ByAccuracy orders the results by the accuracy field.
func ByAccuracy(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAccuracy, opts...).ToFunc()
}

// ByStrength orders the results by the strength field.
func ByStrength(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStrength, opts...).ToFunc()
}

// ByEvasion orders the results by the evasion field.
func ByEvasion(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEvasion, opts...).ToFunc()
}

// ByLuck orders the results by the luck field.
func ByLuck(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLuck, opts...).ToFunc()
}

// BySpeed orders the results by the speed field.
func BySpeed(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSpeed, opts...).ToFunc()
}

// BySystemicpressure orders the results by the systemicpressure field.
func BySystemicpressure(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSystemicpressure, opts...).ToFunc()
}

// ByTorment orders the results by the torment field.
func ByTorment(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTorment, opts...).ToFunc()
}

// ByInsanity orders the results by the insanity field.
func ByInsanity(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInsanity, opts...).ToFunc()
}

// ByLumi orders the results by the lumi field.
func ByLumi(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLumi, opts...).ToFunc()
}

// ByCourage orders the results by the courage field.
func ByCourage(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCourage, opts...).ToFunc()
}

// ByUnderstanding orders the results by the understanding field.
func ByUnderstanding(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUnderstanding, opts...).ToFunc()
}

// ByStatus orders the results by the status field.
func ByStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatus, opts...).ToFunc()
}

// ByStatusChangeYear orders the results by the status_change_year field.
func ByStatusChangeYear(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldStatusChangeYear, opts...).ToFunc()
}

// BySettlementID orders the results by the settlement_id field.
func BySettlementID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSettlementID, opts...).ToFunc()
}

// BySettlementField orders the results by settlement field.
func BySettlementField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newSettlementStep(), sql.OrderByField(field, opts...))
	}
}
func newSettlementStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(SettlementInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, SettlementTable, SettlementColumn),
	)
}

// MarshalGQL implements graphql.Marshaler interface.
func (e Gender) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(e.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (e *Gender) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*e = Gender(str)
	if err := GenderValidator(*e); err != nil {
		return fmt.Errorf("%s is not a valid Gender", str)
	}
	return nil
}

// MarshalGQL implements graphql.Marshaler interface.
func (e Status) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(e.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (e *Status) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*e = Status(str)
	if err := StatusValidator(*e); err != nil {
		return fmt.Errorf("%s is not a valid Status", str)
	}
	return nil
}

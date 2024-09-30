// Code generated by ent, DO NOT EDIT.

package outputgroup

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEQ(FieldName, v))
}

// Incomplete applies equality check predicate on the "incomplete" field. It's identical to IncompleteEQ.
func Incomplete(v bool) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEQ(FieldIncomplete, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldHasSuffix(FieldName, v))
}

// NameIsNil applies the IsNil predicate on the "name" field.
func NameIsNil() predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldIsNull(FieldName))
}

// NameNotNil applies the NotNil predicate on the "name" field.
func NameNotNil() predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNotNull(FieldName))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldContainsFold(FieldName, v))
}

// IncompleteEQ applies the EQ predicate on the "incomplete" field.
func IncompleteEQ(v bool) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldEQ(FieldIncomplete, v))
}

// IncompleteNEQ applies the NEQ predicate on the "incomplete" field.
func IncompleteNEQ(v bool) predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNEQ(FieldIncomplete, v))
}

// IncompleteIsNil applies the IsNil predicate on the "incomplete" field.
func IncompleteIsNil() predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldIsNull(FieldIncomplete))
}

// IncompleteNotNil applies the NotNil predicate on the "incomplete" field.
func IncompleteNotNil() predicate.OutputGroup {
	return predicate.OutputGroup(sql.FieldNotNull(FieldIncomplete))
}

// HasTargetComplete applies the HasEdge predicate on the "target_complete" edge.
func HasTargetComplete() predicate.OutputGroup {
	return predicate.OutputGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, TargetCompleteTable, TargetCompleteColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTargetCompleteWith applies the HasEdge predicate on the "target_complete" edge with a given conditions (other predicates).
func HasTargetCompleteWith(preds ...predicate.TargetComplete) predicate.OutputGroup {
	return predicate.OutputGroup(func(s *sql.Selector) {
		step := newTargetCompleteStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasInlineFiles applies the HasEdge predicate on the "inline_files" edge.
func HasInlineFiles() predicate.OutputGroup {
	return predicate.OutputGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, InlineFilesTable, InlineFilesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasInlineFilesWith applies the HasEdge predicate on the "inline_files" edge with a given conditions (other predicates).
func HasInlineFilesWith(preds ...predicate.TestFile) predicate.OutputGroup {
	return predicate.OutputGroup(func(s *sql.Selector) {
		step := newInlineFilesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasFileSets applies the HasEdge predicate on the "file_sets" edge.
func HasFileSets() predicate.OutputGroup {
	return predicate.OutputGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, false, FileSetsTable, FileSetsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFileSetsWith applies the HasEdge predicate on the "file_sets" edge with a given conditions (other predicates).
func HasFileSetsWith(preds ...predicate.NamedSetOfFiles) predicate.OutputGroup {
	return predicate.OutputGroup(func(s *sql.Selector) {
		step := newFileSetsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.OutputGroup) predicate.OutputGroup {
	return predicate.OutputGroup(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.OutputGroup) predicate.OutputGroup {
	return predicate.OutputGroup(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.OutputGroup) predicate.OutputGroup {
	return predicate.OutputGroup(sql.NotPredicates(p))
}
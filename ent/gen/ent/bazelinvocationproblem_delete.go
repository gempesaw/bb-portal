// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/bazelinvocationproblem"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// BazelInvocationProblemDelete is the builder for deleting a BazelInvocationProblem entity.
type BazelInvocationProblemDelete struct {
	config
	hooks    []Hook
	mutation *BazelInvocationProblemMutation
}

// Where appends a list predicates to the BazelInvocationProblemDelete builder.
func (bipd *BazelInvocationProblemDelete) Where(ps ...predicate.BazelInvocationProblem) *BazelInvocationProblemDelete {
	bipd.mutation.Where(ps...)
	return bipd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (bipd *BazelInvocationProblemDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, bipd.sqlExec, bipd.mutation, bipd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (bipd *BazelInvocationProblemDelete) ExecX(ctx context.Context) int {
	n, err := bipd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (bipd *BazelInvocationProblemDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(bazelinvocationproblem.Table, sqlgraph.NewFieldSpec(bazelinvocationproblem.FieldID, field.TypeInt))
	if ps := bipd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, bipd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	bipd.mutation.done = true
	return affected, err
}

// BazelInvocationProblemDeleteOne is the builder for deleting a single BazelInvocationProblem entity.
type BazelInvocationProblemDeleteOne struct {
	bipd *BazelInvocationProblemDelete
}

// Where appends a list predicates to the BazelInvocationProblemDelete builder.
func (bipdo *BazelInvocationProblemDeleteOne) Where(ps ...predicate.BazelInvocationProblem) *BazelInvocationProblemDeleteOne {
	bipdo.bipd.mutation.Where(ps...)
	return bipdo
}

// Exec executes the deletion query.
func (bipdo *BazelInvocationProblemDeleteOne) Exec(ctx context.Context) error {
	n, err := bipdo.bipd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{bazelinvocationproblem.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (bipdo *BazelInvocationProblemDeleteOne) ExecX(ctx context.Context) {
	if err := bipdo.Exec(ctx); err != nil {
		panic(err)
	}
}
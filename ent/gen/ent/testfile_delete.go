// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
	"github.com/buildbarn/bb-portal/ent/gen/ent/testfile"
)

// TestFileDelete is the builder for deleting a TestFile entity.
type TestFileDelete struct {
	config
	hooks    []Hook
	mutation *TestFileMutation
}

// Where appends a list predicates to the TestFileDelete builder.
func (tfd *TestFileDelete) Where(ps ...predicate.TestFile) *TestFileDelete {
	tfd.mutation.Where(ps...)
	return tfd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (tfd *TestFileDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, tfd.sqlExec, tfd.mutation, tfd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (tfd *TestFileDelete) ExecX(ctx context.Context) int {
	n, err := tfd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (tfd *TestFileDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(testfile.Table, sqlgraph.NewFieldSpec(testfile.FieldID, field.TypeInt))
	if ps := tfd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, tfd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	tfd.mutation.done = true
	return affected, err
}

// TestFileDeleteOne is the builder for deleting a single TestFile entity.
type TestFileDeleteOne struct {
	tfd *TestFileDelete
}

// Where appends a list predicates to the TestFileDelete builder.
func (tfdo *TestFileDeleteOne) Where(ps ...predicate.TestFile) *TestFileDeleteOne {
	tfdo.tfd.mutation.Where(ps...)
	return tfdo
}

// Exec executes the deletion query.
func (tfdo *TestFileDeleteOne) Exec(ctx context.Context) error {
	n, err := tfdo.tfd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{testfile.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (tfdo *TestFileDeleteOne) ExecX(ctx context.Context) {
	if err := tfdo.Exec(ctx); err != nil {
		panic(err)
	}
}
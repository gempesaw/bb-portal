// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/buildbarn/bb-portal/ent/gen/ent/memorymetrics"
	"github.com/buildbarn/bb-portal/ent/gen/ent/predicate"
)

// MemoryMetricsDelete is the builder for deleting a MemoryMetrics entity.
type MemoryMetricsDelete struct {
	config
	hooks    []Hook
	mutation *MemoryMetricsMutation
}

// Where appends a list predicates to the MemoryMetricsDelete builder.
func (mmd *MemoryMetricsDelete) Where(ps ...predicate.MemoryMetrics) *MemoryMetricsDelete {
	mmd.mutation.Where(ps...)
	return mmd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (mmd *MemoryMetricsDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, mmd.sqlExec, mmd.mutation, mmd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (mmd *MemoryMetricsDelete) ExecX(ctx context.Context) int {
	n, err := mmd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (mmd *MemoryMetricsDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(memorymetrics.Table, sqlgraph.NewFieldSpec(memorymetrics.FieldID, field.TypeInt))
	if ps := mmd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, mmd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	mmd.mutation.done = true
	return affected, err
}

// MemoryMetricsDeleteOne is the builder for deleting a single MemoryMetrics entity.
type MemoryMetricsDeleteOne struct {
	mmd *MemoryMetricsDelete
}

// Where appends a list predicates to the MemoryMetricsDelete builder.
func (mmdo *MemoryMetricsDeleteOne) Where(ps ...predicate.MemoryMetrics) *MemoryMetricsDeleteOne {
	mmdo.mmd.mutation.Where(ps...)
	return mmdo
}

// Exec executes the deletion query.
func (mmdo *MemoryMetricsDeleteOne) Exec(ctx context.Context) error {
	n, err := mmdo.mmd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{memorymetrics.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (mmdo *MemoryMetricsDeleteOne) ExecX(ctx context.Context) {
	if err := mmdo.Exec(ctx); err != nil {
		panic(err)
	}
}
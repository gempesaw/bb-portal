// Code generated by ent, DO NOT EDIT.

package build

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the build type in the database.
	Label = "build"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldBuildURL holds the string denoting the build_url field in the database.
	FieldBuildURL = "build_url"
	// FieldBuildUUID holds the string denoting the build_uuid field in the database.
	FieldBuildUUID = "build_uuid"
	// FieldEnv holds the string denoting the env field in the database.
	FieldEnv = "env"
	// EdgeInvocations holds the string denoting the invocations edge name in mutations.
	EdgeInvocations = "invocations"
	// Table holds the table name of the build in the database.
	Table = "builds"
	// InvocationsTable is the table that holds the invocations relation/edge.
	InvocationsTable = "bazel_invocations"
	// InvocationsInverseTable is the table name for the BazelInvocation entity.
	// It exists in this package in order to avoid circular dependency with the "bazelinvocation" package.
	InvocationsInverseTable = "bazel_invocations"
	// InvocationsColumn is the table column denoting the invocations relation/edge.
	InvocationsColumn = "build_invocations"
)

// Columns holds all SQL columns for build fields.
var Columns = []string{
	FieldID,
	FieldBuildURL,
	FieldBuildUUID,
	FieldEnv,
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

// OrderOption defines the ordering options for the Build queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByBuildURL orders the results by the build_url field.
func ByBuildURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBuildURL, opts...).ToFunc()
}

// ByBuildUUID orders the results by the build_uuid field.
func ByBuildUUID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldBuildUUID, opts...).ToFunc()
}

// ByInvocationsCount orders the results by invocations count.
func ByInvocationsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newInvocationsStep(), opts...)
	}
}

// ByInvocations orders the results by invocations terms.
func ByInvocations(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInvocationsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newInvocationsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InvocationsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, InvocationsTable, InvocationsColumn),
	)
}
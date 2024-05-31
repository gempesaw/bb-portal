// Code generated by ent, DO NOT EDIT.

package blob

import (
	"fmt"
	"io"
	"strconv"

	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the blob type in the database.
	Label = "blob"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldURI holds the string denoting the uri field in the database.
	FieldURI = "uri"
	// FieldSizeBytes holds the string denoting the size_bytes field in the database.
	FieldSizeBytes = "size_bytes"
	// FieldArchivingStatus holds the string denoting the archiving_status field in the database.
	FieldArchivingStatus = "archiving_status"
	// FieldReason holds the string denoting the reason field in the database.
	FieldReason = "reason"
	// FieldArchiveURL holds the string denoting the archive_url field in the database.
	FieldArchiveURL = "archive_url"
	// Table holds the table name of the blob in the database.
	Table = "blobs"
)

// Columns holds all SQL columns for blob fields.
var Columns = []string{
	FieldID,
	FieldURI,
	FieldSizeBytes,
	FieldArchivingStatus,
	FieldReason,
	FieldArchiveURL,
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

// ArchivingStatus defines the type for the "archiving_status" enum field.
type ArchivingStatus string

// ArchivingStatusQUEUED is the default value of the ArchivingStatus enum.
const DefaultArchivingStatus = ArchivingStatusQUEUED

// ArchivingStatus values.
const (
	ArchivingStatusQUEUED    ArchivingStatus = "QUEUED"
	ArchivingStatusARCHIVING ArchivingStatus = "ARCHIVING"
	ArchivingStatusSUCCESS   ArchivingStatus = "SUCCESS"
	ArchivingStatusFAILED    ArchivingStatus = "FAILED"
)

func (as ArchivingStatus) String() string {
	return string(as)
}

// ArchivingStatusValidator is a validator for the "archiving_status" field enum values. It is called by the builders before save.
func ArchivingStatusValidator(as ArchivingStatus) error {
	switch as {
	case ArchivingStatusQUEUED, ArchivingStatusARCHIVING, ArchivingStatusSUCCESS, ArchivingStatusFAILED:
		return nil
	default:
		return fmt.Errorf("blob: invalid enum value for archiving_status field: %q", as)
	}
}

// OrderOption defines the ordering options for the Blob queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByURI orders the results by the uri field.
func ByURI(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldURI, opts...).ToFunc()
}

// BySizeBytes orders the results by the size_bytes field.
func BySizeBytes(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSizeBytes, opts...).ToFunc()
}

// ByArchivingStatus orders the results by the archiving_status field.
func ByArchivingStatus(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldArchivingStatus, opts...).ToFunc()
}

// ByReason orders the results by the reason field.
func ByReason(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldReason, opts...).ToFunc()
}

// ByArchiveURL orders the results by the archive_url field.
func ByArchiveURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldArchiveURL, opts...).ToFunc()
}

// MarshalGQL implements graphql.Marshaler interface.
func (e ArchivingStatus) MarshalGQL(w io.Writer) {
	io.WriteString(w, strconv.Quote(e.String()))
}

// UnmarshalGQL implements graphql.Unmarshaler interface.
func (e *ArchivingStatus) UnmarshalGQL(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("enum %T must be a string", val)
	}
	*e = ArchivingStatus(str)
	if err := ArchivingStatusValidator(*e); err != nil {
		return fmt.Errorf("%s is not a valid ArchivingStatus", str)
	}
	return nil
}
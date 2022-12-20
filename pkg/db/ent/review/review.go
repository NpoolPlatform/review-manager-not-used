// Code generated by ent, DO NOT EDIT.

package review

import (
	"entgo.io/ent"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the review type in the database.
	Label = "review"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldDeletedAt holds the string denoting the deleted_at field in the database.
	FieldDeletedAt = "deleted_at"
	// FieldAppID holds the string denoting the app_id field in the database.
	FieldAppID = "app_id"
	// FieldReviewerID holds the string denoting the reviewer_id field in the database.
	FieldReviewerID = "reviewer_id"
	// FieldDomain holds the string denoting the domain field in the database.
	FieldDomain = "domain"
	// FieldObjectID holds the string denoting the object_id field in the database.
	FieldObjectID = "object_id"
	// FieldTrigger holds the string denoting the trigger field in the database.
	FieldTrigger = "trigger"
	// FieldObjectType holds the string denoting the object_type field in the database.
	FieldObjectType = "object_type"
	// FieldState holds the string denoting the state field in the database.
	FieldState = "state"
	// FieldMessage holds the string denoting the message field in the database.
	FieldMessage = "message"
	// Table holds the table name of the review in the database.
	Table = "reviews"
)

// Columns holds all SQL columns for review fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldDeletedAt,
	FieldAppID,
	FieldReviewerID,
	FieldDomain,
	FieldObjectID,
	FieldTrigger,
	FieldObjectType,
	FieldState,
	FieldMessage,
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

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/NpoolPlatform/review-manager/pkg/db/ent/runtime"
//
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() uint32
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() uint32
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() uint32
	// DefaultDeletedAt holds the default value on creation for the "deleted_at" field.
	DefaultDeletedAt func() uint32
	// DefaultAppID holds the default value on creation for the "app_id" field.
	DefaultAppID func() uuid.UUID
	// DefaultReviewerID holds the default value on creation for the "reviewer_id" field.
	DefaultReviewerID func() uuid.UUID
	// DefaultDomain holds the default value on creation for the "domain" field.
	DefaultDomain string
	// DefaultObjectID holds the default value on creation for the "object_id" field.
	DefaultObjectID func() uuid.UUID
	// DefaultTrigger holds the default value on creation for the "trigger" field.
	DefaultTrigger string
	// DefaultObjectType holds the default value on creation for the "object_type" field.
	DefaultObjectType string
	// DefaultState holds the default value on creation for the "state" field.
	DefaultState string
	// DefaultMessage holds the default value on creation for the "message" field.
	DefaultMessage string
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

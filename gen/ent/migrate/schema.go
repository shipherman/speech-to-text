// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// AudiosColumns holds the columns for the "audios" table.
	AudiosColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "path", Type: field.TypeString, Unique: true},
		{Name: "hash", Type: field.TypeString, Unique: true},
		{Name: "text", Type: field.TypeString},
		{Name: "timestamp", Type: field.TypeTime, SchemaType: map[string]string{"postgres": "timestamptz"}},
		{Name: "user_audio", Type: field.TypeInt, Nullable: true},
	}
	// AudiosTable holds the schema information for the "audios" table.
	AudiosTable = &schema.Table{
		Name:       "audios",
		Columns:    AudiosColumns,
		PrimaryKey: []*schema.Column{AudiosColumns[0]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "audios_users_audio",
				Columns:    []*schema.Column{AudiosColumns[5]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.SetNull,
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "email", Type: field.TypeString, Unique: true},
		{Name: "login", Type: field.TypeString, Unique: true},
		{Name: "password", Type: field.TypeString},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		AudiosTable,
		UsersTable,
	}
)

func init() {
	AudiosTable.ForeignKeys[0].RefTable = UsersTable
}

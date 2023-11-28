package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Audio holds the schema definition for the Audio entity.
type Audio struct {
	ent.Schema
}

// Fields of the Audio.
func (Audio) Fields() []ent.Field {
	return []ent.Field{
		field.String("path").
			NotEmpty().Unique(),
		field.String("hash").
			NotEmpty().Unique(),
		field.String("text"),
	}
}

// Edges of the Audio.
func (Audio) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("audio").
			Unique(),
	}
}

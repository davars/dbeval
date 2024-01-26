package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Author holds the schema definition for the Author entity.
type Author struct {
	ent.Schema
}

// Fields of the Author.
func (Author) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.String("name").
			Default(""),
	}
}

// Edges of the Author.
func (Author) Edges() []ent.Edge {
	return nil
}

func (Author) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name"),
	}
}

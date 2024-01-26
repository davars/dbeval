package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"time"
)

// Article holds the schema definition for the Article entity.
type Article struct {
	ent.Schema
}

// Fields of the Article.
func (Article) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Unique(),
		field.String("title").
			Default(""),
		field.String("body").
			Default(""),
		field.Time("published_at").
			Default(time.Now),
	}
}

// Edges of the Article.
func (Article) Edges() []ent.Edge {
	return nil
}

func (Article) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("published_at"),
	}
}

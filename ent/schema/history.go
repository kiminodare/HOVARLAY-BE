package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// History holds the schema definition for the History entity.
type History struct {
	ent.Schema
}

// Fields of the History.
func (History) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(
			func() uuid.UUID {
				id, err := uuid.NewV7()
				if err != nil {
					panic(err)
				}
				return id
			},
		).Immutable().Unique(),
		field.String("text").NotEmpty().SchemaType(map[string]string{dialect.Postgres: "text"}),
		field.String("voice").NotEmpty(),
		field.Float("rate").Default(1).Min(0.1).Max(5),
		field.Float("pitch").Default(1).Min(0).Max(2),
		field.Float("volume").Default(1).Min(0).Max(1),
		field.Time("created_at").Default(func() time.Time { return time.Now() }).StructTag(`json:"createdAt"`),
		field.Time("updated_at").Default(func() time.Time { return time.Now() }).UpdateDefault(func() time.Time { return time.Now() }).StructTag(`json:"updatedAt"`),
	}
}

// Edges of the History.
func (History) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).Ref("histories").Unique(),
	}
}

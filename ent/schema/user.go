package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
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
		field.String("name").NotEmpty(),
		field.String("email").NotEmpty().Unique(),
		field.String("password").NotEmpty(),
		field.Time("created_at").Default(func() time.Time { return time.Now() }).StructTag(`json:"createdAt"`),
		field.Time("updated_at").Default(func() time.Time { return time.Now() }).UpdateDefault(func() time.Time { return time.Now() }).StructTag(`json:"updatedAt"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("histories", History.Type),
	}
}

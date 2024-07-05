package schema

import (
	"context"
	"strings"
	"time"

	ge "github.com/nickheyer/DiscoFlixGo/ent"
	"github.com/nickheyer/DiscoFlixGo/ent/hook"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("email").
			NotEmpty().
			Unique(),
		field.String("password").
			Sensitive().
			NotEmpty(),
		field.Bool("verified").
			Default(false),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Bool("is_admin").
			Default(false),
		field.Strings("roles").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", PasswordToken.Type).
			Ref("user"),
	}
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(
			// Mutate incoming email addresses to be lowercase.
			func(next ent.Mutator) ent.Mutator {
				return hook.UserFunc(func(ctx context.Context, m *ge.UserMutation) (ent.Value, error) {
					if v, exists := m.Email(); exists {
						m.SetEmail(strings.ToLower(v))
					}
					return next.Mutate(ctx, m)
				})
			},
			// Hook on User Create or Update.
			ent.OpCreate|
				ent.OpUpdate|
				ent.OpUpdateOne,
		),
	}
}

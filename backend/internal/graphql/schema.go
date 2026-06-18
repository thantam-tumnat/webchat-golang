package graphql

import (
	"github.com/graphql-go/graphql"
)

// NewSchema ประกอบ Query + Mutation เข้าด้วยกันเป็น schema เดียว
// ผูกแต่ละ field เข้ากับ resolver method ที่เรียก usecase เดิม
func NewSchema(r *Resolver) (graphql.Schema, error) {
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			// rooms: คืนห้องทั้งหมด
			"rooms": &graphql.Field{
				Type:    graphql.NewList(roomType),
				Resolve: r.rooms,
			},
			// messages: ข้อความในห้อง (รับ roomId + page/limit แบบ optional)
			"messages": &graphql.Field{
				Type: graphql.NewList(messageType),
				Args: graphql.FieldConfigArgument{
					"roomId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"page":   &graphql.ArgumentConfig{Type: graphql.Int},
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: r.messages,
			},
		},
	})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createUser": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.createUser,
			},
			"createRoom": &graphql.Field{
				Type: roomType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.createRoom,
			},
			"sendMessage": &graphql.Field{
				Type: messageType,
				Args: graphql.FieldConfigArgument{
					"roomId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"userId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
					"content": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.sendMessage,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
}

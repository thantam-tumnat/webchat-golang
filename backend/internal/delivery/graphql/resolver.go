package graphql

import (
	"chatapp/internal/usecase"

	"github.com/graphql-go/graphql"
)

// Resolver ถือ usecase ทั้งหมด แล้วให้แต่ละ field เรียกใช้
// สังเกตว่า GraphQL เป็นแค่ delivery layer อีกตัว — เรียก usecase ชุดเดียวกับที่ REST ใช้
// business logic จึงไม่ถูก copy ซ้ำ (จุดแข็งของ Clean Architecture)
type Resolver struct {
	userUC    *usecase.UserUsecase
	roomUC    *usecase.RoomUsecase
	messageUC *usecase.MessageUsecase
}

func NewResolver(
	userUC *usecase.UserUsecase,
	roomUC *usecase.RoomUsecase,
	messageUC *usecase.MessageUsecase,
) *Resolver {
	return &Resolver{userUC: userUC, roomUC: roomUC, messageUC: messageUC}
}

// --- Query resolvers ---

func (r *Resolver) rooms(p graphql.ResolveParams) (interface{}, error) {
	return r.roomUC.List(p.Context)
}

func (r *Resolver) messages(p graphql.ResolveParams) (interface{}, error) {
	roomID := uint(p.Args["roomId"].(int))
	page := optInt(p.Args, "page", 1)
	limit := optInt(p.Args, "limit", 50)

	// usecase.List คืน (messages, total, error) — GraphQL ใช้แค่ messages
	msgs, _, err := r.messageUC.List(p.Context, roomID, page, limit)
	return msgs, err
}

// --- Mutation resolvers ---

func (r *Resolver) createUser(p graphql.ResolveParams) (interface{}, error) {
	username := p.Args["username"].(string)
	return r.userUC.CreateOrGet(p.Context, username)
}

func (r *Resolver) createRoom(p graphql.ResolveParams) (interface{}, error) {
	name := p.Args["name"].(string)
	return r.roomUC.Create(p.Context, name)
}

func (r *Resolver) sendMessage(p graphql.ResolveParams) (interface{}, error) {
	roomID := uint(p.Args["roomId"].(int))
	userID := uint(p.Args["userId"].(int))
	content := p.Args["content"].(string)
	return r.messageUC.Send(p.Context, roomID, userID, content)
}

// optInt อ่าน argument ที่เป็น optional ถ้าไม่ส่งมาให้ใช้ค่า default
func optInt(args map[string]interface{}, key string, def int) int {
	if v, ok := args[key]; ok {
		if n, ok := v.(int); ok {
			return n
		}
	}
	return def
}

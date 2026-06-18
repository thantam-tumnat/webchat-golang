package graphql

import (
	"github.com/graphql-go/graphql"
)

// ไฟล์นี้นิยาม "หน้าตา" ของข้อมูลในกราฟ (GraphQL Object types)
// ทั้งหมด map ตรงกับ struct ใน domain โดยอาศัย default resolver ของ graphql-go
// ที่จับคู่ field แบบ case-insensitive (เช่น "roomId" -> struct field "RoomID")
// จึงไม่ต้องเขียน resolver รายตัว ใช้ struct domain เดิมซ้ำได้เลย

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"username":  &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.DateTime},
	},
})

var roomType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Room",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"name":      &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.DateTime},
	},
})

var messageType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Message",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.Int},
		"roomId":    &graphql.Field{Type: graphql.Int},
		"userId":    &graphql.Field{Type: graphql.Int},
		"content":   &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.DateTime},

		// ★ จุดที่ GraphQL ชนะ REST: ขอ message พร้อมข้อมูล user (nested)
		//   ได้ใน request เดียว — ค่า User ถูก preload มากับ message อยู่แล้ว
		"user": &graphql.Field{Type: userType},
	},
})

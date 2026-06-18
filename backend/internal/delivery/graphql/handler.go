package graphql

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// NewHandler สร้าง HTTP handler ของ GraphQL แล้วครอบด้วย adaptor ให้ใช้กับ Fiber ได้
//   - GraphiQL: true  -> เปิด UI ทดสอบ query ในเบราว์เซอร์ที่ GET /graphql
//   - Pretty:   true  -> จัดรูป JSON response ให้อ่านง่าย
func NewHandler(schema graphql.Schema) fiber.Handler {
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	// adaptor.HTTPHandler แปลง net/http handler -> fiber.Handler
	return adaptor.HTTPHandler(h)
}

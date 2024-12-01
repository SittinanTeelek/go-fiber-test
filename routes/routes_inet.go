package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	c "github.com/teelek/go-test/controllers"
)

func InetRoutes(app *fiber.App) {

	api := app.Group("/api") // /api
	v1 := api.Group("/v1")   // /api/v1

	v1Protected := v1.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"gofiber": "21022566",
			"testgo":  "23012023",
		},
	}))

	pf := v1.Group("/pf")

	//สร้างตารางโปรไฟล์ผู้ใช้ผ่านการ automigrate
	pf.Get("/", c.GetProfile)
	pf.Get("/gen", c.GetGen)
	pf.Post("/", c.AddProfile)
	pf.Post("/edit", c.UpdateProfile)
	pf.Post("/remove", c.RemoveProfile)

	//สร้างapi รับค่าตัวเลข ผ่านpath แล้วreturnเป็นค่าfactorialของตัวเลขนั้น
	v1Protected.Get("/fact/:num", c.Factorial)

	//api method POST สมัครสมาชิก ดักฟิลข้อมูลให้ถูกต้อง localhost:3000/api/v1/register และถ้าใส่ข้อมูลไม่ถูกต้องให้โชว์ใส่ข้อมูลผิดพลาด
	v1Protected.Post("/register", c.Register)

	v1Protected.Get("/", c.HelloTest)

	v1Protected.Post("/", c.BodyParserTest)

	v1Protected.Post("/fact", c.QueryTest)

	v1Protected.Get("/user/:name", c.ParramsTest)

	v1Protected.Post("/valid", c.ValidTest)

	v2 := api.Group("/v2")

	v2.Get("/", c.HelloTestV2)

	//สร้างapiขึ้นต้นด้วย api/v3/ (<--ใช้วิธีแบบจัดgroup api)ตามด้วยชื่อเล่นตัวเอง  โดยapiนี้มีการรับ QueryParam ที่ชื่อkeyว่า tax_id นำค่าที่keyเข้าไป(keyได้ทั้งตัวเลขตัวอักษร)แปลงเป็น ascii
	v3 := api.Group("/v3")

	v3.Post("teelek", c.QueryEx)

	//CRUD dogs
	dog := v1Protected.Group("/dog")

	dog.Get("", c.GetDogs)
	dog.Get("/filter50_100", c.GetDog50_100)
	dog.Get("/color", c.GetDogsColor)
	dog.Get("/deleted", c.GetDogsDeleted)
	dog.Get("/filter", c.GetDog)
	dog.Get("/json", c.GetDogsJson)
	dog.Post("/", c.AddDog)
	dog.Put("/:id", c.UpdateDog)
	dog.Delete("/:id", c.RemoveDog)

	//CRUD Companys
	company := v1Protected.Group("/company")

	company.Get("", c.GetCompanys)
	company.Get("/filter", c.GetCompany)
	company.Get("/json", c.GetCompanysJson)
	company.Post("/", c.AddCompany)
	company.Put("/:id", c.UpdateCompany)
	company.Delete("/:id", c.RemoveCompany)

}

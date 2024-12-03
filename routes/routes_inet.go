package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	c "github.com/teelek/go-test/controllers"
)

func InetRoutes(app *fiber.App) {

	auth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			"gofiber": "21022566",
		},
	})

	pfAuth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			"testgo": "23012023",
		},
	})

	api := app.Group("/api") // /api
	v1 := api.Group("/v1")   // /api/v1

	//สร้างapi รับค่าตัวเลข ผ่านpath แล้วreturnเป็นค่าfactorialของตัวเลขนั้น
	fact := v1.Group("/fact", auth)

	fact.Get("/:num", c.Factorial)

	//สร้างapiขึ้นต้นด้วย api/v3/ (<--ใช้วิธีแบบจัดgroup api)ตามด้วยชื่อเล่นตัวเอง  โดยapiนี้มีการรับ QueryParam ที่ชื่อkeyว่า tax_id นำค่าที่keyเข้าไป(keyได้ทั้งตัวเลขตัวอักษร)แปลงเป็น ascii
	v3 := api.Group("/v3", auth)

	v3.Post("/teelek", c.QueryEx)

	//company
	cpn := v1.Group("/company", auth)

	cpn.Get("/", c.GetCompanies)
	cpn.Get("/filter", c.GetCompany)
	cpn.Post("/", c.AddCompany)
	cpn.Put("/:id", c.UpdateCompany)
	cpn.Delete("/:id", c.RemoveCompany)

	//api method POST สมัครสมาชิก ดักฟิลข้อมูลให้ถูกต้อง localhost:3000/api/v1/register และถ้าใส่ข้อมูลไม่ถูกต้องให้โชว์ใส่ข้อมูลผิดพลาด
	v1.Post("/register", c.Register)

	//CRUD dogs
	dog := v1.Group("/dog")

	dog.Get("", c.GetDogs)
	dog.Get("/filter50_100", c.GetDog50_100)
	dog.Get("/color", c.GetDogsColor)
	dog.Get("/deleted", c.GetDogsDeleted)
	dog.Get("/filter", c.GetDog)
	dog.Get("/json", c.GetDogsJson)
	dog.Post("/", c.AddDog)
	dog.Put("/:id", c.UpdateDog)
	dog.Delete("/:id", c.RemoveDog)

	//สร้างตารางโปรไฟล์ผู้ใช้ผ่านการ automigrate
	noAuthProfile := v1.Group("/profile")

	noAuthProfile.Get("/", c.GetProfiles)
	noAuthProfile.Get("/gen", c.GetGen)
	noAuthProfile.Get("/filter", c.GetProfileFilter)
	authProfile := v1.Group("/profile", pfAuth)

	authProfile.Post("/", c.AddProfile)
	authProfile.Put("/:id", c.UpdateProfile)
	authProfile.Delete("/:id", c.RemoveProfile)

}

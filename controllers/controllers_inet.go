package controllers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/teelek/go-test/database"
	m "github.com/teelek/go-test/models"
)

// สร้างapi รับค่าตัวเลข ผ่านpath แล้วreturnเป็นค่าfactorialของตัวเลขนั้น
func Factorial(c *fiber.Ctx) error {
	num, err := strconv.Atoi(c.Params("num"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	fact := 1
	for i := 1; i <= num; i++ {
		fact = fact * i
	}
	return c.JSON(fiber.Map{
		"num":       num,
		"factorail": fact,
	})
}

// สร้างapiขึ้นต้นด้วย api/v3/ (<--ใช้วิธีแบบจัดgroup api)ตามด้วยชื่อเล่นตัวเอง  โดยapiนี้มีการรับ QueryParam ที่ชื่อkeyว่า tax_id นำค่าที่keyเข้าไป(keyได้ทั้งตัวเลขตัวอักษร)แปลงเป็น ascii
func QueryEx(c *fiber.Ctx) error {
	tax_id := c.Query("tax_id")

	var ascii []int
	for _, ch := range tax_id {
		ascii = append(ascii, int(ch))
	}

	result := fmt.Sprintf("ex→ tax_id = %s → %d", tax_id, ascii)

	return c.JSON(result)
}

// api method POST สมัครสมาชิก ดักฟิลข้อมูลให้ถูกต้อง
func Register(c *fiber.Ctx) error {
	newUser := new(m.NewUser)

	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if match, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, newUser.UserName); !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ชื่อผู้ใช้ต้องเป็นตัวอักษรภาษาอังกฤษ (a-z), (A-Z), ตัวเลข (0-9), และเครื่องหมาย (_), (-) เท่านั้น",
		})
	}

	if len(newUser.Password) < 6 || len(newUser.Password) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "รหัสผ่านต้องมีความยาว 6-20 ตัวอักษร",
		})
	}

	if match, _ := regexp.MatchString(`^[a-z0-9-]{2,30}$`, newUser.WebSite); !match {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ชื่อเว็บไซต์ต้องมีความยาว 2-30 ตัวอักษร ต้องเป็นตัวอักษรภาษาอังกฤษตัวเล็ก (a-z), ตัวเลข(0-9) ห้ามใช้เครื่องหมายอักขระพิเศษยกเว้นเครื่องหมายขีด (-) ห้ามเว้นวรรคและห้ามใช้ภาษาไทย",
		})
	}

	validate := validator.New()
	errors := validate.Struct(newUser)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	// ถ้าผ่านการตรวจสอบทั้งหมด
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Registration successful",
		"member":  newUser,
	})
	//return c.JSON(member)

}

// สร้างตารางcompany โดยใช้AutoMigrate โดยที่โครงสร้างcompanyควรจะมีอะไรบ้างใส่มาตามความเหมาะสม และGroupเพิ่มCRUD
func GetCompanies(c *fiber.Ctx) error {
	db := database.DBConn
	var company []m.Companies

	if err := db.Find(&company).Error; err != nil {
		return c.SendStatus(500)
	}
	return c.Status(200).JSON(company)
}

func GetCompany(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var company []m.Companies

	result := db.Find(&company, "company_id = ?", search)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(company)
}

func AddCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Companies

	if err := c.BodyParser(&company); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&company)
	return c.Status(201).JSON(company)
}

func UpdateCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Companies
	id := c.Params("id")

	if err := c.BodyParser(&company); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&company)
	return c.Status(200).JSON(company)
}

func RemoveCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Companies
	id := c.Params("id")

	result := db.Delete(&company, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

// สร้าง api GET ใน group dogs โชว์ข้อมูลที่ถูกลบไปแล้ว ตารางdogs
func GetDogs(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dog

	if err := db.Find(&dogs).Error; err != nil {
		return c.SendStatus(500)
	}
	return c.JSON(dogs)
}

func GetDog(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var dog []m.Dog

	result := db.Find(&dog, "dog_id = ?", search)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(dog)
}

// สร้าง api GET ใน group dogs โชว์ข้อมูลที่ถูกลบไปแล้ว ตารางdogs
func GetDogsDeleted(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dog

	db.Unscoped().Where("deleted_at").Find(&dogs)
	return c.Status(200).JSON(dogs)
}

// สร้างapi GETใหม่ แสดงข้อมูลตารางdogโดย where หา dog_id > 50 แต่น้อยกว่า 100  (gorm)
func GetDog50_100(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dog

	result := db.Where("dog_id > ? AND dog_id < ?", 50, 100).Find(&dogs)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(dogs)
}

func AddDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dog

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&dog)
	return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dog
	id := c.Params("id")

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	result := db.Where("id = ?", id).Updates(&dog)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var dog m.Dog

	result := db.Delete(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func DeletedDog(c *fiber.Ctx) error {
	db := database.DBConn

	var deletedDog []m.Dog

	result := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&deletedDog)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(fiber.StatusOK).JSON(deletedDog)
}

func GetDogsJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dog

	db.Find(&dogs) //10ตัว

	var dataResults []m.DogsRes
	for _, v := range dogs { //1 inet 112 //2 inet1 113
		typeStr := ""
		if v.DogID == 111 {
			typeStr = "red"
		} else if v.DogID == 113 {
			typeStr = "green"
		} else if v.DogID == 999 {
			typeStr = "pink"
		} else {
			typeStr = "no color"
		}

		d := m.DogsRes{
			Name:  v.Name,  //inet
			DogID: v.DogID, //112
			Type:  typeStr, //no color
		}
		dataResults = append(dataResults, d)
		// sumAmount += v.Amount
	}

	r := m.ResultDogsData{
		Data:  dataResults,
		Name:  "golang-test",
		Count: len(dogs), //หาผลรวม,
	}
	return c.Status(200).JSON(r)
}

func GetDogsColor(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dog

	result := db.Find(&dogs)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	colorCounts := map[string]int{
		"red":      0,
		"green":    0,
		"pink":     0,
		"no color": 0,
	}
	var dataResults []m.DogsRes
	for _, dog := range dogs {
		var color string
		switch {
		case dog.DogID >= 10 && dog.DogID <= 50:
			color = "red"
		case dog.DogID >= 100 && dog.DogID <= 150:
			color = "green"
		case dog.DogID >= 200 && dog.DogID <= 250:
			color = "pink"
		default:
			color = "no color"
		}
		colorCounts[color]++
		d := m.DogsRes{
			Name:  dog.Name,
			DogID: dog.DogID,
			Type:  color,
		}
		dataResults = append(dataResults, d)
	}
	r := m.ResultColorData{
		Data:       dataResults,
		Name:       "golang-test",
		Count:      len(dogs), //หาผลรวม,
		SumRed:     colorCounts["red"],
		SumGreen:   colorCounts["green"],
		SumPink:    colorCounts["pink"],
		SumNoColor: colorCounts["no color"],
	}
	return c.Status(200).JSON(r)
}

// สร้างตารางโปรไฟล์ผู้ใช้ผ่านการ automigrate

func GetProfiles(c *fiber.Ctx) error {
	db := database.DBConn
	var profiles []m.Profiles

	result := db.Find(&profiles)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(200).JSON(profiles)
}

func GetProfile(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var profile []m.Profiles

	if search == "" {
		return c.SendStatus(400)
	}

	result := db.Find(&profile, "employee_id = ?", search)

	if result.RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(200).JSON(profile)

}

func AddProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profiles

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	db.Create(&profile)
	return c.Status(201).JSON(profile)
}

func UpdateProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profiles
	id := c.Params("id")

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&profile)
	return c.Status(200).JSON(profile)
}

func RemoveProfile(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var profile m.Profiles

	result := db.Delete(&profile, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func GetGen(c *fiber.Ctx) error {
	db := database.DBConn
	var profiles []m.Profiles

	if err := db.Find(&profiles).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch profile data"})
	}

	genDivider := map[string]int{
		"GenZ":            0,
		"GenY":            0,
		"GenX":            0,
		"Baby Boomer":     0,
		"G.I. Generation": 0,
	}
	var dataResults []m.Profiles
	for _, profiles := range profiles {
		var gen string
		switch {
		case profiles.Age < 24:
			gen = "GenZ"
		case profiles.Age <= 41:
			gen = "GenY"
		case profiles.Age <= 56:
			gen = "GenX"
		case profiles.Age <= 75:
			gen = "Baby Boomer"
		default:
			gen = "G.I. Generation"
		}
		genDivider[gen]++
	}
	result := m.ResultProfileData{
		Data:         dataResults,
		Count:        len(profiles),
		GenZ:         genDivider["GenZ"],
		GenY:         genDivider["GenY"],
		GenX:         genDivider["GenX"],
		BabyBoomer:   genDivider["Baby Boomer"],
		GIGeneration: genDivider["G.I. Generation"],
	}
	return c.Status(200).JSON(result)
}

func GetProfileFilter(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var profiles []m.Profiles

	result := db.Where("employee_id = ? OR name = ? OR last_name = ?", search, search, search).Find(&profiles)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(profiles)
}

func GetProfileByKey(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var profiles []m.Profiles

	result := db.Where("employee_id = ? OR name = ? OR last_name = ?", search, search, search).Find(&profiles)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(profiles)
}

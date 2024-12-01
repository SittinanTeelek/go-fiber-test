package controllers

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/teelek/go-test/database"
	m "github.com/teelek/go-test/models"
)

func HelloTest(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func HelloTestV2(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func BodyParserTest(c *fiber.Ctx) error {
	p := new(m.Person)

	if err := c.BodyParser(p); err != nil {
		return err
	}

	log.Println(p.Name) // john
	log.Println(p.Pass) // doe
	str := p.Name + p.Pass
	return c.JSON(str)
}

func ParramsTest(c *fiber.Ctx) error {

	str := "hello ==> " + c.Params("name")
	return c.JSON(str)
}

func QueryTest(c *fiber.Ctx) error {
	a := c.Query("search")
	str := "my search is  " + a
	return c.JSON(str)
}

func ValidTest(c *fiber.Ctx) error {

	user := new(m.User)
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	errors := validate.Struct(user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	return c.JSON(user)
}

func GetDogs(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs)
	return c.Status(200).JSON(dogs)
}

// สร้าง api GET ใน group dogs โชว์ข้อมูลที่ถูกลบไปแล้ว ตารางdogs
func GetDogsDeleted(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Unscoped().Where("deleted_at").Find(&dogs)
	return c.Status(200).JSON(dogs)
}

// สร้างapi GETใหม่ แสดงข้อมูลตารางdogโดย where หา dog_id > 50 แต่น้อยกว่า 100  (gorm)
func GetDog50_100(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	result := db.Where("dog_id > ? AND dog_id < ?", 50, 100).Find(&dogs)
	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database query failed",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "No dogs found in the specified range",
		})
	}

	return c.Status(200).JSON(&dogs)
}

func GetDog(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var dog []m.Dogs

	result := db.Find(&dog, "dog_id = ?", search)

	// returns found records count, equals `len(users)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(&dog)
}

func AddDog(c *fiber.Ctx) error {
	//twst3
	db := database.DBConn
	var dog m.Dogs

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&dog)
	return c.Status(201).JSON(dog)
}

func UpdateDog(c *fiber.Ctx) error {
	db := database.DBConn
	var dog m.Dogs
	id := c.Params("id")

	if err := c.BodyParser(&dog); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&dog)
	return c.Status(200).JSON(dog)
}

func RemoveDog(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var dog m.Dogs

	result := db.Delete(&dog, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func DeletedDog(c *fiber.Ctx) error {
	db := database.DBConn

	var deletedDog []m.Dogs

	if err := db.Unscoped().Where("deleted_at IS NOT NULL").Find(&deletedDog).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(deletedDog)
}

func GetDogsJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

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

	type ResultData struct {
		Data  []m.DogsRes `json:"data"`
		Name  string      `json:"name"`
		Count int         `json:"count"`
	}
	r := ResultData{
		Data:  dataResults,
		Name:  "golang-test",
		Count: len(dogs), //หาผลรวม,
	}
	return c.Status(200).JSON(r)
}

// สร้างข้อมูลในตารางdog มากกว่า10ตัว(api add dog)GetdogJsonสร้างapi
// ถ้าdog_id อยู่ระหว่าง 10-50 ให้โชว์คำว่า “red”ถ้าdog_id
// อยู่ระหว่าง 100-150 ให้โชว์คำว่า “green”ถ้าdog_id
// อยู่ระหว่าง 200-250 ให้โชว์คำว่า “pink”
// นอกเหนือจากนั้น “no color”
// ผลรวมแต่ละตัว

func GetDogsColor(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

	db.Find(&dogs) //10ตัว

	var dataResults []m.DogsRes
	colorCounts := map[string]int{
		"red":      0,
		"green":    0,
		"pink":     0,
		"no color": 0,
	}

	for _, v := range dogs { //1 inet 112 //2 inet1 113
		typeStr := ""
		if v.DogID >= 10 && v.DogID <= 50 {
			typeStr = "red"
		} else if v.DogID >= 100 && v.DogID <= 150 {
			typeStr = "green"
		} else if v.DogID >= 200 && v.DogID <= 250 {
			typeStr = "pink"
		} else {
			typeStr = "no color"
		}

		colorCounts[typeStr]++

		d := m.DogsRes{
			Name:  v.Name,  //inet
			DogID: v.DogID, //112
			Type:  typeStr, //no color
		}
		dataResults = append(dataResults, d)
		// sumAmount += v.Amount
	}

	sumRed := colorCounts["red"]
	sumGreen := colorCounts["green"]
	sumPink := colorCounts["pink"]
	sumNoColor := colorCounts["no color"]

	type ResultData struct {
		Data       []m.DogsRes `json:"data"`
		Name       string      `json:"name"`
		Count      int         `json:"count"`
		SumRed     int         `json:"sum_red"`
		SumGreen   int         `json:"sum_green"`
		SumPink    int         `json:"sum_pink"`
		SumNoColor int         `json:"sum_nocolo"`
	}

	r := ResultData{
		Data:       dataResults,
		Name:       "golang-test",
		Count:      len(dogs), //หาผลรวม,
		SumRed:     sumRed,
		SumGreen:   sumGreen,
		SumPink:    sumPink,
		SumNoColor: sumNoColor,
	}
	return c.Status(200).JSON(r)
}

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

func Register(c *fiber.Ctx) error {

	newUser := new(m.NewUser)
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	validate := validator.New()
	errors := validate.Struct(newUser)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors.Error())
	}
	return c.JSON(newUser)
}

// สร้างตารางcompany โดยใช้AutoMigrate โดยที่โครงสร้างcompanyควรจะมีอะไรบ้างใส่มาตามความเหมาะสม และGroupเพิ่มCRUD
func GetCompanys(c *fiber.Ctx) error {
	db := database.DBConn
	var company []m.Companys

	db.Find(&company) //delelete = null
	return c.Status(200).JSON(company)
}

func GetCompany(c *fiber.Ctx) error {
	db := database.DBConn
	search := strings.TrimSpace(c.Query("search"))
	var company []m.Companys

	result := db.Find(&company, "company_id = ?", search)

	// returns found records count, equals `len(users)
	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}
	return c.Status(200).JSON(&company)
}

func AddCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Companys

	if err := c.BodyParser(&company); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&company)
	return c.Status(201).JSON(company)
}

func UpdateCompany(c *fiber.Ctx) error {
	db := database.DBConn
	var company m.Companys
	id := c.Params("id")

	if err := c.BodyParser(&company); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&company)
	return c.Status(200).JSON(company)
}

func RemoveCompany(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var company m.Companys

	result := db.Delete(&company, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func GetCompanysJson(c *fiber.Ctx) error {
	db := database.DBConn
	var dogs []m.Dogs

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

	type ResultData struct {
		Data  []m.DogsRes `json:"data"`
		Name  string      `json:"name"`
		Count int         `json:"count"`
	}
	r := ResultData{
		Data:  dataResults,
		Name:  "golang-test",
		Count: len(dogs), //หาผลรวม,
	}
	return c.Status(200).JSON(r)
}

func GetProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profiles []m.Profile

	db.Find(&profiles) //delelete = null
	return c.Status(200).JSON(profiles)
}

func AddProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Create(&profile)
	return c.Status(201).JSON(profile)
}

func UpdateProfile(c *fiber.Ctx) error {
	db := database.DBConn
	var profile m.Profile
	id := c.Params("id")

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(503).SendString(err.Error())
	}

	db.Where("id = ?", id).Updates(&profile)
	return c.Status(200).JSON(profile)
}

func RemoveProfile(c *fiber.Ctx) error {
	db := database.DBConn
	id := c.Params("id")
	var profile m.Profile

	result := db.Delete(&profile, id)

	if result.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.SendStatus(200)
}

func GetGen(c *fiber.Ctx) error {
	db := database.DBConn
	var profile []m.Profile

	db.Find(&profile)

	var dataResults []m.Profile
	gendivider := map[string]int{
		"red":      0,
		"green":    0,
		"pink":     0,
		"no color": 0,
	}

	for _, v := range profile { //1 inet 112 //2 inet1 113
		typeStr := ""
		if v.Age < 24 {
			typeStr = "GenZ"
		} else if v.Age <= 41 {
			typeStr = "GenY"
		} else if v.Age <= 56 {
			typeStr = "GenX"
		} else if v.Age <= 75 {
			typeStr = "Baby Boomer"
		} else {
			typeStr = "G.I. Generation"
		}
		gendivider[typeStr]++
		d := m.Profile{
			EmployeeID: v.EmployeeID,
			Name:       v.Name,
			LastName:   v.LastName,
			BirthDay:   v.BirthDay,
			Age:        v.Age,
			Email:      v.Email,
			Tel:        v.Tel,
		}
		dataResults = append(dataResults, d)
		// sumAmount += v.Amount
	}

	z := gendivider["GenZ"]
	y := gendivider["GenY"]
	x := gendivider["GenX"]
	bb := gendivider["Baby Boomer"]
	gi := gendivider["G.I. Generation"]

	type ResultData struct {
		Data         []m.Profile `json:"data"`
		Count        int         `json:"count"`
		GenZ         int         `json:"genz"`
		GenY         int         `json:"geny"`
		GenX         int         `json:"genx"`
		BabyBoomer   int         `json:"babyboomer"`
		GIGeneration int         `json:"gi_generation"`
	}

	r := ResultData{
		Data:         dataResults,
		Count:        len(profile), //หาผลรวม,
		GenZ:         z,
		GenY:         x,
		GenX:         y,
		BabyBoomer:   bb,
		GIGeneration: gi,
	}
	return c.Status(200).JSON(r)
}

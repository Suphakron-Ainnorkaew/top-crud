package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}

type Employee struct {
	IDemployee int
	Name       string
	Age        string
	Address    string
	Salary     int
}

func main() {
	// เชื่อมต่อกับฐานข้อมูล PostgreSQL
	db, err = sql.Open("postgres", "user=myuser password=mypassword dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ทำการสร้าง router ด้วย Gin
	router := gin.Default()

	// กำหนด route สำหรับหน้าแรก
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		users, err := getUsersFromDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"users": users})
	})
	router.GET("/em", func(c *gin.Context) {
		employees, err := getEmployeeFromDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "emindex.html", gin.H{"employees": employees})
	})
	router.GET("/re", func(c *gin.Context) {
		employees, err := getEmployeeFromDB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.HTML(http.StatusOK, "reindex.html", gin.H{"employees": employees})
	})
	router.POST("/submit", func(c *gin.Context) {
		// อ่านค่าจากฟอร์ม
		IDemployee := c.PostForm("idemployee")
		name := c.PostForm("name")
		age := c.PostForm("age") // ต้องแก้ชื่อ input ให้ถูกต้อง
		address := c.PostForm("address")
		salary := c.PostForm("salary")

		// ดำเนินการบันทึกข้อมูลลงในฐานข้อมูล
		err := saveDataToDB(IDemployee, name, age, address, salary)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// แสดงหน้าเว็บหลังจากที่บันทึกข้อมูลเรียบร้อย
		c.HTML(http.StatusOK, "reindex.html", nil)
	})
	// Endpoint สำหรับลบข้อมูล
	router.POST("/delete", func(c *gin.Context) {
		employeeIDStr := c.PostForm("idemployee")
		employeeID, err := strconv.Atoi(employeeIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
			return
		}

		err = deleteEmployeeFromDB(employeeID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
	})

	// รันเซิร์ฟเวอร์
	router.Run(":8080")
}

func getUsersFromDB() ([]User, error) {
	rows, err := db.Query("SELECT id, username, password, email FROM public.\"User\";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func getEmployeeFromDB() ([]Employee, error) {
	rows, err := db.Query("SELECT idemployee, name, age, address, salary FROM public.\"Employees\";")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Ems []Employee
	for rows.Next() {
		var ems Employee
		err := rows.Scan(&ems.IDemployee, &ems.Name, &ems.Age, &ems.Address, &ems.Salary)
		if err != nil {
			return nil, err
		}
		Ems = append(Ems, ems)
	}
	return Ems, nil
}

// เพิ่มฟังก์ชัน saveDataToDB เพื่อบันทึกข้อมูลลงในฐานข้อมูล
func saveDataToDB(idemployee, name, age, address, salary string) error {
	// ตรวจสอบค่าว่าง และกำหนดค่าเริ่มต้นที่เหมาะสม
	if idemployee == "" {
		idemployee = "0" // หรือค่าเริ่มต้นที่เหมาะสม
	}

	if salary == "" {
		salary = "0" // หรือค่าเริ่มต้นที่เหมาะสม
	}

	// แปลงค่า idemployee และ salary เป็น integer
	idemployeeInt, err := strconv.Atoi(idemployee)
	if err != nil {
		return err
	}

	salaryInt, err := strconv.Atoi(salary)
	if err != nil {
		return err
	}

	// เขียนโค้ดสำหรับการเพิ่มข้อมูลลงในฐานข้อมูล
	_, err = db.Exec("INSERT INTO public.\"Employees\" (idemployee, name, age, address, salary) VALUES ($1, $2, $3, $4, $5)", idemployeeInt, name, age, address, salaryInt)
	if err != nil {
		return err
	}

	return nil
}

func deleteEmployeeFromDB(employeeID int) error {
	_, err := db.Exec("DELETE FROM public.\"Employees\" WHERE idemployee = $1;", employeeID)
	if err != nil {
		return err
	}
	return nil
}

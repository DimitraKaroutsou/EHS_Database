package main

import (
	"log"
	"time"
	"fmt"  //gia to create user

	"gorm.io/driver/postgres"
	"gorm.io/gorm"


	//"database/sql/driver" //gia json
	//"encoding/json"      //gia json
	//"errors"            //gia json

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"


	//"errors"

	//"golang.org/x/exp/errors"

	"github.com/gin-gonic/gin"

	"net/http"   //gia to login


	_ "github.com/lib/pq" // Import the PostgreSQL driver
)


type DeviceType string

type UserTable struct {
    ID       uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Username string       `json:"username"`
	Password string       `json:"password"`
	Email    string       `json:"email"` 
	Country  string       `json:"country"`
	Timezone string       `json:"timezone"`
    Houses   []HouseTable `json:"houses" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Rooms    []RoomTable   `json:"rooms" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
    //Metadata     map[string]interface{} `json:"metadata" gorm:"type:jsonb"`
	// Houses is a collection or array of HouseTable objects
	Permissions   []PermissionTable `json:"permissions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}


type HouseTable struct {
    ID          uuid.UUID     `json:"id" gorm:"type:uuid;primary_key"`
	Address     string        `json:"address"`
	TypeHouse   string        `json:"typehouse"`
	Area        int           `json:"area"`
	Year        int           `json:"year"`
	Floor       string        `json:"floor"`
	EnergyClass string        `json:"energyclass"`
	HeatingType string        `json:"heatingtype"`
	CoolingType string        `json:"coolingtype"`
    UserTableID uuid.UUID     `json:"user_table_id" gorm:”not null;”`
	Rooms       []RoomTable   `json:"rooms" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Devices      []DevicesTable `json:"devices" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//Shelly      []ShellyTable `json:"shelly" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Job         []JobTable    `json:"job" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//Broadlink   []BroadlinkTable `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RoomTable struct {
	ID                    uuid.UUID           `json:"id" gorm:"type:uuid;primary_key"`
	Name                  string              `json:"name"`
	RoomUsage             string              `json:"roomusage"`
	Area                  int                 `json:"area"`
	Windows               int                 `json:"windows"`
	HasAirCondition       string              `json:"hasaircondition"`
	AirConditionBTU       string              `json:"airconditionbtu"`
	AirConditionUsePeriod string              `json:"airconditionuseperiod"`
	AirConditionUsage     string              `json:"airconditionusage"`
	HouseTableID          uuid.UUID           `json:"house_table_id" gorm:”not null;”`
	UserTableID           uuid.UUID           `json:"user_table_id" gorm:”not null;”`
	Devices          []DevicesTable `json:"devices" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//AcController          []AcControllerTable `json:"accontroller" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	//Broadlink             []BroadlinkTable    `json:"broadlink" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Job              []JobTable `json:"job" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type DevicesTable struct{
	ID    uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Name  string     `json:"name" gorm:"not null;`
	Type  DeviceType       `json:"type" gorm:"type:device_type;check:type IN ('Sensor', 'Controller', 'Both')"`     
    HouseTableID uuid.UUID  `json:"house_table_id" gorm:”not null;”`
	RoomTableID  uuid.UUID  `json:"room_table_id"`
	CommunicationProtocol     map[string]interface{} `json:"communication_protocol" gorm:"type:jsonb"`
	Commands       map[string]interface{} `json:"commands" gorm:"type:jsonb"`
	Attributes     map[string]interface{} `json:"attributes" gorm:"type:jsonb"`
	Permissions     []PermissionTable `json:"permissions" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Job                 []JobTable `json:"job" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`	
}

type PermissionTable struct{
	ID    uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	UserTableID      uuid.UUID  `json:"user_table_id" gorm:”not null;”`
	DevicesTableID    uuid.UUID  `json:"devices_table_id" gorm:”not null;”`
}


type JobTable struct {
	ID                  uuid.UUID `gorm:"type:uuid;primary_key"`
	StartTime           string    `json:"starttime"`
	EndTime             string    `json:"endtime"`
	Days                string    `json:"days"`
	Descr               string    `json:"descr"`
	Tag                 uuid.UUID `json:"tag"`
	Timezone            string    `json:"timezone"`
	UserTableID         uuid.UUID `json:"user_table_id" gorm:"not null;"`
	HouseTableID        uuid.UUID `json:"house_table_id" gorm:"not null;"`
	RoomTableID         uuid.UUID `json:"room_table_id"`
	DevicesTableID      uuid.UUID `json:"devices_table_id" gorm:"default:null"`
	//AcControllerTableID uuid.UUID `json:"accontroller_table_id" gorm:"default:null"`
	//AcCommandsTableID   uuid.UUID `json:"accommands_table_id" gorm:"default:null"`
	//ShellyTableID       uuid.UUID `json:"shelly_table_id" gorm:"default:null"`
	JobCommand       map[string]interface{} `json:"job_command" gorm:"type:jsonb"`
	Value               int        `json:"value"`

}

type CacheData struct {
	ID            uuid.UUID      `json:"type:uuid;primary_key"`
	ShellyTableID uuid.UUID      `json:"shelly_table_id" gorm:"default:null"`
	UserTableID   uuid.UUID      `json:"user_table_id" gorm:"not null;"`
	CreationTime  datatypes.Date `json:"time"; gorm:"not null;"`
	Data          datatypes.JSON `json:"data"`
}


var db *gorm.DB

func initDB() {
	//the connecting string
	//dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	//	"localhost", "5432", "admin", "admin", "admin")

	host := "localhost"
	port := "5432"
	user := "admin"
	password := "admin"
	dbname := "admin"

	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"

	var err error
	//opening a connection to our database
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	//creating this table to the database
	db.AutoMigrate(&UserTable{},
		&HouseTable{},
		&RoomTable{},
		&DevicesTable{},
	    &PermissionTable{},
		&JobTable{})
	
    //test for data for the first table
	//var user1 User  //oxi idio onoma me to user ton credentials
	//db.Create(&User{Username: "Dimitra", Password: "12345"})
	//db.Delete(&user1, 1) //fainetai i eggrafi alla fainetai to deleted


//endpoints

// Create a new Gin router with default middleware
r := gin.Default()

// Define a route for the root path,oste sto http://localhost:8080/ na mi vgazei page not found
r.GET("/", func(c *gin.Context) {
 c.JSON(http.StatusOK, gin.H{"message": "Welcome to my application!"})
})

//POST: Submit data to be processed to a specified resource.
// Define an endpoint for creating a user account
r.Any("/users/create/account", createUserAccount)   //defined to handle HTTP POST requests.

// Define an endpoint for user login
r.Any("/users/login", userLogin)

// Define an endpoint for adding a house
r.Any("/house/add", addHouse)

// Define an endpoint for getting all houses
r.Any("/house/getall", getAllHouses)

// Define an endpoint for getting a specific house by ID
r.Any("/house/get/:houseId", getHouseByID)

// Define an endpoint for deleting a specific house by ID
r.DELETE("/house/delete/:houseId", deleteHouseByID)

// Run the server on port 8080,most common,not the same with the database port
r.Run(":8080")

}

func createUserAccount (c *gin.Context) {

	if c.Request.Method == "GET" {
		// Handle GET request (e.g., render an HTML form)
		c.JSON(http.StatusOK, gin.H{"message": "Render your HTML form for user account creation"})
		return
	}

	var newUser UserTable

	 // Try to bind the JSON request body to the newUser variable
	 if err := c.ShouldBindJSON(&newUser); err != nil {
        // If binding fails (invalid JSON), respond with a 400 Bad Request and an error message
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid JSON: %s", err.Error())})
		return
    }

	// Create a new user in the database
	result := db.Create(&newUser)
    
	// If the insertion fails, respond with a 500 Internal Server Error and an error message
	if result.Error != nil {
		// If the insertion fails, respond with a 500 Internal Server Error and an error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
		return
	}
   
	// If everything is successful, respond with a 200 OK and a success message
	c.JSON(http.StatusOK, gin.H{"message": "User account created successfully"})
}


func userLogin(c *gin.Context) {

	if c.Request.Method == "GET" {
		// Handle GET request (e.g., render an HTML form)
		c.JSON(http.StatusOK, gin.H{"message": "Render your HTML form for user account creation"})
		return
	}

	var loginUser UserTable
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Add logic for user login (check credentials, generate token, etc.)
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "User login successful"})
}


// Handler function for adding a house
func addHouse(c *gin.Context) {

	if c.Request.Method == "GET" {
		// Handle GET request (e.g., render an HTML form)
		c.JSON(http.StatusOK, gin.H{"message": "Render your HTML form for user account creation"})
		return
	}

	var newHouse HouseTable
	if err := c.ShouldBindJSON(&newHouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Add logic for adding a new house to the database
    result := db.Create(&newHouse)

	// If the insertion fails, respond with a 500 Internal Server Error and an error message
	if result.Error != nil {
		// If the insertion fails, respond with a 500 Internal Server Error and an error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create house"})
		return
	}
   
	// If everything is successful, respond with a 200 OK and a success message
	c.JSON(http.StatusOK, gin.H{"message": "House added successfully"})

}

func getAllHouses(c *gin.Context) {
    // Query all houses from the database
    var houses []HouseTable
    result := db.Find(&houses)

    // Check for the database query error
    if result.Error != nil {
        // Log the error
        fmt.Println("Database error:", result.Error)

        // Respond with a 500 Internal Server Error and an error message
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve houses from the database"})
        return
    }

    // Respond with the list of houses
    c.JSON(http.StatusOK, houses)
}

func getHouseByID(c *gin.Context) {
    // Get the house ID from the URL parameters
    houseID := c.Param("houseId")

    // Validate and parse the house ID
    parsedHouseID, err := uuid.Parse(houseID)
    if err != nil {
        // Respond with a 400 Bad Request and an error message
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
        return
    }

    // Query the specific house from the database
    var house HouseTable
    result := db.First(&house, parsedHouseID)

    // Check for the database query error
    if result.Error != nil {
        // Log the error
        fmt.Println("Database error:", result.Error)

        // Respond with a 500 Internal Server Error and an error message
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve the house from the database"})
        return
    }

    // Respond with the specific house
    c.JSON(http.StatusOK, house)
}

func deleteHouseByID(c *gin.Context) {
    // Get the house ID from the URL parameters
    houseID := c.Param("houseId")

    // Validate and parse the house ID
    parsedHouseID, err := uuid.Parse(houseID)
    if err != nil {
        // Respond with a 400 Bad Request and an error message
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid house ID"})
        return
    }

    // Delete the specific house from the database
    result := db.Delete(&HouseTable{}, parsedHouseID)

    // Check for the database deletion error
    if result.Error != nil {
        // Log the error
        fmt.Println("Database error:", result.Error)

        // Respond with a 500 Internal Server Error and an error message
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the house from the database"})
        return
    }

    // Respond with a success message
    c.JSON(http.StatusOK, gin.H{"message": "House deleted successfully"})
}




func main() {
	initDB()
}

func (u *UserTable) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	tempPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(tempPassword)
	return
}

func (u *HouseTable) BeforeCreate(tx *gorm.DB) (err error) {
	// Generate a new UUID and assign it to the ID field
	// Check if ID is not provided (zero value)
	if u.ID == uuid.Nil {
		// Set a default UUID value
		u.ID = uuid.New()
	}

	return
}

func (u *RoomTable) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (u *DevicesTable) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (u *PermissionTable) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}


func (u *JobTable) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()

	return
}

func (u *CacheData) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	u.CreationTime = datatypes.Date(time.Now())

	return
}
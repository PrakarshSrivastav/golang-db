package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

const Version = "1.0.0"

type (
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}

	Driver struct {
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		dir     string
		log     Logger
	}
)
type Options struct {
	Logger
}

func New(dir string, options *Options) (*Driver, error) {
	dir = filepath.Clean(dir)
	opts := Options{}
	if options != nil {
		opts = *options
	}
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger((lumber.INFO))
	}

	driver := Driver{
		dir:     dir,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}
	if _, err := os.Stat(dir); err == nil {
		opts.Logger.Debug("Using '%s'(database already exists)\n", dir)
		return &driver, nil
	}

	opts.Logger.Debug("Creating a database at '%s'...\n", dir)
	return &driver, os.MkdirAll(dir, 0755)

}

func (d *Driver) Write() error {

}

func (d *Driver) Read() error {

}

func (d *Driver) ReadAll() {

}

func (d *Driver) Delete() error {

}

func (d *Driver) getOrCreateMutex() *sync.Mutex {

}

type Address struct {
	City    string
	State   string
	Country string
	Pincode json.Number
}

type User struct {
	Name    string
	Age     json.Number
	Contact string
	Company string
	Address Address
}

func main() {
	dir := "./"

	db, err := New(dir, nil)
	if err != nil {
		fmt.Println("Error", err)
	}

	employees := []User{
		{"John", "23", "22333445", "Golden State", Address{"Lucknow", "Uttar Pradesh", "India", "441003"}},
		{"Paul", "25", "22311234", "Lakers", Address{"Jaipur", "Gujarat", "India", "441002"}},
		{"Curry", "30", "22388776", "Golden State", Address{"Dehradun", "Uttrakhand", "India", "441067"}},
		{"Tom", "19", "22388445", "Portland Blazers", Address{"Wuse", "Abuja", "Nigeria", "441878"}},
		{"David", "35", "66554433", "CloudCore", Address{"Hyderabad", "Telangana", "India", "500001"}},
		{"Sophia", "28", "55443322", "NextGen Labs", Address{"Chennai", "Tamil Nadu", "India", "600001"}},
	}

	for _, value := range employees {
		db.Write("users", value.Name, User{
			Name:    value.Name,
			Age:     value.Age,
			Contact: value.Contact,
			Company: value.Company,
			Address: value.Address,
		})
	}
	records, err := db.ReadAll("users")
	if err != nil {
		fmt.Println("Error", err)
	}
	fmt.Println(records)

	allusers := []User{}

	for _, f := range records {
		employeeFound := User{}
		if err := json.Unmarshal([]byte(f), &employeeFound); err != nil {
			fmt.Println("Error", err)
		}
		allusers = append(allusers, employeeFound)
	}
	fmt.Println((allusers))
	// db.Delete("user", "john")
}

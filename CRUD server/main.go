//request data should be sent in JSON format
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

//for GET requests
func fetchData(parsedUrl []string) []byte {
	var rows *sql.Rows
	var err error
	if len(parsedUrl) == 2 {
		rows, err = db.Query("SELECT * FROM "+parsedUrl[0]+" WHERE id=?", parsedUrl[1])
	} else {
		rows, err = db.Query("SELECT * FROM " + parsedUrl[0])
	}

	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println("Failed to get columns", err)
		return nil
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	results := make(map[string]string)

	slice := []map[string]string{}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			results[columns[i]] = value

		}
		slice = append(slice, results)

	}
	json, _ := json.Marshal(slice)

	fmt.Printf("%T\n", json)

	fmt.Println(slice)
	return json

}

//for POST requests
func insertData(data []byte, parsedUrl []string) error {
	statement := ""
	valueString := ""
	colummString := ""
	slice := map[string]string{}
	err := json.Unmarshal(data, &slice)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	fmt.Println(slice)
	for col, value := range slice {
		valueString += "'" + value + "',"
		colummString += col + ","
	}
	fmt.Println(colummString[:len(colummString)-1])
	statement = "INSERT INTO " + parsedUrl[0] + " ( " + colummString[:len(colummString)-1] + " ) VALUES (" + valueString[:len(valueString)-1] + ")"

	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	rows.Close()

	return nil
}

//for PUT requests
func updateData(data []byte, parsedUrl []string) error {
	statement := ""
	resultString := ""
	slice := map[string]string{}
	err := json.Unmarshal(data, &slice)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	for col, value := range slice {
		resultString += col + "='" + value + "',"
	}
	statement = "UPDATE " + parsedUrl[0] + " SET " + resultString[:len(resultString)-1] + " WHERE id=" + parsedUrl[1] + ""
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	rows.Close()

	return nil
}

//for DELETE requests
func deleteData(parsedUrl []string) error {

	statement := "DELETE FROM " + parsedUrl[0] + " WHERE id=" + parsedUrl[1] + ""
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println("error:", err)
		return err
	}
	rows.Close()

	return nil
}

func parseUrl(url string) []string {
	parsed := strings.Split(url, "/")
	return parsed[1:]
}

func handler(w http.ResponseWriter, r *http.Request) {

	db = openConn()
	defer db.Close()
	parsedUrl := parseUrl(r.URL.String())

	switch r.Method {
	case "GET":
		{ //READ
			response := fetchData(parsedUrl)
			w.Write(response)
		}
	case "PUT":
		{
			//UPDATE

			data, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			err := updateData(data, parsedUrl)
			if err != nil {
				fmt.Println("error:", err)
			}
			w.Write([]byte("200 OK"))

		}
	case "POST":
		{
			//CREATE
			data, _ := ioutil.ReadAll(r.Body)
			r.Body.Close()
			err := insertData(data, parsedUrl)
			if err != nil {
				fmt.Println("error:", err)
			}
			w.Write([]byte("200 OK"))

		}
	case "DELETE":
		{
			//DELETE
			err := deleteData(parsedUrl)
			if err != nil {
				fmt.Println("error:", err)
			}
			w.Write([]byte("200 OK"))

		}
	default:
		{
			log.Print("Cannot handle method ", r.Method)
			http.Error(w, "Only GET method", http.StatusNotImplemented)
			return
		}

	}
}

func openConn() *sql.DB {
	db, err := sql.Open("mysql", "root@/golang_db")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {

	http.HandleFunc("/", handler)
	http.ListenAndServe(":9999", nil)
}

/*
MIT License

Copyright (c) 2024 Elliot Michael Keavney

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"tableviewresource"
)

// Function to retrieve databases inside MySQL
func databaseList(dbUsername string, dbPassword string, dbTransport string, dbAddress string, dbPort string, dbTls string) []string {

	// Open database connection
	dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/?tls="+dbTls)
	defer dbConnection.Close()

	// Error
	if err != nil {
		panic("Function databaseList: Is the database online?")
	}

	databaseQuery, err := dbConnection.Query("SHOW DATABASES WHERE `Database` NOT IN ('mysql', 'performance_schema', 'information_schema', 'sys');")
	var databaseListResult []string

	for databaseQuery.Next() {

		var row string

		err = databaseQuery.Scan(&row)

		// Error
		if err != nil {
			panic("Error in database list function")
		}
		databaseListResult = append(databaseListResult, row)
	}
	return databaseListResult
}

// Function to provide table names
func provideTableName(dbConnection *sql.DB, w http.ResponseWriter, dbName string) {

	// SQL query returns table name(s)
	result, err := dbConnection.Query("SELECT TABLE_NAME AS tableName FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA = ?;", dbName)

	var tableName string

	tableInfo(dbConnection, w, dbName)
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<hr class=\"line\">")
	fmt.Fprintf(w, "<br>")

	for result.Next() {

		err = result.Scan(&tableName)

		// Error
		if err != nil {

			panic("SQL query for table names not working")
		}
		showColumn(dbConnection, w, tableName)
	}
}

// Function to retrieve table information
func tableInfo(dbConnection *sql.DB, w http.ResponseWriter, dbName string) {

	// SQL query returns infomation about the table(s)
	result, err := dbConnection.Query("SELECT table_name AS tableName, Table_Type AS tableType, create_time AS createTime, Engine AS engine, TABLE_COLLATION AS tableCollation FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = ?;", dbName)

	// Error
	if err != nil {
		panic("SQL query for table information not working")
	}

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"tableInformation\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><h3>&nbsp &nbsp &nbsp Table Infomation &nbsp &nbsp &nbsp</h3></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp &nbsp Search Table Name &#8680 &nbsp &nbsp</th>")
	fmt.Fprintf(w, "    <th><input type=\"text\" id=\"tableInput\" onkeyup=\"tableFunction()\" placeholder=\"Type to filter table...\" title=\"search\"></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\"table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><b>Table Name</b></th>")
	fmt.Fprintf(w, "    <th><b>Table Type</b></th>")
	fmt.Fprintf(w, "    <th><b>Table Created</b></th>")
	fmt.Fprintf(w, "    <th><b>Table Engine</b></th>")
	fmt.Fprintf(w, "    <th><b>Character Set</b></th>")
	fmt.Fprintf(w, "    <th><b>Go to Table</b></th>")
	fmt.Fprintf(w, "  </tr>")

	for result.Next() {

		var (
			tableName      string
			tableType      string
			createTime     string
			engine         string
			tableCollation string
		)

		err = result.Scan(&tableName, &tableType, &createTime, &engine, &tableCollation)

		// Error
		if err != nil {
			panic("Error in tableInfo function")
		}

		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>"+tableName+"</td>")
		fmt.Fprintf(w, "    <td>"+tableType+"</td>")
		fmt.Fprintf(w, "    <td>"+createTime+"</td>")
		fmt.Fprintf(w, "    <td>"+engine+"</td>")
		characterSet := tableCollation[:len(tableCollation)-10]
		fmt.Fprintf(w, "    <td>"+characterSet+"</td>")
		fmt.Fprintf(w, "    <td><a href=\"#"+tableName+"\" class=\"tableButton\">&#11015</a></td>")
		fmt.Fprintf(w, "  </tr>")
	}

	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "function tableFunction() {")
	fmt.Fprintf(w, "  var input, filter, table, tr, td, i, txtValue;")
	fmt.Fprintf(w, "  input = document.getElementById(\"tableInput\");")
	fmt.Fprintf(w, "  filter = input.value.toUpperCase();")
	fmt.Fprintf(w, "  table = document.getElementById(\"table\");")
	fmt.Fprintf(w, "  tr = table.getElementsByTagName(\"tr\");")
	fmt.Fprintf(w, "  for (i = 0; i < tr.length; i++) {")
	fmt.Fprintf(w, "    td = tr[i].getElementsByTagName(\"td\")[0];")
	fmt.Fprintf(w, "    if (td) {")
	fmt.Fprintf(w, "      txtValue = td.textContent || td.innerText;")
	fmt.Fprintf(w, "      if (txtValue.toUpperCase().indexOf(filter) > -1) {")
	fmt.Fprintf(w, "        tr[i].style.display = \"\";")
	fmt.Fprintf(w, "      } else {")
	fmt.Fprintf(w, "        tr[i].style.display = \"none\";")
	fmt.Fprintf(w, "      }")
	fmt.Fprintf(w, "    }")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "}")
	fmt.Fprintf(w, "</script>")
}

// Function to count rows and retrieve column names
func showColumn(dbConnection *sql.DB, w http.ResponseWriter, tableName string) {

	// SQL query counts rows in table
	result1, result1Err := dbConnection.Query("SELECT COUNT(*) AS rowCount FROM " + tableName + ";")

	// Error
	if result1Err != nil {
		panic("SQL query for row count not working")
	}

	// SQL query returns columns in table
	result2, result2Err := dbConnection.Query("SELECT COLUMN_NAME AS columnName, COLUMN_TYPE AS columnType, IS_NULLABLE AS isNullable, COLUMN_KEY AS columnKey FROM information_schema.COLUMNS WHERE TABLE_NAME = ?;", tableName)

	// Error
	if result2Err != nil {
		panic("SQL query for columns not working")
	}

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\""+tableName+"\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><h3>&nbsp &nbsp &nbsp Column Infomation for Table &nbsp &nbsp &nbsp<br>&nbsp &nbsp &nbsp"+tableName+"&nbsp &nbsp &nbsp</h3></th>")
	fmt.Fprintf(w, "    <th>&nbsp &nbsp &nbsp<a href=\"/#tableInformation\" class=\"tableButton\">&#11014</a>&nbsp &nbsp &nbsp</td>")
	fmt.Fprintf(w, "  </tr>")

	for result1.Next() {

		var rowCount string

		result1Err = result1.Scan(&rowCount)

		// Error
		if result1Err != nil {
			panic("")
		}

		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>&nbsp Total Rows / Records in Table &nbsp</th>")
		fmt.Fprintf(w, "    <th><b>&nbsp"+rowCount+"&nbsp</b></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
	}

	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table>")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th>&nbsp &nbsp Search Column Name &#8680 &nbsp &nbsp</th>")
	fmt.Fprintf(w, "    <th><input type=\"text\" id=\""+tableName+"Input\" onkeyup=\"tableFunction"+tableName+"()\" placeholder=\"Type to filter table "+tableName+"...\" title=\"search\"></th>")
	fmt.Fprintf(w, "  </tr>")
	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<table id=\""+tableName+"Table\">")
	fmt.Fprintf(w, "  <tr>")
	fmt.Fprintf(w, "    <th><b>Column Name</b></th>")
	fmt.Fprintf(w, "    <th><b>Column Type</b></th>")
	fmt.Fprintf(w, "    <th><b>Nullable</b></th>")
	fmt.Fprintf(w, "    <th><b>Key</b></th>")
	fmt.Fprintf(w, "  </tr>")

	for result2.Next() {

		var (
			columnName string
			columnType string
			isNullable string
			columnKey  string
		)

		result2Err = result2.Scan(&columnName, &columnType, &isNullable, &columnKey)

		// Error
		if result2Err != nil {
			panic("Error in showColumn function")
		}

		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <td>"+columnName+"</td>")
		fmt.Fprintf(w, "    <td>"+columnType+"</td>")
		fmt.Fprintf(w, "    <td>"+isNullable+"</td>")
		fmt.Fprintf(w, "    <td>"+columnKey+"</td>")
		fmt.Fprintf(w, "  </tr>")
	}

	fmt.Fprintf(w, "</table>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<hr class=\"line\">")
	fmt.Fprintf(w, "<br>")
	fmt.Fprintf(w, "<script>")
	fmt.Fprintf(w, "function tableFunction"+tableName+"() {")
	fmt.Fprintf(w, "  var input, filter, table, tr, td, i, txtValue;")
	fmt.Fprintf(w, "  input = document.getElementById(\""+tableName+"Input\");")
	fmt.Fprintf(w, "  filter = input.value.toUpperCase();")
	fmt.Fprintf(w, "  table = document.getElementById(\""+tableName+"Table\");")
	fmt.Fprintf(w, "  tr = table.getElementsByTagName(\"tr\");")
	fmt.Fprintf(w, "  for (i = 0; i < tr.length; i++) {")
	fmt.Fprintf(w, "    td = tr[i].getElementsByTagName(\"td\")[0];")
	fmt.Fprintf(w, "    if (td) {")
	fmt.Fprintf(w, "      txtValue = td.textContent || td.innerText;")
	fmt.Fprintf(w, "      if (txtValue.toUpperCase().indexOf(filter) > -1) {")
	fmt.Fprintf(w, "        tr[i].style.display = \"\";")
	fmt.Fprintf(w, "      } else {")
	fmt.Fprintf(w, "        tr[i].style.display = \"none\";")
	fmt.Fprintf(w, "      }")
	fmt.Fprintf(w, "    }")
	fmt.Fprintf(w, "  }")
	fmt.Fprintf(w, "}")
	fmt.Fprintf(w, "</script>")
}

func main() {

	err := godotenv.Load("/etc/tableview/tableview.env")
	if err != nil {
		panic("Error loading tableview.env file for database details")
	}

	//Get the database connection details
	dbUsername := os.Getenv("dbUsername")
	dbPassword := os.Getenv("dbPassword")
	dbTransport := os.Getenv("dbTransport")
	dbAddress := os.Getenv("dbAddress")
	dbPort := os.Getenv("dbPort")
	dbTls := os.Getenv("dbTls")

	//Values allowed for dbTransport Variable
	var allowedTransportValue = []string{"tcp", "udp"}
	validDbTransport := slices.Contains(allowedTransportValue, dbTransport)

	dbPortInt, err := strconv.Atoi(dbPort)
	if err != nil {
		panic("DATABASE PORT MUST BE A NUMBER IN /etc/tableview/tableview.env")
	}

	//Values allowed for dbTls Variable
	var allowedTlsValue = []string{"false", "true"}
	validDbTls := slices.Contains(allowedTlsValue, dbTls)

	//Catch if any errors were made in tableview.env and feed back where to correct error
	if dbUsername == "" {
		panic("DATABASE USERNAME CANNOT BE BLANK IN /etc/tableview/tableview.env")
	} else if dbPassword == "" {
		panic("DATABASE PASSOWRD CANNOT BE BLANK IN /etc/tableview/tableview.env")
	} else if dbTransport == "" {
		panic("DATABASE TRANSPORT OPTION CANNOT BE BLANK IN /etc/tableview/tableview.env")
	} else if validDbTransport == false {
		panic("DATABASE TRANSPORT OPTION MUST BE udp OR tcp IN /etc/tableview/tableview.env")
	} else if dbAddress == "" {
		panic("DATABASE ADDRESS CANNOT BE BLANK IN /etc/tableview/tableview.env")
	} else if dbPortInt <= 0 || dbPortInt >= 65536 {
		panic("DATABASE PORT MUST BE IN THE NUMBER RANGE 1-65535 IN /etc/tableview/tableview.env")
	} else if dbTls == "" {
		panic("DATABASE TLS OPTION CANNOT BE BLANK IN /etc/tableview/tableview.env")
	} else if validDbTls == false {
		panic("DATABASE TRANSPORT OPTION MUST BE false OR true IN /etc/tableview/tableview.env")
	}

	databaseListResult := databaseList(dbUsername, dbPassword, dbTransport, dbAddress, dbPort, dbTls)

	var startHTML string
	startHTML = tableviewresource.StartHTML()

	var endHTML string
	endHTML = tableviewresource.EndHTML()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Fprintf(w, startHTML)
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><a href=\"https://ell.today\" class=\"tableButton externalButton\">Written by Elliot Keavney (Website)</a></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th><a href=\"https://github.com/ellwould/table-view\" class=\"tableButton externalButton\">Table View Source Code (GitHub)</a></th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<table>")
		fmt.Fprintf(w, "  <tr>")
		if dbTls == "false" {
			fmt.Fprintf(w, "<th>The connection between the MySQL server<br>and Table View is unencrypted &#128308</th>")
		} else if dbTls == "true" {
			fmt.Fprintf(w, "<th>The connection between the MySQL server<br>and Table View is encrypted &#128994</th>")
		}
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>MySQL server username: "+dbUsername+"</th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "  <tr>")
		fmt.Fprintf(w, "    <th>MySQL server address: "+dbAddress+" </th>")
		fmt.Fprintf(w, "  </tr>")
		fmt.Fprintf(w, "</table>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<form method=\"POST\" action=\"/\">")
		fmt.Fprintf(w, "  <select id=\"dbName\" name=\"dbName\">")
		fmt.Fprintf(w, "  <option value=>Select a Database</option>")
		for i := 0; i < len(databaseListResult); i++ {
			fmt.Fprintf(w, "<option value="+databaseListResult[i]+">"+databaseListResult[i]+"</option>")
		}
		fmt.Fprintf(w, "  </select>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "  <input type=\"submit\" value=\"submit\" />")
		fmt.Fprintf(w, "</form>")
		fmt.Fprintf(w, "<br>")
		fmt.Fprintf(w, "<br>")

		//Get database name and validate
		form := r.FormValue("dbName")
		var dbName string
		dbName = form
		validDbName := slices.Contains(databaseListResult, dbName)

		if dbName == "" {
			fmt.Fprintf(w, endHTML)
		} else if validDbName == true {

			// Open database connection
			dbConnection, err := sql.Open("mysql", dbUsername+":"+dbPassword+"@"+dbTransport+"("+dbAddress+":"+dbPort+")/"+dbName+"?tls="+dbTls)
			defer dbConnection.Close()

			// Error
			if err != nil {
				panic("Is database online?")
			}

			tableCountQuery, err := dbConnection.Query("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?", dbName)

			// Error
			if err != nil {
				panic("SQL query for tableCount not working")
			}

			for tableCountQuery.Next() {

				var tableCount string
				err = tableCountQuery.Scan(&tableCount)

				// Error
				if err != nil {
					panic("")
				}

				fmt.Fprintf(w, "<table>")
				fmt.Fprintf(w, "  <tr>")
				fmt.Fprintf(w, "    <th><h3>&nbsp &nbsp &nbsp Infomation about Tables in Database "+dbName+"&nbsp &nbsp &nbsp</h3></th>")
				fmt.Fprintf(w, "  </tr>")
				fmt.Fprintf(w, "</table>")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<table>")
				fmt.Fprintf(w, "  <tr>")
				fmt.Fprintf(w, "    <th><b>&nbsp Total Number of Tables in Database "+dbName+"&nbsp<b></th>")
				fmt.Fprintf(w, "    <th><b>&nbsp"+tableCount+"&nbsp</b></th>")
				fmt.Fprintf(w, "  </tr>")
				fmt.Fprintf(w, "</table>")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<br>")
				fmt.Fprintf(w, "<hr class=\"line\">")
				fmt.Fprintf(w, "<br>")
				provideTableName(dbConnection, w, dbName)
				fmt.Fprintf(w, endHTML)
			}
		} else {
			fmt.Fprintf(w, "<table>")
			fmt.Fprintf(w, "  <tr>")
			fmt.Fprintf(w, "    <th><h1>&nbsp &nbsp DATABASE DOES NOT EXIST &nbsp &nbsp</h1></th>")
			fmt.Fprintf(w, "  </tr>")
			fmt.Fprintf(w, "</table>")
			fmt.Fprintf(w, endHTML)
		}
	})

	tvPort := os.Getenv("tvPort")
	tvPortInt, err := strconv.Atoi(tvPort)

	if err != nil {
		panic("TABLE VIEW PORT MUST BE A NUMBER IN /etc/tableview/tableview.env")
	}

	if tvPortInt <= 1023 || tvPortInt >= 49152 {
		panic("TABLE VIEW LISTENING PORT MUST BE IN THE NUMBER RANGE 1024-49151 IN /etc/tableview/tableview.env")
	} else {
		socket := "localhost:" + tvPort
		fmt.Println("Table View is running on port " + socket)

		// Start server on port specified above
		log.Fatal(http.ListenAndServe(socket, nil))
	}
}

// Contributor(s):
// Elliot Michael Keavney

package main

import (
	"fmt"
	"github.com/prometheus/common/log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx"
	"database/sql"
	"time"
	//"strings"
	"text/template"
	"path/filepath"
	"strings"
	"os"
	"gopkg.in/gcfg.v1"
	"io/ioutil"
)

type Conf struct {
	Mysql struct {
		      Host     string
		      Port     string
		      Username string
		      Password string
		      Database string
	      }
	Other struct {
		      Savedir string
	      }
}

var (
	database string
	selectSql string
	putFileName string
	config Conf
)
var defautConfigFile string = "config.ini"

func readConfig() (Config Conf) {
	err := gcfg.ReadFileInto(&Config, defautConfigFile)
	if err != nil {
		fmt.Println("Failed to parse config file: %s", err)
	}
	return Config
}

func writeConfig() (Config Conf) {
	var wireteString = `[Mysql]
host = localhost
port = 3306
username = root
password =
database = Excel
[Other]
savedir = Excel`
	var filename = defautConfigFile

	var d1 = []byte(wireteString)
	ioutil.WriteFile(filename, d1, 0666)
	fmt.Println("写入config")
	config := readConfig()
	fmt.Println(config)
	return config
}

func createSaveDir(dir string) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		logErrorStr := fmt.Sprintf("create save dir: %s  filed", dir)
		log.Error(logErrorStr)
	}
}

func sqlToExcel(w http.ResponseWriter, r *http.Request) {
	selectSql := r.PostFormValue("sql")
	database := r.PostFormValue("db")
	putFileName := r.PostFormValue("name")

	if database == "" {
		database = config.Mysql.Database
	}
	//file name
	if (putFileName == "" ) {
		t := time.Now()
		putFileName = fmt.Sprintf("%d-%02d-%02d_%02d%02d%02d",
			t.Year(), t.Month(), t.Day(),
			t.Hour(), t.Minute(), t.Second())
	}
	if (selectSql != "" ) {
		//db connect
		connectStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
			config.Mysql.Username,
			config.Mysql.Password,
			config.Mysql.Host,
			config.Mysql.Port,
			database)
		db, err := sql.Open("mysql", connectStr)
		if err != nil {
			log.Error(err.Error())
			fmt.Fprintf(w, "{\"code\":10001,\"message\":\"db connect error\"}")
		}
		defer db.Close()

		// Execute the query
		rows, err := db.Query(selectSql)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			log.Error(err.Error()) // proper error handling instead of panic in your app
			fmt.Fprintf(w, "{\"code\":10001,\"message\":\"outpu excel columns error \"}")
		}

		// Make a slice for the values
		values := make([]sql.RawBytes, len(columns))
		new_file := xlsx.NewFile()
		new_sheet, err := new_file.AddSheet("Sheet1")

		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Fetch rows
		titleIndex := 0
		for rows.Next() {

			titleIndex ++
			// get RawBytes from data
			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			// Now do something with the data.
			// Here we just print each column as a string.
			var value string
			if titleIndex == 1 {
				new_row := new_sheet.AddRow()
				for i, col := range values {
					// Here we can check if the value is nil (NULL value)
					if col == nil {
						value = "NULL"
					} else {
						value = string(col)
					}

					new_cell := new_row.AddCell()
					new_cell.Value = columns[i]
				}

			}
			new_row := new_sheet.AddRow()
			for _, col := range values {
				// Here we can check if the value is nil (NULL value)
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}

				new_cell := new_row.AddCell()
				new_cell.Value = value
			}
		}

		savePath := fmt.Sprintf("%s/%s.xlsx", config.Other.Savedir, putFileName)
		new_file.Save(savePath)
		baseDir, _ := filepath.Abs("./")
		saveFullPath := fmt.Sprintf("%s/%s/%s.xlsx", baseDir, config.Other.Savedir, putFileName)
		fmt.Fprintf(w, "{\"code\":10000,\"message\":\"success\",\"name\":\"" + saveFullPath + "\"}")

	} else {
		fmt.Fprintf(w, "{\"code\":10001,\"message\":\"params sql error\"}")
	}

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func index(rw http.ResponseWriter, req *http.Request) {

	langs := strings.Split(req.Header.Get("Accept-Language"), ",")
	locale := strings.ToLower(langs[0])

	t, _ := template.ParseFiles("./view/index.tpl")
	t.Execute(rw, locale)

}

func init() {
	exists, err := PathExists(defautConfigFile)
	log.Info(exists, err)
	if exists == true && err == nil {
		config = readConfig()
	} else {
		config = writeConfig()
	}
	saveDir := config.Other.Savedir
	exists, _ = PathExists(saveDir)
	if exists == false {
		createSaveDir(saveDir)
	}
}

func main() {

	log.Info("server listening on 9010 open url http://localhost:9010")
	//http.Handle("/static/", http.FileServer(http.Dir("static")))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/MysqlToExcel", sqlToExcel)
	err := http.ListenAndServe(":9010", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
package main 

import (
	"golang.org/x/crypto/ssh/terminal"
	log "github.com/sirupsen/logrus"
	"fmt"
	"TD-EXPORT/lib"
	"os"
	"flag"
	"strings"
)

const (
	appName	= "TD-EXPORT"
	desc 	= "Table Definition Export"
	version	= "v0.2.0"
)

func main() {
	var output string 
	flag.StringVar(&output,"output","excel","KTDEX Output(excel / markdown / sql)")
	flag.Parse()

	log.Infof("%s %s",appName,version)
	log.Infof("%s",desc)
	fmt.Println()

	var conn lib.TdexStructure
	var err error

	// DB Endpoint
	conn.Endpoint, err = lib.GetOpt("Endpoint")
	if err != nil {
		log.Errorf("%s",err)
		os.Exit(1)
	}

	// DB Port
	conn.Port, err = lib.GetOpt("Port")
	if err != nil {
		log.Warningf("%s",err)
		log.Infof("Default set Port : 3306")
	}
	if conn.Port == nil {
		conn.Port = lib.PointerStr("3306")
	}

	// DB User(Information_schema Allow user)
	conn.User, err = lib.GetOpt("User")
	if err != nil {
		log.Errorf("%s",err)
		os.Exit(1)
	}

	// DB User Password
	fmt.Printf("Password : ")
	vPass, err := terminal.ReadPassword(0)
	if err != nil {
		log.Errorf("%s",err)
		os.Exit(1)
	}
	passStr := string(vPass)
	conn.Pass = &passStr
	fmt.Println("")

	// Target DB
	conn.TargetDB, err = lib.GetOpt("DB(Seperator , or Space(All))")
	if err != nil {
		log.Warningf("%s",err)
	}

	// Target DB
	conn.ExceptTable, err = lib.GetOpt("Exception Tables(Seperator , or Space(none) / Use wildcard)")
	if err != nil {
		log.Warningf("%s",err)
	}

	// Create Database Object
	conn.Object, err = conn.GetDBObject()
	if err != nil {
		log.Errorf("%s",err)
		os.Exit(1)
	}
	defer conn.Object.Close()
	log.Infof("DB Connect Success")

	// Get Schema Name
	conn.Schema, err = conn.GetSchema()
	if err != nil {
		log.Errorf("%s",err)
		os.Exit(1)
	}

	if len(conn.Schema) < 1 {
		log.Infof("Not in Schema.")
		os.Exit(1)
	}

	var excelz lib.ExcelStructure
	var mdfiles lib.MdStructure
	var sqlfiles lib.SQLStructure
	switch strings.ToUpper(output) {
	case "EXCEL":
		// Create Excel Object and Sheet
		excelz, err = lib.CreateExcel(conn)
		if err != nil {
			log.Errorf("%s",err)
			os.Exit(1)
		}
		log.Infof("Setup Excel Files")	
	case "MARKDOWN":
		// Create Excel Object and Sheet
		mdfiles, err = lib.CreateMarkdown(conn)
		if err != nil {
			log.Errorf("%s",err)
			os.Exit(1)
		}
		log.Infof("Setup Markdown Files")	
	case "SQL":
		// Create Excel Object and Sheet
		sqlfiles, err = lib.CreateSQL(conn)
		if err != nil {
			log.Errorf("%s",err)
			os.Exit(1)
		}
		log.Infof("Setup SQL Files")	
	}

	log.Infof("Get Schema Count : %d",len(conn.Schema))

	// Get Info
	for k, _ := range conn.Schema {
		log.Infof("%s Table Load.",k)
		tables, err := conn.GetTables(k)
		if err != nil {
			log.Errorf("%s",err)
			os.Exit(1)
		}
		log.Infof("%s Table Count : %d",k,len(tables))

		log.Infof("%s Table Column/Index/Const Load",k)
		for i, t := range tables {
			// Exception View Table
			switch strings.ToUpper(output) {
			case "EXCEL", "MARKDOWN":
				if t.General.TableType == "BASE TABLE" {
					// Get t Column Info
					columns, err := t.GetColumns(conn.Object,k)
					if err != nil {
						log.Errorf("%s - %s : %s",k,t,err)
						continue
					}
		
					// Get t Index Info
					indexes, err := t.GetIndexes(conn.Object,k)
					if err != nil {
						log.Errorf("%s - %s : %s",k,t,err)
						continue
					}
					
					// Get t Constraint Info
					consts, err := t.GetConstraints(conn.Object,k)
					if err != nil {
						log.Errorf("%s - %s : %s",k,t,err)
						continue
					}
		
					// Add Table Struct
					tables[i].Columns 		= columns
					tables[i].Indexes 		= indexes
					tables[i].Constraints	= consts
				} else if t.General.TableType == "VIEW"{
					ViewDDL, err := t.GetViewSQL(conn.Object,k)
					if err != nil {
						log.Errorf("%s - %s : %s",k,t,err)
						continue
					}
	
					tables[i].View = ViewDDL
				}
			case "SQL":
				tableDDL, err := t.GetTableDDL(conn.Object,k)
				if err != nil {
					log.Errorf("%s - %s : %s",k,t,err)
					continue
				}

				tables[i].DDL = tableDDL
			}
			
		}
		switch strings.ToUpper(output) {
		case "EXCEL":
			err = excelz.WriteTable(k, tables)
			if err != nil {
				log.Warningf("%s",err)
			}
		case "MARKDOWN":
			err = mdfiles.WriteTable(k, tables)
			if err != nil {
				log.Errorf("%s",err)
				os.Exit(1)
			}
		case "SQL":
			err = sqlfiles.WriteTable(k, tables)
			if err != nil {
				log.Errorf("%s",err)
				os.Exit(1)
			}
		}

		conn.Schema[k] = []lib.PerTable{}
	}



	// File Save Action
	switch strings.ToUpper(output) {
	case "EXCEL":
		// Save Excel
		err = excelz.FileSave(fmt.Sprintf("%s.xlsx",*conn.Endpoint))
		if err != nil {
			log.Errorf("%s",err)
			os.Exit(1)
		}
	}

	log.Infof("Export Complete.")
}





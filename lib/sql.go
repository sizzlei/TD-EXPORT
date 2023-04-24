package lib

import (
	"fmt"
	"os"

)

type SQLStructure struct {
	FilesMap 	map[string]*os.File
}

func CreateSQL(c TdexStructure) (SQLStructure, error) {
	var files SQLStructure
	var err error

	var fileMap map[string]*os.File
	fileMap = make(map[string]*os.File)

	for k, _ := range c.Schema {
		f := fmt.Sprintf("%s(%s).sql",k,*c.Endpoint)
		_ = os.Remove(f)
		fileMap[k], err = os.OpenFile(f,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return files, err
		}	
	}

	files.FilesMap = fileMap

	return files, nil
}

func (q SQLStructure) WriteTable(s string, p []PerTable) error {
	z := q.FilesMap[s]
	// Schema
	Writemd(z,fmt.Sprintf("/* Database : %s */\n",s))

	for _, t := range p {
		// Table 
		Writemd(z,fmt.Sprintf("/* Table : %s */\n",t.TableName))
		Writemd(z,fmt.Sprintf("DROP TABLE IF EXISTS %s;\n",	t.TableName))
		Writemd(z,fmt.Sprintf("%s;\n\n\n",t.DDL))		
	}

	return nil
}
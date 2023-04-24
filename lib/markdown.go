package lib

import (
	"os"
	"fmt"
	"strings"
)

type MdStructure struct {
	FilesMap 	map[string]*os.File
}

func CreateMarkdown(c TdexStructure) (MdStructure,error) {
	var files MdStructure
	var err error

	var fileMap map[string]*os.File
	fileMap = make(map[string]*os.File)

	for k, _ := range c.Schema {
		f := fmt.Sprintf("%s(%s).md",k,*c.Endpoint)
		_ = os.Remove(f)
		fileMap[k], err = os.OpenFile(f,os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return files, err
		}	
	}

	files.FilesMap = fileMap

	return files, nil
}

func (m MdStructure) WriteTable(s string, p []PerTable) error {
	z := m.FilesMap[s]

	// Schema
	Writemd(z,fmt.Sprintf("%s \n",s))
	Writemd(z,"=============\n\n")
	
	// List of Content
	Writemd(z,"## Table List\n")
	for _, t := range p {
		Writemd(z,fmt.Sprintf("- [%s (%s)](#%s)\n ",t.TableName,t.General.Comment.String,strings.ToLower(t.TableName)))
	}
	Writemd(z,"\n")

	for _, t := range p {
		// Table Title
		Writemd(z,fmt.Sprintf("## %s\n",strings.ToLower(t.TableName)))
		Writemd(z,"**Information**\n")

		if t.General.TableType == "BASE TABLE" {
			// Table Info
			Writemd(z,"|Table type|Engine|Row format|Collate|Comment|\n")
			Writemd(z,"|---|---|---|---|---|\n")
			Writemd(z,fmt.Sprintf("|%s|%s|%s|%s|%s|\n\n",
				t.General.TableType,
				t.General.Engine.String,
				t.General.RowFormat.String,
				t.General.Collate.String,
				t.General.Comment.String,
			))

			// Column
			Writemd(z,"**Columns**\n")
			Writemd(z,"|Name|Type|Nullable|Default|Charset|Collation|Key|Extra|Comment|\n")
			Writemd(z,"|---|---|---|---|---|---|---|---|---|\n")
			for _, c := range t.Columns {
				Writemd(z,fmt.Sprintf(
					"|%s|%s|%s|%s|%s|%s|%s|%s|%s|\n",
					c.ColumnName,
					c.ColumnType,
					c.Nullable,
					c.DefaultValue.String,
					c.Charset.String,
					c.Collation.String,
					c.ColumnKey.String,
					c.Extra.String,
					c.Comment.String,
				))
			}
			Writemd(z,"\n")

			// Index
			if len(t.Indexes) > 0 {
				Writemd(z,"**Index**\n")
				for _, i := range t.Indexes {
					if i.IndexUnique == 1 {
						Writemd(z,fmt.Sprintf("- [Normal]%s(%s)\n",i.IndexName,i.IndexCol)) 
					} else {
						Writemd(z,fmt.Sprintf("- [Unique]%s(%s)\n",i.IndexName,i.IndexCol)) 
					}
				}
				Writemd(z,"\n")
			}
			
			// Write Constraint
			if len(t.Constraints) > 0 {
				Writemd(z,"**Constraint**\n")
				for _, o := range t.Constraints {
					Writemd(z,fmt.Sprintf(
						"- %s FOREIGN KEY (%s) REFERNCES %s ON DELETE %s ON UPDATE %s\n",
						o.ConstName,
						o.ConstColumn,
						o.ConstRefer,
						o.ConstDelAction,
						o.ConstUpdAction,

					))
				}
				Writemd(z,"\n")
			}
		} else if t.General.TableType == "VIEW" {
			Writemd(z,"|Table type|Charset|Collate|\n")
			Writemd(z,"|---|---|---|\n")
			Writemd(z,fmt.Sprintf("|%s|%s|%s|\n\n",
				t.General.TableType,
				t.View.Charset,
				t.View.Collate,
			))
			Writemd(z,"**View Create SQL**\n")
			Writemd(z,fmt.Sprintf("\n```%s```\n",t.View.ViewQuery))
			
		}
		Writemd(z," \n")
		
		
	}

	return nil

}

func Writemd(f *os.File, d string) {
	_, _ = f.Write([]byte(d))
}

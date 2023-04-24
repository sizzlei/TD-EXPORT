package lib

import (
	"github.com/xuri/excelize/v2"
	"fmt"
	"strconv"
)

type ExcelStructure struct {
	Object 	*excelize.File
	Sheets	map[string]int
	Style 	map[string]int
}

func CreateExcel(c TdexStructure) (ExcelStructure,error) {
	var excelz ExcelStructure
	var err error

	excelz.Object = excelize.NewFile()

	var sheetMap map[string]int
	sheetMap = make(map[string]int)
	for k, _ := range c.Schema {
		sheetMap[k], err = excelz.Object.NewSheet(k)
		if err != nil {
			return excelz, err
		}
	}

	excelz.Sheets = sheetMap
	excelz.Object.DeleteSheet("Sheet1")

	return excelz, nil
}

func (e ExcelStructure) FileSave(f string) error {
	err := e.Object.SaveAs(f)
	if err != nil {
		return err
	}
	return nil
}

func (e ExcelStructure) WriteTable(s string, p []PerTable) error {
	var err error
	// Sheet Active
	e.Object.SetActiveSheet(e.Sheets[s])
	e.Style, err = e.SetupStyle()
	if err != nil {
		return err
	}

	// Declare
	z := e.Object
	RowNum := 1
	for _, t := range p  {

		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["start"])
		RowNum++
		// Table Name
		z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("C%d",RowNum),fmt.Sprintf("J%d",RowNum))

		// Style
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum),e.Style["title"])

		z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Table name")
		z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),t.TableName)
		RowNum++

		// Comment
		z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("C%d",RowNum),fmt.Sprintf("J%d",RowNum))

		// Style
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum),e.Style["title"])

		z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Description")
		z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),t.General.Comment.String)
		RowNum++

		// Columns
		// Columns Title
		z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum))
		// Style
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
		z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Column Information")
		RowNum++

		if t.General.TableType == "BASE TABLE" {
			z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
			z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"No")
			z.SetCellStr(s, fmt.Sprintf("B%d",RowNum),"Column")
			z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),"Data Type")
			z.SetCellStr(s, fmt.Sprintf("D%d",RowNum),"Nullable")
			z.SetCellStr(s, fmt.Sprintf("E%d",RowNum),"Key")
			z.SetCellStr(s, fmt.Sprintf("F%d",RowNum),"Extra")
			z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),"Collate")
			z.SetCellStr(s, fmt.Sprintf("H%d",RowNum),"Default")
			z.MergeCell(s, fmt.Sprintf("I%d",RowNum),fmt.Sprintf("J%d",RowNum))
			z.SetCellStr(s, fmt.Sprintf("I%d",RowNum),"Comment")
			RowNum++

			for i, c := range t.Columns {
				z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),strconv.Itoa(i))
				z.SetCellStr(s, fmt.Sprintf("B%d",RowNum),c.ColumnName)
				z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),c.ColumnType)
				z.SetCellStr(s, fmt.Sprintf("D%d",RowNum),c.Nullable)
				z.SetCellStr(s, fmt.Sprintf("E%d",RowNum),c.ColumnKey.String)
				z.SetCellStr(s, fmt.Sprintf("F%d",RowNum),c.Extra.String)
				z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),c.Collation.String)
				z.SetCellStr(s, fmt.Sprintf("H%d",RowNum),c.DefaultValue.String)
				z.MergeCell(s, fmt.Sprintf("I%d",RowNum),fmt.Sprintf("J%d",RowNum))
				z.SetCellStr(s, fmt.Sprintf("I%d",RowNum),c.Comment.String)
				RowNum++
			}
			
			if len(t.Indexes) > 0 {
				// Indexes 
				// Index Title
				z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum))
				z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Indexes")
				// Style
				z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
				RowNum++

				z.MergeCell(s, fmt.Sprintf("A%d",RowNum), fmt.Sprintf("B%d",RowNum))
				z.MergeCell(s, fmt.Sprintf("C%d",RowNum), fmt.Sprintf("F%d",RowNum))
				z.MergeCell(s, fmt.Sprintf("G%d",RowNum), fmt.Sprintf("J%d",RowNum))
				z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Index Type")
				z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),"Index Name")
				z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),"Columns")
				z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
				RowNum++

				for _, x := range t.Indexes {
					z.MergeCell(s, fmt.Sprintf("A%d",RowNum), fmt.Sprintf("B%d",RowNum))
					z.MergeCell(s, fmt.Sprintf("C%d",RowNum), fmt.Sprintf("F%d",RowNum))
					z.MergeCell(s, fmt.Sprintf("G%d",RowNum), fmt.Sprintf("J%d",RowNum))
					if x.IndexUnique == 1 {
						z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Normal Index")
					} else {
						z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Unique Index")
					}
					z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),x.IndexName)
					z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),x.IndexCol)
					RowNum++
				}
			}
			
			
			if len(t.Constraints) > 0 {
				// Contraints
				// Contraint Title
				z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum))
				z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Constraint")
				// Style
				z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
				RowNum++

				z.MergeCell(s, fmt.Sprintf("A%d",RowNum), fmt.Sprintf("C%d",RowNum))
				z.MergeCell(s, fmt.Sprintf("E%d",RowNum), fmt.Sprintf("H%d",RowNum))
				z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Constraint Name")
				z.SetCellStr(s, fmt.Sprintf("D%d",RowNum),"Column")
				z.SetCellStr(s, fmt.Sprintf("E%d",RowNum),"Referance")
				z.SetCellStr(s, fmt.Sprintf("I%d",RowNum),"ON DELETE")
				z.SetCellStr(s, fmt.Sprintf("J%d",RowNum),"ON UPDATE")
				RowNum++

				for _, o := range t.Constraints {
					z.MergeCell(s, fmt.Sprintf("A%d",RowNum), fmt.Sprintf("C%d",RowNum))
					z.MergeCell(s, fmt.Sprintf("E%d",RowNum), fmt.Sprintf("H%d",RowNum))
					z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),o.ConstName)
					z.SetCellStr(s, fmt.Sprintf("D%d",RowNum),o.ConstColumn)
					z.SetCellStr(s, fmt.Sprintf("E%d",RowNum),o.ConstRefer)
					z.SetCellStr(s, fmt.Sprintf("I%d",RowNum),o.ConstDelAction)
					z.SetCellStr(s, fmt.Sprintf("J%d",RowNum),o.ConstUpdAction)
					RowNum++
				}
			}
		} else if t.General.TableType == "VIEW" {
			z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum))
			z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"View Create SQL")
			z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
			RowNum++

			z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum))
			z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),t.View.ViewQuery)
			RowNum++
		}
		
		z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum))
		z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Table Information")
		// Style
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["title"])
		RowNum++

		// Table Info
		z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("C%d",RowNum),fmt.Sprintf("D%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("E%d",RowNum),fmt.Sprintf("F%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("G%d",RowNum),fmt.Sprintf("J%d",RowNum))

		z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Engine")
		z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),t.General.Engine.String)
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum),e.Style["title"])

		z.SetCellStr(s, fmt.Sprintf("E%d",RowNum),"Row Format")
		z.SetCellStyle(s, fmt.Sprintf("E%d",RowNum),fmt.Sprintf("F%d",RowNum),e.Style["title"])
		z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),t.General.RowFormat.String)
		RowNum++

		z.MergeCell(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("C%d",RowNum),fmt.Sprintf("D%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("E%d",RowNum),fmt.Sprintf("F%d",RowNum))
		z.MergeCell(s, fmt.Sprintf("G%d",RowNum),fmt.Sprintf("J%d",RowNum))

		z.SetCellStr(s, fmt.Sprintf("A%d",RowNum),"Table Type")
		z.SetCellStr(s, fmt.Sprintf("C%d",RowNum),t.General.TableType)
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("B%d",RowNum),e.Style["title"])
		z.SetCellStr(s, fmt.Sprintf("E%d",RowNum),"Collation")
		if t.General.TableType == "BASE TABLE" {
			z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),t.General.Collate.String)
		} else if t.General.TableType == "VIEW" {
			z.SetCellStr(s, fmt.Sprintf("G%d",RowNum),t.View.Collate)
		}
		z.SetCellStyle(s, fmt.Sprintf("E%d",RowNum),fmt.Sprintf("F%d",RowNum),e.Style["title"])
		RowNum++
		z.SetCellStyle(s, fmt.Sprintf("A%d",RowNum),fmt.Sprintf("J%d",RowNum),e.Style["end"])
		RowNum++
		RowNum++

	}
	
	return nil
}

func (e ExcelStructure) SetupStyle() (map[string]int,error) {
	var Style map[string]int
	var err error
	Style = make(map[string]int)

	Style["title"], err = e.Object.NewStyle(
		&excelize.Style{
			Border: []excelize.Border{
				{Type: "left",Color: "#000000",Style: 1},
				{Type: "right",Color: "#000000",Style: 1},
				{Type: "top",Color: "#000000",Style: 1},
				{Type: "bottom",Color: "#000000",Style: 1},
			},
			Font: &excelize.Font{Bold:true,Color:"#ffffff"},
			Fill: excelize.Fill{Type:"pattern",Color:[]string{"#000000"},Pattern:1},
			// Alignment: &excelize.Alignment{Horizontal: "center"},
		},
	)
	if err != nil {
		return Style,err
	}

	Style["start"], err = e.Object.NewStyle(
		&excelize.Style{
			Border: []excelize.Border{
				{Type: "bottom",Color: "#000000",Style: 1},
			},
		},
	)
	if err != nil {
		return Style,err
	}

	Style["end"], err = e.Object.NewStyle(
		&excelize.Style{
			Border: []excelize.Border{
				{Type: "top",Color: "#000000",Style: 1},
			},
		},
	)
	if err != nil {
		return Style,err
	}

	
	return Style, nil
}
package excel

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

//Cell value type Constatnts
const (
	CellTypeBool    = 0
	CellTypeInt     = 1
	CellTypeString  = 2
	CellTypeFloat   = 3
	CellTypeFormula = 4
)

//Style Constants
const (
	StyleNone = -1
)

//ExcelizeWriter ...
type ExcelizeWriter interface {
	ExcelizeWrite(x *excelize.File, sheetName string, axis string) (string, error)
	ExcelizeClear(x *excelize.File, sheetName string, axis string) error
	Range(axis string) (string, error)
}

//Cell ...
type Cell struct {
	val     interface{}
	typ     int
	styleID int
	colRef  *Column
}

//Range ...
func (c Cell) Range(axis string) string {
	if c.val == nil {
		return ""
	}
	return axis + ":" + axis
}

//ExcelizeClear ...
func (c Cell) ExcelizeClear(x *excelize.File, sheetName string, axis string) error {
	return x.SetCellValue(sheetName, axis, "")

}

//ExcelizeWrite ...
func (c Cell) ExcelizeWrite(x *excelize.File, sheetName string, axis string) (string, error) {
	///TODO, handle specifying prcision and bitsize
	///TODO, impliment hyperling cell value and Default string cell value

	if c.val != nil {

		var err error
		switch c.typ {
		case CellTypeString:
			err = x.SetCellStr(sheetName, axis, c.val.(string))
		case CellTypeBool:
			err = x.SetCellBool(sheetName, axis, c.val.(bool))
		case CellTypeFloat:
			err = x.SetCellFloat(sheetName, axis, c.val.(float64), 2, 64)
		case CellTypeFormula:
			err = x.SetCellFormula(sheetName, axis, c.val.(string))
		case CellTypeInt:
			err = x.SetCellInt(sheetName, axis, c.val.(int))
		default:
			err = x.SetCellValue(sheetName, axis, c.val)
		}
		if err != nil {
			return "", err
		}
	}

	style := c.getEffectiveStyleID()

	if style != StyleNone {
		err := x.SetCellStyle(sheetName, axis, axis, style)
		if err != nil {
			return "", err
		}
	}

	return axis + ":" + axis, nil

}
func (c Cell) getEffectiveStyleID() int {
	if c.styleID != StyleNone {
		return c.styleID
	}
	if c.colRef != nil {
		return c.colRef.getEffectiveStyleID()
	}
	return StyleNone
}

//NewCell ...
func NewCell(typ int, value interface{}, styleID int) (Cell, error) {
	c := Cell{
		typ:     typ,
		val:     value,
		styleID: styleID,
	}
	if value == nil {
		return c, nil
	}

	switch typ {
	case CellTypeBool:
		if _, ok := value.(bool); !ok {
			return Cell{}, fmt.Errorf("value: %v is not of type bool", value)
		}
	case CellTypeFloat:
		switch v := value.(type) {
		case float32:
			c.val = float64(v)
		case float64:
		default:
			return Cell{}, fmt.Errorf("value: %v is not floating point type", value)
		}
	case CellTypeInt:
		switch v := value.(type) {
		case int:
		case int8:
			c.val = int(v)
		case int16:
			c.val = int(v)
		case int32:
			c.val = int(v)
		case int64:
			c.val = int(v)
		default:
			return Cell{}, fmt.Errorf("value: %v is not integer type", value)

		}
	case CellTypeFormula, CellTypeString:
		if _, ok := value.(string); !ok {
			return Cell{}, fmt.Errorf("value: %v is not a string", value)
		}
	}
	return c, nil

}

//Column ...
type Column struct {
	hdr        string
	vals       map[int]Cell
	styleID    int
	colStyleID int
	colTabRef  *ColTable
}

func (col Column) getEffectiveStyleID() int {
	if col.styleID != StyleNone {
		return col.styleID
	}
	if col.colTabRef != nil {
		return col.colTabRef.getEffectiveStyleID()
	}
	return StyleNone
}

//Range ...
func (col Column) Range(axis string) (string, error) {
	if len(col.vals) == 0 {
		if col.hdr == "" {
			return "", nil
		}
		return axis + ":" + axis, nil

	}

	c, r, err := excelize.CellNameToCoordinates(axis)

	if err != nil {
		return "", err
	}
	maxRowOffset := 0
	for offset := range col.vals {
		maxRowOffset = max(maxRowOffset, offset)
	}

	bottomRowIndex := maxRowOffset + r

	bottomCellName, err := excelize.CoordinatesToCellName(c, bottomRowIndex)
	if err != nil {
		return "", err
	}

	return axis + ":" + bottomCellName, nil
}

//ExcelizeClear ...
func (col Column) ExcelizeClear(x *excelize.File, sheetName string, axis string) error {
	for _, c := range col.vals {
		err := c.ExcelizeClear(x, sheetName, axis)
		if err != nil {
			return err
		}
	}
	return nil
}

//ExcelizeWrite ..
func (col Column) ExcelizeWrite(x *excelize.File, sheetName string, axis string) (string, error) {
	//log.Printf("writing column value %v to axis %s", col.vals, axis)
	colRange, err := col.Range(axis)
	if err != nil {
		return "", err
	}

	if colRange == ":" {
		return colRange, nil
	}

	if col.hdr != "" {
		//log.Printf("writing column heade value %s", col.hdr)
		err := x.SetCellStr(sheetName, axis, col.hdr)
		if err != nil {
			return "", err
		}
	}

	axisColName, axisRowNumber, err := excelize.SplitCellName(axis)
	if err != nil {
		return "", err
	}

	if col.colStyleID != StyleNone {
		err = x.SetColStyle(sheetName, axisColName, col.colStyleID)
		if err != nil {
			return "", err
		}
	}

	if len(col.vals) != 0 {
		for offset, cell := range col.vals {
			cellAxis, err := excelize.JoinCellName(axisColName, axisRowNumber+offset)
			if err != nil {
				return "", err
			}
			// if col.styleId != StyleNone {
			// 	err = x.SetCellStyle(sheetName, cellAxis, cellAxis, col.styleId)
			// 	if err != nil {
			// 		return "", err
			// 	}
			// }
			//log.Printf("writing cell value %v at offset %d to axis %s", cell.val, offset, cellAxis)
			_, err = cell.ExcelizeWrite(x, sheetName, cellAxis)
			if err != nil {
				return "", err
			}
		}
	}

	return colRange, nil

}

//NewColumn ...
func NewColumn(header string, cells map[int]Cell, style, colStyle int) (Column, error) {
	val := cells
	if cells == nil {
		val = make(map[int]Cell)
	}

	for k := range val {
		if k < 0 {
			return Column{}, fmt.Errorf("One of the Cells has negative offset %d", k)
		}

		if header != "" && k == 0 {
			return Column{}, fmt.Errorf("One of the Cells has offset 0 which is also the offset for the header")
		}
	}

	return Column{
		hdr:        header,
		vals:       val,
		styleID:    style,
		colStyleID: colStyle,
	}, nil
}

//NewColumnFromSlice ...
func NewColumnFromSlice(header string, cells []Cell, style, colStyle int) Column {

	cellsMap := make(map[int]Cell)
	i := 0

	if header != "" {
		i++
	}

	for _, v := range cells {
		cellsMap[i] = v
		i++
	}
	c, _ := NewColumn(header, cellsMap, style, colStyle)
	return c
}

//SetHeader ...
func (col *Column) SetHeader(header string) error {
	if col.hdr != "" {
		col.hdr = header
		return nil
	}
	if _, exist := col.vals[0]; exist {
		return fmt.Errorf("Cannot Set header there is a cell value in offset 0")
	}

	col.hdr = header
	return nil
}

//AddCell ...
func (col *Column) AddCell(c Cell, offset int) error {

	if col.HasCellAtOffset(offset) {
		return fmt.Errorf("Error. A cell already exist at that offcset: %d", offset)
	}
	return col.SetCell(c, offset)
}

//SetCell ...
func (col *Column) SetCell(c Cell, offset int) error {
	if offset < 0 {
		return fmt.Errorf("Error: Cell has negative offset %d", offset)
	}
	c.colRef = col
	col.vals[offset] = c
	return nil
}

//AppendCell ...
func (col *Column) AppendCell(c Cell) int {
	maxOffset := -1 //same reason as AppendCol of ColTable
	for k := range col.vals {
		maxOffset = max(maxOffset, k)
	}
	maxOffset++
	if col.hdr != "" && maxOffset == 0 {
		maxOffset++
	}
	col.SetCell(c, maxOffset)
	return maxOffset
}

//DeleteCell ...
func (col *Column) DeleteCell(offset int) {
	delete(col.vals, offset)
}

//HasCellAtOffset ...
func (col *Column) HasCellAtOffset(offset int) bool {
	_, v := col.vals[offset]
	return v
}

//ColTable ...
type ColTable struct {
	vals        map[int]Column
	style       int ///TODO figure how to propagate this to cells as default
	tableFormat string
}

func (t ColTable) getEffectiveStyleID() int {
	return t.style
}

//ExcelizeWrite ...
func (t *ColTable) ExcelizeWrite(x *excelize.File, sheetName string, axis string) (string, error) {
	c, r, err := excelize.CellNameToCoordinates(axis)
	if err != nil {
		return "", err
	}
	tableRange, err := t.Range(axis)
	if err != nil {
		return "", err
	}
	if tableRange == ":" {
		return ":", nil
	}

	topleft, bottomRight := R2c(tableRange)

	for colOffset, col := range t.vals {
		colAxis, err := excelize.CoordinatesToCellName(c+colOffset, r)
		if err != nil {
			return "", err
		}
		_, err = col.ExcelizeWrite(x, sheetName, colAxis)
		if err != nil {
			return "", err
		}

	}
	if t.tableFormat != "" {
		err = x.AddTable(sheetName, topleft, bottomRight, t.tableFormat)
		if err != nil {
			return "", err
		}
	}
	return t.Range(axis)
}

//ExcelizeClear ...
func (t *ColTable) ExcelizeClear(x *excelize.File, sheetName string, axis string) error {
	c, r, err := excelize.CellNameToCoordinates(axis)
	if err != nil {
		return err
	}
	for colOffset, col := range t.vals {
		colAxis, err := excelize.CoordinatesToCellName(c+colOffset, r)
		if err != nil {
			return err
		}
		err = col.ExcelizeClear(x, sheetName, colAxis)
		if err != nil {
			return err
		}
	}
	return nil
}

//Range ...
func (t *ColTable) Range(axis string) (string, error) {
	if len(t.vals) == 0 {
		return "", nil
	}
	c, r, err := excelize.CellNameToCoordinates(axis)
	if err != nil {
		return "", err
	}
	maxColOffset := 0
	maxRowIndex := 0

	for colOffset, col := range t.vals {
		maxColOffset = max(maxColOffset, colOffset)
		colAxis, err := excelize.CoordinatesToCellName(c+colOffset, r)
		if err != nil {
			return "", err
		}
		colRange, err := col.Range(colAxis)
		if err != nil {
			return "", err
		}
		_, colBottomCelName := R2c(colRange)
		_, rowIndex, err := excelize.SplitCellName(colBottomCelName)
		if err != nil {
			return "", err
		}

		maxRowIndex = max(maxRowIndex, rowIndex)

	}
	maxCellColName, err := excelize.ColumnNumberToName(c + maxColOffset)
	if err != nil {
		return "", err
	}
	bottomRightCellName, err := excelize.JoinCellName(maxCellColName, maxRowIndex)
	if err != nil {
		return "", err
	}

	return axis + ":" + bottomRightCellName, nil
}

//NewColTable ...
func NewColTable(cols map[int]Column) (*ColTable, error) {

	t := ColTable{
		vals: cols,
	}

	if cols == nil {
		t.vals = make(map[int]Column)
	}

	//var maxColOffset int
	for offset := range t.vals {
		if offset < 0 {
			return &ColTable{}, fmt.Errorf("Error map keys (Column offset) cannot be negative: found %d", offset)
		}
		//maxColOffset = max(maxColOffset, offset)
	}

	return &t, nil

}

//AddCol ...
func (t *ColTable) AddCol(col Column, offset int) error {
	if _, exist := t.vals[offset]; exist {
		return fmt.Errorf("Cannot add column at offset %d. a column already exists at that offset", offset)
	}
	t.SetCol(col, offset)
	return nil
}

//SetCol ...
func (t *ColTable) SetCol(col Column, offset int) {
	t.vals[offset] = col
	col.colTabRef = t
}

//AppendCol ...
func (t *ColTable) AppendCol(col Column) int {
	maxKey := -1 //to take care of empty map. else nothing will have offset 0
	for k := range t.vals {
		maxKey = max(maxKey, k)
	}
	maxKey++
	t.SetCol(col, maxKey)
	return maxKey
}

//AddCell ...
func (t *ColTable) AddCell(c Cell, colOffset, rowOffset int, create bool) error {
	col, exist := t.vals[colOffset]
	if !exist {
		if !create {
			return fmt.Errorf("Cannot add Cell at offset %d BY %d. no column exists at offset %d to add the cell", colOffset, rowOffset, colOffset)
		}
		err := t.addEmptyColumn("Col "+strconv.Itoa(123), colOffset)
		if err != nil {
			return err
		}
		col = t.vals[colOffset]
	}

	return col.AddCell(c, rowOffset)
}

//SetCell ...
func (t *ColTable) SetCell(c Cell, colOffset, rowOffset int, create bool) error {
	col, exist := t.vals[colOffset]
	if !exist {
		if !create {
			return fmt.Errorf("Cannot set Cell at offset %d BY %d. no column exists at offset %d to add the cell", colOffset, rowOffset, colOffset)
		}
		err := t.addEmptyColumn("Col "+strconv.Itoa(123), colOffset)
		if err != nil {
			return err
		}
		col = t.vals[colOffset]

	}

	return col.SetCell(c, rowOffset)
}

//AppendCell ...
func (t *ColTable) AppendCell(c Cell, colOffset int) int {
	col, exist := t.vals[colOffset]
	if !exist {
		t.addEmptyColumn("Col "+strconv.Itoa(123), colOffset)
		col = t.vals[colOffset]
	}
	return col.AppendCell(c)
}

//GetCol ...
func (t *ColTable) GetCol(offset int) (Column, error) {
	if col, exist := t.vals[offset]; exist {
		return col, nil
	}
	return Column{}, fmt.Errorf("Error: No column at offset %d in table", offset)

}

//SetColHdr ...
func (t *ColTable) SetColHdr(hdr string, offset int) error {
	col, exist := t.vals[offset]
	if !exist {
		return fmt.Errorf("Cannot set Column header at offset %d no column exists at offset", offset)
	}
	return col.SetHeader(hdr)
}

//SetTableFormat ...
func (t *ColTable) SetTableFormat(format string) {
	t.tableFormat = format
}

//GetColOffset ...
func (t *ColTable) GetColOffset(hdr string) (int, error) {

	for k, v := range t.vals {
		if strings.EqualFold(hdr, v.hdr) {
			return k, nil
		}
	}
	return -1, fmt.Errorf("Error: No column with header '%s' in table", hdr)

}

func (t *ColTable) addEmptyColumn(header string, offset int) error {
	//colName := "Col " + strconv.Itoa(123)

	col, err := NewColumn(header, nil, StyleNone, StyleNone)
	if err != nil {
		return err
	}

	err = t.AddCol(col, offset)
	if err != nil {
		return err
	}
	return nil
}

func (t *ColTable) print() {
	for c := 0; c <= 4; c++ {
		if c >= len(t.vals) {
			break
		}
		buff := t.vals[c]
		for r := 0; r <= 4; r++ {
			if r >= len(buff.vals) {
				break
			}
			fmt.Printf("C%d:R%d - %v\t\t", c, r, buff.vals[r].val)
		}
		fmt.Println()
	}
}

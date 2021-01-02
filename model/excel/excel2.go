package excel

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

//FitColWidth ...
func FitColWidth(x *excelize.File, sheetName string, area string, factor float64) error {

	topLeft, bottomRight := R2c(area)
	c1, r1, err := excelize.SplitCellName(topLeft)
	if err != nil {
		return err
	}

	c2, r2, err := excelize.SplitCellName(bottomRight)
	if err != nil {
		return err
	}

	c1n, err := excelize.ColumnNameToNumber(c1)
	if err != nil {
		return err
	}

	c2n, err := excelize.ColumnNameToNumber(c2)

	for i := c1n; i <= c2n; i++ {
		maxChar := 0
		for j := r1; j <= r2; j++ {
			axis, err := excelize.CoordinatesToCellName(i, j)
			if err != nil {
				return err
			}
			v, err := x.GetCellValue(sheetName, axis)
			if err != nil {
				return err
			}
			maxChar = max(maxChar, len(v))
		}
		c, err := excelize.ColumnNumberToName(i)
		if err != nil {
			return err
		}
		//log.Printf("Col: %s, Max Char: %d", c, maxChar)
		err = x.SetColWidth(sheetName, c, c, factor*float64(maxChar))
		if err != nil {
			return err
		}

	}

	return nil

}

func FitSheetName(sheetName string) string {
	actualSheetName := sheetName
	if len(sheetName) > 31 {

		actualSheetName = sheetName[0:28] + "..."
	}
	return actualSheetName
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Dollarize(axis string, col bool, row bool) (string, error) {
	re := regexp.MustCompile(`^\$?([A-Za-z]+)\$?([0-9]+)$`)
	res := re.FindStringSubmatch(axis)
	if len(res) != 3 {
		return "", fmt.Errorf("Invalid Cell Label: %s", axis)
	}
	c := res[1]
	r := res[2]
	if col {
		c = "$" + c
	}
	if row {
		r = "$" + r
	}
	return c + r, nil
}

func OffsetJump(axis string, cOffset, rOffset int) (string, error) {
	c, r, err := excelize.CellNameToCoordinates(axis)
	if err != nil {
		return "", err
	}

	return excelize.CoordinatesToCellName(c+cOffset, r+rOffset)

}

//R2c ...
func R2c(r string) (string, string) {
	if r == "" {
		return "", ""
	}
	comps := strings.Split(r, ":")
	return comps[0], comps[1]

}

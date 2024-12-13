package table

import (
	"fmt"
	"strings"
	"strconv"
)

// Define the size of the table
const TableSize = 21

// ChartRow represents a single row in the chart
type ChartRow struct {
	Key   string
	Value [13]string
}

// Chart represents the main chart structure
type Chart struct {
	Name    string
	Rows    [TableSize]ChartRow
	NextRow int
}

// NewChart initializes a new chart with the given name
func NewChart(name string) *Chart {
	chart := &Chart{Name: name}
	for i := range chart.Rows {
		chart.Rows[i].Key = "--"
		for j := range chart.Rows[i].Value {
			chart.Rows[i].Value[j] = "---"
		}
	}
	return chart
}

//
func (c *Chart) GetRowCount() int {
	return c.NextRow
}

// Insert adds a key-value pair to the chart
func (c *Chart) Insert(key string, up int, value string) {
	index := c.getRow(key)
	if index < 0 {
		index = c.NextRow
		c.NextRow++
		c.Rows[index].Key = strings.ToUpper(key)
	}
	c.Rows[index].Value[up] = strings.ToUpper(value)
}

// GetValue retrieves a value from the chart
func (c *Chart) GetValue(key string, up int) string {
	index := c.getRow(key)
	if index < 0 {
		fmt.Printf("Cannot find value in %s for %s vs %d\n", c.Name, key, up)
		panic("Key not found")
	}
	return c.Rows[index].Value[up]
}

//
func (c *Chart) GetValueByTotal(total, up int) string {
	return c.GetValue(strconv.Itoa(total), up);
}

// Print prints the entire chart to the console
func (c *Chart) Print() {
	fmt.Println(c.Name)
	fmt.Println("--------2-----3-----4-----5-----6-----7-----8-----9-----T-----J-----Q-----K-----A---")
	for i := 0; i < c.NextRow; i++ {
		row := c.Rows[i]
		fmt.Printf("%2s : ", row.Key)
		for _, value := range row.Value {
			fmt.Printf("%4s, ", value)
		}
		fmt.Println()
	}
	fmt.Println("------------------------------------------------------------------------------------")
}

// getRow finds the index of the row for the given key
func (c *Chart) getRow(key string) int {
	keyUpper := strings.ToUpper(key)
	for i := 0; i < c.NextRow; i++ {
		if c.Rows[i].Key == keyUpper {
			return i
		}
	}
	return -1
}


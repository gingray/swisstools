package ui

import (
	"os"
	"os/exec"

	"github.com/gingray/swisstools/pkg/common"
	"github.com/olekukonko/tablewriter"
)

type TableView struct {
}

func NewTableView() *TableView {
	return &TableView{}
}

func (t *TableView) Show(view *common.DataView) error {
	table := tablewriter.NewTable(os.Stdout)
	defer table.Close()
	table.Header(view.Keys)
	var items [][]string
	for _, row := range view.Rows {
		var singleRow []string
		for _, value := range view.Keys {
			singleRow = append(singleRow, row[value])
		}
		items = append(items, singleRow)

	}
	table.Bulk(items)
	table.Footer(view.Keys)
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	err := table.Render()
	return err

}

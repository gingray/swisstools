package ui

import (
	"github.com/gingray/swisstools/pkg/common"
	"github.com/olekukonko/tablewriter"
	"os"
	"os/exec"
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
		for _, value := range view.Keys {
			items = append(items, []string{row[value]})
		}
	}
	table.Bulk(items)
	table.Footer(view.Keys)
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	err := table.Render()
	return err

}

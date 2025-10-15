package ui

import (
	"github.com/olekukonko/tablewriter"
	"os"
	"os/exec"
)

type TableView struct {
}

func (t *TableView) Show(rows [][]string) error {
	table := tablewriter.NewTable(os.Stdout)
	defer table.Close()
	table.Header([]string{"Url", "Title", "Status", "Created"})
	table.Bulk(rows)
	table.Footer([]string{"Url", "Title", "Status", "Created"})
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	err := table.Render()
	return err
}

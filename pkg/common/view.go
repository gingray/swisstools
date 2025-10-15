package common

type ViewRecords interface {
	Show(rows [][]string) error
}

type DataView struct {
	Keys []string
	Rows []map[string]string
}

func (d *DataView) NewDataView() *DataView {
	return &DataView{Rows: make([]map[string]string, 0), Keys: make([]string, 0)}
}

func (d *DataView) AddKey(key string) {
	d.Keys = append(d.Keys, key)
}

func (d *DataView) AddRow(row map[string]string) {
	d.Rows = append(d.Rows, row)
}

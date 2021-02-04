package commands

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"sort"
	"strings"

	"github.com/uhppoted/uhppoted-api/config"
)

type Doors []string

func getDoors(conf *config.Config) (Doors, error) {
	type door struct {
		door  string
		index uint32
	}

	displayOrder := strings.Split(conf.WildApricot.DisplayOrder.Doors, ",")

	set := map[string]door{}
	for _, device := range conf.Devices {
		for _, d := range device.Doors {
			normalised := normalise(d)
			if _, ok := set[normalised]; ok {
				return nil, fmt.Errorf("WARN  Duplicate door in configuration (%v)", d)
			}

			index := uint32(math.MaxUint32)
			for i := range displayOrder {
				name := normalise(displayOrder[i])
				if normalised == name {
					index = uint32(i + 1)
					break
				}
			}

			set[normalised] = door{
				door:  clean(d),
				index: index,
			}
		}
	}

	doors := []door{}
	for _, d := range set {
		doors = append(doors, d)
	}

	sort.SliceStable(doors, func(i, j int) bool { return doors[i].door < doors[j].door })
	sort.SliceStable(doors, func(i, j int) bool { return doors[i].index < doors[j].index })

	list := []string{}
	for _, d := range doors {
		list = append(list, d.door)
	}

	return list, nil
}

func (doors *Doors) MarshalText() ([]byte, error) {
	return doors.MarshalTextIndent("")
}

func (doors *Doors) MarshalTextIndent(indent string) ([]byte, error) {
	header, data := doors.asTable()
	table := [][]string{}

	table = append(table, header)
	table = append(table, data...)

	var b bytes.Buffer

	if len(table) > 0 {
		widths := make([]int, len(table[0]))
		for _, row := range table {
			for i, field := range row {
				if len(field) > widths[i] {
					widths[i] = len(field)
				}
			}
		}

		for i := 1; i < len(widths); i++ {
			widths[i-1] += 1
		}

		for _, row := range table {
			fmt.Fprintf(&b, "%s", indent)
			for i, field := range row {
				fmt.Fprintf(&b, "%-*v", widths[i], field)
			}
			fmt.Fprintln(&b)
		}
	}

	return b.Bytes(), nil
}

func (doors *Doors) ToTSV(f io.Writer) error {
	header, data := doors.asTable()

	w := csv.NewWriter(f)
	w.Comma = '\t'

	w.Write(header)
	for _, row := range data {
		w.Write(row)
	}

	w.Flush()

	return nil
}

func (doors *Doors) asTable() ([]string, [][]string) {
	header := []string{
		"Door",
	}

	data := [][]string{}

	if doors != nil {
		for _, d := range *doors {
			row := []string{
				fmt.Sprintf("%v", d),
			}

			data = append(data, row)
		}
	}

	return header, data
}

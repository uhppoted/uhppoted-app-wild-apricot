package commands

import (
	"fmt"
	"math"
	"sort"
	"strings"

	api "github.com/uhppoted/uhppoted-lib/acl"
	"github.com/uhppoted/uhppoted-lib/config"
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

			// ignore blank doors - no use to man or beast
			if normalised == "" {
				continue
			}

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

func (doors *Doors) AsTable() *api.Table {
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

	table := api.Table{
		Header:  header,
		Records: data,
	}

	return &table
}

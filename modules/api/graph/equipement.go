package graph

import (
	"io"
	"time"

	"github.com/sapk/GoWatch/modules/rrd"
)

//EquipementPing execute a ping and return informations
func EquipementPing(ID, duration string, out io.Writer) error {
	now := time.Now()
	switch duration {
	case "minute":
		out.Write(rrd.GraphPing(ID, "5 minutes", now.Add(-time.Minute*5), now))
	case "hour":
		out.Write(rrd.GraphPing(ID, "Hourly", now.Add(-time.Hour), now))
	case "day":
		out.Write(rrd.GraphPing(ID, "Daily", now.Add(-time.Hour*24), now))
	case "week":
		out.Write(rrd.GraphPing(ID, "Weekly", now.Add(-time.Hour*24*7), now))
	}
	return nil
}

package graph

import (
	"io"
	"time"

	"github.com/sapk/GoWatch/modules/collector"
)

//EquipementPing send a image to the out buffer
func EquipementPing(ID, duration string, out io.Writer) error {
	now := time.Now()
	switch duration {
	case "minute":
		out.Write(collector.GraphPing(ID, "5 minutes", now.Add(-time.Minute*5), now))
	case "hour":
		out.Write(collector.GraphPing(ID, "Hourly", now.Add(-time.Hour), now))
	case "day":
		out.Write(collector.GraphPing(ID, "Daily", now.Add(-time.Hour*24), now))
	case "week":
		out.Write(collector.GraphPing(ID, "Weekly", now.Add(-time.Hour*24*7), now))
	}
	return nil
}

//EquipementPingData send ping in database
func EquipementPingData(ID, duration string) []collector.Ping {
	now := time.Now()
	switch duration {
	case "minute":
		return collector.DataPing(ID, now.Add(-time.Minute*5), now)
	case "hour":
		return collector.DataPing(ID, now.Add(-time.Hour), now)
	case "day":
		return collector.DataPing(ID, now.Add(-time.Hour*24), now)
	case "week":
		return collector.DataPing(ID, now.Add(-time.Hour*24*7), now)
	}
	return []collector.Ping{}
}

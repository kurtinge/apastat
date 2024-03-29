package collector

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type ApacheCollector struct {
	Host string
	Uri  string
}

func NewApacheCollector(host *string) Collector {
	if host == nil {
		*host = "localhost"
	}
	return &ApacheCollector{
		Host: *host,
		Uri:  "/server-status",
	}
}

func (c *ApacheCollector) GetStats() (*ServerStatus, error) {
	resp, err := http.Get("http://" + c.Host + c.Uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parseApacheServerStatus(resp.Body)
}

func (c *ApacheCollector) GetStatsFromFile() (*ServerStatus, error) {
	file, err := os.Open("serverstatus.html")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return c.parseApacheServerStatus(file)
}

func (c *ApacheCollector) parseApacheServerStatus(htmlBody io.Reader) (*ServerStatus, error) {
	tkn := html.NewTokenizer(htmlBody)

	serverStatus := &ServerStatus{}

	var isTd bool
	var isTh bool
	var isDt bool

	var columns []string
	var currentColumn int
	foundServerStatus := false

	var serverSlot Slot

	for {
		tt := tkn.Next()

		switch {

		case tt == html.ErrorToken:
			return serverStatus, nil

		case tt == html.StartTagToken:

			t := tkn.Token()
			if t.Data == "b" {
				continue
			}

			isTd = t.Data == "td"
			isTh = t.Data == "th"
			isDt = t.Data == "dt"

			if t.Data == "tr" {
				serverSlot = Slot{}
				currentColumn = 0
			}

			if t.Data == "table" {
				if t.Attr[0].Key == "border" {
					foundServerStatus = true
				}
			}

		case tt == html.EndTagToken:
			t := tkn.Token()
			if foundServerStatus && t.Data == "tr" {

				if serverSlot.Mode != "" {
					serverStatus.ServerSlots = append(serverStatus.ServerSlots, serverSlot)
				}
			}

			if foundServerStatus && t.Data == "td" {
				currentColumn++
			}

			if foundServerStatus && t.Data == "table" {
				return serverStatus, nil
			}
			isTd = false
			isTh = false
			isDt = false

		case tt == html.TextToken:

			t := tkn.Token()

			if foundServerStatus && isTh {
				columns = append(columns, strings.TrimSpace(t.Data))
			}

			if foundServerStatus && isTd {
				switch columns[currentColumn] {
				case "Srv":
					serverSlot.ServerSlot = strings.TrimSpace(serverSlot.ServerSlot + t.Data)
				case "PID":
					serverSlot.Pid, _ = strconv.Atoi(t.Data)
				case "M":
					mode := strings.TrimSpace(t.Data)
					if mode != "" {
						serverSlot.Mode = ServerMode(mode)
					}
				case "CPU":
					serverSlot.Cpu, _ = strconv.ParseFloat(t.Data, 64)
				case "SS":
					serverSlot.SecondsSinceRequest, _ = strconv.Atoi(t.Data)
				case "Client":
					serverSlot.Client = t.Data
				case "VHost":
					serverSlot.Vhost = strings.TrimSpace(fmt.Sprintf("%s %s", serverSlot.Vhost, t.Data))
				case "Protocol":
					serverSlot.Protocol = t.Data
				case "Request":
					serverSlot.Request = t.Data
				}
			}

			if isDt {
				infoRow := strings.Split(t.Data, ":")
				if len(infoRow) == 2 {
					if infoRow[0] == "Server uptime" {
						serverStatus.Uptime = infoRow[1]
					}
				} else {
					infoRow := strings.Split(t.Data, "-")
					firstColumns := strings.Split(infoRow[0], " ")
					if len(firstColumns) > 1 {
						if firstColumns[1] == "requests/sec" {
							rqs, err := strconv.Atoi(infoRow[1])
							if err == nil {
								serverStatus.RequestSec = rqs
							}
						}
					}
				}

			}
		}
	}
}

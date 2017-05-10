// Parts of this file is totally copied from https://github.com/iopred/bruxism/blob/master/statsplugin/statsplugin.go
// This applies: https://github.com/iopred/bruxism/blob/master/LICENSE
// Thanks to Mister Christopher Rhodes for this! :D

package bot_reactions

import (
	"bytes"
	"fmt"
	"runtime"
	"text/tabwriter"
	"time"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
)

var statsStartTime = time.Now()

func getDurationString(duration time.Duration) string {
	return fmt.Sprintf(
		"%0.2d:%02d:%02d",
		int(duration.Hours()),
		int(duration.Minutes())%60,
		int(duration.Seconds())%60,
	)
}

type stats struct {
	Trigger string
}

func (s *stats) Help() string {
	return "I have statistics!"
}

func (s *stats) HelpDetail(m *discordgo.Message) string {
	return s.Help()
}

func (s *stats) Reaction(m *discordgo.Message, a *discordgo.Member, update bool) string {
	stats := runtime.MemStats{}
	runtime.ReadMemStats(&stats)

	w := &tabwriter.Writer{}
	buf := &bytes.Buffer{}

	hostname,err := os.Hostname()
	pwd,err := os.Getwd()

	var execmode = "Named pipe"

	w.Init(buf, 0, 4, 0, ' ', 0)

	fmt.Fprintf(w, "```\n")
	
	if err == nil {
		fmt.Fprintf(w, "Hostname: \t%s\n", hostname)
	}

        if err == nil {
                fmt.Fprintf(w, "Working directory: \t%s\n", pwd)
        }

	fmt.Fprintf(w, "DiscordGo: \t%s\n", discordgo.VERSION)
	fmt.Fprintf(w, "Go: \t%s\n", runtime.Version())
	// fmt.Fprintf(w, "Ruby version: \t%s\n", os.Getenv("RUBY_VERSION"))
	fmt.Fprintf(w, "Uptime: \t%s\n", getDurationString(time.Now().Sub(statsStartTime)))
	fmt.Fprintf(w, "Memory used: \t%s / %s (%s garbage collected)\n", humanize.Bytes(stats.Alloc), humanize.Bytes(stats.Sys), humanize.Bytes(stats.TotalAlloc))
	fmt.Fprintf(w, "Concurrent tasks: \t%d\n", runtime.NumGoroutine())
	fmt.Fprintf(w, "Execution mode: \t%s\n", execmode)
	fmt.Fprintf(w, "\n```")

	w.Flush()
	out := buf.String()
	return out
}

func init() {
	_stats := &stats{
		Trigger: "stats",
	}
	addReaction(_stats.Trigger, _stats)
}


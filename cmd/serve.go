package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/kekscode/calendar-events-exporter/pkg/calendar"
	exporter "github.com/kekscode/calendar-events-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve starts the exporter serving metrics",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		urls, err := cmd.Flags().GetStringArray("icalendar-urls")
		if err != nil {
			log.Printf("error: icalendar target list of URLs is not valid: %v", err)
		}

		// Trim whitespaces
		for i := range urls {
			urls[i] = strings.TrimSpace(urls[i])
		}

		store, err := calendar.NewEventStore("ical", urls)
		if err != nil {
			log.Printf("Error loading calendar monitor: %v", store)
		}

		// Main loop
		ticker := time.NewTicker(2 * time.Second)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					store.Update()
					for _, e := range store.Events() {
						log.Printf("%v\n", e)
					}
					log.Printf("Tick at", t)
				}

			}
		}()

		prometheus.MustRegister(exporter.NewExporter())
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9310", nil)

		done <- true
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	serveCmd.Flags().StringArrayP("icalendar-urls", "u", []string{"file:///calendar.ics"}, "URL location of the ics file to monitor. This flag may be repeated for different targets.")
}

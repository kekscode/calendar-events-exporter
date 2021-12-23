package cmd

import (
	"net/http"
	"strings"
	"time"

	"github.com/kekscode/calendar-events-exporter/pkg/calendar"
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

		// Create the calendar events store
		store, err := calendar.NewEventStore("ical", urls)
		if err != nil {
			log.Printf("Error loading calendar monitor: %v", store)
		}

		// Main loop
		ticker := time.NewTicker(5 * time.Second)
		done := make(chan bool)
		go func() {
			for {
				select {
				case <-done:
					return
				case t := <-ticker.C:
					// TODO: Add metrics for the number of events in the store
					// TODO: Add metrics for calendar sources
					// TODO: Add metrics for event duration in minutes

					// Update the store
					store.Update()

					// Upsert metrics to prometheus
					for _, e := range store.Events() {
						// Create metric for event
						var evt = prometheus.NewGauge(
							prometheus.GaugeOpts{
								Name: "calendar_event_info",
								Help: "Info on a calendar event",
								ConstLabels: prometheus.Labels{
									"id":          e.ID,
									"summary":     e.Summary,
									"description": e.Description,
									"location":    e.Location,
									"time_start":  e.StartTime.String(),
									"time_end":    e.EndTime.String(),
									"updated":     t.String(),
								},
							},
						)
						// Set metric to 1 like recommended for software versions:
						// https://www.robustperception.io/exposing-the-software-version-to-prometheus
						evt.Set(1.0)

						// Register metric
						prometheus.Register(evt)
					}
				}
			}
		}()

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

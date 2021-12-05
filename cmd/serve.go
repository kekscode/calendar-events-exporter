package cmd

import (
	"net/http"
	"time"

	"github.com/kekscode/calendar-events-exporter/pkg/calendar"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

/*
var (
	eventInfo = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "calendar_event_info",
		Help: "Info on a calendar event",
	})
)

func init() {
	// Expose and set to a fixed number
	// See: https://www.robustperception.io/exposing-the-software-version-to-prometheus
	prometheus.MustRegister(eventInfo)
	eventInfo.Set(1.0)
}
*/

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
					for _, e := range store.GetEvents() {
						log.Printf("%v+\n", e)
					}

					// Extract this in to a function returning all updated metrics
					//var eventIds []prometheus.Gauge
					//					if err := prometheus.Register(
					//						prometheus.NewGauge(
					//							prometheus.GaugeOpts{
					//								Name:        "calendar_event_info",
					//								Help:        "Info on a calendar event",
					//								ConstLabels: prometheus.Labels{
					//									//"uid":         store.Events[1].GetProperty(ics.ComponentPropertyUniqueId).Value,
					//									//"summary":     store.Events[1].GetProperty(ics.ComponentPropertySummary).Value,
					//									//"description": store.Events[1].GetProperty(ics.ComponentPropertyDescription).Value,
					//									//"location":    store.Events[1].GetProperty(ics.ComponentPropertyLocation).Value,
					//									//"dstart":      store.Events[1].GetProperty(ics.ComponentPropertyDtStart).Value,
					//									//"dend": store.Events[1].GetProperty(ics.ComponentPropertyDtEnd).Value,
					//								},
					//							}),
					//					); err != nil {
					//						log.Printf("Could not register metrics: %v\n", err)
					//					}

					//for _, e := range store.Events {
					//	eventIds = append(eventIds, prometheus.NewGauge(
					//		prometheus.GaugeOpts{
					//			Name:        "calendar_event_info",
					//			Help:        "Info on a calendar event",
					//			ConstLabels: prometheus.Labels{"uid": e.GetProperty("UID").IANAToken},
					//		}))
					//}

					//// Register generated metrics
					//for _, e := range eventIds {
					//	if err := prometheus.Register(e); err != nil {
					//		log.Printf("Could not register metrics: %v\n", err)
					//	}
					//}

					log.Printf("Tick at", t)

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

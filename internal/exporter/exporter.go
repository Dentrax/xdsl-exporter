package exporter

import (
	"3e8.eu/go/dsl"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	Namespace = "xdsl"
)

type Exporter struct {
	dsl    dsl.Client
	logger log.Logger

	// via go-dsl
	// see: https://github.com/janh/go-dsl/blob/690a62b79cd43d01b5f10fe2ef0d1a8a2b3f00f7/models/status.go#L13-L77
	state                    *prometheus.Desc
	mode                     *prometheus.Desc
	uptime                   *prometheus.Desc
	farEndInventory          *prometheus.Desc
	nearEndInventory         *prometheus.Desc
	downstreamActualRate     *prometheus.Desc
	upstreamActualRate       *prometheus.Desc
	downstreamAttainableRate *prometheus.Desc
	upstreamAttainableRate   *prometheus.Desc
	// downstreamMinimumErrorFreeThroughput *prometheus.Desc
	// upstreamMinimumErrorFreeThroughput   *prometheus.Desc
	// downstreamBitswapEnabled             *prometheus.Desc
	// upstreamBitswapEnabled               *prometheus.Desc
	// downstreamSeamlessRateAdaption       *prometheus.Desc
	// upstreamSeamlessRateAdaption         *prometheus.Desc
	// downstreamInterleavingDelay          *prometheus.Desc
	// upstreamInterleavingDelay            *prometheus.Desc
	// downstreamImpulseNoiseProtection     *prometheus.Desc
	// upstreamImpulseNoiseProtection       *prometheus.Desc
	// downstreamRetransmissionEnabled      *prometheus.Desc
	// upstreamRetransmissionEnabled        *prometheus.Desc
	// downstreamVectoringState             *prometheus.Desc
	// upstreamVectoringState               *prometheus.Desc
	// downstreamAttenuation                *prometheus.Desc
	// upstreamAttenuation                  *prometheus.Desc
	// downstreamSNRMargin                  *prometheus.Desc
	// upstreamSNRMargin                    *prometheus.Desc
	// downstreamPower                      *prometheus.Desc
	// upstreamPower                        *prometheus.Desc
	// downstreamRTXTXCount                 *prometheus.Desc
	// upstreamRTXTXCount                   *prometheus.Desc
	// downstreamRTXCCount                  *prometheus.Desc
	// upstreamRTXCCount                    *prometheus.Desc
	// downstreamRTXUCCount                 *prometheus.Desc
	// upstreamRTXUCCount                   *prometheus.Desc
	// downstreamFECCount                   *prometheus.Desc
	// upstreamFECCount                     *prometheus.Desc
	// downstreamCRCCount                   *prometheus.Desc
	// upstreamCRCCount                     *prometheus.Desc
	// downstreamESCount                    *prometheus.Desc
	// upstreamESCount                      *prometheus.Desc
	// downstreamSESCount                   *prometheus.Desc
	// upstreamSESCount                     *prometheus.Desc
}

func New(dsl dsl.Client, logger log.Logger) *Exporter {
	return &Exporter{
		dsl:    dsl,
		logger: logger,
		state: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "state"),
			"State of the DSL modem.",
			[]string{"state"},
			nil,
		),
		mode: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "mode"),
			"Mode of the DSL modem.",
			[]string{"mode"},
			nil,
		),
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "uptime"),
			"Uptime of the DSL modem.",
			[]string{"uptime"},
			nil,
		),
		farEndInventory: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "remote"),
			"Far end inventory name of the manufacturer",
			[]string{"vendor", "version"},
			nil,
		),
		nearEndInventory: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "modem"),
			"Near end inventory name of the manufacturer.",
			[]string{"vendor", "version"},
			nil,
		),
		downstreamActualRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "actual_rate_downstream"),
			"Actual rate of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamActualRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "actual_rate_upstream"),
			"Actual rate of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamAttainableRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "attainable_rate_downstream"),
			"Attainable rate of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamAttainableRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "", "attainable_rate_upstream"),
			"Attainable rate of upstream.",
			[]string{"unit"},
			nil,
		),
	}
}

func (e *Exporter) Describe(descs chan<- *prometheus.Desc) {
	descs <- e.state
	descs <- e.mode
	descs <- e.uptime
	descs <- e.farEndInventory
	descs <- e.nearEndInventory
	descs <- e.downstreamActualRate
	descs <- e.upstreamActualRate
	descs <- e.downstreamAttainableRate
	descs <- e.upstreamActualRate
}

func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	level.Debug(e.logger).Log("msg", "collecting metrics...")

	if err := e.getDataFromClient(metrics); err != nil {
		level.Error(e.logger).Log("msg", "could not get data", "err", err.Error()) //nolint:errcheck
	}
}

func (e *Exporter) getDataFromClient(metrics chan<- prometheus.Metric) error {
	if err := e.dsl.UpdateData(); err != nil {
		return err
	}

	status := e.dsl.Status()

	metrics <- prometheus.MustNewConstMetric(e.state, prometheus.UntypedValue, 1, status.State.String())
	metrics <- prometheus.MustNewConstMetric(e.mode, prometheus.UntypedValue, 1, status.Mode.String())
	metrics <- prometheus.MustNewConstMetric(e.uptime, prometheus.UntypedValue, 1, status.Uptime.String())
	metrics <- prometheus.MustNewConstMetric(e.farEndInventory, prometheus.UntypedValue, 1, status.FarEndInventory.Vendor, status.FarEndInventory.Version)
	metrics <- prometheus.MustNewConstMetric(e.nearEndInventory, prometheus.UntypedValue, 1, status.NearEndInventory.Vendor, status.NearEndInventory.Version)
	metrics <- prometheus.MustNewConstMetric(e.downstreamActualRate, prometheus.GaugeValue, float64(status.DownstreamActualRate.Int), status.DownstreamActualRate.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamActualRate, prometheus.GaugeValue, float64(status.UpstreamActualRate.Int), status.UpstreamActualRate.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamAttainableRate, prometheus.GaugeValue, float64(status.DownstreamAttainableRate.Int), status.DownstreamAttainableRate.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamAttainableRate, prometheus.GaugeValue, float64(status.UpstreamAttainableRate.Int), status.UpstreamAttainableRate.Unit())

	return nil
}

func (e *Exporter) CloseClient() {
	e.dsl.Close()
}

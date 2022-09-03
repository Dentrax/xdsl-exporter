package exporter

import (
	"3e8.eu/go/dsl"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	rtop "github.com/rapidloop/rtop/pkg/client"
	"github.com/rapidloop/rtop/pkg/types"
	"strconv"
)

const (
	Namespace     = "xdsl"
	SubsystemDsl  = "dsl"
	SubsystemRtop = "rtop"
)

type Exporter struct {
	dsl    dsl.Client
	rtop   *rtop.Client
	logger log.Logger

	// via go-dsl
	// see: https://github.com/janh/go-dsl/blob/690a62b79cd43d01b5f10fe2ef0d1a8a2b3f00f7/models/status.go#L13-L77
	state                                *prometheus.Desc
	mode                                 *prometheus.Desc
	uptime                               *prometheus.Desc
	farEndInventory                      *prometheus.Desc
	nearEndInventory                     *prometheus.Desc
	downstreamActualRate                 *prometheus.Desc
	upstreamActualRate                   *prometheus.Desc
	downstreamAttainableRate             *prometheus.Desc
	upstreamAttainableRate               *prometheus.Desc
	downstreamMinimumErrorFreeThroughput *prometheus.Desc
	upstreamMinimumErrorFreeThroughput   *prometheus.Desc
	downstreamBitswapEnabled             *prometheus.Desc
	upstreamBitswapEnabled               *prometheus.Desc
	downstreamSeamlessRateAdaption       *prometheus.Desc
	upstreamSeamlessRateAdaption         *prometheus.Desc
	downstreamInterleavingDelay          *prometheus.Desc
	upstreamInterleavingDelay            *prometheus.Desc
	downstreamImpulseNoiseProtection     *prometheus.Desc
	upstreamImpulseNoiseProtection       *prometheus.Desc
	downstreamRetransmissionEnabled      *prometheus.Desc
	upstreamRetransmissionEnabled        *prometheus.Desc
	downstreamVectoringState             *prometheus.Desc
	upstreamVectoringState               *prometheus.Desc
	downstreamAttenuation                *prometheus.Desc
	upstreamAttenuation                  *prometheus.Desc
	downstreamSNRMargin                  *prometheus.Desc
	upstreamSNRMargin                    *prometheus.Desc
	downstreamPower                      *prometheus.Desc
	upstreamPower                        *prometheus.Desc
	downstreamRTXTXCount                 *prometheus.Desc
	upstreamRTXTXCount                   *prometheus.Desc
	downstreamRTXCCount                  *prometheus.Desc
	upstreamRTXCCount                    *prometheus.Desc
	downstreamRTXUCCount                 *prometheus.Desc
	upstreamRTXUCCount                   *prometheus.Desc
	downstreamFECCount                   *prometheus.Desc
	upstreamFECCount                     *prometheus.Desc
	downstreamCRCCount                   *prometheus.Desc
	upstreamCRCCount                     *prometheus.Desc
	downstreamESCount                    *prometheus.Desc
	upstreamESCount                      *prometheus.Desc
	downstreamSESCount                   *prometheus.Desc
	upstreamSESCount                     *prometheus.Desc

	// via rtop
	rtopInfo         *prometheus.Desc
	rtopLoad1        *prometheus.Desc
	rtopLoad5        *prometheus.Desc
	rtopLoad15       *prometheus.Desc
	rtopLoadRunning  *prometheus.Desc
	rtopLoadTotal    *prometheus.Desc
	rtopCPUUser      *prometheus.Desc
	rtopCPUSystem    *prometheus.Desc
	rtopCPUNice      *prometheus.Desc
	rtopCPUIdle      *prometheus.Desc
	rtopCPUIOWait    *prometheus.Desc
	rtopCPUIRQ       *prometheus.Desc
	rtopCPUSoftIRQ   *prometheus.Desc
	rtopCPUSteal     *prometheus.Desc
	rtopCPUGuest     *prometheus.Desc
	rtopMEMTotal     *prometheus.Desc
	rtopMEMFree      *prometheus.Desc
	rtopMEMUsed      *prometheus.Desc
	rtopMEMBuffers   *prometheus.Desc
	rtopMEMCached    *prometheus.Desc
	rtopMEMSwapFree  *prometheus.Desc
	rtopMEMSwapTotal *prometheus.Desc
	rtopFSTotal      *prometheus.Desc
	rtopFSUsed       *prometheus.Desc
	rtopFSFree       *prometheus.Desc
	rtopNETRx        *prometheus.Desc
	rtopNETTx        *prometheus.Desc
}

func New(dsl dsl.Client, rtop *rtop.Client, logger log.Logger) *Exporter {
	return &Exporter{
		dsl:    dsl,
		rtop:   rtop,
		logger: logger,
		state: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "state"),
			"State of the DSL modem.",
			[]string{"state"},
			nil,
		),
		mode: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "mode"),
			"Mode of the DSL modem.",
			[]string{"mode"},
			nil,
		),
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "uptime"),
			"Uptime of the DSL modem.",
			[]string{"uptime"},
			nil,
		),
		farEndInventory: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "modem_manufacturer_far"),
			"Far end inventory name of the manufacturer",
			[]string{"vendor", "version"},
			nil,
		),
		nearEndInventory: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "modem_manufacturer_near"),
			"Near end inventory name of the manufacturer.",
			[]string{"vendor", "version"},
			nil,
		),
		downstreamActualRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "actual_rate_downstream"),
			"Actual rate of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamActualRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "actual_rate_upstream"),
			"Actual rate of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamAttainableRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "attainable_rate_downstream"),
			"Attainable rate of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamAttainableRate: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "attainable_rate_upstream"),
			"Attainable rate of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamMinimumErrorFreeThroughput: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "minimum_error_free_throughput_downstream"),
			"Minimum error free throughput of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamMinimumErrorFreeThroughput: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "minimum_error_free_throughput_upstream"),
			"Minimum error free throughput of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamBitswapEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "bitswap_enabled_downstream"),
			"Bitswap enabled of downstream.",
			nil,
			nil,
		),
		upstreamBitswapEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "bitswap_enabled_upstream"),
			"Bitswap enabled of upstream.",
			nil,
			nil,
		),
		downstreamSeamlessRateAdaption: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "seamless_rate_adaption_downstream"),
			"Seamless rate adaption of downstream.",
			nil,
			nil,
		),
		upstreamSeamlessRateAdaption: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "seamless_rate_adaption_upstream"),
			"Seamless rate adaption of upstream.",
			nil,
			nil,
		),
		downstreamInterleavingDelay: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "interleaving_delay_downstream"),
			"Interleaving delay of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamInterleavingDelay: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "interleaving_delay_upstream"),
			"Interleaving delay of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamImpulseNoiseProtection: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "impulse_noise_protection_downstream"),
			"Impulse noise protection of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamImpulseNoiseProtection: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "impulse_noise_protection_upstream"),
			"Impulse noise protection of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamRetransmissionEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "retransmission_enabled_downstream"),
			"Retransmission enabled of downstream.",
			nil,
			nil,
		),
		upstreamRetransmissionEnabled: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "retransmission_enabled_upstream"),
			"Retransmission enabled of upstream.",
			nil,
			nil,
		),
		downstreamVectoringState: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "vectoring_state_downstream"),
			"Vectoring state of downstream.",
			[]string{"state"},
			nil,
		),
		upstreamVectoringState: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "vectoring_state_upstream"),
			"Vectoring state of upstream.",
			[]string{"state"},
			nil,
		),
		downstreamAttenuation: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "attenuation_downstream"),
			"Attenuation of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamAttenuation: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "attenuation_upstream"),
			"Attenuation of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamSNRMargin: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "snr_margin_downstream"),
			"SNR margin of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamSNRMargin: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "snr_margin_upstream"),
			"SNR margin of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamPower: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "power_downstream"),
			"Power of downstream.",
			[]string{"unit"},
			nil,
		),
		upstreamPower: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "power_upstream"),
			"Power of upstream.",
			[]string{"unit"},
			nil,
		),
		downstreamRTXTXCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "rtxtx_count_downstream"),
			"RTTX TX count of downstream.",
			nil,
			nil,
		),
		upstreamRTXTXCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "rtxtx_count_upstream"),
			"RTTX TX count of upstream.",
			nil,
			nil,
		),
		downstreamRTXCCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "rtxcc_count_downstream"),
			"RTXCC count of downstream.",
			nil,
			nil,
		),
		upstreamRTXCCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "rtxcc_count_upstream"),
			"RTXCC count of upstream.",
			nil,
			nil,
		),
		downstreamRTXUCCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "rtxucc_count_downstream"),
			"RTXUCC count of downstream.",
			nil,
			nil,
		),
		upstreamRTXUCCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "rtxucc_count_upstream"),
			"RTXUCC count of upstream.",
			nil,
			nil,
		),
		downstreamFECCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "fec_count_downstream"),
			"FEC count of downstream.",
			nil,
			nil,
		),
		upstreamFECCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "fec_count_upstream"),
			"FEC count of upstream.",
			nil,
			nil,
		),
		downstreamCRCCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "crc_count_downstream"),
			"CRC count of downstream.",
			nil,
			nil,
		),
		upstreamCRCCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "crc_count_upstream"),
			"CRC count of upstream.",
			nil,
			nil,
		),
		downstreamESCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "es_count_downstream"),
			"ES count of downstream.",
			nil,
			nil,
		),
		upstreamESCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "es_count_upstream"),
			"ES count of upstream.",
			nil,
			nil,
		),
		downstreamSESCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "ses_count_downstream"),
			"SES count of downstream.",
			nil,
			nil,
		),
		upstreamSESCount: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemDsl, "ses_count_upstream"),
			"SES count of upstream.",
			nil,
			nil,
		),
		// rtop
		rtopInfo: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "info"),
			"Information about the host.",
			[]string{"hostname", "uptime"},
			nil,
		),
		rtopLoad1: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "load1"),
			"Load1 of the host.",
			nil,
			nil,
		),
		rtopLoad5: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "load5"),
			"Load5 of the host.",
			nil,
			nil,
		),
		rtopLoad15: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "load15"),
			"Load15 of the host.",
			nil,
			nil,
		),
		rtopLoadTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "load_total"),
			"LoadTotal of the host.",
			nil,
			nil,
		),
		rtopLoadRunning: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "load_running"),
			"LoadRunning of the host.",
			nil,
			nil,
		),
		rtopCPUUser: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_user"),
			"CPU user of the host.",
			nil,
			nil,
		),
		rtopCPUSystem: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_system"),
			"CPU system of the host.",
			nil,
			nil,
		),
		rtopCPUNice: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_nice"),
			"CPU nice of the host.",
			nil,
			nil,
		),
		rtopCPUIdle: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_idle"),
			"CPU idle of the host.",
			nil,
			nil,
		),
		rtopCPUIOWait: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_iowait"),
			"CPU iowait of the host.",
			nil,
			nil,
		),
		rtopCPUIRQ: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_irq"),
			"CPU irq of the host.",
			nil,
			nil,
		),
		rtopCPUSoftIRQ: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_softirq"),
			"CPU softirq of the host.",
			nil,
			nil,
		),
		rtopCPUSteal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_steal"),
			"CPU steal of the host.",
			nil,
			nil,
		),
		rtopCPUGuest: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "cpu_guest"),
			"CPU guest of the host.",
			nil,
			nil,
		),
		rtopMEMTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_total"),
			"Total memory of the host.",
			nil,
			nil,
		),
		rtopMEMFree: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_free"),
			"Free memory of the host.",
			nil,
			nil,
		),
		rtopMEMUsed: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_used"),
			"Used memory of the host.",
			nil,
			nil,
		),
		rtopMEMBuffers: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_buffers"),
			"Buffers memory of the host.",
			nil,
			nil,
		),
		rtopMEMCached: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_cached"),
			"Cached memory of the host.",
			nil,

			nil,
		),
		rtopMEMSwapFree: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_swap_free"),
			"Free swap memory of the host.",
			nil,
			nil,
		),
		rtopMEMSwapTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "mem_swap_total"),
			"Total swap memory of the host.",
			nil,
			nil,
		),
		rtopFSTotal: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "fs_total"),
			"Total filesystems of the host.",
			[]string{"mount"},
			nil,
		),
		rtopFSUsed: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "fs_used"),
			"Used filesystem of the host.",
			[]string{"mount"},
			nil,
		),
		rtopFSFree: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "fs_free"),
			"Free filesystem of the host.",
			[]string{"mount"},
			nil,
		),
		rtopNETRx: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "net_rx"),
			"Total received bytes of the network.",
			[]string{"interface", "ipv4", "ipv6"},
			nil,
		),
		rtopNETTx: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, SubsystemRtop, "net_tx"),
			"Total transmitted bytes of the network.",
			[]string{"interface", "ipv4", "ipv6"},
			nil,
		),
	}
}

func (e *Exporter) Describe(descs chan<- *prometheus.Desc) {
	// dsl
	descs <- e.state
	descs <- e.mode
	descs <- e.uptime
	descs <- e.farEndInventory
	descs <- e.nearEndInventory
	descs <- e.downstreamActualRate
	descs <- e.upstreamActualRate
	descs <- e.downstreamAttainableRate
	descs <- e.upstreamActualRate
	descs <- e.downstreamMinimumErrorFreeThroughput
	descs <- e.upstreamMinimumErrorFreeThroughput
	descs <- e.downstreamBitswapEnabled
	descs <- e.upstreamBitswapEnabled
	descs <- e.downstreamSeamlessRateAdaption
	descs <- e.upstreamSeamlessRateAdaption
	descs <- e.downstreamInterleavingDelay
	descs <- e.upstreamInterleavingDelay
	descs <- e.downstreamImpulseNoiseProtection
	descs <- e.upstreamImpulseNoiseProtection
	descs <- e.downstreamRetransmissionEnabled
	descs <- e.upstreamRetransmissionEnabled
	descs <- e.downstreamVectoringState
	descs <- e.upstreamVectoringState
	descs <- e.downstreamAttenuation
	descs <- e.upstreamAttenuation
	descs <- e.downstreamSNRMargin
	descs <- e.upstreamSNRMargin
	descs <- e.downstreamPower
	descs <- e.upstreamPower
	descs <- e.downstreamRTXTXCount
	descs <- e.upstreamRTXTXCount
	descs <- e.downstreamRTXCCount
	descs <- e.upstreamRTXCCount
	descs <- e.downstreamRTXUCCount
	descs <- e.upstreamRTXUCCount
	descs <- e.downstreamFECCount
	descs <- e.upstreamFECCount
	descs <- e.downstreamCRCCount
	descs <- e.upstreamCRCCount
	descs <- e.downstreamESCount
	descs <- e.upstreamESCount
	descs <- e.downstreamSESCount
	descs <- e.upstreamSESCount

	// rtop
	descs <- e.rtopInfo
	descs <- e.rtopLoad1
	descs <- e.rtopLoad5
	descs <- e.rtopLoad15
	descs <- e.rtopLoadRunning
	descs <- e.rtopLoadTotal
	descs <- e.rtopCPUUser
	descs <- e.rtopCPUSystem
	descs <- e.rtopCPUNice
	descs <- e.rtopCPUIdle
	descs <- e.rtopCPUIOWait
	descs <- e.rtopCPUIRQ
	descs <- e.rtopCPUSoftIRQ
	descs <- e.rtopCPUSteal
	descs <- e.rtopCPUGuest
	descs <- e.rtopMEMTotal
	descs <- e.rtopMEMFree
	descs <- e.rtopMEMUsed
	descs <- e.rtopMEMBuffers
	descs <- e.rtopMEMCached
	descs <- e.rtopMEMSwapFree
	descs <- e.rtopMEMSwapTotal
	descs <- e.rtopFSTotal
	descs <- e.rtopFSUsed
	descs <- e.rtopFSFree
	descs <- e.rtopNETRx
	descs <- e.rtopNETTx
}

func (e *Exporter) Collect(metrics chan<- prometheus.Metric) {
	level.Debug(e.logger).Log("msg", "collecting metrics...")

	if err := e.getDataFromClients(metrics); err != nil {
		level.Error(e.logger).Log("msg", "could not get data", "err", err.Error()) //nolint:errcheck
	}
}

func (e *Exporter) getDataFromClients(metrics chan<- prometheus.Metric) error {
	if err := e.getDataFromDsl(metrics); err != nil {
		return err
	}
	if err := e.getDataFromRtop(metrics); err != nil {
		return err
	}
	return nil
}

func (e *Exporter) getDataFromDsl(metrics chan<- prometheus.Metric) error {
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
	metrics <- prometheus.MustNewConstMetric(e.downstreamMinimumErrorFreeThroughput, prometheus.GaugeValue, float64(status.DownstreamMinimumErrorFreeThroughput.Int), status.DownstreamMinimumErrorFreeThroughput.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamMinimumErrorFreeThroughput, prometheus.GaugeValue, float64(status.UpstreamMinimumErrorFreeThroughput.Int), status.UpstreamMinimumErrorFreeThroughput.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamBitswapEnabled, prometheus.GaugeValue, boolToFloat64(status.DownstreamBitswapEnabled.Bool))
	metrics <- prometheus.MustNewConstMetric(e.upstreamBitswapEnabled, prometheus.GaugeValue, boolToFloat64(status.UpstreamBitswapEnabled.Bool))
	metrics <- prometheus.MustNewConstMetric(e.downstreamSeamlessRateAdaption, prometheus.GaugeValue, boolToFloat64(status.DownstreamSeamlessRateAdaption.Bool))
	metrics <- prometheus.MustNewConstMetric(e.upstreamSeamlessRateAdaption, prometheus.GaugeValue, boolToFloat64(status.UpstreamSeamlessRateAdaption.Bool))
	metrics <- prometheus.MustNewConstMetric(e.downstreamInterleavingDelay, prometheus.GaugeValue, status.DownstreamInterleavingDelay.Float, status.DownstreamInterleavingDelay.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamInterleavingDelay, prometheus.GaugeValue, status.UpstreamInterleavingDelay.Float, status.UpstreamInterleavingDelay.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamImpulseNoiseProtection, prometheus.GaugeValue, status.DownstreamImpulseNoiseProtection.Float, status.DownstreamImpulseNoiseProtection.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamImpulseNoiseProtection, prometheus.GaugeValue, status.UpstreamImpulseNoiseProtection.Float, status.UpstreamImpulseNoiseProtection.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamRetransmissionEnabled, prometheus.GaugeValue, boolToFloat64(status.DownstreamRetransmissionEnabled.Bool))
	metrics <- prometheus.MustNewConstMetric(e.upstreamRetransmissionEnabled, prometheus.GaugeValue, boolToFloat64(status.UpstreamRetransmissionEnabled.Bool))
	metrics <- prometheus.MustNewConstMetric(e.downstreamVectoringState, prometheus.GaugeValue, 1, status.DownstreamVectoringState.Value())
	metrics <- prometheus.MustNewConstMetric(e.upstreamVectoringState, prometheus.GaugeValue, 1, status.UpstreamVectoringState.Value())
	metrics <- prometheus.MustNewConstMetric(e.downstreamAttenuation, prometheus.GaugeValue, status.DownstreamAttenuation.Float, status.DownstreamAttenuation.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamAttenuation, prometheus.GaugeValue, status.UpstreamAttenuation.Float, status.UpstreamAttenuation.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamSNRMargin, prometheus.GaugeValue, status.DownstreamSNRMargin.Float, status.DownstreamSNRMargin.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamSNRMargin, prometheus.GaugeValue, status.UpstreamSNRMargin.Float, status.UpstreamSNRMargin.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamPower, prometheus.GaugeValue, status.DownstreamPower.Float, status.DownstreamPower.Unit())
	metrics <- prometheus.MustNewConstMetric(e.upstreamPower, prometheus.GaugeValue, status.UpstreamPower.Float, status.UpstreamPower.Unit())
	metrics <- prometheus.MustNewConstMetric(e.downstreamRTXTXCount, prometheus.GaugeValue, float64(status.DownstreamRTXTXCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamRTXTXCount, prometheus.GaugeValue, float64(status.UpstreamRTXTXCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.downstreamRTXCCount, prometheus.GaugeValue, float64(status.DownstreamRTXCCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamRTXCCount, prometheus.GaugeValue, float64(status.UpstreamRTXCCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.downstreamRTXUCCount, prometheus.GaugeValue, float64(status.DownstreamRTXUCCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamRTXUCCount, prometheus.GaugeValue, float64(status.UpstreamRTXUCCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.downstreamFECCount, prometheus.GaugeValue, float64(status.DownstreamFECCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamFECCount, prometheus.GaugeValue, float64(status.UpstreamFECCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.downstreamCRCCount, prometheus.GaugeValue, float64(status.DownstreamCRCCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamCRCCount, prometheus.GaugeValue, float64(status.UpstreamCRCCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.downstreamESCount, prometheus.GaugeValue, float64(status.DownstreamESCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamESCount, prometheus.GaugeValue, float64(status.UpstreamESCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.downstreamSESCount, prometheus.GaugeValue, float64(status.DownstreamSESCount.Int))
	metrics <- prometheus.MustNewConstMetric(e.upstreamSESCount, prometheus.GaugeValue, float64(status.UpstreamSESCount.Int))

	return nil
}

func (e *Exporter) getDataFromRtop(metrics chan<- prometheus.Metric) error {
	stats, err := e.rtop.GetStats()
	if err != nil {
		return err
	}

	metrics <- prometheus.MustNewConstMetric(e.rtopInfo, prometheus.GaugeValue, 1, stats.Hostname, stats.Uptime.String())

	metrics <- prometheus.MustNewConstMetric(e.rtopLoad1, prometheus.GaugeValue, stringToFloat64(stats.Loads.Load1))
	metrics <- prometheus.MustNewConstMetric(e.rtopLoad5, prometheus.GaugeValue, stringToFloat64(stats.Loads.Load5))
	metrics <- prometheus.MustNewConstMetric(e.rtopLoad15, prometheus.GaugeValue, stringToFloat64(stats.Loads.Load15))
	metrics <- prometheus.MustNewConstMetric(e.rtopLoadRunning, prometheus.GaugeValue, stringToFloat64(stats.Loads.RunningProcs))
	metrics <- prometheus.MustNewConstMetric(e.rtopLoadTotal, prometheus.GaugeValue, stringToFloat64(stats.Loads.TotalProcs))

	metrics <- prometheus.MustNewConstMetric(e.rtopCPUUser, prometheus.GaugeValue, float64(stats.CPU.User))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUSystem, prometheus.GaugeValue, float64(stats.CPU.System))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUNice, prometheus.GaugeValue, float64(stats.CPU.Nice))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUIdle, prometheus.GaugeValue, float64(stats.CPU.Idle))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUIOWait, prometheus.GaugeValue, float64(stats.CPU.IOWait))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUIRQ, prometheus.GaugeValue, float64(stats.CPU.IRQ))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUSoftIRQ, prometheus.GaugeValue, float64(stats.CPU.SoftIRQ))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUSteal, prometheus.GaugeValue, float64(stats.CPU.Steal))
	metrics <- prometheus.MustNewConstMetric(e.rtopCPUGuest, prometheus.GaugeValue, float64(stats.CPU.Guest))

	metrics <- prometheus.MustNewConstMetric(e.rtopMEMTotal, prometheus.GaugeValue, float64(stats.MEM.Total))
	metrics <- prometheus.MustNewConstMetric(e.rtopMEMFree, prometheus.GaugeValue, float64(stats.MEM.Free))
	metrics <- prometheus.MustNewConstMetric(e.rtopMEMUsed, prometheus.GaugeValue, float64(stats.MEM.Used()))
	metrics <- prometheus.MustNewConstMetric(e.rtopMEMCached, prometheus.GaugeValue, float64(stats.MEM.Cached))
	metrics <- prometheus.MustNewConstMetric(e.rtopMEMBuffers, prometheus.GaugeValue, float64(stats.MEM.Buffers))
	metrics <- prometheus.MustNewConstMetric(e.rtopMEMSwapFree, prometheus.GaugeValue, float64(stats.MEM.SwapFree))
	metrics <- prometheus.MustNewConstMetric(e.rtopMEMSwapTotal, prometheus.GaugeValue, float64(stats.MEM.SwapTotal))

	e.calculateNET(stats, e.rtopNETRx, e.rtopNETTx, metrics)
	e.calculateFS(stats, e.rtopFSTotal, e.rtopFSUsed, e.rtopFSFree, metrics)

	return nil
}

func (e *Exporter) calculateNET(stats types.Stats, descRx *prometheus.Desc, descTx *prometheus.Desc, metrics chan<- prometheus.Metric) {
	for k, v := range stats.NetInterface {
		metrics <- prometheus.MustNewConstMetric(descRx, prometheus.GaugeValue, float64(v.Rx), k, v.IPv4, v.IPv6)
		metrics <- prometheus.MustNewConstMetric(descTx, prometheus.GaugeValue, float64(v.Tx), k, v.IPv4, v.IPv6)
	}
}

func (e *Exporter) calculateFS(stats types.Stats, descTotal *prometheus.Desc, descUsed *prometheus.Desc, descFree *prometheus.Desc, metrics chan<- prometheus.Metric) {
	for _, fs := range stats.FSInfos {
		metrics <- prometheus.MustNewConstMetric(descTotal, prometheus.GaugeValue, float64(fs.Total), fs.MountPoint)
		metrics <- prometheus.MustNewConstMetric(descUsed, prometheus.GaugeValue, float64(fs.Used), fs.MountPoint)
		metrics <- prometheus.MustNewConstMetric(descFree, prometheus.GaugeValue, float64(fs.Free), fs.MountPoint)
	}
}

func (e *Exporter) CloseClient() {
	e.dsl.Close()
}

func boolToFloat64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func stringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return f
}

/*
Copyright © 2022 Furkan Türkal

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/Dentrax/xdsl-exporter/internal/rtop"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Dentrax/xdsl-exporter/internal/config"
	"github.com/Dentrax/xdsl-exporter/internal/dsl"
	"github.com/Dentrax/xdsl-exporter/internal/exporter"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfg     = config.Config{}
	cfgFile string
	cmd     = &cobra.Command{
		Use:           "xdsl-exporter",
		Short:         "A Prometheus Exporter for your rusty xDSL Modem",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func Execute() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVar(&cfg.ListenAddress, "listen-address", ":9090", "Address on which to expose metrics and web interface.")
	cmd.PersistentFlags().StringVar(&cfg.MetricsPath, "metrics-path", "/metrics", "Path under which to expose metrics.")
	cmd.PersistentFlags().StringVar(&cfg.KnownHostsPath, "known-hosts-path", "~/.ssh/known_hosts", "Path to your known_hosts file.")
	cmd.PersistentFlags().StringVar(&cfg.TargetHost, "target-host", "192.168.1.1", "Hostname or IP address of the target xDSL Modem")
	cmd.PersistentFlags().IntVar(&cfg.TargetPort, "target-port", 22, "Port of the target xDSL Modem")
	cmd.PersistentFlags().StringVar(&cfg.TargetUser, "target-user", "admin", "Host user")
	cmd.PersistentFlags().StringVar(&cfg.TargetPassword, "target-password", "", "Host password")
	cmd.PersistentFlags().StringVar(&cfg.TargetSSHKeyPath, "target-ssh-key-path", "", "Path to the SSH key to use for authentication")
	cmd.PersistentFlags().StringVar(&cfg.TargetSSHPassphrase, "target-ssh-passphrase", "", "Passphrase to use for the SSH key")
	cmd.PersistentFlags().StringVar(&cfg.TargetClient, "target-client", "", strings.Join(dsl.GetSupportedClients(), ","))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".xdsl-exporter")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func run() error {
	promlogConfig := &promlog.Config{
		Level: &promlog.AllowedLevel{},
	}
	promlogConfig.Level.Set("info")
	logger := promlog.New(promlogConfig)

	prometheus.MustRegister(version.NewCollector("xdsl_exporter"))

	if err := cfg.Check(); err != nil {
		return fmt.Errorf("config check: %w", err)
	}

	dslClient, err := dsl.New(cfg)
	if err != nil {
		return err
	}

	rtopClient, err := rtop.New(cfg)
	if err != nil {
		return err
	}

	exporter := exporter.New(dslClient, rtopClient, logger)
	prometheus.MustRegister(exporter)

	http.Handle(cfg.MetricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
            	<html>
            	<head><title>xDSL Exporter Metrics</title></head>
            	<body>
            	<p><a href='` + cfg.MetricsPath + `'>Metrics</a></p>
            	</body>
            	</html>
				`))
	})

	go func() {
		level.Info(logger).Log("msg", "Listening on address", "address", cfg.ListenAddress) //nolint:errcheck
		srv := &http.Server{
			Addr:              cfg.ListenAddress,
			ReadHeaderTimeout: 60 * time.Second,
		}
		if err := web.ListenAndServe(srv, "", logger); err != nil {
			level.Error(logger).Log("msg", "Error running HTTP server", "err", err) //nolint:errcheck
			os.Exit(1)
		}
	}()
	done := make(chan struct{})
	go func() {
		level.Info(logger).Log("msg", "Listening signals...") //nolint:errcheck
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		exporter.CloseClient()
		close(done)
	}()

	<-done
	level.Info(logger).Log("msg", "Shutting down...") //nolint:errcheck

	return nil
}

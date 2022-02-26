package cmd

import (
	"context"
	"log"
	"os"

	"github.com/SaltFishPr/redis-viewer/internal/config"
	"github.com/SaltFishPr/redis-viewer/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "redis-viewer",
	Short: "view redis data in terminal.",
	Long:  `Redis Viewer is a tool to view redis data in terminal.`,
	Run: func(cmd *cobra.Command, args []string) {
		config.LoadConfig()
		cfg := config.GetConfig()

		var rdb redis.Cmdable
		switch cfg.Mode {
		case "sentinel":
			rdb = redis.NewFailoverClient(
				&redis.FailoverOptions{
					MasterName:    cfg.MasterName,
					SentinelAddrs: cfg.SentinelAddrs,
					Username:      cfg.Username,
					Password:      cfg.Password,
				},
			)
		case "cluster":
			rdb = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    cfg.ClusterAddrs,
				Username: cfg.Username,
				Password: cfg.Password,
			})
		default:
			rdb = redis.NewClient(
				&redis.Options{
					Addr:     cfg.Addr,
					Username: cfg.Username,
					Password: cfg.Password,
					DB:       cfg.DB,
				},
			)
		}

		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			log.Fatal("connect to redis failed: ", err)
		}

		p := tea.NewProgram(tui.New(rdb), tea.WithAltScreen(), tea.WithMouseCellMotion())
		if err := p.Start(); err != nil {
			log.Fatal("start failed: ", err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.redis-viewer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

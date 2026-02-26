package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/MaikelVeen/voice-agent/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	argVerbose = "verbose"

	envPrefix = "VOICE"
)

var rootCmd = &cobra.Command{
	Use:          "voice-agent <command>",
	Short:        "voice-agent is a CLI to make coding agent speak",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		return initializeConfig(cmd)
	},
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if viper.GetBool("version") {
			fmt.Println(version.Get().String())
			return nil
		}
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(NewVersionCommand().Command)

	pflags := rootCmd.PersistentFlags()

	pflags.BoolP(argVerbose, "v", false, "Enable verbose output")
	_ = viper.BindPFlag(argVerbose, pflags.Lookup(argVerbose))

	pflags.Bool("version", false, "Show version information")
	_ = viper.BindPFlag("version", pflags.Lookup("version"))
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			_ = cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

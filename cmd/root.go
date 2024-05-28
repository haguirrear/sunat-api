/*
Copyright © 2023 Hector Aguirre <hector.aguirre.arista@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var ConfigData Config

type Config struct {
	User         string
	Password     string
	ClientID     string
	ClientSecret string
	AuthBaseURL  string
	BaseURL      string
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sunat",
	Short: "CLI app para interactuar con la API REST de SUNAT",
	Long:  `sunatapi es una app de terminal que interactua con la API REST de SUNAT.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("root called")
	// },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.sunatapi.yaml)")
	RootCmd.PersistentFlags().StringP("user", "u", "", "Usuario (RUC+Usuario SOL)")
	RootCmd.PersistentFlags().StringP("password", "p", "", "Clave SOL")
	RootCmd.PersistentFlags().String("client-id", "", "Client Id para el uso de la API de SUNAT")
	RootCmd.PersistentFlags().String("client-secret", "", "Client Secret para el uso de la API de SUNAT")
	RootCmd.PersistentFlags().String("auth-url", "https://api-seguridad.sunat.gob.pe", "URL base para el endpoint de obtener Token")
	RootCmd.PersistentFlags().String("base-url", "https://api-cpe.sunat.gob.pe", "URL base para las apis de SUNAT")

	viper.BindPFlag("user", RootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", RootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("clientid", RootCmd.PersistentFlags().Lookup("client-id"))
	viper.BindPFlag("clientsecret", RootCmd.PersistentFlags().Lookup("client-secret"))
	viper.BindPFlag("authbaseurl", RootCmd.PersistentFlags().Lookup("auth-url"))
	viper.BindPFlag("baseurl", RootCmd.PersistentFlags().Lookup("base-url"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// // Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name ".sunatapi" (without extension).
		currDir, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(currDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".sunatapi")
	}

	viper.SetEnvPrefix("SUNAT")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	parseAndValidateConfig()

}

func parseAndValidateConfig() {
	if err := viper.Unmarshal(&ConfigData); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading configuration: %v\n", err)
		os.Exit(1)
	}

	var unset []string

	if ConfigData.User == "" {
		unset = append(unset, "User")
	}

	if ConfigData.Password == "" {
		unset = append(unset, "Password")
	}

	if ConfigData.ClientID == "" {
		unset = append(unset, "ClientID")
	}

	if ConfigData.ClientSecret == "" {
		unset = append(unset, "ClientSecret")
	}

	if ConfigData.AuthBaseURL == "" {
		unset = append(unset, "AuthBaseURL")
	}

	if ConfigData.BaseURL == "" {
		unset = append(unset, "BaseURL")
	}

	if len(unset) > 0 {
		fmt.Fprintf(os.Stderr, "error: required configurations not set\n")
		for _, name := range unset {
			fmt.Fprintf(os.Stderr, "- %s\n", name)
		}
	}
}

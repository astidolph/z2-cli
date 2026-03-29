package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/strava-cli/internal/auth"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with Strava",
	Long: `Connect strava-cli to your Strava account.

You'll need a Strava API application. Create one at https://www.strava.com/settings/api
Set the "Authorization Callback Domain" to "localhost".`,
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter your Strava Client ID: ")
		clientID, _ := reader.ReadString('\n')
		clientID = strings.TrimSpace(clientID)

		fmt.Print("Enter your Strava Client Secret: ")
		clientSecret, _ := reader.ReadString('\n')
		clientSecret = strings.TrimSpace(clientSecret)

		if clientID == "" || clientSecret == "" {
			return fmt.Errorf("client ID and secret are required")
		}

		config := &auth.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		}
		if err := auth.SaveConfig(config); err != nil {
			return fmt.Errorf("could not save config: %w", err)
		}

		token, err := auth.Authenticate(clientID, clientSecret)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		if err := auth.SaveToken(token); err != nil {
			return fmt.Errorf("could not save token: %w", err)
		}

		fmt.Println("Successfully authenticated with Strava!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

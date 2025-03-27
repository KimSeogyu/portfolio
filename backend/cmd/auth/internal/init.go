/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package internal

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kimseogyu/portfolio/backend/cmd/auth/internal/config"
	"github.com/spf13/cobra"
)

var (
	dbType      string
	dbHost      string
	dbPort      int
	dbUser      string
	dbPassword  string
	dbName      string
	dbSchema    string
	grpcPort    int
	gatewayPort int
	outputPath  string
	force       bool
	interactive bool
)

func promptForInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func promptForConfirmation(prompt string) bool {
	for {
		fmt.Print(prompt)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Please answer 'y' or 'n'")
	}
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	Long: `Initialize a new configuration file.
This command will create a config file with the specified configuration.
In interactive mode, you will be prompted for each configuration value.
In non-interactive mode, you can provide values via command line flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if config file exists
		if _, err := os.Stat(outputPath); err == nil {
			if !force {
				if interactive {
					if !promptForConfirmation(fmt.Sprintf("Config file '%s' already exists. Do you want to overwrite it? (y/n): ", outputPath)) {
						fmt.Println("Operation cancelled.")
						return
					}
				} else {
					fmt.Printf("Error: Config file '%s' already exists. Use --force to overwrite or --interactive for interactive mode.\n", outputPath)
					return
				}
			}
		}

		var cfg *config.Config
		if interactive {
			// Interactive mode
			fmt.Println("Enter database configuration (press Enter to use default values):")
			defaultConfig := config.Local()

			if dbType == "" {
				dbType = promptForInput(fmt.Sprintf("Database type (postgres/test_postgres) [%s]: ", defaultConfig.DBConfig.DBType))
				if dbType == "" {
					dbType = string(defaultConfig.DBConfig.DBType)
				}
			}

			if dbHost == "" {
				dbHost = promptForInput(fmt.Sprintf("Database host [%s]: ", defaultConfig.DBConfig.DB.Host))
				if dbHost == "" {
					dbHost = defaultConfig.DBConfig.DB.Host
				}
			}

			if dbPort == 0 {
				dbPortStr := promptForInput(fmt.Sprintf("Database port [%d]: ", defaultConfig.DBConfig.DB.Port))
				if dbPortStr == "" {
					dbPort = defaultConfig.DBConfig.DB.Port
				} else {
					var err error
					dbPort, err = strconv.Atoi(dbPortStr)
					if err != nil {
						fmt.Printf("Invalid port number: %v\n", err)
						return
					}
				}
			}

			if dbUser == "" {
				dbUser = promptForInput(fmt.Sprintf("Database user [%s]: ", defaultConfig.DBConfig.DB.User))
				if dbUser == "" {
					dbUser = defaultConfig.DBConfig.DB.User
				}
			}

			if dbPassword == "" {
				dbPassword = promptForInput(fmt.Sprintf("Database password [%s]: ", defaultConfig.DBConfig.DB.Password))
				if dbPassword == "" {
					dbPassword = defaultConfig.DBConfig.DB.Password
				}
			}

			if dbName == "" {
				dbName = promptForInput(fmt.Sprintf("Database name [%s]: ", defaultConfig.DBConfig.DB.DBName))
				if dbName == "" {
					dbName = defaultConfig.DBConfig.DB.DBName
				}
			}

			if dbSchema == "" {
				dbSchema = promptForInput(fmt.Sprintf("Database schema [%s]: ", defaultConfig.DBConfig.DB.Schema))
				if dbSchema == "" {
					dbSchema = defaultConfig.DBConfig.DB.Schema
				}
			}

			if grpcPort == 0 {
				grpcPortStr := promptForInput(fmt.Sprintf("gRPC port [%d]: ", defaultConfig.GRPCConfig.GrpcPort))
				if grpcPortStr == "" {
					grpcPort = defaultConfig.GRPCConfig.GrpcPort
				} else {
					var err error
					grpcPort, err = strconv.Atoi(grpcPortStr)
					if err != nil {
						fmt.Printf("Invalid port number: %v\n", err)
						return
					}
				}
			}

			if gatewayPort == 0 {
				gatewayPortStr := promptForInput(fmt.Sprintf("Gateway port [%d]: ", defaultConfig.GRPCConfig.GatewayPort))
				if gatewayPortStr == "" {
					gatewayPort = defaultConfig.GRPCConfig.GatewayPort
				} else {
					var err error
					gatewayPort, err = strconv.Atoi(gatewayPortStr)
					if err != nil {
						fmt.Printf("Invalid port number: %v\n", err)
						return
					}
				}
			}

		} else {
			// Non-interactive mode: set default values if not provided
			defaultConfig := config.Local()
			if dbType == "" {
				dbType = string(defaultConfig.DBConfig.DBType)
			}
			if dbHost == "" {
				dbHost = defaultConfig.DBConfig.DB.Host
			}
			if dbPort == 0 {
				dbPort = defaultConfig.DBConfig.DB.Port
			}
			if dbUser == "" {
				dbUser = defaultConfig.DBConfig.DB.User
			}
			if dbPassword == "" {
				dbPassword = defaultConfig.DBConfig.DB.Password
			}
			if dbName == "" {
				dbName = defaultConfig.DBConfig.DB.DBName
			}
			if dbSchema == "" {
				dbSchema = defaultConfig.DBConfig.DB.Schema
			}
			if grpcPort == 0 {
				grpcPort = defaultConfig.GRPCConfig.GrpcPort
			}
			if gatewayPort == 0 {
				gatewayPort = defaultConfig.GRPCConfig.GatewayPort
			}
		}

		// Create config based on user input
		cfg = &config.Config{
			DBConfig: config.DBConfig{
				DBType: config.DBType(dbType),
				DB: config.DB{
					Host:     dbHost,
					Port:     dbPort,
					User:     dbUser,
					Password: dbPassword,
					DBName:   dbName,
					Schema:   dbSchema,
				},
			},
			GRPCConfig: config.GRPCConfig{
				GrpcPort:    grpcPort,
				GatewayPort: gatewayPort,
			},
			CacheConfig: config.CacheConfig{
				RedisAddrs: []string{
					"portfolio-backend-redis-node-0:6379",
					"portfolio-backend-redis-node-1:6379",
					"portfolio-backend-redis-node-2:6379",
					"portfolio-backend-redis-node-3:6379",
					"portfolio-backend-redis-node-4:6379",
					"portfolio-backend-redis-node-5:6379",
				},
			},
		}

		// Validate config
		if err := cfg.Validate(); err != nil {
			fmt.Printf("Error validating config: %v\n", err)
			return
		}

		// Create output directory if it doesn't exist
		outputDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			return
		}

		// Marshal config to YAML
		data, err := cfg.MarshalYAML()
		if err != nil {
			fmt.Printf("Error marshaling config: %v\n", err)
			return
		}

		// Write config to file
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Printf("Error writing config file: %v\n", err)
			return
		}

		fmt.Printf("Successfully created config file at: %s\n", outputPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Database configuration flags
	initCmd.Flags().StringVarP(&dbType, "db-type", "t", "", "Database type (postgres or test_postgres)")
	initCmd.Flags().StringVarP(&dbHost, "db-host", "H", "", "Database host")
	initCmd.Flags().IntVarP(&dbPort, "db-port", "p", 0, "Database port")
	initCmd.Flags().StringVarP(&dbUser, "db-user", "u", "", "Database user")
	initCmd.Flags().StringVarP(&dbPassword, "db-password", "w", "", "Database password")
	initCmd.Flags().StringVarP(&dbName, "db-name", "n", "", "Database name")
	initCmd.Flags().StringVarP(&dbSchema, "db-schema", "s", "public", "Database schema")
	initCmd.Flags().IntVarP(&grpcPort, "grpc-port", "g", 0, "gRPC server port")
	initCmd.Flags().StringVarP(&outputPath, "output", "o", ".config/auth.config.yaml", "Output config file path")
	initCmd.Flags().BoolVarP(&force, "force", "f", false, "Force overwrite of existing config file")
	initCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Run in interactive mode")
}

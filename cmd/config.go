package cmd

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SNS      SNSConfig      `yaml:"sns"`
	SQS      SQSConfig      `yaml:"sqs"`
	S3Bucket S3BucketConfig `yaml:"s3bucket"`
	Client   ClientConfig   `yaml:"client"`
}

type ClientConfig struct {
	PushFrequency string `yaml:"push-frequency"`
}

type S3BucketConfig struct {
	Region string `yaml:"region"`
	Name   string `yaml:"name"`
	S3Path string `yaml:"s3path"`
}

type SNSConfig struct {
	TopicARN string `yaml:"topic-arn"`
	Region   string `yaml:"region"`
}

type SQSConfig struct {
	Region   string `yaml:"region"`
	QueueARN string `yaml:"queue-arn"`
	QueueURL string `yaml:"queue-url"`
}

func NewDefaultConfig() *Config {
	return &Config{
		SNS: SNSConfig{
			TopicARN: "arn:aws:sns:us-west-2:123456789012:example-topic",
			Region:   "us-west-2",
		},
		SQS: SQSConfig{
			Region:   "us-west-2",
			QueueARN: "arn:aws:sqs:us-west-2:123456789012",
			QueueURL: "https://sqs.us-west-2.amazonaws.com/193048895737/somename",
		},
		S3Bucket: S3BucketConfig{
			Region: "us-west-2",
			Name:   "mybucket",
			S3Path: "jetplam.yaml",
		},
		Client: ClientConfig{
			PushFrequency: "1m",
		},
	}
}

// SetDefaultConfigValues sets the default values for the Config struct
func SetDefaultConfigValues(cfg *Config) {
	defaultConfig := NewDefaultConfig()
	cfg.SNS = defaultConfig.SNS
	cfg.SQS = defaultConfig.SQS
	cfg.S3Bucket = defaultConfig.S3Bucket
	cfg.Client = defaultConfig.Client
}

// WriteDefaultConfigToFile writes the default Config struct to a YAML file if the file doesn't exist
func WriteDefaultConfigToFile(cfg *Config, filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		defaultConfigData, err := yaml.Marshal(cfg)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %v", err)
		}

		err = os.WriteFile(filename, defaultConfigData, 0o644)
		if err != nil {
			return fmt.Errorf("failed to write default config file: %v", err)
		}
	}
	return nil
}

func ReadConfigFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

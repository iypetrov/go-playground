package sdnotify

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

// typeStr is the YAML key used to reference this extension in collector configs.
const typeStr = "sdnotify"

// NewFactory returns the factory that the OCB-generated components.go registers.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		component.MustNewType(typeStr),
		createDefaultConfig,
		createExtension,
		component.StabilityLevelAlpha,
	)
}

func createDefaultConfig() component.Config {
	return &Config{FailIfNotSupervised: false}
}

func createExtension(_ context.Context, set extension.Settings, cfg component.Config) (extension.Extension, error) {
	return newSDNotify(cfg.(*Config), set.Logger), nil
}

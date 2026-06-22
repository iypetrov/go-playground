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
	// All advanced features off by default: today's "just send READY=1 and
	// STOPPING=1 if NOTIFY_SOCKET is set" behaviour is preserved bit-for-bit.
	// HealthcheckV2 defaults to the conventional sibling ID so that flipping
	// `deep_healthcheck: true` works without also restating the dependency.
	return &Config{
		HealthcheckV2: component.MustNewID("healthcheckv2"),
	}
}

func createExtension(_ context.Context, set extension.Settings, cfg component.Config) (extension.Extension, error) {
	return newSDNotify(cfg.(*Config), set.Logger), nil
}

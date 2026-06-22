package sdnotify

import (
	"testing"

	"go.opentelemetry.io/collector/component"
)

func TestConfig_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name:    "default zero value is valid",
			cfg:     Config{},
			wantErr: false,
		},
		{
			name: "deep_healthcheck without healthcheckv2 is invalid",
			cfg: Config{
				DeepHealthcheck: true,
			},
			wantErr: true,
		},
		{
			name: "deep_healthcheck with healthcheckv2 is valid",
			cfg: Config{
				DeepHealthcheck: true,
				HealthcheckV2:   component.MustNewID("healthcheckv2"),
			},
			wantErr: false,
		},
		{
			name: "watchdog and unset_environment don't require healthcheckv2",
			cfg: Config{
				EnableWatchdog:   true,
				UnsetEnvironment: true,
			},
			wantErr: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.cfg.Validate()
			if (err != nil) != tc.wantErr {
				t.Fatalf("Validate(): err=%v, wantErr=%v", err, tc.wantErr)
			}
		})
	}
}

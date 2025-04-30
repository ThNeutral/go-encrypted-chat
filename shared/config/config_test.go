package config_test

import (
	"chat/shared/config"
	"chat/shared/utils"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantType  config.AppType
		wantAddr  *net.TCPAddr
		expectErr bool
	}{
		{
			name:     "valid server config",
			args:     []string{"--server", "--addr=localhost:8080"},
			wantType: config.SERVER,
			wantAddr: utils.Must(net.ResolveTCPAddr("tcp", "localhost:8080")),
		},
		{
			name:     "valid client config",
			args:     []string{"--client", "--addr=127.0.0.1:9090"},
			wantType: config.CLIENT,
			wantAddr: utils.Must(net.ResolveTCPAddr("tcp", "127.0.0.1:9090")),
		},
		{
			name:      "missing app type",
			args:      []string{"--addr=localhost:9090"},
			expectErr: true,
		},
		{
			name:      "missing address",
			args:      []string{"--server"},
			expectErr: true,
		},
		{
			name:      "both client and server",
			args:      []string{"--client", "--server", "--addr=localhost:9090"},
			expectErr: true,
		},
		{
			name:      "invalid address",
			args:      []string{"--client", "--addr=%%%"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := config.NewConfig()
			require.NotNil(t, inst)
			err := inst.Init(tt.args)

			if tt.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, inst.Type(), tt.wantType)

			if tt.wantAddr != nil {
				require.Equal(t, inst.Addr(), tt.wantAddr)
			}
		})
	}
}

package redisfsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitaliy-ukiru/fsm-telebot/v2"
)

func TestStorage_generateKey(t *testing.T) {
	s := &Storage{pref: StorageSettings{Prefix: "test"}}
	type args struct {
		target  fsm.StorageKey
		keyType keyType
		keys    []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "simple key",
			args: args{
				keyType: dataKey,
				keys:    []string{"myKey"},
			},
			want: "test:0:0:0:data:myKey",
		},
		{
			name: "multiple key",
			args: args{
				keyType: dataKey,
				keys:    []string{"multiple", "keys"},
			},
			want: "test:0:0:0:data:multiple:keys",
		},
		{
			name: "with thread",
			args: args{
				keyType: stateKey,
				target: fsm.StorageKey{
					ThreadID: 1,
				},
			},
			want: "test:0:0:0:1:state",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(
				t,
				tt.want,
				s.generateKey(tt.args.target, tt.args.keyType, tt.args.keys...),
				"generateKey(%v, %v, %v)",
				tt.args.target,
				tt.args.keyType,
				tt.args.keys,
			)
		})
	}
	t.Run("key type", func(t *testing.T) {
		assert.NotEqual(
			t,
			s.generateKey(fsm.StorageKey{}, stateKey),
			s.generateKey(fsm.StorageKey{}, dataKey, "state"),
			"state key and data[state] equals",
		)
	})
}

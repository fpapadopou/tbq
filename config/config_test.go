package config

import (
	"errors"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		envs      map[string]string
		want      Config
		wantError error
	}{
		{
			"no env set",
			nil,
			Config{
				SourceType:    "redis",
				RedisHost:     "localhost",
				RedisPort:     "6379",
				RedisPassword: "",
				RedisDatabase: 0,
			},
			nil,
		},
		{
			"env set",
			map[string]string{
				"TBQ_SOURCE": "foo",
			},
			Config{
				SourceType:    "foo",
				RedisHost:     "localhost",
				RedisPort:     "6379",
				RedisPassword: "",
				RedisDatabase: 0,
			},
			nil,
		},
		{
			"invalid env",
			map[string]string{
				"TBQ_REDIS_DB": "not a number",
			},
			Config{},
			errors.New("strconv.Atoi: parsing \"not a number\": invalid syntax"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up env.
			os.Clearenv()
			for k, v := range tt.envs {
				err := os.Setenv(k, v)
				if err != nil {
					t.Errorf("setting env %s to value %s failed: %v", k, v, err)
				}
			}
			// Execute.
			got, gotError := New()
			if got != tt.want {
				t.Errorf("test failed for New() - got %v, want %v", got, tt.want)
			}

			// Validate.
			if gotError != nil && gotError.Error() != tt.wantError.Error() {
				t.Errorf("test failed for New() - got error %v, want error %v", gotError, tt.wantError)
			}
		})
	}

}

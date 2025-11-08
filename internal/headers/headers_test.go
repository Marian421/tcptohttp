package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		want     Headers
		wantN    int
		wantDone bool
		wantErr  bool
	}{
		{
			name:     "Valid single header",
			input:    []byte("Host: localhost:42069\r\n\r\n"),
			want:     map[string]string{"host": "localhost:42069"},
			wantN:    25,
			wantDone: true,
			wantErr:  false,
		}, {
			name:     "Spacing before header",
			input:    []byte("  Host: localhost:42069\r\n\r\n"),
			want:     nil,
			wantN:    0,
			wantDone: false,
			wantErr:  true,
		}, {
			name:     "Spacing between header and colon",
			input:    []byte("host : localhost:42069\r\n\r\n"),
			want:     nil,
			wantN:    0,
			wantDone: false,
			wantErr:  true,
		}, {
			name:     "Trims the whitespaces",
			input:    []byte("Host:    localhost:42069   \r\n\r\n"),
			want:     map[string]string{"host": "localhost:42069"},
			wantN:    31,
			wantDone: true,
			wantErr:  false,
		}, {
			name:     "Not done with the parsing",
			input:    []byte("Host:    localhost:42069   \r\n"),
			want:     map[string]string{"host": "localhost:42069"},
			wantN:    29,
			wantDone: false,
			wantErr:  false,
		}, {
			name:     "Validates field-name charachters",
			input:    []byte("H©st: localhost:42069\r\n\r\n"),
			want:     nil,
			wantN:    0,
			wantDone: false,
			wantErr:  true,
		}, {
			name:     "Multiple values for a header",
			input:    []byte("Set-person: person1\r\nSet-person: person2\r\n\r\n"),
			want:     map[string]string{"set-person": "person1, person2"},
			wantN:    44,
			wantDone: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHeaders()
			n, done, err := h.Parse(tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, tt.wantN, n)
				assert.Equal(t, tt.wantDone, done)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, h)
			assert.Equal(t, tt.want, h)
			assert.Equal(t, tt.wantN, n)
			assert.Equal(t, tt.wantDone, done)
		})
	}
}

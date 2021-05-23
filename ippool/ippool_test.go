package ippool

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestIPPool(t *testing.T) {
	require := require.New(t)
	tests := []*Target{
		{
			Host:     "8.8.8.8",
			Port:     53,
			Counter:  3,
			Interval: time.Second,
			Timeout:  time.Second,
		},
		{
			Host:     "1.1.1.1",
			Port:     53,
			Counter:  2,
			Interval: time.Second,
			Timeout:  time.Second,
		},
		{
			Host:     "2.2.2.2",
			Port:     53,
			Counter:  1,
			Interval: time.Second,
			Timeout:  time.Second * 2,
		},
	}
	p := NewIPPool(time.Second * 8)
	for _, test := range tests {
		p.AddTarget(test)
	}
	require.Equal(len(p.targets), len(tests))
	require.Equal(len(p.Results()), 0)

	ctx := context.Background()
	p.Start(ctx)
	time.Sleep(time.Second * 10)
	p.Reset()
	require.Equal(len(p.targets), 0)
	require.Equal(len(p.Results()), len(tests))

	tests[0].Host = "114.114.114.114"
	tests[1].Host = "5.5.5.5"
	tests[2].Host = "1.2.4.8"
	for _, test := range tests {
		p.AddTarget(test)
	}
	time.Sleep(time.Second * 15)
	fastest, err := p.GetFastestTarget()
	require.NoError(err)
	require.Equal(fastest, tests[0])
	random, err := p.GetRandomTarget()
	require.NoError(err)
	require.Equal(random, tests[0])
	p.Stop(ctx)
}

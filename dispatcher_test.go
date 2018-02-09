package workerpool

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type WpTestSuite struct {
	suite.Suite
}

func TestWpTestSuite(t *testing.T) {
	suite.Run(t, new(WpTestSuite))
}

func (suite *WpTestSuite) TestPerformsWork() {
	var workers uint32 = 10
	values := make([]*int, workers)
	for i := 0; i < len(values); i++ {
		v := 0
		values[i] = &v
	}

	d := NewDispatcher(workers)
	d.Start()

	for _, i := range values {
		t := i
		d.Submit(func() {
			doWork(t)
		})
	}

	d.Stop()

	for _, i := range values {
		fmt.Println(*i)
		require.Equal(suite.T(), 1, *i, "Work was not finished")
	}
}

func doWork(i *int) {
	*i = 1
}

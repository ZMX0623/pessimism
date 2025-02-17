package service_test

import (
	"context"
	"fmt"
	"testing"

	svc "github.com/base-org/pessimism/internal/api/service"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetHealth(t *testing.T) {
	ctrl := gomock.NewController(t)

	var tests = []struct {
		name        string
		description string
		function    string

		constructionLogic func() testSuite
		testLogic         func(*testing.T, testSuite)
	}{
		{
			name:        "Get Health Success",
			description: "",
			function:    "ProcessInvariantRequest",

			constructionLogic: func() testSuite {
				cfg := svc.Config{}
				ts := createTestSuite(ctrl, cfg)

				ts.mockEthClientInterface.EXPECT().
					DialContext(context.Background(), gomock.Any()).
					Return(nil).
					AnyTimes()

				ts.mockService.EXPECT().
					CheckETHRPCHealth(gomock.Any()).
					Return(true).
					AnyTimes()

				ts.mockEthClientInterface.EXPECT().
					HeaderByNumber(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes()

				return ts
			},

			testLogic: func(t *testing.T, ts testSuite) {
				hc := ts.apiSvc.CheckHealth()

				assert.True(t, hc.Healthy)
				assert.True(t, hc.ChainConnectionStatus.IsL2Healthy)
				assert.True(t, hc.ChainConnectionStatus.IsL1Healthy)

			},
		},
		{
			name:        "Get Unhealthy Response",
			description: "Emulates unhealthy rpc endpoints",
			function:    "ProcessInvariantRequest",

			constructionLogic: func() testSuite {
				cfg := svc.Config{}
				ts := createTestSuite(ctrl, cfg)

				ts.mockEthClientInterface.EXPECT().
					DialContext(gomock.Any(), gomock.Any()).
					Return(testErr1()).
					AnyTimes()

				ts.mockService.EXPECT().
					CheckETHRPCHealth(gomock.Any()).
					Return(false).
					AnyTimes()

				ts.mockEthClientInterface.EXPECT().
					HeaderByNumber(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes()

				return ts
			},

			testLogic: func(t *testing.T, ts testSuite) {
				hc := ts.apiSvc.CheckHealth()
				assert.False(t, hc.Healthy)
				assert.False(t, hc.ChainConnectionStatus.IsL2Healthy)
				assert.False(t, hc.ChainConnectionStatus.IsL1Healthy)
			},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			testMeta := tc.constructionLogic()
			tc.testLogic(t, testMeta)
		})

	}

}

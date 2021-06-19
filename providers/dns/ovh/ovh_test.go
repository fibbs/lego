package ovh

import (
	"testing"
	"time"

	"github.com/fibbs/lego/v4/platform/tester"
	"github.com/stretchr/testify/require"
)

const envDomain = envNamespace + "DOMAIN"

var envTest = tester.NewEnvTest(
	EnvEndpoint,
	EnvApplicationKey,
	EnvApplicationSecret,
	EnvConsumerKey).
	WithDomain(envDomain)

func TestNewDNSProvider(t *testing.T) {
	testCases := []struct {
		desc     string
		envVars  map[string]string
		expected string
	}{
		{
			desc: "success",
			envVars: map[string]string{
				EnvEndpoint:          "ovh-eu",
				EnvApplicationKey:    "B",
				EnvApplicationSecret: "C",
				EnvConsumerKey:       "D",
			},
		},
		{
			desc: "missing credentials",
			envVars: map[string]string{
				EnvEndpoint:          "",
				EnvApplicationKey:    "",
				EnvApplicationSecret: "",
				EnvConsumerKey:       "",
			},
			expected: "ovh: some credentials information are missing: OVH_ENDPOINT,OVH_APPLICATION_KEY,OVH_APPLICATION_SECRET,OVH_CONSUMER_KEY",
		},
		{
			desc: "missing endpoint",
			envVars: map[string]string{
				EnvEndpoint:          "",
				EnvApplicationKey:    "B",
				EnvApplicationSecret: "C",
				EnvConsumerKey:       "D",
			},
			expected: "ovh: some credentials information are missing: OVH_ENDPOINT",
		},
		{
			desc: "missing invalid endpoint",
			envVars: map[string]string{
				EnvEndpoint:          "foobar",
				EnvApplicationKey:    "B",
				EnvApplicationSecret: "C",
				EnvConsumerKey:       "D",
			},
			expected: "ovh: unknown endpoint 'foobar', consider checking 'Endpoints' list of using an URL",
		},
		{
			desc: "missing application key",
			envVars: map[string]string{
				EnvEndpoint:          "ovh-eu",
				EnvApplicationKey:    "",
				EnvApplicationSecret: "C",
				EnvConsumerKey:       "D",
			},
			expected: "ovh: some credentials information are missing: OVH_APPLICATION_KEY",
		},
		{
			desc: "missing application secret",
			envVars: map[string]string{
				EnvEndpoint:          "ovh-eu",
				EnvApplicationKey:    "B",
				EnvApplicationSecret: "",
				EnvConsumerKey:       "D",
			},
			expected: "ovh: some credentials information are missing: OVH_APPLICATION_SECRET",
		},
		{
			desc: "missing consumer key",
			envVars: map[string]string{
				EnvEndpoint:          "ovh-eu",
				EnvApplicationKey:    "B",
				EnvApplicationSecret: "C",
				EnvConsumerKey:       "",
			},
			expected: "ovh: some credentials information are missing: OVH_CONSUMER_KEY",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			defer envTest.RestoreEnv()
			envTest.ClearEnv()

			envTest.Apply(test.envVars)

			p, err := NewDNSProvider()

			if test.expected == "" {
				require.NoError(t, err)
				require.NotNil(t, p)
				require.NotNil(t, p.config)
				require.NotNil(t, p.client)
				require.NotNil(t, p.recordIDs)
			} else {
				require.EqualError(t, err, test.expected)
			}
		})
	}
}

func TestNewDNSProviderConfig(t *testing.T) {
	testCases := []struct {
		desc              string
		apiEndpoint       string
		applicationKey    string
		applicationSecret string
		consumerKey       string
		expected          string
	}{
		{
			desc:              "success",
			apiEndpoint:       "ovh-eu",
			applicationKey:    "B",
			applicationSecret: "C",
			consumerKey:       "D",
		},
		{
			desc:     "missing credentials",
			expected: "ovh: credentials missing",
		},
		{
			desc:              "missing api endpoint",
			apiEndpoint:       "",
			applicationKey:    "B",
			applicationSecret: "C",
			consumerKey:       "D",
			expected:          "ovh: credentials missing",
		},
		{
			desc:              "missing invalid api endpoint",
			apiEndpoint:       "foobar",
			applicationKey:    "B",
			applicationSecret: "C",
			consumerKey:       "D",
			expected:          "ovh: unknown endpoint 'foobar', consider checking 'Endpoints' list of using an URL",
		},
		{
			desc:              "missing application key",
			apiEndpoint:       "ovh-eu",
			applicationKey:    "",
			applicationSecret: "C",
			consumerKey:       "D",
			expected:          "ovh: credentials missing",
		},
		{
			desc:              "missing application secret",
			apiEndpoint:       "ovh-eu",
			applicationKey:    "B",
			applicationSecret: "",
			consumerKey:       "D",
			expected:          "ovh: credentials missing",
		},
		{
			desc:              "missing consumer key",
			apiEndpoint:       "ovh-eu",
			applicationKey:    "B",
			applicationSecret: "C",
			consumerKey:       "",
			expected:          "ovh: credentials missing",
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			config := NewDefaultConfig()
			config.APIEndpoint = test.apiEndpoint
			config.ApplicationKey = test.applicationKey
			config.ApplicationSecret = test.applicationSecret
			config.ConsumerKey = test.consumerKey

			p, err := NewDNSProviderConfig(config)

			if test.expected == "" {
				require.NoError(t, err)
				require.NotNil(t, p)
				require.NotNil(t, p.config)
				require.NotNil(t, p.client)
				require.NotNil(t, p.recordIDs)
			} else {
				require.EqualError(t, err, test.expected)
			}
		})
	}
}

func TestLivePresent(t *testing.T) {
	if !envTest.IsLiveTest() {
		t.Skip("skipping live test")
	}

	envTest.RestoreEnv()
	provider, err := NewDNSProvider()
	require.NoError(t, err)

	err = provider.Present(envTest.GetDomain(), "", "123d==")
	require.NoError(t, err)
}

func TestLiveCleanUp(t *testing.T) {
	if !envTest.IsLiveTest() {
		t.Skip("skipping live test")
	}

	envTest.RestoreEnv()
	provider, err := NewDNSProvider()
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	err = provider.CleanUp(envTest.GetDomain(), "", "123d==")
	require.NoError(t, err)
}

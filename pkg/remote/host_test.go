package remote

import "testing"

func TestParseHostError(t *testing.T) {
  testCases := []string{
    "localhost:",
    "[::1%2]:",
    "[::1%2",
    "user@",
    "user@:23",
    "port:IO",
    "",
  }

  for _, input := range testCases {
    if got, err := ParseHost(input); err == nil {
      t.Fatalf("expected error, got '%v'", got)
    }
  }
}

func TestParseHost(t *testing.T) {
  testCases := map[string]Host{
    "localhost": Host{
      Host: "localhost",
    },
    "only.host": Host{
      Host: "only.host",
    },
    "host-with.port:23": Host{
      Host: "host-with.port",
      Port: 23,
    },
    "127.0.0.1": Host{
      Host: "127.0.0.1",
    },
    "127.1.2.3:23": Host{
      Host: "127.1.2.3",
      Port: 23,
    },
    "[::1%2]": Host{
      Host: "::1%2",
    },
    "[2001:0db8:85a3:0000:0000:8a2e:0370:7334%2]:23": Host{
      Host: "2001:0db8:85a3:0000:0000:8a2e:0370:7334%2",
      Port: 23,
    },
    "user@localhost": Host{
      User: "user",
      Host: "localhost",
    },
    "user@user.host": Host{
      User: "user",
      Host: "user.host",
    },
    "user@user-host.port:23": Host{
      User: "user",
      Host: "user-host.port",
      Port: 23,
    },
    "user@127.2.4.6": Host{
      User: "user",
      Host: "127.2.4.6",
    },
    "user@127.127.127.127:23": Host{
      User: "user",
      Host: "127.127.127.127",
      Port: 23,
    },
    "user@[::ffff%2]": Host{
      User: "user",
      Host: "::ffff%2",
    },
    "user@[::f00d%2]:23": Host{
      User: "user",
      Host: "::f00d%2",
      Port: 23,
    },
  }

  for input, want := range testCases {
    got, err := ParseHost(input)

    if err != nil {
      t.Fatalf("ParseHost: %v", err)
    }

    compare(t, &want, got)
  }
}

func compare(t *testing.T, want, got *Host) {
	if want.User != got.User {
		t.Errorf("got User=%q; want %q", got.User, want.User)
	}

	if want.Port != got.Port {
		t.Errorf("got Port=%q; want %q", got.Port, want.Port)
	}

	if want.Host != got.Host {
		t.Errorf("got Host=%q; want %q", got.Host, want.Host)
	}
}

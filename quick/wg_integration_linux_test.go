package quick

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func TestWgBinAppliesInterfaceAndDeviceConfig(t *testing.T) {
	wgBin := os.Getenv("WG_QUICK_OP_TEST_WGBIN")
	if wgBin == "" {
		t.Skip("WG_QUICK_OP_TEST_WGBIN is not set")
	}

	iface := fmt.Sprintf("wqop%d", os.Getpid())
	if len(iface) > 15 {
		iface = iface[:15]
	}
	_, err := netlink.LinkByName(iface)
	require.Error(t, err, "temporary interface already exists: %s", iface)

	privateKey, err := wgtypes.GeneratePrivateKey()
	require.NoError(t, err)
	peerPrivateKey, err := wgtypes.GeneratePrivateKey()
	require.NoError(t, err)
	presharedKey, err := wgtypes.GeneratePrivateKey()
	require.NoError(t, err)

	configText := fmt.Sprintf(`[Interface]
PrivateKey = %s
ListenPort = 0
FwMark = 0x2345
MTU = 1377
Table = off
WgBin = %s

[Peer]
PublicKey = %s
PresharedKey = %s
AllowedIPs = 192.0.2.1/32, 2001:db8::1/128
PersistentKeepalive = 17
`, privateKey.String(), wgBin, peerPrivateKey.PublicKey().String(), presharedKey.String())
	cfg := &Config{}
	require.NoError(t, cfg.UnmarshalText([]byte(configText)))

	t.Cleanup(func() {
		link, linkErr := netlink.LinkByName(iface)
		if linkErr == nil {
			require.NoError(t, netlink.LinkDel(link))
		}
	})

	require.NoError(t, Sync(cfg, iface, zerolog.Nop()))
	link, err := netlink.LinkByName(iface)
	require.NoError(t, err)
	require.Equal(t, cfg.MTU, link.Attrs().MTU)

	device, err := client.Device(iface)
	require.NoError(t, err)
	require.Equal(t, 0x2345, device.FirewallMark)
	require.Equal(t, privateKey.PublicKey(), device.PublicKey)
	require.Len(t, device.Peers, 1)
	require.Equal(t, peerPrivateKey.PublicKey(), device.Peers[0].PublicKey)
	require.Equal(t, presharedKey, device.Peers[0].PresharedKey)
	require.Equal(t, "17s", device.Peers[0].PersistentKeepaliveInterval.String())
	require.Len(t, device.Peers[0].AllowedIPs, 2)
	require.Equal(t, "192.0.2.1/32", device.Peers[0].AllowedIPs[0].String())
	require.Equal(t, "2001:db8::1/128", device.Peers[0].AllowedIPs[1].String())
}

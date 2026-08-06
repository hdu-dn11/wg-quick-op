package quick

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
)

func TestSyncLinkMTUSetsConfiguredMTU(t *testing.T) {
	oldLinkSetMTU := linkSetMTU
	t.Cleanup(func() {
		linkSetMTU = oldLinkSetMTU
	})

	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{
		Name: "wg-test",
		MTU:  1280,
	}}
	var called int
	linkSetMTU = func(got netlink.Link, mtu int) error {
		called++
		assert.Equal(t, link, got)
		assert.Equal(t, 1420, mtu)
		return nil
	}

	err := SyncLinkMTU(&Config{MTU: 1420}, link, zerolog.Nop())

	assert.NoError(t, err)
	assert.Equal(t, 1, called)
	assert.Equal(t, 1420, link.Attrs().MTU)
}

func TestSyncLinkMTUSkipsWhenAlreadyConfigured(t *testing.T) {
	oldLinkSetMTU := linkSetMTU
	t.Cleanup(func() {
		linkSetMTU = oldLinkSetMTU
	})

	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{
		Name: "wg-test",
		MTU:  1420,
	}}
	linkSetMTU = func(netlink.Link, int) error {
		t.Fatal("LinkSetMTU should not be called when MTU is already configured")
		return nil
	}

	err := SyncLinkMTU(&Config{MTU: 1420}, link, zerolog.Nop())

	assert.NoError(t, err)
}

func TestSyncLinkMTUSkipsWhenUnset(t *testing.T) {
	oldLinkSetMTU := linkSetMTU
	t.Cleanup(func() {
		linkSetMTU = oldLinkSetMTU
	})

	link := &netlink.Dummy{LinkAttrs: netlink.LinkAttrs{
		Name: "wg-test",
		MTU:  1280,
	}}
	linkSetMTU = func(netlink.Link, int) error {
		t.Fatal("LinkSetMTU should not be called when MTU is unset")
		return nil
	}

	err := SyncLinkMTU(&Config{}, link, zerolog.Nop())

	assert.NoError(t, err)
	assert.Equal(t, 1280, link.Attrs().MTU)
}

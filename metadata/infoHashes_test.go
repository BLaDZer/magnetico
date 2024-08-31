package metadata

import (
	"net"
	"reflect"
	"testing"
)

func TestInfoHashes_Push(t *testing.T) {
	t.Parallel()

	ih := &infoHashes{
		infoHashes:  make(map[[20]byte][]net.TCPAddr),
		maxNLeeches: 2,
	}

	infoHash := [20]byte{1, 2, 3, 4, 5, 6}
	peerAddresses := []net.TCPAddr{
		{IP: net.ParseIP("1.0.0.1"), Port: 443},
		{IP: net.ParseIP("1.0.0.2"), Port: 1337},
		{IP: net.ParseIP("1.0.0.2"), Port: 1337},
		{IP: net.ParseIP("1.0.0.3"), Port: 6969},
		{IP: net.ParseIP("1.0.0.4"), Port: 8080},
	}

	ih.push(infoHash, peerAddresses)

	expected := []net.TCPAddr{
		{IP: net.ParseIP("1.0.0.2"), Port: 1337},
		{IP: net.ParseIP("1.0.0.3"), Port: 6969},
	}

	actual := ih.infoHashes[infoHash]

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected infoHashes[%v] to be %v, but got %v", infoHash, expected, actual)
	}
}

func TestInfoHashes_Pop(t *testing.T) {
	t.Parallel()

	ih := &infoHashes{
		infoHashes:  make(map[[20]byte][]net.TCPAddr),
		maxNLeeches: 2,
	}

	infoHash := [20]byte{1, 2, 3, 4, 5, 6}
	peerAddresses := []net.TCPAddr{
		{IP: net.ParseIP("1.0.0.1"), Port: 443},
		{IP: net.ParseIP("1.0.0.2"), Port: 1337},
		{IP: net.ParseIP("1.0.0.3"), Port: 6969},
	}

	ih.infoHashes[infoHash] = peerAddresses

	expected := &peerAddresses[0]
	actual := ih.pop(infoHash)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected pop(%v) to return %v, but got %v", infoHash, expected, actual)
	}

	if len(ih.infoHashes[infoHash]) != 2 {
		t.Errorf("Expected infoHashes[%v] to have length 2, but got %d", infoHash, len(ih.infoHashes[infoHash]))
	}
}

func TestInfoHashes_Pop_Nil(t *testing.T) {
	t.Parallel()

	ih := &infoHashes{
		infoHashes:  make(map[[20]byte][]net.TCPAddr),
		maxNLeeches: 2,
	}

	infoHash := [20]byte{1, 2, 3, 4, 5, 6}

	actual := ih.pop(infoHash)

	if actual != nil {
		t.Errorf("Expected pop(%v) to return nil, but got %v", infoHash, actual)
	}
}

func Test_infoHashes_isAllowedFilter(t *testing.T) {
	t.Parallel()

	_, ipnet, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		t.Error(err)
	}
	ih := &infoHashes{
		filterPeers: []net.IPNet{
			*ipnet,
		},
	}

	tests := []struct {
		peer net.TCPAddr
		want bool
	}{
		{
			peer: net.TCPAddr{
				IP:   net.ParseIP("127.0.0.1"),
				Port: 5678,
			},
			want: true,
		},
		{
			peer: net.TCPAddr{
				IP:   net.ParseIP("192.168.1.1"),
				Port: 6789,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.peer.String(), func(t *testing.T) {
			if got := ih.isAllowed(tt.peer); got != tt.want {
				t.Errorf("infoHashes.isAllowed() = %v, want %v", got, tt.want)
			}
		})
	}
}

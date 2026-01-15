package content

import (
	"encoding/binary"
	"time"

	"github.com/skerkour/stdx-go/crypto/blake3"
)

func HashPageMetadata(pageType PageType, path string, date time.Time, sendAsNewsletter bool, language string, title string, description string, tags []string) [32]byte {
	var hash [32]byte

	hasher := blake3.New(32, nil)
	hasher.Write([]byte(pageType))
	hasher.Write([]byte(path))
	binary.Write(hasher, binary.LittleEndian, date.UnixMilli())
	binary.Write(hasher, binary.LittleEndian, sendAsNewsletter)
	hasher.Write([]byte(language))
	hasher.Write([]byte(title))
	hasher.Write([]byte(description))
	for _, tag := range tags {
		hasher.Write([]byte(tag))
	}

	hasher.Sum(hash[:0])

	return hash
}

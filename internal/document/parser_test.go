package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsedSource_Blocks(t *testing.T) {
	data := []byte(`Test involving nested blocks:

1. Install the linkerd CLI

    ` + "```" + `bash
    curl https://run.linkerd.io/install | sh
    ` + "```" + `

1. Install Linkerd2

    ` + "```" + `bash
    linkerd install | kubectl apply -f -
    ` + "```" + `
`)
	source := NewSource(data)
	blocks := source.Parse().Blocks()
	assert.Equal(t, 5, len(blocks))
}

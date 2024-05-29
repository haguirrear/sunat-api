package main

import (
	_ "embed"

	"github.com/haguirrear/sunatapi/cmd"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante/consultar"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante/enviar"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante/procesar"
)

//go:embed version
var version string

func main() {
	cmd.Execute(version)
}

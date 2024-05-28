package main

import (
	"github.com/haguirrear/sunatapi/cmd"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante/consultar"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante/enviar"
	_ "github.com/haguirrear/sunatapi/cmd/comprobante/procesar"
)

func main() {
	cmd.Execute()
}

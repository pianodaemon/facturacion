package dal

import (
	"bytes"
	"encoding/json"
)

type (

	// Define the Receptor struct
	Receptor struct {
		UID string `json:"UID"`
	}

	// Define the Traslado struct (for taxes)
	Traslado struct {
		Base       float64 `json:"Base"`
		Impuesto   string  `json:"Impuesto"`
		TipoFactor string  `json:"TipoFactor"`
		TasaOCuota string  `json:"TasaOCuota"`
		Importe    float64 `json:"Importe"`
	}

	// Define the Local struct (for local taxes)
	Local struct {
		Base       float64 `json:"Base"`
		Impuesto   string  `json:"Impuesto"`
		TipoFactor string  `json:"TipoFactor"`
		TasaOCuota string  `json:"TasaOCuota"`
		Importe    float64 `json:"Importe"`
	}

	// Define the Impuestos struct (for both Traslados and Locales)
	Impuestos struct {
		Traslados []Traslado `json:"Traslados"`
		Locales   []Local    `json:"Locales"`
	}

	// Define the Concepto struct
	Concepto struct {
		ClaveProdServ string    `json:"ClaveProdServ"`
		Cantidad      float64   `json:"Cantidad"`
		ClaveUnidad   string    `json:"ClaveUnidad"`
		Unidad        string    `json:"Unidad"`
		ValorUnitario float64   `json:"ValorUnitario"`
		Descripcion   string    `json:"Descripcion"`
		Impuestos     Impuestos `json:"Impuestos"`
	}

	// Define the Payload struct
	Payload struct {
		Receptor      Receptor   `json:"Receptor"`
		TipoDocumento string     `json:"TipoDocumento"`
		Conceptos     []Concepto `json:"Conceptos"`
		UsoCFDI       string     `json:"UsoCFDI"`
		Serie         string     `json:"Serie"`
		FormaPago     string     `json:"FormaPago"`
		MetodoPago    string     `json:"MetodoPago"`
		Moneda        string     `json:"Moneda"`
		EnviarCorreo  bool       `json:"EnviarCorreo"`
	}

	IssueHeaders struct {
		FPlugin    string
		FApiKey    string
		FSecretKey string
	}
)

type PgSQLPayloadProvider struct {
	prefacturaID string
}

func NewPgSQLPayloadProvider(sourceID string) *PgSQLPayloadProvider {
	return &PgSQLPayloadProvider{
		prefacturaID: sourceID,
	}
}

func (p *PgSQLPayloadProvider) fetchReceptor() Receptor {
	return Receptor{
		UID: "66ec233077430",
	}
}

func (p *PgSQLPayloadProvider) fetchConceptos() []Concepto {
	return []Concepto{
		{
			ClaveProdServ: "81112101",
			Cantidad:      1,
			ClaveUnidad:   "E48",
			Unidad:        "Unidad de servicio",
			ValorUnitario: 229.9,
			Descripcion:   "Desarrollo a la medida",
			Impuestos: Impuestos{
				Traslados: []Traslado{
					{
						Base:       229.9,
						Impuesto:   "002",
						TipoFactor: "Tasa",
						TasaOCuota: "0.16",
						Importe:    36.784,
					},
				},
				Locales: []Local{
					{
						Base:       229.9,
						Impuesto:   "ISH",
						TipoFactor: "Tasa",
						TasaOCuota: "0.03",
						Importe:    6.897,
					},
				},
			},
		},
	}
}

func (p *PgSQLPayloadProvider) fetchMoneda() string {
	return "MXN"
}

func (p *PgSQLPayloadProvider) fetchFormaPago() string {
	return "03"
}

func (p *PgSQLPayloadProvider) fetchMetodoPago() string {
	return "PUE"
}

func (p *PgSQLPayloadProvider) fetchUsoCFDI() string {
	return "S01"
}

func (p *PgSQLPayloadProvider) fetchSerie() string {
	return "60659"
}

func (p *PgSQLPayloadProvider) DumpBuffer() (*bytes.Buffer, error) {

	payload := &Payload{
		Receptor:      p.fetchReceptor(),   // Call the interface method for "Receptor"
		TipoDocumento: "factura",           // Static value
		Conceptos:     p.fetchConceptos(),  // Call the interface method for "Conceptos"
		UsoCFDI:       p.fetchUsoCFDI(),    // Call the interface method for "UsoCFDI"
		Serie:         p.fetchSerie(),      // Call the interface method for "Serie"
		FormaPago:     p.fetchFormaPago(),  // Call the interface method for "FormaPago"
		MetodoPago:    p.fetchMetodoPago(), // Call the interface method for "MetodoPago"
		Moneda:        p.fetchMoneda(),     // Call the interface method for "Moneda"
		EnviarCorreo:  false,               // Static value
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(payloadBytes), nil
}

func (p *PgSQLPayloadProvider) GetIssueHeaders() (*IssueHeaders, error) {
	return &IssueHeaders{
		FPlugin:    "9d4095c8f7ed5785cb14c0e3b033eeb8252416ed",
		FApiKey:    "JDJ5JDEwJFlEUFR2aEExa0FyVHQxQWRmWkJvTE9nSlRWZlNDVjhVd3RGYWJIZmNjY1BNYUo5R3F0Mk5P",
		FSecretKey: "JDJ5JDEwJE5vOUM1OGo1RGJMT0k5cVpqVzkxbHVITU50N2RZNzh0ZzBoQXh6RGJlcVFiY0NjcTIxOGpH",
	}, nil
}

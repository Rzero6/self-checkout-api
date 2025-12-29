package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Rzero6/self-checkout-api/config"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type InvoiceData struct {
	OrderID     string
	PaymentType string
	Status      string
	Date        string
	Items       []InvoiceItem
	Total       string
}
type InvoiceItem struct {
	Name     string
	Price    string
	Qty      string
	Subtotal string
}

const headerSize = 10
const bodySize = 8

func GeneratePDF(data InvoiceData) (string, error) {

	logo := os.Getenv("APP_LOGO_PATH")
	email := os.Getenv("APP_EMAIL_ADDRESS")
	appName := os.Getenv("APP_NAME")

	cfg := config.SetupPDF()
	m := maroto.New(cfg)

	m.RegisterHeader(addHeader(logo, appName))
	m.AddRow(5)
	m.AddAutoRow(
		text.NewCol(12, "INVOICE", props.Text{
			Style: fontstyle.Bold,
			Align: align.Center,
			Size:  16,
		}),
	)
	m.AddRow(3)
	m.AddAutoRow(
		text.NewCol(12, data.OrderID,
			props.Text{Style: fontstyle.Bold, Align: align.Center, Size: 20}),
	)
	addDivider(m)
	addPaymentDetail(m, data)
	m.AddRow(10)
	addItemList(m, data)
	m.RegisterFooter(addFooter(email))
	document, err := m.Generate()
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("invoice_%s.pdf", data.OrderID)
	filePath := filepath.Join(os.TempDir(), fileName)

	err = document.Save(filePath)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// Row in milimeter
// Col in 12 grid bootstrap system
func addHeader(logo, appName string) core.Row {
	return row.New(20).Add(
		image.NewFromFileCol(
			6,
			logo,
			props.Rect{Center: false, Percent: 100},
		),
		col.New(6).Add(
			text.New(appName, props.Text{
				Align: align.Right,
				Size:  10,
			}),
		),
	)
}

func addDivider(m core.Maroto) {
	m.AddRow(5)
	m.AddRow(2, line.NewCol(12, props.Line{Color: getDarkEmeraldColor(), Thickness: 1}))
	m.AddRow(5)
}

func addPaymentDetail(m core.Maroto, data InvoiceData) {
	header := row.New(headerSize).Add(
		col.New(1),
		text.NewCol(4, "PAYMENT METHOD", props.Text{
			Align: align.Left, Size: 12, Style: fontstyle.Bold,
			Color: getWhiteColor(),
		}),
		text.NewCol(3, "STATUS", props.Text{
			Align: align.Center, Size: 12, Style: fontstyle.Bold,
			Color: getWhiteColor(),
		}),
		text.NewCol(3, "DATE", props.Text{
			Align: align.Right, Size: 12, Style: fontstyle.Bold,
			Color: getWhiteColor(),
		}),
		col.New(1),
	)
	header.WithStyle(&props.Cell{
		BackgroundColor: getDarkEmeraldColor(),
	})

	body := row.New(bodySize).Add(
		col.New(1),
		text.NewCol(4, data.PaymentType, props.Text{Align: align.Left, Size: 11}),
		text.NewCol(3, data.Status, props.Text{Align: align.Center, Size: 11}),
		text.NewCol(3, data.Date, props.Text{Align: align.Right, Size: 11}),
		col.New(1),
	)
	body.WithStyle(&props.Cell{
		BackgroundColor: &props.Color{Red: 224, Green: 245, Blue: 240},
	})

	m.AddRows(header)
	m.AddRows(body)
}

func addItemList(m core.Maroto, data InvoiceData) {
	header := data.GetHeader()
	m.AddRows(header)

	for i, item := range data.Items {
		row := item.GetContent(i)
		m.AddRows(row)
	}

	m.AddRow(5)
	r := row.New(headerSize).Add(
		col.New(6),
		col.New(2).Add(
			text.New("TOTAL", props.Text{
				Align: align.Left,
				Size:  14,
				Style: fontstyle.Bold,
				Color: getWhiteColor()}),
		),
		col.New(3).Add(
			text.New(data.Total, props.Text{
				Align: align.Right,
				Size:  14,
				Style: fontstyle.Bold,
				Color: getWhiteColor(),
			}),
		),
		col.New(1),
	)
	r.WithStyle(&props.Cell{
		BackgroundColor: getDarkEmeraldColor(),
	})
	m.AddRows(r)
}

func (o InvoiceData) GetHeader() core.Row {
	r := row.New(headerSize).Add(
		col.New(1),
		text.NewCol(5, "ITEM", props.Text{Align: align.Left, Size: 12, Style: fontstyle.Bold, Color: getWhiteColor()}),
		text.NewCol(2, "PRICE", props.Text{Align: align.Center, Size: 12, Style: fontstyle.Bold, Color: getWhiteColor()}),
		text.NewCol(1, "QTY", props.Text{Align: align.Center, Size: 12, Style: fontstyle.Bold, Color: getWhiteColor()}),
		text.NewCol(2, "SUBTOTAL", props.Text{Align: align.Right, Size: 12, Style: fontstyle.Bold, Color: getWhiteColor()}),
		col.New(1),
	)
	r.WithStyle(&props.Cell{
		BackgroundColor: getDarkEmeraldColor(),
	})
	return r
}
func (o InvoiceItem) GetContent(i int) core.Row {
	r := row.New(bodySize).Add(
		col.New(1),
		text.NewCol(5, o.Name, props.Text{Align: align.Left, Size: 11}),
		text.NewCol(2, o.Price, props.Text{Align: align.Center, Size: 11}),
		text.NewCol(1, o.Qty, props.Text{Align: align.Center, Size: 11}),
		text.NewCol(2, o.Subtotal, props.Text{Align: align.Right, Size: 11}),
		col.New(1),
	)
	if i%2 == 0 {
		r.WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 224, Green: 245, Blue: 240},
		})
	} else {
		r.WithStyle(&props.Cell{
			BackgroundColor: &props.Color{Red: 233, Green: 233, Blue: 233},
		})
	}
	return r
}

func addFooter(email string) core.Row {
	return row.New(20).Add(
		col.New(12).Add(
			text.New("Thank you for your purchase!\nFor questions, contact "+email, props.Text{
				Align: align.Center,
				Color: getGrayColor(),
				Size:  8,
			}),
		),
	)
}

func getWhiteColor() *props.Color {
	return &props.Color{
		Red:   255,
		Green: 255,
		Blue:  255,
	}
}
func getDarkEmeraldColor() *props.Color {
	return &props.Color{
		Red:   0,
		Green: 127,
		Blue:  97,
	}
}

func getGrayColor() *props.Color {
	return &props.Color{
		Red:   160,
		Green: 160,
		Blue:  160,
	}
}

package internal

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/urfave/cli/v2"
	"os"
	"text/tabwriter"
)

func PushCommand(client *http.Client) *Cmd{
	flags := []cli.Flag {
		&cli.StringFlag{
			Name: "phone",
			Aliases: []string{"p"},
			Usage: "phone number to send push request",
		},
		&cli.Int64Flag{
			Name:        "amount",
			Aliases:     []string{"a"},
			Usage:       "amount for push pay",
		},
		&cli.StringFlag{
			Name: "reference",
			Aliases: []string{"ref","r"},
			Usage: "push pay message/description",
		},
	}
	return &Cmd{
		ApiClient:   client,
		RequestType: airtel.USSDPush,
		Name:        "push",
		Usage:       "send ussd push request",
		Description: "send ussd push request",
		Flags:       flags,
		SubCommands: nil,
	}
}

func pushResponsePrintOut(response api.PushPayResponse)  {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer func(w *tabwriter.Writer) {
		err := w.Flush()
		if err != nil {
			fmt.Printf("error while closing tabwriter: %v\n", err)
		}
	}(w)


	_, _ = fmt.Fprintf(w, "\n %s\t", "PUSH RESPONSE")
	_, _ = fmt.Fprintf(w, "\n %s\t", "--------------")

	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "ID:", response.ID)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Error:", response.Error)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Result Code:", response.ResultCode)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Status Code:", response.StatusCode)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Status:", response.Status)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Status Message:", response.StatusMessage)
	_, _ = fmt.Fprintf(w, "\n %s\t%t\t", "Success:", response.Success)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Error Desc:", response.ErrorDescription)

	_, _ = fmt.Fprintf(w, "\n")
}
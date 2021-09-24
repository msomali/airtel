package internal

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/internal/models"
	"os"
	"text/tabwriter"
)

func TokenCommand(client *http.Client) *Cmd{
	return &Cmd{
		ApiClient:   client,
		RequestType: airtel.Authorization,
		Name:        "token",
		Usage:       "retrieve access token",
		Description: "retrieve access token",
		Flags:       nil,
		SubCommands: nil,
	}
}

func tokenResponsePrintOut(response models.AirtelAuthResponse)  {
	// initialize tabwriter
	w := new(tabwriter.Writer)

	w.Init(os.Stdout, 8, 8, 0, '\t', 0)

	defer func(w *tabwriter.Writer) {
		err := w.Flush()
		if err != nil {
			fmt.Printf("error while closing tabwriter: %v\n", err)
		}
	}(w)

	_, _ = fmt.Fprintf(w, "\n %s\t", "ACCESS TOKEN RESPONSE")
	_, _ = fmt.Fprintf(w, "\n %s\t", "---------------------")

	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Access Token:", response.AccessToken)
	_, _ = fmt.Fprintf(w, "\n %s\t%d\t", "Expires In:", response.ExpiresIn)
	_, _ = fmt.Fprintf(w, "\n %s\t%s\t", "Token type:", response.TokenType)
	_, _ = fmt.Fprintf(w, "\n")
}

/*
 * MIT License
 *
 * Copyright (c) 2021 TECHCRAFT TECHNOLOGIES CO LTD.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package internal

import (
	"fmt"
	"github.com/techcraftlabs/airtel"
	"github.com/techcraftlabs/airtel/api/http"
	"github.com/techcraftlabs/airtel/internal/models"
	"os"
	"text/tabwriter"
)

func TokenCommand(client *http.Client) *Cmd {
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

func tokenResponsePrintOut(response models.AirtelAuthResponse) {
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

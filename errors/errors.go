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

package errors

var _ error = (*Code)(nil)

type (
	Code struct {
		ID   string
		Desc string
		Err  error
	}
)

func (c *Code) Error() string {
	return c.Desc
}

func (e *Code) Unwrap() error { return e.Err }

const (
	ESB000001 = "Something went wrong"
	ESB000004 = "An error occurred while initiating the payment"
	ESB000008 = "Field validation"
	ESB000011 = "Failed"
	ESB000014 = "An error occurred while fetching the transaction status"
	ESB000033 = "Invalid MSISDN Length. MSISDN Length should be {0}"
	ESB000034 = "Invalid Country Command"
	ESB000035 = "Invalid Currency CodeName"
	ESB000036 = "Invalid MSISDN Length. MSISDN Length should be {0} and should start with 0"
	ESB000039 = "Vendor is not configured to do transaction in the country {0}"
)

var (
	Codes = map[string]string{
		"ESB000001": ESB000001,
		"ESB000004": ESB000004,
		"ESB000008": ESB000008,
		"ESB000011": ESB000011,
		"ESB000014": ESB000014,
		"ESB000033": ESB000033,
		"ESB000034": ESB000034,
		"ESB000035": ESB000035,
		"ESB000036": ESB000036,
		"ESB000039": ESB000039,
	}
)

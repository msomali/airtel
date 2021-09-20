package internal

import "errors"

const errStatusCodeMargin = 400

var DoErr = errors.New("result code is above or equal to 400")

package controllers

import "io"

type BodyReader func(io.Reader) ([]byte, error)

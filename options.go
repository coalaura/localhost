package main

import "fmt"

type Options struct {
	Certificate string
	Directory   string
	Index       string
	Key         string
	LiveReload  bool
	Open        bool
	Port        int
	Redirect    bool
	Verbose     bool
}

func NewOptions() Options {
	return Options{
		Directory: ".",
	}
}

func (o *Options) HasSSL() bool {
	return o.Certificate != "" && o.Key != ""
}

func (o *Options) GetPort() int {
	if o.Port == 0 {
		if o.HasSSL() {
			o.Port = 443
		} else {
			o.Port = 80
		}
	}

	return o.Port
}

func (o *Options) GetHost() string {
	var (
		host string
		port int
	)

	if o.HasSSL() {
		port = 443
		host = "https://localhost"
	} else {
		port = 80
		host = "http://localhost"
	}

	if actual := o.GetPort(); actual != port {
		host = fmt.Sprintf("%s:%d", host, actual)
	}

	return host
}

func (o *Options) GetIndex(def string) string {
	if o.Index != "" {
		return o.Index
	}

	return def
}

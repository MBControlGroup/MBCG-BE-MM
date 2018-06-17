package main

import (
    "os"
    "github.com/MBControlGroup/MBCG-BE-MM/service"
    flag "github.com/spf13/pflag"
)

const (
    PORT string = "9300"
)

func main() {
    port := os.Getenv("PORT")
    if len(port) == 0 {
        port = PORT
    }

    pPort := flag.StringP("port", "p", PORT, "PORT for httpd listening")
    flag.Parse()
    if len(*pPort) != 0 {
        port = *pPort
    }

    server := service.NewServer()
    server.Run(":" + port)
}

package main

import (
	"fmt"
	// "log"
	// "net/http"
	// "os"
	// "os/signal"
	// "runtime"
	// "syscall"
	// "text/tabwriter"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/rds"

	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/common/version"
	// "gopkg.in/alecthomas/kingpin.v2"
)


func main() {
	fmt.Printf("Starting aws-rds-exporter: start")
	fmt.Printf("\n")

	sess, err := session.NewSession(&aws.Config{Region: aws.String("us-east-1")})

	if err != nil {
		fmt.Printf("Error getting session")
	}

	if sess == nil {
		fmt.Printf("Error getting session after first error")
	} else {
		sType := reflect.TypeOf(sess)
		fmt.Printf(sType.String())
	}

	fmt.Printf("\n")
}





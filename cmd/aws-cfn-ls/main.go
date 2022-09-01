package main

import (
	"fmt"
	"os"

	"github.com/hangxie/aws-utils/cloudformation"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage:", os.Args[0], "stack-name")
		os.Exit(1)
	}
	resourceList, err := cloudformation.ListResources(os.Args[1])
	if err != nil {
		fmt.Println("failed to list resources:", err.Error())
		os.Exit(1)
	}

	for _, r := range resourceList {
		fmt.Println(r.Type, r.Id)
	}
}

package main

import (
	"fmt"
	"os"
	"sort"

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

	sort.Slice(resourceList, lessFunc(resourceList))
	for _, r := range resourceList {
		fmt.Println(r.Type, r.Id)
	}
}

func lessFunc(list []cloudformation.AwsResource) func(int, int) bool {
	return func(i, j int) bool {
		if list[i].Type == list[j].Type {
			return list[i].Id < list[j].Type
		}
		return list[i].Type < list[j].Type
	}
}

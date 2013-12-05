package wmi

import (
	"fmt"
	"testing"
)

func TestQuery(t *testing.T) {
	var dst []Win32_Process
	err := Query("SELECT * FROM Win32_Process", &dst)
	if err != nil {
		t.Fatal(err)
	}
	/*
	b, err := json.MarshalIndent(dst[:10], "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	//b = nil
	fmt.Println(string(b))
	//*/
	for k, v := range dst {
		fmt.Printf("%v %+v\n", k, v)
	}
}

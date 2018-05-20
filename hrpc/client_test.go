package hrpc

import (
	"testing"
	"fmt"
	"time"
	"log"
)

func TestNewClient(t *testing.T) {
	option := ClientOption{
		Debug: true,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 5 * time.Second,
		HeartbeatInterval: 7 * time.Second,
	}

	client, err := Dial("tcp", "127.0.0.1:6789", option)
	if err != nil {
		t.Fatal("dialing", err)
	}
	defer client.Close()

	// Synchronous calls
	for i := 0; i < 6; i++ {
		args := &Args{7, 8}
		reply := new(Reply)
		err = client.Call("Arith.Add", args, reply)
		if err != nil {
			log.Println(err.Error())
			t.Errorf("Add: expected no error but got string %q", err.Error())
		}
		if reply.C != args.A+args.B {
			t.Errorf("Add: expected %d got %d", reply.C, args.A+args.B)
		}
		log.Println(fmt.Sprintf("Arith.Add: %v, %v", args, reply))

		time.Sleep(6 * time.Second)
	}
}

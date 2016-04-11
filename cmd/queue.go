// Copyright © 2016 Genome Research Limited
// Author: Sendu Bala <sb10@sanger.ac.uk>.
//
//  This file is part of VRPipe.
//
//  VRPipe is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  VRPipe is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with VRPipe. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	// "bufio"
	"fmt"
	"github.com/sb10/vrpipe/jobqueue"
	"github.com/sb10/vrpipe/queue"
	"github.com/spf13/cobra"
	// "github.com/ugorji/go/codec"
	"log"
	"net"
	"runtime"
	// "strings"
	// "strconv"
	// "time"
	"os"
	"time"
)

var queues map[string]*queue.Queue

// queueCmd represents the queue command
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "temp playground for queue implementations",
	Long:  `don't use this`,
	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(runtime.NumCPU())
		queues = make(map[string]*queue.Queue)

		l, err := net.Listen("tcp", ":11301")
		if err != nil {
			log.Fatal("Error listening:", err)
		}
		defer l.Close()
		defer os.Remove("/lustre/scratch116/vr/user/sb10/tmp/socket")

		for {
			// listen for an incoming connection.
			netConn, err := l.Accept()
			if err != nil {
				log.Fatal("Error accepting: ", err)
			}

			// handle connections in a new goroutine.
			go handleClient(netConn)
		}
	},
}

func init() {
	RootCmd.AddCommand(queueCmd)
	// queueCmd.Flags().StringVar(&enqueue, "enqueue", "", "Add a job to the queue")
	// queueCmd.Flags().BoolVar(&dequeue, "dequeue", false, "Get a job from the queue")
}

func handleClient(netConn net.Conn) {
	c := jobqueue.New(netConn)
	netConn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	defer netConn.Close()

	// the first thing a client needs to do on connecting is send the desired
	// queue name, which we get here
	err := c.HandleQueue()
	if err != nil {
		fmt.Println(err)
		return
	}

	// now we loop for up to 5mins waiting for some other command
	for {
		err := c.HandleCmd()
		if err != nil {
			if err != jobqueue.ErrClose {
				fmt.Println(err)
			}
			return
		}
		netConn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	}
}

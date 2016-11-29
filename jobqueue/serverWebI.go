// Copyright © 2016 Genome Research Limited
// Author: Sendu Bala <sb10@sanger.ac.uk>.
//
//  This file is part of wr.
//
//  wr is free software: you can redistribute it and/or modify
//  it under the terms of the GNU Lesser General Public License as published by
//  the Free Software Foundation, either version 3 of the License, or
//  (at your option) any later version.
//
//  wr is distributed in the hope that it will be useful,
//  but WITHOUT ANY WARRANTY; without even the implied warranty of
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU Lesser General Public License for more details.
//
//  You should have received a copy of the GNU Lesser General Public License
//  along with wr. If not, see <http://www.gnu.org/licenses/>.

package jobqueue

// This file contains the web interface code of the server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

// jstatusReq is what the status webpage sends us to ask for info about jobs.
// The possible Requests are:
// current = get count info for every job in every RepGroup in the cmds queue.
// details = get example job details for jobs in the RepGroup, grouped by having
//           the same Status, Exitcode and FailReason.
// retry = retry the buried jobs with the given RepGroup, ExitCode and FailReason.
type jstatusReq struct {
	Key        string // sending Key means "give me detailed info about this single job"
	RepGroup   string // sending RepGroup means "send me limited info about the jobs with this RepGroup"
	State      string // A Job.State to limit RepGroup by
	Exitcode   int
	FailReason string
	All        bool // If false, retry mode will act on a single random matching job, instead of all of them
	Request    string
}

// jstatus is the job info we send to the status webpage (only real difference
// to Job is that the times are seconds instead of *time.Duration... *** not
// really sure if we really need this and should just give the webpage Jobs
// directly instead).
type jstatus struct {
	Key          string
	RepGroup     string
	Cmd          string
	State        string
	Cwd          string
	ExpectedRAM  int
	ExpectedTime float64
	Cores        int
	PeakRAM      int
	Exited       bool
	Exitcode     int
	FailReason   string
	Pid          int
	Host         string
	Walltime     float64
	CPUtime      float64
	StdErr       string
	StdOut       string
	Env          []string
	Attempts     uint32
	Similar      int
}

// webInterfaceStatic is a http handler for our static documents in static.go
// (which in turn come from the static folder in the git repository). static.go
// is auto-generated by:
// $ esc -pkg jobqueue -prefix static -private -o jobqueue/static.go static
func webInterfaceStatic(w http.ResponseWriter, r *http.Request) {
	// our home page is /status.html
	path := r.URL.Path
	if path == "/" || path == "/status" {
		path = "/status.html"
	}

	// during development, to avoid having to rebuild and restart manager on
	// every change to a file in static dir, do:
	// $ esc -pkg jobqueue -prefix $GOPATH/src/github.com/VertebrateResequencing/wr/static -private -o jobqueue/static.go $GOPATH/src/github.com/VertebrateResequencing/wr/static
	// and set the boolean to true. Don't forget to rerun esc without the abs
	// paths and change the boolean back to false before any commit!
	doc, err := _escFSByte(false, path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	if strings.HasPrefix(path, "/js") {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	} else if strings.HasPrefix(path, "/css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	} else if strings.HasPrefix(path, "/fonts") {
		if strings.HasSuffix(path, ".eot") {
			w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
		} else if strings.HasSuffix(path, ".svg") {
			w.Header().Set("Content-Type", "image/svg+xml")
		} else if strings.HasSuffix(path, ".ttf") {
			w.Header().Set("Content-Type", "application/x-font-truetype")
		} else if strings.HasSuffix(path, ".woff") {
			w.Header().Set("Content-Type", "application/font-woff")
		} else if strings.HasSuffix(path, ".woff2") {
			w.Header().Set("Content-Type", "application/font-woff2")
		}
	} else if strings.HasSuffix(path, "favicon.ico") {
		w.Header().Set("Content-Type", "image/x-icon")
	}

	w.Write(doc)
}

// webSocket upgrades a http connection to a websocket
func webSocket(w http.ResponseWriter, r *http.Request) (conn *websocket.Conn, ok bool) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	ok = true
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		ok = false
	}
	return
}

// webInterfaceStatusWS reads from and writes to the websocket on the status
// webpage
func webInterfaceStatusWS(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, ok := webSocket(w, r)
		if !ok {
			log.Println("failed to set up websocket at", r.Host)
			return
		}

		writeMutex := &sync.Mutex{}

		// go routine to read client requests and respond to them
		go func(conn *websocket.Conn) {
			// log panics and die
			defer s.logPanic("jobqueue websocket client handling", true)

			for {
				req := jstatusReq{}
				err := conn.ReadJSON(&req)
				if err != nil { // probably the browser was refreshed, breaking conn
					break
				}

				q, existed := s.qs["cmds"]
				if !existed {
					continue
				}

				switch {
				case req.Key != "":
					jobs, _, errstr := s.getJobsByKeys(q, []string{req.Key}, true, true)
					if errstr == "" && len(jobs) == 1 {
						stderr, _ := jobs[0].StdErr()
						stdout, _ := jobs[0].StdOut()
						env, _ := jobs[0].Env()
						status := jstatus{
							Key:          jobs[0].key(),
							RepGroup:     jobs[0].RepGroup,
							Cmd:          jobs[0].Cmd,
							State:        jobs[0].State,
							Cwd:          jobs[0].Cwd,
							ExpectedRAM:  jobs[0].RAM,
							ExpectedTime: jobs[0].Time.Seconds(),
							Cores:        jobs[0].Cores,
							PeakRAM:      jobs[0].PeakRAM,
							Exited:       jobs[0].Exited,
							Exitcode:     jobs[0].Exitcode,
							FailReason:   jobs[0].FailReason,
							Pid:          jobs[0].Pid,
							Host:         jobs[0].Host,
							Walltime:     jobs[0].Walltime.Seconds(),
							CPUtime:      jobs[0].CPUtime.Seconds(),
							StdErr:       stderr,
							StdOut:       stdout,
							Env:          env,
							Attempts:     jobs[0].Attempts,
						}
						writeMutex.Lock()
						err = conn.WriteJSON(status)
						writeMutex.Unlock()
						if err != nil {
							break
						}
					}
				case req.Request != "":
					switch req.Request {
					case "current":
						// get all current jobs
						jobs := s.getJobsCurrent(q, 0, "", false, false)
						writeMutex.Lock()
						err := webInterfaceStatusSendGroupStateCount(conn, "+all+", jobs)
						if err != nil {
							writeMutex.Unlock()
							break
						}

						// for each different RepGroup amongst these jobs,
						// send the job state counts
						repGroups := make(map[string][]*Job)
						for _, job := range jobs {
							repGroups[job.RepGroup] = append(repGroups[job.RepGroup], job)
						}
						failed := false
						for repGroup, jobs := range repGroups {
							complete, _, qerr := s.getCompleteJobsByRepGroup(repGroup)
							if qerr != "" {
								failed = true
								break
							}
							jobs = append(jobs, complete...)
							err := webInterfaceStatusSendGroupStateCount(conn, repGroup, jobs)
							if err != nil {
								failed = true
								break
							}
						}
						writeMutex.Unlock()
						if failed {
							break
						}
					case "details":
						// *** probably want to take the count as a req option,
						// so user can request to see more than just 1 job per
						// State+Exitcode+FailReason
						jobs, _, errstr := s.getJobsByRepGroup(q, req.RepGroup, 1, req.State, true, true)
						if errstr == "" && len(jobs) > 0 {
							writeMutex.Lock()
							failed := false
							for _, job := range jobs {
								stderr, _ := job.StdErr()
								stdout, _ := job.StdOut()
								env, _ := job.Env()
								status := jstatus{
									Key:          job.key(),
									RepGroup:     req.RepGroup, // not job.RepGroup, since we want to return the group the user asked for, not the most recent group the job was made for
									Cmd:          job.Cmd,
									State:        job.State,
									Cwd:          job.Cwd,
									ExpectedRAM:  job.RAM,
									ExpectedTime: job.Time.Seconds(),
									Cores:        job.Cores,
									PeakRAM:      job.PeakRAM,
									Exited:       job.Exited,
									Exitcode:     job.Exitcode,
									FailReason:   job.FailReason,
									Pid:          job.Pid,
									Host:         job.Host,
									Walltime:     job.Walltime.Seconds(),
									CPUtime:      job.CPUtime.Seconds(),
									Attempts:     job.Attempts,
									Similar:      job.Similar,
									StdErr:       stderr,
									StdOut:       stdout,
									Env:          env,
								}
								err = conn.WriteJSON(status)
								if err != nil {
									failed = true
									break
								}
							}
							writeMutex.Unlock()
							if failed {
								break
							}
						}
					case "retry":
						s.rpl.RLock()
						for key := range s.rpl.lookup[req.RepGroup] {
							item, err := q.Get(key)
							if err != nil {
								break
							}
							stats := item.Stats()
							if stats.State == "bury" {
								job := item.Data.(*Job)
								if job.Exitcode == req.Exitcode && job.FailReason == req.FailReason {
									err := q.Kick(key)
									if err != nil {
										break
									}
									job.UntilBuried = job.Retries + 1
									if !req.All {
										break
									}
								}
							}
						}
						s.rpl.RUnlock()
					case "remove":
						s.rpl.RLock()
						var toDelete []string
						for key := range s.rpl.lookup[req.RepGroup] {
							item, err := q.Get(key)
							if err != nil {
								break
							}
							stats := item.Stats()
							if stats.State == "bury" || stats.State == "delay" || stats.State == "dependent" {
								job := item.Data.(*Job)
								if job.Exitcode == req.Exitcode && job.FailReason == req.FailReason {
									// we can't allow the removal of jobs that
									// have dependencies, as *queue would regard
									// that as satisfying the dependency and
									// downstream jobs would start
									hasDeps, err := q.HasDependents(key)
									if err != nil || hasDeps {
										continue
									}

									err = q.Remove(key)
									if err != nil {
										break
									}
									if err == nil {
										s.db.deleteLiveJob(key)
										toDelete = append(toDelete, key)
										if stats.State == "delay" {
											s.decrementGroupCount(job.schedulerGroup, q)
										}
									}
									if !req.All {
										break
									}
								}
							}
						}
						for _, key := range toDelete {
							delete(s.rpl.lookup[req.RepGroup], key)
						}
						s.rpl.RUnlock()
					default:
						continue
					}
				default:
					continue
				}
			}
		}(conn)

		// go routine to push changes to the client
		go func(conn *websocket.Conn) {
			// log panics and die
			defer s.logPanic("jobqueue websocket status updating", true)

			statusReceiver := s.statusCaster.Join()
			for status := range statusReceiver.In {
				writeMutex.Lock()
				err := conn.WriteJSON(status)
				writeMutex.Unlock()
				if err != nil {
					break
				}
			}
			statusReceiver.Close()
		}(conn)
	}
}

// webInterfaceStatusSendGroupStateCount sends the per-repgroup state counts
// to the status webpage websocket
func webInterfaceStatusSendGroupStateCount(conn *websocket.Conn, repGroup string, jobs []*Job) (err error) {
	queueCounts := make(map[string]int)
	for _, job := range jobs {
		var subQueue string
		switch job.State {
		case "delayed":
			subQueue = "delay"
		case "reserved", "running":
			subQueue = "run"
		case "buried":
			subQueue = "bury"
		default:
			subQueue = job.State
		}
		queueCounts[subQueue]++
	}
	for to, count := range queueCounts {
		err = conn.WriteJSON(&jstateCount{repGroup, "new", to, count})
		if err != nil {
			return
		}
	}
	return
}

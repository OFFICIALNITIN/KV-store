package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/OFFICIALNITIN/KV-store/internal/store"
)

func HandleConnection(conn net.Conn, kv *store.KVStore) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		params := strings.Fields(strings.TrimSpace(line))
		if len(params) == 0 {
			continue
		}

		switch strings.ToUpper(params[0]) {
		case "SET":
			if len(params) < 3 {
				conn.Write([]byte("ERR: SET <key> <value> [ttl_secs]\n"))
				continue
			}
			ttl := 0 * time.Second
			if len(params) == 4 {
				if sec, err := strconv.Atoi(params[3]); err == nil {
					ttl = time.Duration(sec) * time.Second
				}
			}
			kv.Set(params[1], params[2], ttl)
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(params) < 2 {
				conn.Write([]byte("ERR: GET <key>\n"))
				continue
			}
			val, found := kv.Get(params[1])
			if !found {
				conn.Write([]byte("(nil)\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("%v\n", val)))
			}

		case "DELETE":
			if len(params) < 2 {
				conn.Write([]byte("ERR: DELETE <key>\n"))
				continue
			}
			kv.Delete(params[1])
			conn.Write([]byte("OK\n"))

		case "EXIT":
			return

		default:
			conn.Write([]byte("ERR: Unknown command\n"))
		}

	}
}

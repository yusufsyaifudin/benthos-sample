package calc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"go.uber.org/multierr"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type CMDCalculateServer struct{}

var _ cli.Command = (*CMDCalculateServer)(nil)

func (c *CMDCalculateServer) Help() string {
	return "Similar like calculate-file-number cmd but using server and distributed lock (redis)"
}

func (c *CMDCalculateServer) Run(args []string) int {
	type Argument struct {
		Port      int    `short:"p" long:"port" description:"Server port"`
		RedisAddr string `long:"redis-addr" description:"Redis address"`
		MutexKey  string `long:"mutex-key" description:"Mutex key"`
		RedisKey  string `long:"redis-key" description:"Redis key to save final value"`
	}

	var argsVal Argument
	args, err := flags.ParseArgs(&argsVal, args)
	if err != nil {
		err = fmt.Errorf("failed parsing flag: %w", err)
		log.Println(err)
		return -1
	}

	client := goredislib.NewClient(&goredislib.Options{
		Addr: argsVal.RedisAddr,
	})

	pool := goredis.NewPool(client)
	rs := redsync.New(pool)
	lock := rs.NewMutex(argsVal.MutexKey)

	handlerCfg := handlerConfig{
		redisCacheName: argsVal.RedisKey,
		cache:          client,
		mutex:          lock,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", c.handler(handlerCfg))

	// wrap mux with our logger. this will
	var handler http.Handler = mux
	handler = logRequestHandler(mux)

	var errChan = make(chan error, 1)
	go func() {
		log.Printf("starting http on port: %d\n", argsVal.Port)
		if _err := http.ListenAndServe(fmt.Sprintf(":%d", argsVal.Port), handler); _err != nil {
			errChan <- fmt.Errorf("http server error: %w", _err)
		}
	}()

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChan:
		log.Println("got an interrupt")
	case err := <-errChan:
		if err != nil {
			log.Fatalln("error while running server")
		}
	}

	return 0
}

func (c *CMDCalculateServer) Synopsis() string {
	return "Similar like calculate-file-number cmd but using server and distributed lock (redis)"
}

type handlerConfig struct {
	redisCacheName string
	cache          *goredislib.Client
	mutex          *redsync.Mutex
}

func (c *CMDCalculateServer) handler(cfg handlerConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Method != http.MethodPost {
			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("method not POST, got %s", r.Method),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		if r.Body == nil {
			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("body is nil"),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		type ReqData struct {
			MathString string `json:"mathString"`
		}

		defer func() {
			if _err := r.Body.Close(); _err != nil {
				log.Printf("error close request body: %s\n", _err)
				return
			}
		}()

		reqData := ReqData{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&reqData)
		if err != nil {
			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("failed read request body: %s", err),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		// parse line string
		line := strings.TrimSpace(reqData.MathString)
		op := line[0]
		val := line[1:]
		currValBigInt := new(big.Int)
		_, err = fmt.Sscan(val, currValBigInt)
		if err != nil {
			err = fmt.Errorf("cannot parse string to big int %s: %w", val, err)

			resp, _ := json.Marshal(map[string]interface{}{
				"error": err,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		// validate valid operator
		switch op {
		case '*', '/', '+', '-':
		default:
			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("unknown operator: %s", string(op)),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		// read from redis to lock and calculate
		// lock context
		err = cfg.mutex.LockContext(ctx)
		if err != nil {
			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("error lock redis %s", err),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		// read from redis
		finalVal, err := cfg.cache.Get(ctx, cfg.redisCacheName).Result()
		if errors.Is(err, goredislib.Nil) {
			log.Println("this is first time call")

			finalVal = "0" // for the first time, the final value is 0
			err = nil      // discard error
		}

		if err != nil {
			// Release the lock so other processes or threads can obtain a lock.
			if ok, _err := cfg.mutex.UnlockContext(ctx); !ok || _err != nil {
				err = fmt.Errorf("cannot unlock redis: %w: %s", err, _err)
			}

			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("error get current value from redis %s", err),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		finalValBigInt := new(big.Int)
		_, err = fmt.Sscan(finalVal, finalValBigInt)
		if err != nil {
			// Release the lock so other processes or threads can obtain a lock.
			if ok, _err := cfg.mutex.UnlockContext(ctx); !ok || _err != nil {
				err = fmt.Errorf("cannot unlock redis: %w: %s", err, _err)
			}

			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("error convert to big int %s", err),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		switch op {
		case '*':
			finalValBigInt = finalValBigInt.Mul(finalValBigInt, currValBigInt)
		case '/':
			finalValBigInt = finalValBigInt.Div(finalValBigInt, currValBigInt)
		case '+':
			finalValBigInt = finalValBigInt.Add(finalValBigInt, currValBigInt)
		case '-':
			finalValBigInt = finalValBigInt.Sub(finalValBigInt, currValBigInt)
		}

		err = cfg.cache.Set(ctx, cfg.redisCacheName, finalValBigInt.String(), 0).Err()
		if err != nil {
			// Release the lock so other processes or threads can obtain a lock.
			if ok, _err := cfg.mutex.UnlockContext(ctx); !ok || _err != nil {
				err = fmt.Errorf("cannot unlock redis: %w: %s", err, _err)
			}

			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("error set current value to redis %s", err),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		// Release the lock so other processes or threads can obtain a lock.
		ok, _err := cfg.mutex.UnlockContext(ctx)
		if !ok || _err != nil {
			resp, _ := json.Marshal(map[string]interface{}{
				"error": fmt.Sprintf("cannot unlock redis: %s", _err),
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(resp)
			return
		}

		resp, _ := json.Marshal(map[string]interface{}{
			"currentValue": finalValBigInt.String(),
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resp)
		return
	}
}

func logRequestHandler(h http.Handler) http.Handler {
	var toSimpleMap = func(h http.Header) map[string]string {
		out := map[string]string{}
		for k, v := range h {
			out[k] = strings.Join(v, " ")
		}

		return out
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		t0 := time.Now()

		var (
			err        error
			reqBody    []byte
			reqBodyErr error
			reqBodyObj interface{}
		)

		if r.Body != nil {
			reqBody, reqBodyErr = ioutil.ReadAll(r.Body)
			if reqBodyErr != nil {
				err = multierr.Append(err, fmt.Errorf("error read request body: %w", reqBodyErr))
				reqBody = []byte("")
			}

			r.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		if _err := json.Unmarshal(reqBody, &reqBodyObj); _err == nil {
			reqBody = []byte("")
		}

		// call the original http.Handler we're wrapping
		httpResp := httptest.NewRecorder()
		h.ServeHTTP(httpResp, r)

		for key, val := range httpResp.Header() {
			w.Header().Set(key, strings.Join(val, " "))
		}

		w.WriteHeader(httpResp.Code)
		_, _ = w.Write(httpResp.Body.Bytes())

		respBody := httpResp.Body.Bytes()
		var respObj interface{}
		if _err := json.Unmarshal(respBody, &respObj); _err == nil {
			respBody = []byte("")
		}

		errStr := ""
		if err != nil {
			errStr = err.Error()
		}

		type HTTPData struct {
			Header     map[string]string `json:"header,omitempty"`
			DataObject interface{}       `json:"data_object,omitempty"`
			DataString string            `json:"data_string,omitempty"`
		}

		type AccessLogData struct {
			Path        string   `json:"path,omitempty"`
			Request     HTTPData `json:"request,omitempty"`
			Response    HTTPData `json:"response,omitempty"`
			Error       string   `json:"error,omitempty"`
			ElapsedTime int64    `json:"elapsed_time,omitempty"`
		}

		// log outgoing request
		logData, _ := json.Marshal(AccessLogData{
			Path: r.URL.String(),
			Request: HTTPData{
				Header:     toSimpleMap(r.Header),
				DataObject: reqBodyObj,
				DataString: string(reqBody),
			},
			Response: HTTPData{
				Header:     toSimpleMap(httpResp.Header()),
				DataObject: respObj,
				DataString: string(respBody),
			},
			Error:       errStr,
			ElapsedTime: time.Since(t0).Milliseconds(),
		})

		log.Println(string(logData))
	}

	// http.HandlerFunc wraps a function so that it
	// implements http.Handler interface
	return http.HandlerFunc(fn)
}

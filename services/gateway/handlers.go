package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
)

type GraphQLRequest struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponse struct {
	Data   interface{}    `json:"data,omitempty"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

func supergraphHandler(subgraphURLs map[string]string, rdb *goredis.Client) http.HandlerFunc {
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"errors":[{"message":"failed to read body"}]}`, http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var gqlReq GraphQLRequest
		if err := json.Unmarshal(body, &gqlReq); err != nil {
			http.Error(w, `{"errors":[{"message":"invalid graphql request"}]}`, http.StatusBadRequest)
			return
		}

		requestID := uuid.New().String()
		w.Header().Set("X-Request-Id", requestID)
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		tenantID := r.Header.Get("X-Tenant-Id")
		userID := r.Header.Get("X-User-Id")
		appID := r.Header.Get("X-App-Id")
		branchID := r.Header.Get("X-Branch-Id")

		results := make(map[string]*GraphQLResponse)
		var mu sync.Mutex
		var wg sync.WaitGroup

		for name, url := range subgraphURLs {
			wg.Add(1)
			go func(svcName, svcURL string) {
				defer wg.Done()

				gqlBody, _ := json.Marshal(gqlReq)
				targetReq, _ := http.NewRequestWithContext(r.Context(), http.MethodPost, svcURL, bytes.NewReader(gqlBody))
				targetReq.Header.Set("Content-Type", "application/json")
				targetReq.Header.Set("Authorization", authHeader)
				targetReq.Header.Set("X-Tenant-Id", tenantID)
				targetReq.Header.Set("X-User-Id", userID)
				targetReq.Header.Set("X-App-Id", appID)
				targetReq.Header.Set("X-Branch-Id", branchID)
				targetReq.Header.Set("X-Request-Id", requestID)

				resp, err := client.Do(targetReq)
				if err != nil {
					mu.Lock()
					results[svcName] = &GraphQLResponse{
						Errors: []GraphQLError{{Message: fmt.Sprintf("%s upstream error: %v", svcName, err)}},
					}
					mu.Unlock()
					return
				}
				defer resp.Body.Close()

				var gqlResp GraphQLResponse
				if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
					mu.Lock()
					results[svcName] = &GraphQLResponse{
						Errors: []GraphQLError{{Message: svcName + " decode error"}},
					}
					mu.Unlock()
					return
				}

				mu.Lock()
				results[svcName] = &gqlResp
				mu.Unlock()
			}(name, url)
		}

		wg.Wait()

		merged := mergeResults(results)
		json.NewEncoder(w).Encode(merged)
	}
}

func mergeResults(results map[string]*GraphQLResponse) map[string]interface{} {
	merged := make(map[string]interface{})
	var allErrors []interface{}

	for svc, resp := range results {
		if resp.Data != nil {
			if dataMap, ok := resp.Data.(map[string]interface{}); ok {
				for k, v := range dataMap {
					merged[k] = v
				}
			}
		}
		for _, e := range resp.Errors {
			if len(e.Message) > 12 && (e.Message[:12] == "Unknown type" || e.Message[:12] == "Cannot query") {
				continue
			}
			allErrors = append(allErrors, map[string]interface{}{
				"message": e.Message,
				"service": svc,
			})
		}
	}

	result := map[string]interface{}{
		"data": merged,
	}
	if len(allErrors) > 0 {
		result["errors"] = allErrors
	}
	return result
}

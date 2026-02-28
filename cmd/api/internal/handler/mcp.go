package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// jsonRPCRequest represents a JSON-RPC 2.0 request.
type jsonRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// jsonRPCResponse represents a JSON-RPC 2.0 response.
type jsonRPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      any       `json:"id"`
	Result  any       `json:"result,omitempty"`
	Error   *rpcError `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MCP handles POST /api/v1/mcp JSON-RPC 2.0 requests.
// Note: JSON-RPC 2.0 defines its own error format (not RFC 9457)
// because JSON-RPC clients expect {jsonrpc, id, error} responses.
func MCP(w http.ResponseWriter, r *http.Request) {
	var req jsonRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      nil,
			Error:   &rpcError{Code: -32700, Message: fmt.Sprintf("parse error: %v", err)},
		})
		return
	}

	resp := routeRPC(req)
	writeJSON(w, http.StatusOK, resp)
}

func routeRPC(req jsonRPCRequest) jsonRPCResponse {
	switch req.Method {
	case "tools/list":
		return handleToolsList(req)
	case "tools/call":
		return handleToolsCall(req)
	default:
		return jsonRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &rpcError{Code: -32601, Message: "method not found"},
		}
	}
}

func handleToolsList(req jsonRPCRequest) jsonRPCResponse {
	tools := []map[string]string{
		{"name": "analyze", "description": "Run OpenAPI spec analysis"},
		{"name": "models", "description": "List supported AI models"},
	}
	return jsonRPCResponse{JSONRPC: "2.0", ID: req.ID, Result: tools}
}

func handleToolsCall(req jsonRPCRequest) jsonRPCResponse {
	return jsonRPCResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]string{"status": "stub", "message": "tool execution not yet implemented"},
	}
}

package sharedtypes

// HandlerRequest represents the client request for a specific chat or embeddings operation.
type HandlerRequest struct {
	Adapter             string            `json:"adapter"` // "chat", "embeddings"
	InstructionGuid     string            `json:"instructionGuid"`
	ModelIds            []string          `json:"modelIds"`                   // optional model ids to define a set of specific models to be used for this request
	Data                interface{}       `json:"data"`                       // for embeddings, this can be a string or []string; for chat, only string is allowed
	Images              []string          `json:"images"`                     // List of images in base64 format
	ChatRequestType     string            `json:"chatRequestType"`            // "summary", "code", "keywords", "general"; only relevant if "adapter" is "chat"
	DataStream          bool              `json:"dataStream"`                 // only relevant if "adapter" is "chat"
	MaxNumberOfKeywords uint32            `json:"maxNumberOfKeywords"`        // only relevant if "chatRequestType" is "keywords"
	IsConversation      bool              `json:"isConversation"`             // only relevant if "chatRequestType" is "code"
	ConversationHistory []HistoricMessage `json:"conversationHistory"`        // only relevant if "isConversation" is true
	GeneralContext      string            `json:"generalContext"`             // any added context you might need
	MsgContext          string            `json:"msgContext"`                 // any added context you might need
	SystemPrompt        string            `json:"systemPrompt"`               // only relevant if "chatRequestType" is "general"
	ModelOptions        ModelOptions      `json:"modelOptions,omitempty"`     // only relevant if "adapter" is "chat"
	EmbeddingOptions    EmbeddingOptions  `json:"embeddingOptions,omitempty"` // only relevant if "adapter" is "embeddings"
}

// HandlerResponse represents the LLM Handler response for a specific request.
type HandlerResponse struct {
	// Common properties
	InstructionGuid string `json:"instructionGuid"`
	Type            string `json:"type"` // "info", "error", "chat", "embeddings"

	// Chat properties
	IsLast           *bool   `json:"isLast,omitempty"`
	Position         *uint32 `json:"position,omitempty"`
	InputTokenCount  *int    `json:"inputTokenCount,omitempty"`
	OutputTokenCount *int    `json:"outputTokenCount,omitempty"`
	ChatData         *string `json:"chatData,omitempty"`

	// Embeddings properties
	EmbeddedData   interface{} `json:"embeddedData,omitempty"`   // []float32 or [][]float32; for BAAI/bge-m3 these are dense vectors
	LexicalWeights interface{} `json:"lexicalWeights,omitempty"` // map[uint]float32 or []map[uint]float32; only for BAAI/bge-m3
	ColbertVecs    interface{} `json:"colbertVecs,omitempty"`    // [][]float32 or [][][]float32; only for BAAI/bge-m3

	// Error properties
	Error *ErrorResponse `json:"error,omitempty"`

	// Info properties
	InfoMessage *string `json:"infoMessage,omitempty"`
}

// ErrorResponse represents the error response sent to the client when something fails during the processing of the request.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TransferDetails holds communication channels for the websocket listener and writer.
type TransferDetails struct {
	ResponseChannel chan HandlerResponse
	RequestChannel  chan HandlerRequest
}

// HistoricMessage represents a past chat messages.
type HistoricMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"` // image in base64 format
}

// OpenAIOption represents an option for an OpenAI API call.
type ModelOptions struct {
	FrequencyPenalty *float32 `json:"frequencyPenalty,omitempty"`
	MaxTokens        *int32   `json:"maxTokens,omitempty"`
	PresencePenalty  *float32 `json:"presencePenalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	Temperature      *float32 `json:"temperature,omitempty"`
	TopP             *float32 `json:"topP,omitempty"`
}

// EmbeddingsOptions represents the options for an embeddings request.
type EmbeddingOptions struct {
	ReturnDense   bool `json:"returnDense"`   // defines if the response should include dense vectors; only for BAAI/bge-m3
	ReturnSparse  bool `json:"returnSparse"`  // defines if the response should include lexical weights; only for BAAI/bge-m3
	ReturnColbert bool `json:"returnColbert"` // defines if the response should include colbert vectors; only for BAAI/bge-m3
}

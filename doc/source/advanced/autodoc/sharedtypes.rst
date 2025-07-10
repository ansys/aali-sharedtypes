.. _sharedtypes:

sharedtypes
===========

This package provides functionality for sharedtypes.

Functions
---------

.. function:: ExtractSessionContext

   func ExtractSessionContext(ctx \*logging.ContextMap, msg \[\]byte) (SessionContext, error)

   SetSessionContext sets the SessionContext struct from the JSON payload  Parameters: - msg: the JSON payload  Returns: - SessionContext: the SessionContext struct


Types
-----

.. type:: DbFilters

   type DbFilters struct

   DbFilters represents the filters for the database.


.. type:: DbArrayFilter

   type DbArrayFilter struct

   DbArrayFilter represents the filter for an array field in the database.


.. type:: DbJsonFilter

   type DbJsonFilter struct

   DbJsonFilter represents the filter for a JSON field in the database.


.. type:: DbData

   type DbData struct

   DbData represents the data stored in the database.


.. type:: DbResponse

   type DbResponse struct

   DbResponse represents the response from the database.


.. type:: DBListCollectionsOutput

   type DBListCollectionsOutput struct

   DBListCollectionsOutput represents the output of listing collections in the database.


.. type:: GeneralNeo4jQueryInput

   type GeneralNeo4jQueryInput struct

   GeneralNeo4jQueryInput represents the input for executing a Neo4j query.


.. type:: GeneralNeo4jQueryOutput

   type GeneralNeo4jQueryOutput struct

   GeneralNeo4jQueryOutput represents the output of executing a Neo4j query.


.. type:: Neo4jResponse

   type Neo4jResponse struct

   neo4jResponse represents the response from the Neo4j query.


.. type:: DbAddDataInput

   type DbAddDataInput struct

   DbAddDataInput represents the input for adding data to the database.


.. type:: DbAddDataOutput

   type DbAddDataOutput struct

   DbAddDataOutput represents the output of adding data to the database.


.. type:: DbCreateCollectionInput

   type DbCreateCollectionInput struct

   DbCreateCollectionInput represents the input for creating a collection in the database.


.. type:: DbCreateCollectionOutput

   type DbCreateCollectionOutput struct

   DbCreateCollectionOutput represents the output of creating a collection in the database.


.. type:: HandlerRequest

   type HandlerRequest struct

   HandlerRequest represents the client request for a specific chat or embeddings operation.


.. type:: HandlerResponse

   type HandlerResponse struct

   HandlerResponse represents the LLM Handler response for a specific request.


.. type:: ErrorResponse

   type ErrorResponse struct

   ErrorResponse represents the error response sent to the client when something fails during the processing of the request.


.. type:: TransferDetails

   type TransferDetails struct

   TransferDetails holds communication channels for the websocket listener and writer.


.. type:: HistoricMessage

   type HistoricMessage struct

   HistoricMessage represents a past chat messages.


.. type:: ModelOptions

   type ModelOptions struct

   OpenAIOption represents an option for an OpenAI API call.


.. type:: EmbeddingOptions

   type EmbeddingOptions struct

   EmbeddingsOptions represents the options for an embeddings request.


.. type:: FunctionDefinition

   type FunctionDefinition struct

   FunctionDefinition is a struct that contains the id, name, description, package, inputs and outputs of a function


.. type:: FlowKitPythonFunction

   type FlowKitPythonFunction struct

   FlowKitPythonFunction is a struct that contains the name, path, description, inputs, outputs and definitions of a FlowKit-Python function


.. type:: FunctionDefinitionShort

   type FunctionDefinitionShort struct

   FunctionDefinitionShort is equivalent to FunctionDefinition but without the API key (used for aali-agent rest API)


.. type:: FunctionInput

   type FunctionInput struct

   FunctionInput is a struct that contains the name, type, go type and options of a function input


.. type:: FunctionOutput

   type FunctionOutput struct

   FunctionOutput is a struct that contains the name, type and go type of a function output


.. type:: FilledInputOutput

   type FilledInputOutput struct

   FilledInputOutput is a struct that contains the name, go type and value of a filled input/output


.. type:: MaterialLlmCriterion

   type MaterialLlmCriterion struct

   Represents a criterion returned from the llm


.. type:: MaterialCriterionWithGuid

   type MaterialCriterionWithGuid struct

   Represents a criterion with its GUID


.. type:: MaterialAttribute

   type MaterialAttribute struct

   Represents a defined material attribute with its name and GUID.


.. type:: ExecRequest

   type ExecRequest struct

   ExecRequest represents the requests that can be sent to aali-exec


.. type:: ExecutionInstruction

   type ExecutionInstruction struct

   ExecutionInstruction contain an array of strings that represent the code to be executed in aali-exec


.. type:: ExecResponse

   type ExecResponse struct

   ExecResponse represents the response that aali-exec sends back


.. type:: ExecutionDetails

   type ExecutionDetails struct

   ExecutionDetails represents the details of the execution


.. type:: FileDetails

   type FileDetails struct

   FileDetails contain parts of a file that is being sent


.. type:: SessionContext

   type SessionContext struct

   Message represents the JSON message you are expecting


.. type:: ConversationHistoryMessage

   type ConversationHistoryMessage struct

   ConversationHistoryMessage is a structure that contains the message ID, role, content, and images of a conversation history message.


.. type:: Feedback

   type Feedback struct

   Feedback is a structure that contains the conversation history, message ID, and feedback options of a workflow feedback.


.. type:: DataExtractionDocumentData

   type DataExtractionDocumentData struct

   DataExtractionDocumentData represents the data extracted from a document.


.. type:: CodeGenerationElement

   type CodeGenerationElement struct


.. type:: CodeGenerationType

   type CodeGenerationType string

   Enum values for CodeGenerationType


.. type:: XMLMemberExample

   type XMLMemberExample struct


.. type:: XMLMemberExampleCode

   type XMLMemberExampleCode struct


.. type:: XMLMemberParam

   type XMLMemberParam struct


.. type:: CodeGenerationExample

   type CodeGenerationExample struct


.. type:: CodeGenerationUserGuideSection

   type CodeGenerationUserGuideSection struct


.. type:: AnsysGPTDefaultFields

   type AnsysGPTDefaultFields struct

   DefaultFields represents the default fields for the user query.


.. type:: ACSSearchResponse

   type ACSSearchResponse struct

   ACSSearchResponse represents the response from the ACS search.


.. type:: AnsysGPTCitation

   type AnsysGPTCitation struct

   AnsysGPTCitation represents the citation from the AnsysGPT.


.. type:: AnsysGPTRetrieverModuleChunk

   type AnsysGPTRetrieverModuleChunk struct

   RetrieverModuleChunk represents a chunk of data context received from the retriever module.

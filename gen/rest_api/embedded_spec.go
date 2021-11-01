// Code generated by go-swagger; DO NOT EDIT.

package rest_api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "title": "dtaservice.proto",
    "version": "version not set"
  },
  "paths": {
    "/v1/document/transform": {
      "post": {
        "tags": [
          "DTAServer"
        ],
        "summary": "Request to transform a plain text document",
        "operationId": "TransformDocument",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dtaserviceDocumentRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceTransformDocumentResponse"
            }
          }
        }
      }
    },
    "/v1/document/transform-pipe": {
      "post": {
        "tags": [
          "DTAServer"
        ],
        "operationId": "TransformPipe",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dtaserviceTransformPipeRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceTransformPipeResponse"
            }
          }
        }
      }
    },
    "/v1/service/list": {
      "get": {
        "tags": [
          "DTAServer"
        ],
        "operationId": "ListServices",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceListServicesResponse"
            }
          }
        }
      }
    },
    "/v1/service/options": {
      "get": {
        "tags": [
          "DTAServer"
        ],
        "operationId": "Options",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceOptionsResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "dtaserviceDocumentRequest": {
      "type": "object",
      "title": "The request message containing the document to be transformed",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "file_name": {
          "type": "string"
        },
        "options": {
          "type": "object"
        },
        "service_name": {
          "type": "string"
        }
      }
    },
    "dtaserviceListServicesResponse": {
      "type": "object",
      "properties": {
        "services": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "dtaserviceOptionsResponse": {
      "type": "object",
      "properties": {
        "services": {
          "type": "string"
        }
      }
    },
    "dtaservicePipeService": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "options": {
          "type": "object"
        }
      }
    },
    "dtaserviceTransformDocumentResponse": {
      "type": "object",
      "title": "The response message containing the transformed message",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "error": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "output": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "dtaserviceTransformPipeRequest": {
      "type": "object",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "file_name": {
          "type": "string"
        },
        "pipeService": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dtaservicePipeService"
          }
        }
      }
    },
    "dtaserviceTransformPipeResponse": {
      "type": "object",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "error": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "last_transformer": {
          "type": "string"
        },
        "output": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "protobufNullValue": {
      "description": "` + "`" + `NullValue` + "`" + ` is a singleton enumeration to represent the null value for the\n` + "`" + `Value` + "`" + ` type union.\n\n The JSON representation for ` + "`" + `NullValue` + "`" + ` is JSON ` + "`" + `null` + "`" + `.\n\n - NULL_VALUE: Null value.",
      "type": "string",
      "default": "NULL_VALUE",
      "enum": [
        "NULL_VALUE"
      ]
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "swagger": "2.0",
  "info": {
    "title": "dtaservice.proto",
    "version": "version not set"
  },
  "paths": {
    "/v1/document/transform": {
      "post": {
        "tags": [
          "DTAServer"
        ],
        "summary": "Request to transform a plain text document",
        "operationId": "TransformDocument",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dtaserviceDocumentRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceTransformDocumentResponse"
            }
          }
        }
      }
    },
    "/v1/document/transform-pipe": {
      "post": {
        "tags": [
          "DTAServer"
        ],
        "operationId": "TransformPipe",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/dtaserviceTransformPipeRequest"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceTransformPipeResponse"
            }
          }
        }
      }
    },
    "/v1/service/list": {
      "get": {
        "tags": [
          "DTAServer"
        ],
        "operationId": "ListServices",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceListServicesResponse"
            }
          }
        }
      }
    },
    "/v1/service/options": {
      "get": {
        "tags": [
          "DTAServer"
        ],
        "operationId": "Options",
        "responses": {
          "200": {
            "description": "A successful response.",
            "schema": {
              "$ref": "#/definitions/dtaserviceOptionsResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "dtaserviceDocumentRequest": {
      "type": "object",
      "title": "The request message containing the document to be transformed",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "file_name": {
          "type": "string"
        },
        "options": {
          "type": "object"
        },
        "service_name": {
          "type": "string"
        }
      }
    },
    "dtaserviceListServicesResponse": {
      "type": "object",
      "properties": {
        "services": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "dtaserviceOptionsResponse": {
      "type": "object",
      "properties": {
        "services": {
          "type": "string"
        }
      }
    },
    "dtaservicePipeService": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string"
        },
        "options": {
          "type": "object"
        }
      }
    },
    "dtaserviceTransformDocumentResponse": {
      "type": "object",
      "title": "The response message containing the transformed message",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "error": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "output": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "dtaserviceTransformPipeRequest": {
      "type": "object",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "file_name": {
          "type": "string"
        },
        "pipeService": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/dtaservicePipeService"
          }
        }
      }
    },
    "dtaserviceTransformPipeResponse": {
      "type": "object",
      "properties": {
        "document": {
          "type": "string",
          "format": "byte"
        },
        "error": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "last_transformer": {
          "type": "string"
        },
        "output": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "protobufNullValue": {
      "description": "` + "`" + `NullValue` + "`" + ` is a singleton enumeration to represent the null value for the\n` + "`" + `Value` + "`" + ` type union.\n\n The JSON representation for ` + "`" + `NullValue` + "`" + ` is JSON ` + "`" + `null` + "`" + `.\n\n - NULL_VALUE: Null value.",
      "type": "string",
      "default": "NULL_VALUE",
      "enum": [
        "NULL_VALUE"
      ]
    }
  }
}`))
}

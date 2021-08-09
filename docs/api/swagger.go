package api

var JSON = `{
  "swagger" : "2.0",
  "info" : {
    "version" : "1.0.0",
    "title" : "CheckMyMole API"
  },
  "host" : "prod.api.checkmoleapp.demo-redisys.com",
  "tags" : [ {
    "name" : "User"
  }, {
    "name" : "Body parts"
  }, {
    "name" : "Questions"
  }, {
    "name" : "Lesions"
  }, {
    "name" : "Requests"
  } ],
  "schemes" : [ "https" ],
  "paths" : {
    "/users/me" : {
      "get" : {
        "tags" : [ "User" ],
        "summary" : "Get the details of the current user",
        "parameters" : [ ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Account"
            }
          }
        }
      }
    },
    "/body-parts" : {
      "get" : {
        "tags" : [ "Body parts" ],
        "summary" : "List all body parts ordered by \"order\"",
        "parameters" : [ ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/BodyPart"
              }
            }
          }
        }
      },
      "post" : {
        "tags" : [ "Body parts" ],
        "summary" : "Create a new body part",
        "parameters" : [ {
          "in" : "body",
          "name" : "body",
          "description" : "Body part to be created",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/BodyPart"
          }
        } ],
        "responses" : {
          "201" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/BodyPart"
            }
          },
          "400" : {
            "description" : "Invalid input"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "body-parts.write" ]
        } ]
      }
    },
    "/body-parts/{id}" : {
      "put" : {
        "tags" : [ "Body parts" ],
        "summary" : "Update an existing body part",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the part to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Body part to be updated",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/BodyPart"
          }
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/BodyPart"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied or invalid JSON input passed"
          },
          "404" : {
            "description" : "Body part not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "body-parts.write" ]
        } ]
      },
      "delete" : {
        "tags" : [ "Body parts" ],
        "summary" : "Delete an existing body part",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the part to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/BodyPart"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied"
          },
          "404" : {
            "description" : "Body part not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "body-parts.write" ]
        } ]
      }
    },
    "/questions" : {
      "get" : {
        "tags" : [ "Questions" ],
        "summary" : "List all questions ordered by \"order\"",
        "parameters" : [ ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Question"
              }
            }
          }
        }
      },
      "post" : {
        "tags" : [ "Questions" ],
        "summary" : "Create a new question",
        "parameters" : [ {
          "in" : "body",
          "name" : "body",
          "description" : "Question to be created",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Question"
          }
        } ],
        "responses" : {
          "201" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Question"
            }
          },
          "400" : {
            "description" : "Invalid input"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "questions.write" ]
        } ]
      }
    },
    "/questions/{id}" : {
      "put" : {
        "tags" : [ "Questions" ],
        "summary" : "Update an existing question",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the question to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Question to be updated",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Question"
          }
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Question"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied or invalid JSON input passed"
          },
          "404" : {
            "description" : "Question not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "questions.write" ]
        } ]
      },
      "delete" : {
        "tags" : [ "Questions" ],
        "summary" : "Delete an existing question",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the question to delete",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Question"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied"
          },
          "404" : {
            "description" : "Question not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "questions.write" ]
        } ]
      }
    },
    "/users/me/lesions" : {
      "get" : {
        "tags" : [ "Lesions" ],
        "summary" : "List all of user's lesions",
        "parameters" : [ {
          "name" : "include_body_parts",
          "in" : "query",
          "description" : "Should body parts be included (as \"body_part\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_last_reports",
          "in" : "query",
          "description" : "Should last reports be included (as \"last_report\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_last_requests",
          "in" : "query",
          "description" : "Should last reports be included (as \"last_report\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "offset",
          "in" : "query",
          "description" : "How many rows should be skipped",
          "required" : false,
          "type" : "integer",
          "default" : 0
        }, {
          "name" : "limit",
          "in" : "query",
          "description" : "How many rows should be returned",
          "required" : false,
          "type" : "integer",
          "default" : 50
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Lesion"
              }
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      },
      "post" : {
        "tags" : [ "Lesions" ],
        "summary" : "Create a new lesion",
        "parameters" : [ {
          "in" : "body",
          "name" : "body",
          "description" : "Lesion to be created",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Lesion"
          }
        } ],
        "responses" : {
          "201" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Lesion"
            }
          },
          "400" : {
            "description" : "Invalid input"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      }
    },
    "/users/me/lesions/{id}" : {
      "put" : {
        "tags" : [ "Lesions" ],
        "summary" : "Update an existing lesion",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the lesion to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Lesion to be updated",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Lesion"
          }
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Lesion"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied or invalid JSON input passed"
          },
          "404" : {
            "description" : "Lesion not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      },
      "delete" : {
        "tags" : [ "Lesions" ],
        "summary" : "Delete an existing lesion",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the lesion to delete",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Lesion"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied"
          },
          "404" : {
            "description" : "Lesion not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      }
    },
    "/users/me/requests" : {
      "get" : {
        "tags" : [ "Requests" ],
        "summary" : "List all of user's requests",
        "parameters" : [ {
          "name" : "offset",
          "in" : "query",
          "description" : "How many rows should be skipped",
          "required" : false,
          "type" : "integer",
          "default" : 0
        }, {
          "name" : "limit",
          "in" : "query",
          "description" : "How many rows should be returned",
          "required" : false,
          "type" : "integer",
          "default" : 50
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Request"
              }
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      },
      "post" : {
        "tags" : [ "Requests" ],
        "summary" : "Create a new request",
        "parameters" : [ {
          "in" : "body",
          "name" : "body",
          "description" : "Request to be updated. You can pass a slice of report IDs as the \"reports\" field if you'd like to change their request_id field in bulk. The status field of the reports gets automatically changed to submitted.",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Request"
          }
        } ],
        "responses" : {
          "201" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Request"
            }
          },
          "400" : {
            "description" : "Invalid input"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      }
    },
    "/users/me/requests/{id}" : {
      "get" : {
        "tags" : [ "Requests" ],
        "summary" : "Gets a single request with related objects",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the part to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "name" : "include_reports",
          "in" : "query",
          "description" : "Should reports be included (as \"reports[]\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "temp_urls",
          "in" : "query",
          "description" : "Should reports.photos be signed?",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_lesions",
          "in" : "query",
          "description" : "Should lesions be included (as \"reports[].lesion\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_answers",
          "in" : "query",
          "description" : "Should answers be included (as \"reports[].answers\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_questions",
          "in" : "query",
          "description" : "Should questions be included (as \"reports[].answers[].question\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "skip_answer",
          "in" : "query",
          "description" : "Whether answer_text should be skipped",
          "required" : false,
          "type" : "boolean"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Request"
            }
          },
          "404" : {
            "description" : "Request not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.read" ]
        } ]
      },
      "put" : {
        "tags" : [ "Requests" ],
        "summary" : "Update an existing request",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the part to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Request to be updated. You can pass a slice of report IDs as the \"reports\" field if you'd like to change their request_id field in bulk. The status field of the reports gets automatically changed to submitted.",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Request"
          }
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Request"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied or invalid JSON input passed"
          },
          "404" : {
            "description" : "Request not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      },
      "delete" : {
        "tags" : [ "Requests" ],
        "summary" : "Delete an existing request",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the request to delete",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Request"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied"
          },
          "404" : {
            "description" : "Request not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      }
    },
    "/users/me/lesions/{lesion}/reports" : {
      "get" : {
        "tags" : [ "Reports" ],
        "summary" : "List all of the reports in the lesion",
        "parameters" : [ {
          "name" : "lesion",
          "in" : "path",
          "description" : "ID of the lesion",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "name" : "temp_urls",
          "in" : "query",
          "description" : "Should the photo URLs be signed?",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_answers",
          "in" : "query",
          "description" : "Should answers be included (as \"answers\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_questions",
          "in" : "query",
          "description" : "Should questions be included (as \"answers[]\" -> \"question\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "offset",
          "in" : "query",
          "description" : "How many rows should be skipped",
          "required" : false,
          "type" : "integer",
          "default" : 0
        }, {
          "name" : "limit",
          "in" : "query",
          "description" : "How many rows should be returned",
          "required" : false,
          "type" : "integer",
          "default" : 50
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Report"
              }
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      },
      "post" : {
        "tags" : [ "Reports" ],
        "summary" : "Create a new report",
        "parameters" : [ {
          "name" : "lesion",
          "in" : "path",
          "description" : "ID of the lesion",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Report to be created",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Report"
          }
        } ],
        "responses" : {
          "201" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Report"
            }
          },
          "400" : {
            "description" : "Invalid input"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      }
    },
    "/users/me/lesions/{lesion}/reports/{id}" : {
      "put" : {
        "tags" : [ "Reports" ],
        "summary" : "Update an existing report",
        "parameters" : [ {
          "name" : "lesion",
          "in" : "path",
          "description" : "ID of the lesion",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the report to update",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Report to be updated",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/Report"
          }
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Report"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied or invalid JSON input passed"
          },
          "404" : {
            "description" : "Report not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      },
      "delete" : {
        "tags" : [ "Reports" ],
        "summary" : "Delete an existing report",
        "parameters" : [ {
          "name" : "lesion",
          "in" : "path",
          "description" : "ID of the lesion",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        }, {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the report to delete",
          "required" : true,
          "type" : "integer",
          "format" : "int64"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Report"
            }
          },
          "400" : {
            "description" : "Invalid ID supplied"
          },
          "404" : {
            "description" : "Report not found"
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "owned.lesions.write" ]
        } ]
      }
    },
    "/requests" : {
      "get" : {
        "tags" : [ "Requests" ],
        "summary" : "List all of the requests in the system",
        "parameters" : [ {
          "name" : "account_id",
          "in" : "query",
          "description" : "Filter requests by this particular account",
          "required" : false,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "status",
          "in" : "query",
          "description" : "Filter requests with this particular status",
          "required" : false,
          "type" : "string"
        }, {
          "name" : "include_accounts",
          "in" : "query",
          "description" : "Should accounts be included (as the \"account\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "offset",
          "in" : "query",
          "description" : "How many rows should be skipped",
          "required" : false,
          "type" : "integer",
          "default" : 0
        }, {
          "name" : "limit",
          "in" : "query",
          "description" : "How many rows should be returned",
          "required" : false,
          "type" : "integer",
          "default" : 50
        }, {
          "name" : "skip_answer",
          "in" : "query",
          "description" : "Whether answer_text should be skipped",
          "required" : false,
          "type" : "boolean",
          "default" : false
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Request"
              }
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "requests.read" ]
        } ]
      }
    },
    "/requests/{id}" : {
      "get" : {
        "tags" : [ "Requests" ],
        "summary" : "Get a single request by its ID",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the request to fetch",
          "required" : true,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "include_account",
          "in" : "query",
          "description" : "Should the account be included (as the \"account\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "skip_answer",
          "in" : "query",
          "description" : "Whether answer_text should be skipped",
          "required" : false,
          "type" : "boolean"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Report"
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "requests.read" ]
        } ]
      },
      "put" : {
        "tags" : [ "Requests" ],
        "summary" : "Allows doctors to update a request (including changing answer and state to answered)",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the request to update",
          "required" : true,
          "type" : "string",
          "format" : "uuid"
        }, {
          "in" : "body",
          "name" : "body",
          "description" : "Request to be updated. It can contain a \"notify_msg\" to send a push notification to user's channel.",
          "required" : true,
          "schema" : {
            "$ref" : "#/definitions/RequestWithReports"
          }
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Report"
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "requests.respond" ]
        } ]
      }
    },
    "/lesions" : {
      "get" : {
        "tags" : [ "Lesions" ],
        "summary" : "List all of the lesions in the system",
        "parameters" : [ {
          "name" : "account_id",
          "in" : "query",
          "description" : "Filter lesions by this particular account_id",
          "required" : false,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "body_part_id",
          "in" : "query",
          "description" : "Filter lesions with this particular body_part_id",
          "required" : false,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "include_body_parts",
          "in" : "query",
          "description" : "Should body_parts be included (as the \"body_part\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "offset",
          "in" : "query",
          "description" : "How many rows should be skipped",
          "required" : false,
          "type" : "integer",
          "default" : 0
        }, {
          "name" : "limit",
          "in" : "query",
          "description" : "How many rows should be returned",
          "required" : false,
          "type" : "integer",
          "default" : 50
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Lesion"
              }
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "lesions.read" ]
        } ]
      }
    },
    "/lesions/{id}" : {
      "get" : {
        "tags" : [ "Lesions" ],
        "summary" : "Get a single lesion by its ID",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the lesion to fetch",
          "required" : true,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "include_body_part",
          "in" : "query",
          "description" : "Should the body_part be included (as the \"body_part\" field)",
          "required" : false,
          "type" : "boolean"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Lesion"
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "lesions.read" ]
        } ]
      }
    },
    "/reports" : {
      "get" : {
        "tags" : [ "Reports" ],
        "summary" : "List all of the reports in the system",
        "parameters" : [ {
          "name" : "request_id",
          "in" : "query",
          "description" : "Filter reports by this particular request_id",
          "required" : false,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "lesion_id",
          "in" : "query",
          "description" : "Filter reports by this particular lesion_id",
          "required" : false,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "temp_urls",
          "in" : "query",
          "description" : "Should photo URLs be signed?",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_lesions",
          "in" : "query",
          "description" : "Should lesions be included (as the \"lesion\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_body_parts",
          "in" : "query",
          "description" : "Should body_parts be included (as the \"body_part\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_answers",
          "in" : "query",
          "description" : "Should answers be included (as the \"answers\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_questions",
          "in" : "query",
          "description" : "Should questions be included (as the \"answers[].question\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "offset",
          "in" : "query",
          "description" : "How many rows should be skipped",
          "required" : false,
          "type" : "integer",
          "default" : 0
        }, {
          "name" : "limit",
          "in" : "query",
          "description" : "How many rows should be returned",
          "required" : false,
          "type" : "integer",
          "default" : 50
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "type" : "array",
              "items" : {
                "$ref" : "#/definitions/Report"
              }
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "reports.read" ]
        } ]
      }
    },
    "/reports/{id}" : {
      "get" : {
        "tags" : [ "Reports" ],
        "summary" : "Get a single report by its ID",
        "parameters" : [ {
          "name" : "id",
          "in" : "path",
          "description" : "ID of the report to fetch",
          "required" : true,
          "type" : "string",
          "format" : "uuid"
        }, {
          "name" : "include_lesion",
          "in" : "query",
          "description" : "Should the lesion be included (as the \"lesion\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_body_part",
          "in" : "query",
          "description" : "Should the body_part be included (as the \"body_part\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_answers",
          "in" : "query",
          "description" : "Should answers be included (as the \"answers\" field)",
          "required" : false,
          "type" : "boolean"
        }, {
          "name" : "include_questions",
          "in" : "query",
          "description" : "Should questions be included (as the \"answers[].question\" field)",
          "required" : false,
          "type" : "boolean"
        } ],
        "responses" : {
          "200" : {
            "description" : "successful operation",
            "schema" : {
              "$ref" : "#/definitions/Report"
            }
          }
        },
        "security" : [ {
          "checkmoleapp_auth" : [ "reports.read" ]
        } ]
      }
    }
  },
  "securityDefinitions" : {
    "checkmoleapp_auth" : {
      "type" : "oauth2",
      "authorizationUrl" : "https://auth.checkmoleapp.demo-redisys.com/login",
      "flow" : "implicit",
      "scopes" : {
        "body-parts.write" : "create, update and delete body parts",
        "questions.write" : "create, update and delete questions",
        "owned.lesions.read" : "read lesions, reports and requests owned by the user",
        "owned.lesions.write" : "modify lesions, reports and requests owned by the user",
        "requests.read" : "read all requests in the system",
        "reports.read" : "read all reports in the system",
        "requests.respond" : "update all requests in the system",
        "lesions.read" : "read all lesions in the system"
      }
    }
  },
  "definitions" : {
    "Account" : {
      "type" : "object",
      "required" : [ "email", "gender", "name", "phone" ],
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "name" : {
          "type" : "string",
          "example" : "John Doe",
          "description" : "Full name of the user"
        },
        "email" : {
          "type" : "string",
          "example" : "hello@world.org",
          "description" : "User's email"
        },
        "phone" : {
          "type" : "string",
          "description" : "User's phone number"
        },
        "gender" : {
          "type" : "string",
          "description" : "Whatever user writes in the gender field"
        },
        "created_at" : {
          "type" : "integer"
        },
        "updated_at" : {
          "type" : "integer"
        }
      },
      "example" : {
        "gender" : "gender",
        "updated_at" : 6,
        "phone" : "phone",
        "name" : "John Doe",
        "created_at" : 0,
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "email" : "hello@world.org"
      }
    },
    "BodyPart" : {
      "type" : "object",
      "required" : [ "image", "name", "order" ],
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "name" : {
          "type" : "string",
          "example" : "Forehead"
        },
        "displayed" : {
          "type" : "boolean",
          "description" : "Whether it should be displayed in the UI"
        },
        "image" : {
          "type" : "string",
          "example" : "https://example.com/some-photo.png",
          "description" : "URL to the image that will be used in the body part selector"
        },
        "order" : {
          "type" : "integer",
          "format" : "int64",
          "example" : 10,
          "description" : "Body parts will be sorted in the application according to this order."
        },
        "parent" : {
          "type" : "string",
          "format" : "uuid"
        }
      },
      "example" : {
        "displayed" : true,
        "image" : "https://example.com/some-photo.png",
        "parent" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "name" : "Forehead",
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "order" : 10
      }
    },
    "Question" : {
      "type" : "object",
      "required" : [ "answers", "displayed", "name", "order", "type" ],
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "name" : {
          "type" : "string",
          "example" : "Do you smoke?"
        },
        "type" : {
          "type" : "string",
          "description" : "Question type"
        },
        "answers" : {
          "type" : "object",
          "description" : "Data describing possible answers",
          "properties" : { }
        },
        "displayed" : {
          "type" : "boolean",
          "example" : false,
          "description" : "Whether the question should be displayed"
        },
        "order" : {
          "type" : "integer",
          "format" : "int64",
          "example" : 10,
          "description" : "Questions will be sorted in the application according to this order."
        },
        "created_at" : {
          "type" : "integer"
        },
        "updated_at" : {
          "type" : "integer"
        }
      },
      "example" : {
        "displayed" : false,
        "updated_at" : 6,
        "name" : "Do you smoke?",
        "answers" : "{}",
        "created_at" : 0,
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "type" : "type",
        "order" : 10
      }
    },
    "Lesion" : {
      "type" : "object",
      "required" : [ "account_id", "body_part_id", "body_part_location", "name" ],
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "account_id" : {
          "type" : "string",
          "format" : "uuid",
          "description" : "ID of the account that owns the lesion"
        },
        "name" : {
          "type" : "string",
          "example" : "Mole on my forehead"
        },
        "body_part_id" : {
          "type" : "string",
          "format" : "uuid",
          "description" : "ID of the body part"
        },
        "body_part_location" : {
          "type" : "string",
          "description" : "URL to the photo containing the marked location"
        },
        "created_at" : {
          "type" : "integer"
        },
        "updated_at" : {
          "type" : "integer"
        }
      },
      "example" : {
        "account_id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "updated_at" : 6,
        "body_part_id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "body_part_location" : "body_part_location",
        "name" : "Mole on my forehead",
        "created_at" : 0,
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91"
      }
    },
    "Request" : {
      "type" : "object",
      "required" : [ "account_id", "status" ],
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "account_id" : {
          "type" : "string",
          "format" : "uuid",
          "description" : "ID of the account that owns the request"
        },
        "status" : {
          "type" : "string",
          "description" : "draft, submitted or answered"
        },
        "answer_text" : {
          "type" : "string",
          "description" : "answer from the doctor"
        },
        "answered_by" : {
          "type" : "string",
          "description" : "name of the doctor"
        },
        "answered_at" : {
          "type" : "integer",
          "description" : "when it was answered"
        },
        "created_at" : {
          "type" : "integer"
        },
        "updated_at" : {
          "type" : "integer"
        }
      },
      "example" : {
        "answer_text" : "answer_text",
        "account_id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "updated_at" : 1,
        "answered_by" : "answered_by",
        "answered_at" : 0,
        "created_at" : 6,
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "status" : "status"
      }
    },
    "RequestWithReports" : {
      "type" : "object",
      "required" : [ "status" ],
      "properties" : {
        "status" : {
          "type" : "string",
          "description" : "draft, submitted or answered"
        },
        "answer_text" : {
          "type" : "string",
          "description" : "answer from the doctor"
        },
        "answered_by" : {
          "type" : "string",
          "description" : "name of the doctor"
        },
        "answered_at" : {
          "type" : "integer",
          "description" : "when it was answered"
        },
        "reports" : {
          "type" : "array",
          "items" : {
            "$ref" : "#/definitions/RequestWithReports_reports"
          }
        }
      },
      "example" : {
        "answer_text" : "answer_text",
        "reports" : [ {
          "consultation_result" : "consultation_result",
          "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91"
        }, {
          "consultation_result" : "consultation_result",
          "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91"
        } ],
        "answered_by" : "answered_by",
        "answered_at" : 0,
        "status" : "status"
      }
    },
    "Report" : {
      "type" : "object",
      "required" : [ "lesion_id", "photos" ],
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "request_id" : {
          "type" : "string",
          "format" : "uuid",
          "description" : "ID of the request that the report belongs to"
        },
        "lesion_id" : {
          "type" : "string",
          "format" : "uuid",
          "description" : "ID of the lesion that this report is about"
        },
        "photos" : {
          "type" : "array",
          "items" : {
            "type" : "string",
            "description" : "URL of the photo"
          }
        },
        "status" : {
          "type" : "string",
          "description" : "Any value is accepted"
        },
        "consultation_result" : {
          "type" : "string"
        },
        "answers" : {
          "type" : "array",
          "items" : {
            "$ref" : "#/definitions/Report_answers"
          }
        }
      },
      "example" : {
        "lesion_id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "consultation_result" : "consultation_result",
        "answers" : [ {
          "question" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
          "answer" : "{}"
        }, {
          "question" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
          "answer" : "{}"
        } ],
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "request_id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "photos" : [ "photos", "photos" ],
        "status" : "status"
      }
    },
    "RequestWithReports_reports" : {
      "properties" : {
        "id" : {
          "type" : "string",
          "format" : "uuid"
        },
        "consultation_result" : {
          "type" : "string",
          "description" : "Report-specific consultation result"
        }
      },
      "example" : {
        "consultation_result" : "consultation_result",
        "id" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91"
      }
    },
    "Report_answers" : {
      "properties" : {
        "question" : {
          "type" : "string",
          "format" : "uuid",
          "description" : "ID of the question"
        },
        "answer" : {
          "type" : "object",
          "description" : "JSON value with the answer",
          "properties" : { }
        }
      },
      "example" : {
        "question" : "046b6c7f-0b8a-43b9-b35d-6489e6daee91",
        "answer" : "{}"
      }
    }
  }
}`

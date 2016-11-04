package capture

const (
	schema = `{
		"$schema": "http://json-schema.org/draft-04/schema#",
		"title": "Campaign POST request body",
	    "type": "object",
	    "properties": {
	        "campaignName": {
	            "type": "string",
	            "description": "Campaign reference name"
	        },
	        "campaignVersion": {
	            "type": "string",
	            "description": "Campaign semantic version number"
	        },
	        "sessionFingerprint": {
	            "type": "string",
	            "description": "Campaign entrant browser session hash"
	        }, 
	        "submitAction": {
	            "type": "string",
	            "description": "Action to be taken on entry submission",
	            "enum": ["email", "store", ""]
	        },
	        "entrant": {
	            "type": "object",
	            "description": "Campaign entrant information",
	            "properties": {
	                "title": {
	                    "type": "string",
	                    "description": "Campaign entrant title"
	                },
	                "firstName": {
	                    "type": "string",
	                    "description": "Campaign entrant forename"
	                },
	                "lastName": {
	                    "type": "string",
	                    "description": "Campaign entrant surname"
	                },
	                "emailAddress": {
	                    "type": "string",
	                    "description": "Campaign entrant email address"
	                },
	                "phoneNumber": {
	                    "type": "string",
	                    "description": "Campaign entrant phone number"
	                },
	                "birthDate": {
	                    "type": "string",
	                    "description": "Campaign entrant date of birth"
	                }
	            },
	            "additionalProperties": false,
	            "required": ["title", "firstName", "lastName", "emailAddress"]
	        },
	        "permissions" : {
	        	"type": "object",
	        	"properties": {
	        		"OptInEmail": {
	                    "type": "boolean",
	                    "description": "Marketing optin indicator for email"
	                },
	                "optInPhone": {
	                    "type": "boolean",
	                    "description": "Marketing optin indicator for telephone"
	                },
	                "OptInSms": {
	                    "type": "boolean",
	                    "description": "Marketing optin indicator for SMS"
	                },
	                "OptInPost": {
	                    "type": "boolean",
	                    "description": "Marketing optin indicator for post"
	                }
	        	},
	        	"additionalProperties": false
	        },
	        "form": {
	            "type": "array",
	            "description": "Array of campaign questions and responses",
	            "items": {
	                "type": "object",
	                "description": "Campaign question and response",
	                "properties": {
	                    "key": {
	                        "type": "string",
	                        "description": "Campaign question text"
	                    },
	                    "value": {
	                        "type": "string",
	                        "description": "Campaign question response text"
	                    }
	                },
	                "additionalProperties": false,
	                "required": ["key", "value"]
	            }
	        },
	        "tags": {
	            "type": "array",
	            "description": "Array of arbitrary key value pairs",
	            "items": {
	                "type": "object",
	                "description": "Optional extra information",
	                "properties": {
	                    "key": {
	                        "type": "string"
	                    },
	                    "value": {
	                        "type": "string"
	                    }
	                },
	                "additionalProperties": false,
	                "required": ["key", "value"]
	            }
	        }
	    },
	    "additionalProperties": false,
	    "required": ["campaignName", "campaignVersion", "entrant", "form"]
	}`
)

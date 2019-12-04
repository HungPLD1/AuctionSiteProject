// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag at
// 2019-12-03 22:06:54.930086 +0700 +07 m=+0.046875101

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/categories": {
            "get": {
                "description": "Search categories by id, return all by default, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of categories, if empty then return all",
                        "name": "id",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Categories"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/item": {
            "get": {
                "description": "Search item by query, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Name of the item (or part of it)",
                        "name": "name",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Item Categories by number",
                        "name": "categories",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Items"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/item/:id": {
            "get": {
                "description": "Get the item's informations and pictures by ID, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "the item ID number",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Items"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Login by JSON form, return a jwt session token in JSON form",
                "parameters": [
                    {
                        "description": "username",
                        "name": "userid",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "password",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session token",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/logs/:id": {
            "get": {
                "description": "Get Bid Session Logs by session ID, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.BidSessionLog"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/profile": {
            "get": {
                "description": "Show user profile, return user general profile in JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserCommon"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "401": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            },
            "put": {
                "description": "Modify/Update user profile, return message in JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "401": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/review/:id": {
            "get": {
                "description": "Show review of User by User ID, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.UserReview"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/session/:id": {
            "get": {
                "description": "Show Session information by ID, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.BidSession"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/signup": {
            "post": {
                "description": "Register new Account in JSON form, return a jwt session token in JSON form",
                "parameters": [
                    {
                        "description": "username",
                        "name": "userid",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "password",
                        "name": "password",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Session token",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/wishlist": {
            "get": {
                "description": "Show user WishList, return a JSON form",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Items"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "401": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        },
        "/wishlist/:id": {
            "post": {
                "description": "Add new item to wishlist, return a JSON message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Item id to be added to wishlist",
                        "name": "itemid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "401": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            },
            "delete": {
                "description": "Remove item from wishlist, return a JSON message",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Item id to be removed from wishlist",
                        "name": "itemid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Success message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "400": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "401": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    },
                    "500": {
                        "description": "Error message",
                        "schema": {
                            "type": "body"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.BidSession": {
            "type": "object",
            "properties": {
                "enddate": {
                    "type": "string"
                },
                "itemid": {
                    "type": "integer"
                },
                "sellerid": {
                    "type": "string"
                },
                "sessionid": {
                    "type": "integer"
                },
                "sessionstatus": {
                    "type": "string"
                },
                "startdate": {
                    "type": "string"
                }
            }
        },
        "model.BidSessionLog": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "number"
                },
                "biddate": {
                    "type": "string"
                },
                "sessionid": {
                    "type": "integer"
                },
                "userid": {
                    "type": "string"
                }
            }
        },
        "model.Categories": {
            "type": "object",
            "properties": {
                "categoriesName": {
                    "type": "string"
                },
                "categoriesid": {
                    "type": "integer"
                }
            }
        },
        "model.Items": {
            "type": "object",
            "properties": {
                "categoriesid": {
                    "type": "integer"
                },
                "imagelink": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "itemDescription": {
                    "type": "string"
                },
                "itemName": {
                    "type": "string"
                },
                "itemSaleStatus": {
                    "type": "string"
                },
                "itemcondition": {
                    "type": "string"
                },
                "itemid": {
                    "type": "integer"
                }
            }
        },
        "model.UserCommon": {
            "type": "object",
            "properties": {
                "accesslevel": {
                    "type": "integer"
                },
                "address": {
                    "type": "string"
                },
                "gender": {
                    "description": "UserBirth \t\t\ttime.Time ` + "`" + `gorm:\"type:date\" json:\"birthdate\"` + "`" + `",
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phone": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                },
                "userid": {
                    "type": "string"
                }
            }
        },
        "model.UserReview": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "score": {
                    "type": "integer"
                },
                "targetid": {
                    "type": "string"
                },
                "writerid": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}

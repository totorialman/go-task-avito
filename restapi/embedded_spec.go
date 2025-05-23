// Code generated by go-swagger; DO NOT EDIT.

package restapi

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
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Сервис для управления ПВЗ и приемкой товаров",
    "title": "backend service",
    "version": "1.0.0"
  },
  "paths": {
    "/dummyLogin": {
      "post": {
        "summary": "Получение тестового токена",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "role"
              ],
              "properties": {
                "role": {
                  "type": "string",
                  "enum": [
                    "employee",
                    "moderator"
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Успешная авторизация",
            "schema": {
              "$ref": "#/definitions/Token"
            }
          },
          "400": {
            "description": "Неверный запрос",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/login": {
      "post": {
        "summary": "Авторизация пользователя",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "email",
                "password"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "password": {
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Успешная авторизация",
            "schema": {
              "$ref": "#/definitions/Token"
            }
          },
          "401": {
            "description": "Неверные учетные данные",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/products": {
      "post": {
        "summary": "Добавление товара в текущую приемку (только для сотрудников ПВЗ)",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "type",
                "pvzId"
              ],
              "properties": {
                "pvzId": {
                  "type": "string",
                  "format": "uuid"
                },
                "type": {
                  "type": "string",
                  "enum": [
                    "электроника",
                    "одежда",
                    "обувь"
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Товар добавлен",
            "schema": {
              "$ref": "#/definitions/Product"
            }
          },
          "400": {
            "description": "Неверный запрос или нет активной приемки",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pvz": {
      "get": {
        "summary": "Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией",
        "parameters": [
          {
            "type": "string",
            "format": "date-time",
            "name": "startDate",
            "in": "query"
          },
          {
            "type": "string",
            "format": "date-time",
            "name": "endDate",
            "in": "query"
          },
          {
            "minimum": 1,
            "type": "integer",
            "default": 1,
            "name": "page",
            "in": "query"
          },
          {
            "maximum": 30,
            "minimum": 1,
            "type": "integer",
            "default": 10,
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Список ПВЗ",
            "schema": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "pvz": {
                    "$ref": "#/definitions/PVZ"
                  },
                  "receptions": {
                    "type": "array",
                    "items": {
                      "type": "object",
                      "properties": {
                        "products": {
                          "type": "array",
                          "items": {
                            "$ref": "#/definitions/Product"
                          }
                        },
                        "reception": {
                          "$ref": "#/definitions/Reception"
                        }
                      }
                    }
                  }
                }
              }
            }
          }
        }
      },
      "post": {
        "summary": "Создание ПВЗ (только для модераторов)",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/PVZ"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "ПВЗ создан",
            "schema": {
              "$ref": "#/definitions/PVZ"
            }
          },
          "400": {
            "description": "Неверный запрос",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pvz/{pvzId}/close_last_reception": {
      "post": {
        "summary": "Закрытие последней открытой приемки товаров в рамках ПВЗ",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "name": "pvzId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Приемка закрыта",
            "schema": {
              "$ref": "#/definitions/Reception"
            }
          },
          "400": {
            "description": "Неверный запрос или приемка уже закрыта",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pvz/{pvzId}/delete_last_product": {
      "post": {
        "summary": "Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "name": "pvzId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Товар удален"
          },
          "400": {
            "description": "Неверный запрос, нет активной приемки или нет товаров для удаления",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/receptions": {
      "post": {
        "summary": "Создание новой приемки товаров (только для сотрудников ПВЗ)",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "pvzId"
              ],
              "properties": {
                "pvzId": {
                  "type": "string",
                  "format": "uuid"
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Приемка создана",
            "schema": {
              "$ref": "#/definitions/Reception"
            }
          },
          "400": {
            "description": "Неверный запрос или есть незакрытая приемка",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/register": {
      "post": {
        "summary": "Регистрация пользователя",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "email",
                "password",
                "role"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "password": {
                  "type": "string"
                },
                "role": {
                  "type": "string",
                  "enum": [
                    "employee",
                    "moderator"
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Пользователь создан",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "400": {
            "description": "Неверный запрос",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "PVZ": {
      "type": "object",
      "required": [
        "city"
      ],
      "properties": {
        "city": {
          "type": "string",
          "enum": [
            "Москва",
            "Санкт-Петербург",
            "Казань"
          ]
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "registrationDate": {
          "type": "string",
          "format": "date-time"
        }
      }
    },
    "Product": {
      "type": "object",
      "required": [
        "type",
        "receptionId"
      ],
      "properties": {
        "dateTime": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "receptionId": {
          "type": "string",
          "format": "uuid"
        },
        "type": {
          "type": "string",
          "enum": [
            "электроника",
            "одежда",
            "обувь"
          ]
        }
      }
    },
    "Reception": {
      "type": "object",
      "required": [
        "dateTime",
        "pvzId",
        "status"
      ],
      "properties": {
        "dateTime": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "pvzId": {
          "type": "string",
          "format": "uuid"
        },
        "status": {
          "type": "string",
          "enum": [
            "in_progress",
            "close"
          ]
        }
      }
    },
    "Token": {
      "type": "string"
    },
    "User": {
      "type": "object",
      "required": [
        "email",
        "role"
      ],
      "properties": {
        "email": {
          "type": "string",
          "format": "email"
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "role": {
          "type": "string",
          "enum": [
            "employee",
            "moderator"
          ]
        }
      }
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
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "Сервис для управления ПВЗ и приемкой товаров",
    "title": "backend service",
    "version": "1.0.0"
  },
  "paths": {
    "/dummyLogin": {
      "post": {
        "summary": "Получение тестового токена",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "role"
              ],
              "properties": {
                "role": {
                  "type": "string",
                  "enum": [
                    "employee",
                    "moderator"
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Успешная авторизация",
            "schema": {
              "$ref": "#/definitions/Token"
            }
          },
          "400": {
            "description": "Неверный запрос",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/login": {
      "post": {
        "summary": "Авторизация пользователя",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "email",
                "password"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "password": {
                  "type": "string"
                }
              }
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Успешная авторизация",
            "schema": {
              "$ref": "#/definitions/Token"
            }
          },
          "401": {
            "description": "Неверные учетные данные",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/products": {
      "post": {
        "summary": "Добавление товара в текущую приемку (только для сотрудников ПВЗ)",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "type",
                "pvzId"
              ],
              "properties": {
                "pvzId": {
                  "type": "string",
                  "format": "uuid"
                },
                "type": {
                  "type": "string",
                  "enum": [
                    "электроника",
                    "одежда",
                    "обувь"
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Товар добавлен",
            "schema": {
              "$ref": "#/definitions/Product"
            }
          },
          "400": {
            "description": "Неверный запрос или нет активной приемки",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pvz": {
      "get": {
        "summary": "Получение списка ПВЗ с фильтрацией по дате приемки и пагинацией",
        "parameters": [
          {
            "type": "string",
            "format": "date-time",
            "name": "startDate",
            "in": "query"
          },
          {
            "type": "string",
            "format": "date-time",
            "name": "endDate",
            "in": "query"
          },
          {
            "minimum": 1,
            "type": "integer",
            "default": 1,
            "name": "page",
            "in": "query"
          },
          {
            "maximum": 30,
            "minimum": 1,
            "type": "integer",
            "default": 10,
            "name": "limit",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "description": "Список ПВЗ",
            "schema": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/GetPvzOKBodyItems0"
              }
            }
          }
        }
      },
      "post": {
        "summary": "Создание ПВЗ (только для модераторов)",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/PVZ"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "ПВЗ создан",
            "schema": {
              "$ref": "#/definitions/PVZ"
            }
          },
          "400": {
            "description": "Неверный запрос",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pvz/{pvzId}/close_last_reception": {
      "post": {
        "summary": "Закрытие последней открытой приемки товаров в рамках ПВЗ",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "name": "pvzId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Приемка закрыта",
            "schema": {
              "$ref": "#/definitions/Reception"
            }
          },
          "400": {
            "description": "Неверный запрос или приемка уже закрыта",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/pvz/{pvzId}/delete_last_product": {
      "post": {
        "summary": "Удаление последнего добавленного товара из текущей приемки (LIFO, только для сотрудников ПВЗ)",
        "parameters": [
          {
            "type": "string",
            "format": "uuid",
            "name": "pvzId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "Товар удален"
          },
          "400": {
            "description": "Неверный запрос, нет активной приемки или нет товаров для удаления",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/receptions": {
      "post": {
        "summary": "Создание новой приемки товаров (только для сотрудников ПВЗ)",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "pvzId"
              ],
              "properties": {
                "pvzId": {
                  "type": "string",
                  "format": "uuid"
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Приемка создана",
            "schema": {
              "$ref": "#/definitions/Reception"
            }
          },
          "400": {
            "description": "Неверный запрос или есть незакрытая приемка",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          },
          "403": {
            "description": "Доступ запрещен",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    },
    "/register": {
      "post": {
        "summary": "Регистрация пользователя",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "type": "object",
              "required": [
                "email",
                "password",
                "role"
              ],
              "properties": {
                "email": {
                  "type": "string",
                  "format": "email"
                },
                "password": {
                  "type": "string"
                },
                "role": {
                  "type": "string",
                  "enum": [
                    "employee",
                    "moderator"
                  ]
                }
              }
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Пользователь создан",
            "schema": {
              "$ref": "#/definitions/User"
            }
          },
          "400": {
            "description": "Неверный запрос",
            "schema": {
              "$ref": "#/definitions/Error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "message": {
          "type": "string"
        }
      }
    },
    "GetPvzOKBodyItems0": {
      "type": "object",
      "properties": {
        "pvz": {
          "$ref": "#/definitions/PVZ"
        },
        "receptions": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/GetPvzOKBodyItems0ReceptionsItems0"
          }
        }
      }
    },
    "GetPvzOKBodyItems0ReceptionsItems0": {
      "type": "object",
      "properties": {
        "products": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Product"
          }
        },
        "reception": {
          "$ref": "#/definitions/Reception"
        }
      }
    },
    "PVZ": {
      "type": "object",
      "required": [
        "city"
      ],
      "properties": {
        "city": {
          "type": "string",
          "enum": [
            "Москва",
            "Санкт-Петербург",
            "Казань"
          ]
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "registrationDate": {
          "type": "string",
          "format": "date-time"
        }
      }
    },
    "Product": {
      "type": "object",
      "required": [
        "type",
        "receptionId"
      ],
      "properties": {
        "dateTime": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "receptionId": {
          "type": "string",
          "format": "uuid"
        },
        "type": {
          "type": "string",
          "enum": [
            "электроника",
            "одежда",
            "обувь"
          ]
        }
      }
    },
    "Reception": {
      "type": "object",
      "required": [
        "dateTime",
        "pvzId",
        "status"
      ],
      "properties": {
        "dateTime": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "pvzId": {
          "type": "string",
          "format": "uuid"
        },
        "status": {
          "type": "string",
          "enum": [
            "in_progress",
            "close"
          ]
        }
      }
    },
    "Token": {
      "type": "string"
    },
    "User": {
      "type": "object",
      "required": [
        "email",
        "role"
      ],
      "properties": {
        "email": {
          "type": "string",
          "format": "email"
        },
        "id": {
          "type": "string",
          "format": "uuid"
        },
        "role": {
          "type": "string",
          "enum": [
            "employee",
            "moderator"
          ]
        }
      }
    }
  }
}`))
}

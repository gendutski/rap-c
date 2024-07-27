# Unit API Contract

1. Get Unit List<br>
    - Path: **/api/unit/list**
    - Method: **Get** 
    - Authorization: **Bearer <token>**
    - Request:
    ```json
    {
        "name": "<string>",
        "sortField": "<name|createdAt>",
        "descendingOrder": <bool>,
        "limit": <int>,
        "page": <int>
    }
    ```
    - Ok Response:
    ```json
    {
        "units": [
            {
                "name": "<string>",
                "createdAt": "<timestamp>",
                "createdBy": "<string>"
            }
        ],
        "request": {
            "name": "<string>",
            "sortField": "<name|createdAt>",
            "descendingOrder": <bool>,
            "limit": <int>,
            "page": <int>
        }
    }
    ```

2. Get Unit Total<br>
    - Path: **/api/unit/total**
    - Method: **Get** 
    - Authorization: **Bearer <token>**
    - Request:
    ```json
    {
        "name": "<string>",
    }
    ```
    - Ok Response:
    ```json
    {
        "total": <int>,
        "request": {
            "name": "<string>",
            "sortField": "<name|createdAt>",
            "descendingOrder": <bool>,
            "limit": <int>,
            "page": <int>
        }
    }
    ```

3. Create Unit<br>
    Create unique unit for measurement
    - Path: **/api/unit/create**
    - Method: **Post** 
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "name": "<string>",
    }
    ```
    - Ok Response:
    ```json
    {
        "name": "<string>",
        "createdAt": "<timestamp>",
        "createdBy": "<string>"
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "name": [
                    // name field must not empty
                    {"tag": "required", "param": ""},
                    // name field max length is 30
                    {"tag": "max", "param": "30"}
                ]
            }
        }
        ```
        - Duplicate Request (http status 400)
        ```json
        {
            "code": 400006,
            "message": "duplicate unit name, `<name>` is already in use"
        }
        ```
4. Delete Unit<br>
    Delete unused measurement unit
    - Path: **/api/unit/delete**
    - Method: **Delete** 
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "name": "<string>",
    }
    ```
    - Ok Response:
    ```json
    {
        "status": "ok",
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "name": [
                    // name field must not empty
                    {"tag": "required", "param": ""},
                ]
            }
        }
        ```
        - Forbidden Request (http status 403)
        ```json
        {
            "code": 403004,
            "message": "cannot delete used units"
        }
        ```
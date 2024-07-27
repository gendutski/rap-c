# User API Contract

1. Get User By Username
    - Path: **/api/user/detail/{username}**
    - Method: **Get** 
    - Authorization: **Bearer <token>**
    - Ok Response:
    ```json
    {
        "username": "<string>",
        "fullName": "<string>",
        "email": "<string>",
        "passwordMustChange": <bool>,
        "disabled": <bool>,
        "isGuest": <bool>,
        "createdAt": "<timestamp>",
        "createdBy": "<string>",
        "updatedAt": "<timestamp>",
        "updatedBy": "<string>"
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "username": [
                    // username field must not empty
                    {"tag": "required", "param": ""},
                ]
            }
        }
        ```
        - User not found (http status 404)
        ```json
        {
            "code": 404002,
            "message": "user with `username` = `<string>` not found"
        }

2. Get User List
    List of users that will be sorted based on the fields entered in sortField(default `createdAt`) in `asc` or `desc` (if `descendingOrder` = `true`)
    - Path: **/api/user/list**
    - Method: **Get** 
    - Authorization: **Bearer token**
    - Request:
    ```json
    {
        "username": "<string>",
        "email": "<string>",
        "fullName": "<string>",
        "show": "<all|active|not active>",
        "sortField": "<username|fullName|email|createdAt|updatedAt>",
        "descendingOrder": <bool>,
        "limit": <int>,
        "page": <int>
    }
    ```
    - Ok Response:
    ```json
    {
        "users": [
            {
                "username": "<string>",
                "fullName": "<string>",
                "email": "<string>",
                "passwordMustChange": <bool>,
                "disabled": <bool>,
                "isGuest": <bool>,
                "createdAt": "<timestamp>",
                "createdBy": "<string>",
                "updatedAt": "<timestamp>",
                "updatedBy": "<string>"
            }
        ],
        "request": {
            "username": "<string>",
            "email": "<string>",
            "fullName": "<string>",
            "show": "<all|active|not active>",
            "sortField": "<username|fullName|email|createdAt|updatedAt>",
            "descendingOrder": <bool>,
            "limit": <int>,
            "page": <int>
        }
    }
    ```

3. Get User Total
    Total users based on request
    - Path: **/api/user/total**
    - Method: **Get** 
    - Authorization: **Bearer token**
    - Request:
    ```json
    {
        "username": "<string>",
        "email": "<string>",
        "fullName": "<string>",
        "show": "<all|active|not active>"
    }
    ```
    - Response:
    ```json
    {
        "total": <int>,
        "request": {
            "username": "<string>",
            "email": "<string>",
            "fullName": "<string>",
            "show": "<all|active|not active>",
            "sortField": "<username|fullName|email|role|createdAt|updatedAt>",
            "descendingOrder": <bool>,
            "limit": <int>,
            "page": <int>
        }
    }
    ```

4. Create User<br>
    Add new users, auto generate passwords, and send welcome emails to users. 
    The password must be changed by the user concerned the first time they log in
    - Path: **/api/user/create**
    - Method: **Post** 
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "username": "<string>",
        "fullName": "<string>",
        "email": "<email>"
    }
    ```
    - Ok Response:
    ```json
    {
        "username": "<string>",
        "fullName": "<string>",
        "email": "<string>",
        "passwordMustChange": <bool>,
        "disabled": <bool>,
        "isGuest": <bool>,
        "createdAt": "<timestamp>",
        "createdBy": "<string>",
        "updatedAt": "<timestamp>",
        "updatedBy": "<string>"
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "email": [
                    // email field must not empty
                    {"tag": "required", "param": ""},
                    // email field must be a valid email address
                    {"tag": "email", "param": ""}
                ],
                "fullName": [
                    // fullName field must not empty
                    {"tag": "required", "param": ""}
                ],
                "username": [
                    // username field must not empty
                    {"tag": "required", "param": ""}
                    // username field max length is 30
                    {"tag": "max", "param": "30"}
                    // username field not valid (only accept alphanumeric, period, dash, and underscore)
                    {"tag": "username", "param": ""}
                ]
            }
        }
        ```
        - User duplicate email (http status 400)
        ```json
        {
            "code": 400001,
            "message": "duplicate email, email '<email>` is already in use"
        }
        ```
        - User duplicate username (http status 400)
        ```json
        {
            "code": 400002,
            "message": "duplicate username, username '<username>` is already in use"
        }
        ```

5. Update User<br>
    Update username, fullname, email and password, from current user
    - Path: **/api/user/update**
    - Method: **Put** 
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "username": "<string>",
        "fullName": "<string>",
        "email": "<email>",
        "password": "<string>",
        "confirmPassword": "<string>"
    }
    ```
    - Ok Response:
    ```json
    {
        "username": "<string>",
        "fullName": "<string>",
        "email": "<string>",
        "passwordMustChange": <bool>,
        "disabled": <bool>,
        "isGuest": <bool>,
        "createdAt": "<timestamp>",
        "createdBy": "<string>",
        "updatedAt": "<timestamp>",
        "updatedBy": "<string>"
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "username": [
                    // username field max length is 30
                    {"tag": "max", "param": "30"}
                    // username field not valid (only accept alphanumeric, period, dash, and underscore)
                    {"tag": "username", "param": ""}
                ],
                "email": [
                    // email field must be a valid email address
                    {"tag": "email", "param": ""}
                ],
                "password": [
                    // password field must not less than 8 characters
                    {"tag": "min", "param": "8"}
                ],
                "confirmPassword": [
                    // confirmPassword field must not empty if password field not empty
                    {"tag": "required_with", "param": "Password"},
                    // confirmPassword field not match with password field
                    {"tag": "eqfield", "param": "Password"}
                ]
            }
        }
        ```

6. Enable/Disable User<br>
    Enable or disable other user
    - Path: **/api/user/active-status**
    - Method: **Put** 
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "username": "<string>",
        "disabled": <bool>
    }
    ```
    - Ok Response:
    ```json
    {
        "username": "<string>",
        "fullName": "<string>",
        "email": "<string>",
        "passwordMustChange": <bool>,
        "disabled": <bool>,
        "isGuest": <bool>,
        "createdAt": "<timestamp>",
        "createdBy": "<string>",
        "updatedAt": "<timestamp>",
        "updatedBy": "<string>"
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "username": [
                    // username field must not empty
                    {"tag": "required"}
                ],
            }
        }
        ```
        - Deactivating Inactive User (http status 400)
        ```json
        {
            "code": 400004,
            "message": "try deactivating inactive users"
        }
        ```
        - Activating Active User (http status 400)
        ```json
        {
            "code": 400005,
            "message": "try activating active users"
        }
        ```

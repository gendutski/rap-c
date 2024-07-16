# Auth API Contract

1. Login
    - Path: **/api/login**
    - Method: **Post**
    - Payload:
    ```json
    {
        "email": "<email>",
        "password": "<string>"
    }
    ```
    - Ok Response:
    ```json
    {
        "token": "<string>",
        "user": {
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
                "password": [
                    // password field must not empty
                    {"tag": "required", "param": ""}
                ]
            }
        }
        ```
        - Attempt Login Failed (http status 401)
        ```json
        {
            "code": 401001,
            "message": "wrong email or password"
        }
        ```

2. Guest Login
    - Path: **/api/guest-login**
    - Method: **Post**
    - Ok Response:
    ```json
    {
        "token": "<string>",
        "user": {
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
    }
    ```
    - Error (non internal service error) Response:
        - Attempt Guest Login Forbidden (http status 403)
        ```json
        {
            "code": 403001,
            "message": "guest login is disabled"
        }
        ```
        - Non Guest Attempt Guest Login (http status 401)
        ```json
        {
            "code": 401004,
            "message": "cannot login as guest"
        }
        ```
        - Guest login not found (http status 401)
        ```json
        {
            "code": 401001,
            "message": "wrong email or password"
        }
        ```

3. Renew Password (must change password)
    - Path: **/api/renew-password**
    - Method: **Put**
    - Authorization: Bearer **<token>**
    - Payload:
    ```json
    {
        "password": "<string>",
        "confirmPassword": "<string>"
    }
    ```
    - Ok Response:
    ```json
    {
        "status": "ok"
    }
    ```
    - Error (non internal service error) Response:
        - Validator Bad Request (http status 400)
        ```json
        {
            "code": 400999,
            "message": {
                "password": [
                    // password field must not empty
                    {"tag": "required", "param": ""}
                ],
                "confirmPassword": [
                    // confirmPassword field must not empty
                    {"tag": "required", "param": ""},
                    // confirmPassword field not match with password field
                    {"tag": "eqfield", "param": "Password"}
                ]
            }
        }
        ```
        - Unchanged password (http status 400)
        ```json
        {
            "code": 400003,
            "message": "cannot use same password"
        }
        ```

4. Request reset password (forgot password)
    - Path: **/api/request-reset**
    - Method: **Post**
    - Payload:
    ```json
    {
        "email": "<email>"
    }
    ```
    - Ok Response:
    ```json
    {
        "status": "email for request reset password has been sent"
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
                ]
            }
        }
        ```
        - Email not found (http status 404)
        ```json
        {
            "code": 404002,
            "message": "user with `email` = `<email>` not found"
        }
        ```

5. Reset password (get token from email)
    - Path: **/api/reset-password**
    - Method: **Post**
    - Payload:
    ```json
    {
        "email": "<email>",
        "token": "<string>",
        "password": "<string>",
        "confirmPassword": "<string>"
    }
    ```
    - Ok Response:
    ```json
    {
        "status": "email for request reset password has been sent"
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
                "token": [
                    // token field must not empty
                    {"tag": "required", "param": ""},
                ],
                "password": [
                    // password field must not empty
                    {"tag": "required", "param": ""}
                ],
                "confirmPassword": [
                    // confirmPassword field must not empty
                    {"tag": "required", "param": ""},
                    // confirmPassword field not match with password field
                    {"tag": "eqfield", "param": "Password"}
                ]
            }
        }
        ```
        - Reset token not found or expired (http status 404)
        ```json
        {
            "code": 404001,
            "message": "request reset password not found"
        }
        ```
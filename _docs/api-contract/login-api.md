# Login API Contract

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
    - Response:
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

# User API Contract

1. Renew Password
    - Path: **/api/user/renew-password**
    - Method: **Put**
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "password": "<email>",
        "confirmPassword": "<string>"
    }
    ```
    - Response:
    ```json
    {
        "status": "ok"
    }
    ```

2. Create User<br>
    Menambah user baru, auto generate password, dan mengirim welcome email ke user. 
    Password harus diganti user bersangkutan saat pertama kali login
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
    - Response:
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

3. Get User By Username<br>
    Mendapatkan single user berdasarkan username
    - Path: **/api/user/detail/{username}**
    - Method: **Get** 
    - Authorization: **Bearer token**
    - Response:
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

4. Get User List<br>
    Daftar user yang akan di sort berdasarkan field yang diinput di sortField(default `createdAt`) secara `asc` atau `desc` (jika `descendingOrder` = `true`)
    - Path: **/api/user/list**
    - Method: **Get** 
    - Authorization: **Bearer token**
    - Request:
    ```json
    {
        "userName": "<string>",
        "email": "<string>",
        "fullName": "<string>",
        "show": "<all|active|not active>",
        "sortField": "<username|fullName|email|createdAt|updatedAt>",
        "descendingOrder": <bool>,
        "limit": <int>,
        "page": <int>
    }
    ```
    - Response:
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
            "userName": "<string>",
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

5. Get User Total<br>
    Total user berdasarkan request
    - Path: **/api/user/total**
    - Method: **Get** 
    - Authorization: **Bearer token**
    - Request:
    ```json
    {
        "userName": "<string>",
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
            "userName": "<string>",
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
6. Update User<br>
    Update username, password, fullname dari user yang bersangkutan
    - Path: **/api/user/update**
    - Method: **Put** 
    - Authorization: **Bearer token non guest**
    - Payload:
    ```json
    {
        "username": "<string>",
        "password": "<string>",
        "fullName": "<string>"
    }
    ```
    - Response:
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
7. Enable/Disable User<br>
    Aktifkan atau non aktifkan user lain
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
    - Response:
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

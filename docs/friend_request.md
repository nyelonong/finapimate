**FINMATE User Friend Request**
----
  User friend request list.

* **URL**

  /v1/user/friends/request

* **Method:**

  `POST`

* **JSON Body**

        {
            "user_id": 123,
        }


* **Success Response:**

  * **Code:** 200 <br />
    **Content:**

        [
            {
                "friend_id": 1,
                "user_id_a": 123,
                "user_id_b": 456,
                "create_time": "2006-01-02T15:04:05Z",
                "status": 1,
                "user_profile": {
                    "email": "zaki.afrani@gmail.com",
                    "name": "Muhammad Zaki Al-Afrani",
                    "password": "6542b775e3263c27e321b929f52fc6e0",
                    "gender": 1,
                    "birth_date": "2006-01-02T15:04:05Z",
                    "nik": "12345678903545",
                    "nik_valid": 1,
                    "msisdn": "024456647647"
                    "create_time": "2006-01-02T15:04:05Z"
                }

            },
            {
                "friend_id": 2,
                "user_id_a": 123,
                "user_id_b": 7746,
                "create_time": "2006-01-02T15:04:05Z",
                "status": 1,
                "user_profile": {
                    "email": "zaki.afrani@gmail.com",
                    "name": "Muhammad Zaki Al-Afrani",
                    "password": "6542b775e3263c27e321b929f52fc6e0",
                    "gender": 1,
                    "birth_date": "2006-01-02T15:04:05Z",
                    "nik": "12345678903545",
                    "nik_valid": 1,
                    "msisdn": "024456647647"
                    "create_time": "2006-01-02T15:04:05Z"
                }

            }
        ]

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**

        {
            "message": "Invalid Body Request."
        }


* **Sample Call:**

  `curl -X POST -H "Content-Type: application/json" -d '{"user_id": 123}' http://localhost/v1/user/friends/request`

**FINMATE User Approve Friends**
----
  User approving friends.

* **URL**

  /v1/user/friends/approve

* **Method:**

  `POST`

* **JSON Body**

        [
            {
                "friend_id": 1,
                "user_id_a": 123,
                "user_id_b": 456
            },
            {
                "friend_id": 2,
                "user_id_a": 123,
                "user_id_b": 7746
            }
        ]


    * user_id_a is yours
    * user_id_b is your beloved friend

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**

        [
            {
                "friend_id": 1,
                "user_id_a": 123,
                "user_id_b": 456
            },
            {
                "friend_id": 2,
                "user_id_a": 123,
                "user_id_b": 7746
            }
        ]

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**

        {
            "message": "Invalid Body Request."
        }


* **Sample Call:**

  `curl -X POST -H "Content-Type: application/json" -d '[{"friend_id": 2, "user_id_a": 123, "user_id_b": 456}, {"friend_id": 2, "user_id_a": 123, "user_id_b": 7746}]' http://localhost/v1/user/friends/approve`

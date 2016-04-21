**FINMATE User Add Friends**
----
  User adding friends.

* **URL**

  /v1/user/addfriends

* **Method:**

  `POST`

* **JSON Body**

        [
            {
                "user_id_a": 123,
                "user_id_b": 456
            },
            {
                "user_id_a": 123,
                "user_id_b": 7746
            },
            {
                "user_id_a": 123,
                "user_id_b": 434
            }
        ]


        * user_id_a is yours
        * user_id_b is your beloved friend

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**

        [
            {
                "user_id_a": 123,
                "user_id_b": 456
            },
            {
                "user_id_a": 123,
                "user_id_b": 7746
            },
            {
                "user_id_a": 123,
                "user_id_b": 434
            }
        ]

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**

        {
            "message": "Invalid Body Request."
        }


* **Sample Call:**

  `curl -X POST -H "Content-Type: application/json" -d '[{"user_id_a": 123, "user_id_b": 456}, {"user_id_a": 123, "user_id_b": 7746}]' http://localhost/v1/user/addfriends`

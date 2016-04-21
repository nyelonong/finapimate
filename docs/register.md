**FINMATE User Register**
----
  User registration.

* **URL**

  /v1/register

* **Method:**

  `POST`

* **JSON Body**
        {
            "email": "zaki.afrani@gmail.com",
            "name": "Muhammad Zaki Al-Afrani",
            "password": "6542b775e3263c27e321b929f52fc6e0",
            "gender": 1,
            "birth_date": "2006-01-02T15:04:05Z",
            "nik": "12345678903545",
            "msisdn": "024456647647"
        }
<i>birth_date is timestamp with RFC3339 format</i>

* **Success Response:**

  * **Code:** 200 <br />
    **Content:**
        {
            "email": "zaki.afrani@gmail.com",
            "name": "Muhammad Zaki Al-Afrani",
            "password": "6542b775e3263c27e321b929f52fc6e0",
            "gender": 1,
            "birth_date": "2006-01-02T15:04:05Z07:00",
            "nik": "12345678903545",
            "msisdn": "024456647647"
        }

* **Error Response:**

  * **Code:** 400 <br />
    **Content:**

        {
            "message": "Invalid Body Request."
        }


* **Sample Call:**

  `curl -X POST -H "Content-Type: application/json" -d '{"email": "zaki.afrani@gmail.com", "name": "Muhammad Zaki Al-Afrani", "password": "6542b775e3263c27e321b929f52fc6e0", "gender": 1, "birth_date": "2006-01-02T15:04:05Z07:00", "nik": "12345678903545", "msisdn": "024456647647"}' http://localhost/v1/register`

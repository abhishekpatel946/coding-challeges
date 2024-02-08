
FAKE_USER_DATA = {
    "user_id": 1,
    "username": "johndoe",
    "password": "password123",
    "name": "John Doe",
    "email": "john.doe@example.com",
    "address": "123 Main St",
    "city": "Indore",
    "state": "MP",
    "country": "India",
    "postcode": "62701",
    "phone_number": "+1-202-555-0123",
    "birthdate": "1985-10-25",
    "gender": "Male",
    "job": "Software Engineer",
    "company": "Acme Corporation",
    "website": "http://www.example.com"
}


class MockDB:
    def __init__(self):
        pass

    def get_db():
        return FAKE_USER_DATA

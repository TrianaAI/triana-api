# Triana API

## Error Body Format
When an error occurs, the API will send an error HTTP status code along with a JSON response body. The format of the error response will be as follows:
```json
{
    "message": "Error message",
    "other data": "Additional data if needed"
}
```
Other data will be included for the following cases:
- **Validation Error**: If the request body is not valid, the API will return a 422 status code with a message indicating the validation error.
```json
{
    "message": "Error message",
    "details": {
        "Field name": "Error type",
    }
}
```
## APIs
### **POST** /register
Register a new user. The request body should contain the user's email and password. The response will include a success message and the user (id, email and name).
- **Example Request**:
```json
{
    "name": "new name",
    "email" : "new@gmail.com",
    "nationality": "indonesia",
    "dob": "2025/04/01",
    "weight": 2,
    "height": 165.6,
    "heartrate": 98.6,
    "bodytemp": 35.5
}
```

- **Example Response**:
```json
{
    "message": "User registered successfully",
    "user": {
        "id": "some-uuid",
        "name": "new name",
        "email" : "new@gmail.com"
    },
}
```

### **POST** /verify-otp
Verify the OTP sent to the user's email. To avoid complexity, send the same data as the register request along with the user inputted OTP. The response will return a success message and the new session ID.
- **Example Request**:
```json
{
    "name": "new 20000",
    "email" : "new@gmail.com",
    "nationality": "indonesia",
    "dob": "2025/04/01",
    "weight": 2,
    "height": 165.6,
    "heartrate": 98.6,
    "bodytemp": 35.5,
    "OTP": "132267"
}
```
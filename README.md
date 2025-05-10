# Triana API

Backend API for Triana, built with Go and Gin.

## Setup

1. Clone the repository:

```bash
git clone git@github.com:BeeCodingAI/triana-api.git
```

2. Navigate to the project directory:

```bash
cd triana-api
```

3. Install the required dependencies:

```bash
go mod tidy
```

4. Set up the environment variables. Create a `.env` file in the root directory and add follow .env.example.
5. Run the application:

```bash
go run main.go
```

or use `air` for live reloading during development:

```bash
air
```

6. The API will be available at `http://localhost:8080`, go to main.go to change the port.

## Error Body Format

When an error occurs, the API will send an error HTTP status code along with a JSON response body. The format of the error response will be as follows:

```json
{
  "message": "Error message",
  "other data": "Additional data if needed"
}
```

Other data will be included for the following cases:

- **Validation Error**: If the request body is not valid, the API will return a 400 status code with a message indicating the validation error.

```json
{
  "message": "Error message",
  "details": {
    "Field name": "Error type"
  }
}
```

## APIs

### **POST** /register

Register a new user. The request body should contain the user's email and password. The response will include a success message and the user (id, email and name).

- **Example Request**:

```json
{
    "name": "Mario",
    "email" : "new@gmail.com",
    "nationality": "Indonesian",
    "dob": "2004-04-01", // use YYYY-MM-DD format
    "weight": 40.2,
    "gender": "male",
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
    "name": "Mario",
    "email": "new@gmail.com"
  }
}
```

### **POST** /verify-otp

Verify the OTP sent to the user's email. To avoid complexity, send the same data as the register request along with the user inputted OTP. The response will return a success message and the new session ID.

- **Example Request**:

```json
{
  "name": "Mario",
  "email": "new@gmail.com",
  "nationality": "Indonesian",
  "dob": "2004-04-01", // use YYYY-MM-DD format
  "weight": 2,
  "gender": "male",
  "height": 165.6,
  "heartrate": 98.6,
  "bodytemp": 35.5,
  "OTP": "132267"
}
```

### **POST** /session/:id

This route is for user to send a new message to the chat session. The request body should contain only the new message. The response includes the reply from the AI, the session ID, message and the next action to be taken.

- **Example Request**:

```json
{
  "new_message": "What is my name and age?"
}
```

- **Example Response**:

```json
{
  "message": "Chat history updated successfully",
  "next_action": "CONTINUE_CHAT",
  "reply": "Hello! How can I help you today, Mario? I see you're Indonesian, 21 years, 1 month, and 4 days old. I also have your weight (40 kg), height (160 cm), heart rate (98.6 bpm), and body temperature (35.5°C). Is there anything specific you'd like to discuss?\n",
  "session_id": "464879a4-07c8-4be0-9b7d-5ee1d05d5e23"
}
```

### **POST** /session/:id/diagnose

This route allows a doctor to add a diagnosis to a specific session. The request body should contain the diagnosis string. The response confirms whether the diagnosis was saved successfully.

- **Example Request**:

```json
{
  "diagnosis": "Patient has mild fever."
}
```

- **Example Response**:

```json
{
  "message": "Diagnosis saved successfully"
}
```

### **GET** /session/:id

Get the data for a specific session. The response will include the session ID, user data, and the chat history (messages in ASC order).

- **Example Response**:

```json
{
  "id": "464879a4-07c8-4be0-9b7d-5ee1d05d5e23",
  "user_id": "89497617-d9e7-4403-bd39-cd7a2554cea0",
  "user": {
    "id": "89497617-d9e7-4403-bd39-cd7a2554cea0",
    "name": "Mario",
    "email": "new@gmail.com",
    "nationality": "Indonesian",
    "dob": "2004-04-01T00:00:00Z",
    "created_at": "2025-05-05T09:11:51.066136Z",
    "updated_at": "2025-05-05T09:18:09.772594Z",
    "sessions": null
  },
  "weight": 40,
  "height": 160,
  "heartrate": 98.6,
  "bodytemp": 35.5,
  "messages": [
    {
      "id": "0a33a470-6440-48e2-9faa-cd74d74dbd0c",
      "role": "user",
      "content": "Hello!",
      "session_id": "464879a4-07c8-4be0-9b7d-5ee1d05d5e23",
      "created_at": "2025-05-05T09:31:31.040801Z",
      "updated_at": "2025-05-05T09:31:31.040801Z"
    },
    {
      "id": "bf98b87d-f656-4996-a111-f6a3202cf074",
      "role": "triana",
      "content": "Hello! How can I help you today, Mario? I see you're Indonesian, 21 years, 1 month, and 4 days old. I also have your weight (40 kg), height (160 cm), heart rate (98.6 bpm), and body temperature (35.5°C). Is there anything specific you'd like to discuss?\n",
      "session_id": "464879a4-07c8-4be0-9b7d-5ee1d05d5e23",
      "created_at": "2025-05-05T09:31:31.041801Z",
      "updated_at": "2025-05-05T09:31:31.041801Z"
    }
  ],
  "prediagnosis": "",
  "created_at": "2025-05-05T09:18:09.76029Z",
  "updated_at": "2025-05-05T09:18:09.76029Z"
}
```

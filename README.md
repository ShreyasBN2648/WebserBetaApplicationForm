    
# Nyxborn Beta Tester Application WebServer

This is a simple web server for handling beta tester applications for the video game Nyxborn. The server provides two main endpoints and stores submitted applications in a MongoDB collection.


## Features

- Home page (/) with basic information about the beta program
- Application form (/form.html) for potential beta testers to submit their information
- Form data storage in MongoDB for later review
- Confirmation message upon successful submission


## Endpoints

- Home Page (/)
    - Returns basic information about the Nyxborn beta testing program
- Application Form (/form.html)
    - Displays the beta tester application form (as shown in the image)
    - Handles form submissions via POST request
    - Stores submitted data in MongoDB
    - Returns a thank you message upon successful submission

## Form Fields

The application form collects the following information:

- Full Name 
- Email
- Age Range (dropdown selection)
- Location
- Primary GPU
- Primary CPU

## Data Storage

Submitted applications are stored in a MongoDB collection with the following structure:

```
{
  name: String,
  email: String,
  age-Range: String,
  location: String,
  gpu: String,
  cpu: String,
  submissionDate: Date
}
```


## Installation

clone Beta Application Web Server with git bash

```bash
  git clone https://github.com/ShreyasBN2648/WebserBetaApplicationForm.git
```
```
  cd WebserBetaApplicationForm/cmd
```
```
  go run main.go
```
    
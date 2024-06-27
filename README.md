# Bingoscape

Bingoscape is an oldschool Runescape clan bingo management tool built with Go, Templ, htmx, and Tailwind CSS. This tool helps you organize and manage bingo events for your clan, providing an intuitive and interactive interface for tracking progress and results.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Introduction

Bingoscape is a streamlined tool for managing oldschool Runescape clan bingo events, emphasizing simplicity and ease of use. Designed to handle one instance per clan, Bingoscape offers two types of login: management and bingo team. Management users have persistent accounts, enabling long-term administration, while team logins are temporary, created for the duration of a single bingo event.

### Key Features

- **Simplicity and Focus**: Bingoscape is designed for ease of use, with a single instance per clan. Future updates will introduce a more comprehensive clan system.
- **Flexible Login System**: Two distinct types of logins cater to different needs. Management users have ongoing accounts, while team logins are temporary, making it easy to set up and manage events.
- **Reusable Templates**: Save time and effort by storing bingo tasks as templates. These can be reused in future bingo events, ensuring consistency and reducing setup time.
- **Submission Review**: Management users can review team submissions, accepting them or marking them as needing action. This process includes the ability to comment on submissions, providing clear feedback and communication.
- **Detailed Submission Views**: View submissions on a per-tile and per-team basis, offering a comprehensive overview of progress and performance.

## Technologies Used

- **Go**: Backend logic
- **Templ**: HTML templating
- **htmx**: For interactivity without too much JavaScript
- **Tailwind CSS**: For styling and responsive design
- **esbuild**: For javascript bundling

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- [Go](https://golang.org/dl/) (version 1.22.x+)
- [Node.js](https://nodejs.org/) and npm (for Tailwind CSS and esbuild)
- [GNU Make](https://www.gnu.org/software/make/)
- [PostgreSQL](https://www.postgresql.org/) 

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/kaffeed/bingoscape.git
   cd bingoscape
   ```

2. Install dependencies for Tailwind CSS and esbuild:
   ```bash
   npm install
   ```

3. Build the program with:
   ```bash
   make build
   ```
4. Create the Database and Database User
   First, connect to the database:
   ```bash
   psql -U postgres
   ```

   then, create the database. For this example we will use the user bingoscapeuser and the database bingoscape.
   ```sql
   create database bingoscape;
   create user bingoscapeuser with encrypted password '<password>';
   grant all privileges on database bingoscape to bingoscapeuser;
   ```

   depending on your postgres version (15+) you will also have to grant the user access to the schema public.
   While still in psql, connect to your database like so: 
   ```
   \c bingoscape
   ```
   and afterwards run 
   ```sql
   grant all on schema public to bingoscapeuser;
   ```

### Running the Application

<mark>CAVE: You need to create the database and the user on your own, the schema is migrated on the initial application startup.</mark>

1. Build the Go application:
   ```bash
   make build
   ```

2. Create a .env file in the folder you're running the application from according to the [.env.example file](.env.example):

3. Run the application:
   ```bash
   ./bin/bingoscape
   ```
4. In the same folder is a second executable, called mgmt that also requires the environment from the last step - you can use this one to create your initial management user like this:

   ```bash
   ./bin/mgmt add -u <username> -p <password>
   ```
5. Open your web browser and go to localhost:PORT as specified in HTTP_LISTEN_ADDR .env variable to see the application.

### Running a local development server

1. To run bingoscape in local development mode, execute the following make step:
   ```bash
   make dev
   ```
   This starts the local web server, watches the assets for changes and hot-reloads the application on the default templ proxy port which is usually 7331

## Contributing

Contributions are welcome! Please fork the repository and create a pull request with your changes. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions or feedback, feel free to reach out:

- GitHub: [kaffeed](https://github.com/kaffeed)

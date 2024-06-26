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
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Introduction

Bingoscape is designed for oldschool Runescape clans to create, manage, and track bingo events. It leverages modern web technologies to provide a seamless and interactive user experience.

## Features

- Create and customize bingo boards
- Progress tracking on a team / tile base
- User-friendly interface with htmx and Tailwind CSS

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

### Running the Application

1. Build the Go application:
   ```bash
   make build
   ```

2. Create a .env file in bin/ according to the [.env.example file](.env.example):

3. Run the application:
   ```bash
   ./bin/bingoscape
   ```

4. Open your web browser and go to localhost:<PORT> to see the application.

## Usage

- Create a bingo event: Fill out the details of your bingo event including the number of participants, board size, and event duration.
- Manage bingo boards: Customize the bingo boards by adding, editing, or removing tasks.
- Track progress: View real-time updates of the participants' progress.

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

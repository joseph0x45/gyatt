# The GyATT stack
Go (y) Alpine Templ Tailwind

## How to use
Just download the script and run it. Or you can also pipe it directly to sh
```sh
curl https://raw.githubusercontent.com/joseph0x45/gyatt/refs/heads/main/gyatt.sh | sh -s <project_name>
```

## What it does
The script creates an opinionated folder structure.
- components package: Contains templ components
- handlers package: Contains handlers for HTTP requests
- db package: Contains all the code handling interaction with a database 
- models package: Contains domain specific types
- static folder: Contains static files such as CSS, and Alpine library. This folder's content is embedded into the binary.
- Makefile: Contains basic targets for building the app

## Dependencies
- TailwindCSS CLI from the [releases](https://github.com/tailwindlabs/tailwindcss/releases/tag/v4.1.14) 
